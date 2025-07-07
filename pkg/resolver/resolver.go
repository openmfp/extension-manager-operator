package resolver

// THIS CODE WILL BE UPDATED WITH SCHEMA CHANGES. PREVIOUS IMPLEMENTATION FOR SCHEMA CHANGES WILL BE KEPT IN THE COMMENT SECTION. IMPLEMENTATION FOR UNCHANGED SCHEMA WILL BE KEPT.

import (
	"context"

	kcpv1alpha1 "github.com/kcp-dev/kcp/sdk/apis/apis/v1alpha1"
	corev1alpha1 "github.com/openmfp/extension-manager-operator/api/v1alpha1"
	"github.com/openmfp/extension-manager-operator/pkg/graph"
	"github.com/openmfp/extension-manager-operator/pkg/graph/model"
	"sigs.k8s.io/controller-runtime/pkg/client"
	mcmanager "sigs.k8s.io/multicluster-runtime/pkg/manager"
)

type Resolver struct {
	mgr                 mcmanager.Manager
	providerWorkspaceID string
}

var (
	LabelContentFor = "ui.platform-mesh.io/content-for"
	LabelEntity     = "ui.platform-mesh.ui/entity"
)

// ContentConfigurations is the resolver for the contentConfigurations field.
func (r *queryResolver) ContentConfigurations(ctx context.Context, path string) ([]*model.ContentConfiguration, error) {
	// TODO: exchange path with logicalcluster.Name
	cluster, err := r.mgr.GetCluster(ctx, path)
	if err != nil {
		return nil, err
	}

	clusterClient := cluster.GetClient()

	// TODO: get the proper entity type
	// collect all ContentConfigurations from the current workspace
	var contentConfigurations corev1alpha1.ContentConfigurationList
	err = clusterClient.List(ctx, &contentConfigurations, client.MatchingLabels{
		LabelEntity: "account",
	})
	if err != nil {
		return nil, err
	}

	var apiBindings kcpv1alpha1.APIBindingList
	err = clusterClient.List(ctx, &apiBindings)
	if err != nil {
		return nil, err
	}

	for _, binding := range apiBindings.Items {
		// TODO: exchange path with logicalcluster.Name
		exportWs, err := r.mgr.GetCluster(ctx, binding.Spec.Reference.Export.Path)
		if err != nil {
			return nil, err
		}

		var exportContentConfigs corev1alpha1.ContentConfigurationList
		err = exportWs.GetClient().List(ctx, &exportContentConfigs, client.MatchingLabels{
			LabelContentFor: binding.Spec.Reference.Export.Name,
			LabelEntity:     "account",
		})
		if err != nil {
			return nil, err
		}

		contentConfigurations.Items = append(contentConfigurations.Items, exportContentConfigs.Items...)
	}

	providerWs, err := r.mgr.GetCluster(ctx, r.providerWorkspaceID)
	if err != nil {
		return nil, err
	}

	var providerContentConfigs corev1alpha1.ContentConfigurationList
	err = providerWs.GetClient().List(ctx, &providerContentConfigs, client.MatchingLabels{
		LabelEntity: "account",
	})
	if err != nil {
		return nil, err
	}
	contentConfigurations.Items = append(contentConfigurations.Items, providerContentConfigs.Items...)

	var result []*model.ContentConfiguration
	for _, cc := range contentConfigurations.Items {
		mc := &model.ContentConfiguration{}
		if mc.Spec.RemoteConfiguration != nil {
			mc.Spec.RemoteConfiguration = &model.RemoteConfiguration{
				ContentType: cc.Spec.RemoteConfiguration.ContentType,
				URL:         cc.Spec.RemoteConfiguration.URL,
				InternalURL: &cc.Spec.RemoteConfiguration.InternalUrl,
				Authentication: &model.Authentication{
					Type: &cc.Spec.RemoteConfiguration.Authentication.Type,
					SecretRef: &model.LocalObjectReference{
						Name: cc.Spec.RemoteConfiguration.Authentication.SecretRef.Name,
					},
				},
			}
		}

		if cc.Spec.InlineConfiguration != nil {
			mc.Spec.InlineConfiguration = &model.InlineConfiguration{
				ContentType: cc.Spec.InlineConfiguration.ContentType,
				Content:     cc.Spec.InlineConfiguration.Content,
			}
		}

		result = append(result, mc)
	}

	return result, nil
}

// Query returns graph.QueryResolver implementation.
func (r *Resolver) Query() graph.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
