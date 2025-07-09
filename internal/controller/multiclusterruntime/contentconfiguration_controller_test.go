package multiclusterruntime

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/kcp-dev/multicluster-provider/apiexport"
	openmfpconfig "github.com/platform-mesh/golang-commons/config"
	"github.com/platform-mesh/golang-commons/logger"
	"golang.org/x/sync/errgroup"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	ctrl "sigs.k8s.io/controller-runtime"
	mcmanager "sigs.k8s.io/multicluster-runtime/pkg/manager"
	"sigs.k8s.io/yaml"

	apisv1alpha1 "github.com/kcp-dev/kcp/sdk/apis/apis/v1alpha1"
	"github.com/kcp-dev/kcp/sdk/apis/core"
	tenancyv1alpha1 "github.com/kcp-dev/kcp/sdk/apis/tenancy/v1alpha1"
	topologyv1alpha1 "github.com/kcp-dev/kcp/sdk/apis/topology/v1alpha1"
	"github.com/kcp-dev/logicalcluster/v3"
	clusterclient "github.com/kcp-dev/multicluster-provider/client"
	"github.com/kcp-dev/multicluster-provider/envtest"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"

	"github.com/openmfp/extension-manager-operator/api/v1alpha1"
	"github.com/openmfp/extension-manager-operator/internal/config"
)

const (
	defaultTestTimeout  = 15 * time.Second
	defaultTickInterval = 250 * time.Millisecond
	defaultNamespace    = "default"
)

var (
	env       *envtest.Environment
	kcpConfig *rest.Config
)

type ContentConfigurationTestSuite struct {
	suite.Suite
	cli                clusterclient.ClusterClient
	provider, consumer logicalcluster.Path
	consumerWS         *tenancyv1alpha1.Workspace
	ctx                context.Context
	cancel             context.CancelFunc
	g                  *errgroup.Group
}

func init() {
	runtime.Must(v1alpha1.AddToScheme(scheme.Scheme))
	runtime.Must(apisv1alpha1.AddToScheme(scheme.Scheme))
	runtime.Must(tenancyv1alpha1.AddToScheme(scheme.Scheme))
	runtime.Must(topologyv1alpha1.AddToScheme(scheme.Scheme))

}

func (suite *ContentConfigurationTestSuite) SetupSuite() {
	logConfig := logger.DefaultConfig()
	logConfig.NoJSON = true
	logConfig.Name = "ContentConfigurationControllerTestSuite"
	log, err := logger.New(logConfig)
	ctrl.SetLogger(log.Logr())
	suite.Require().NoError(err, "failed to create logger %v", err)
	suite.ctx, suite.cancel = context.WithCancel(context.Background())
	// Prevent the metrics listener being created
	metricsserver.DefaultBindAddress = "0"

	env = &envtest.Environment{}
	env.BinaryAssetsDirectory = "../../../bin"
	err = os.Setenv("PRESERVE", "true")
	suite.Require().NoError(err, "failed to set PRESERVE environment variable")
	kcpConfig, err = env.Start()
	suite.Require().NoError(err, "failed to start envtest environment")

	suite.cli, err = clusterclient.New(kcpConfig, client.Options{})
	suite.Require().NoError(err, "failed to create cluster client")
	_, suite.provider = envtest.NewWorkspaceFixture(suite.T(), suite.cli, core.RootCluster.Path(), envtest.WithNamePrefix("provider"))
	suite.consumerWS, suite.consumer = envtest.NewWorkspaceFixture(suite.T(), suite.cli, core.RootCluster.Path(), envtest.WithNamePrefix("consumer"))

	// Prepare apiexports and resource schema
	suite.loadFromFile("../../../test/setup/apiresourceschema-contentconfigurations.core.openmfp.io.yaml", suite.provider)
	suite.loadFromFile("../../../test/setup/apiexport-core.openmfp.io.yaml", suite.provider)

	// Create apiexportendpointslice
	aes := &apisv1alpha1.APIExportEndpointSlice{
		ObjectMeta: metav1.ObjectMeta{
			Name: "core.openmfp.io",
		},
		Spec: apisv1alpha1.APIExportEndpointSliceSpec{
			APIExport: apisv1alpha1.ExportBindingReference{
				Name: "core.openmfp.io",
				Path: suite.provider.String(),
			},
		},
	}
	suite.cli.Cluster(suite.provider).Create(suite.ctx, aes) //nolint:errcheck

	ab := &apisv1alpha1.APIBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: "core.openmfp.io",
		},
		Spec: apisv1alpha1.APIBindingSpec{
			Reference: apisv1alpha1.BindingReference{
				Export: &apisv1alpha1.ExportBindingReference{
					Name: "core.openmfp.io",
					Path: suite.provider.String(),
				},
			},
		},
	}
	err = suite.cli.Cluster(suite.consumer).Create(suite.ctx, ab)
	suite.Require().NoError(err, "failed to create APIBinding for core.openmfp.io in consumer workspace")

	suite.Eventually(func() bool {
		getErr := suite.cli.Cluster(suite.consumer).Get(suite.ctx, types.NamespacedName{Name: "core.openmfp.io"}, ab)
		return getErr == nil && ab.Status.Phase == apisv1alpha1.APIBindingPhaseBound
	}, 10*time.Second, 100*time.Millisecond, "APIBinding for core.openmfp.io in consumer workspace did not become ready")

	// lookup api export
	err = suite.cli.Cluster(suite.provider).Get(suite.ctx, types.NamespacedName{Name: "core.openmfp.io"}, aes)
	suite.Require().NoError(err, "failed to get APIExport for core.openmfp.io in consumer workspace")

	cfg := rest.CopyConfig(kcpConfig)
	cfg.Host = aes.Status.APIExportEndpoints[0].URL
	provider, err := apiexport.New(cfg, apiexport.Options{})
	suite.Require().NoError(err, "failed to create APIExport client for core.openmfp.io in consumer workspace")

	mgr, err := mcmanager.New(cfg, provider, mcmanager.Options{Logger: log.Logr()})
	suite.Require().NoError(err, "failed to create APIExport client for core.openmfp.io in consumer workspace")

	operatorCfg := config.OperatorConfig{}
	operatorCfg.Subroutines.ContentConfiguration.Enabled = true
	rec := NewContentConfigurationReconciler(log, mgr, operatorCfg)

	err = rec.SetupWithManager(mgr, &openmfpconfig.CommonServiceConfig{}, log)
	suite.Require().NoError(err, "failed to setup ContentConfiguration reconciler with manager")

	var groupContext context.Context
	suite.g, groupContext = errgroup.WithContext(suite.ctx)
	suite.g.Go(func() error {
		return provider.Run(groupContext, mgr)
	})
	suite.g.Go(func() error {
		return mgr.Start(groupContext)
	})
}

func (suite *ContentConfigurationTestSuite) loadFromFile(filePath string, workspace logicalcluster.Path) {
	data, err := os.ReadFile(filePath)
	require.NoError(suite.T(), err, "failed to read file %s", filePath)

	var u unstructured.Unstructured
	err = yaml.Unmarshal(data, &u.Object)
	require.NoError(suite.T(), err, "failed to unmarshal file %s", filePath)

	err = suite.cli.Cluster(workspace).Create(suite.ctx, &u)
	require.NoError(suite.T(), err, "failed to create resource %s", filePath)
}

func (suite *ContentConfigurationTestSuite) TearDownSuite() {
	suite.cancel()
	suite.g.Wait() //nolint:errcheck
	env.Stop()     //nolint:errcheck
}

func (suite *ContentConfigurationTestSuite) TestProcessContentConfiguration() {

	//Given
	var err error
	testContext := context.Background()
	name := "example-content-configuration"
	cc := &v1alpha1.ContentConfiguration{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: v1alpha1.ContentConfigurationSpec{
			InlineConfiguration: &v1alpha1.InlineConfiguration{
				ContentType: "json",
				Content: `{
								"name": "overview",
								"creationTimestamp": "2022-05-17T11:37:17Z",
								"luigiConfigFragment": {
								  "data": {
									"nodes": [
									  {
										"entityType": "global",
										"pathSegment": "home",
										"hideFromNav": false,
										"order": 1,
										"label": "Home",
										"icon": "account",
										"defineEntity": {
										  "id": "main"
										},
										"children": [
										  {
											"pathSegment": "overview",
											"label": "Overview",
											"icon": "home",
											"defineEntity": {
											  "id": "overview"
											},
											"compound": {
											  "renderer": {
												"use": "grid",
												"config": {
												  "columns": "1fr 1fr 1fr 1fr"
												}
											  }
											}
										  }
										]
									  }
									]
								  }
								}
							  }`,
			},
		},
	}

	// When
	err = suite.cli.Cluster(suite.consumer).Create(testContext, cc)
	suite.Require().NoError(err)

	// Wait for workspace creation and ready
	updatedCC := &v1alpha1.ContentConfiguration{}
	suite.Assert().Eventually(func() bool {
		err := suite.cli.Cluster(suite.consumer).Get(testContext, types.NamespacedName{Name: name}, updatedCC)
		readyCondition := meta.FindStatusCondition(updatedCC.Status.Conditions, "Ready")
		return err == nil && readyCondition != nil && readyCondition.Status == metav1.ConditionTrue
	}, defaultTestTimeout, defaultTickInterval)

	suite.verifyCondition(updatedCC.Status.Conditions, "Ready", metav1.ConditionTrue, "Complete")
	suite.verifyCondition(updatedCC.Status.Conditions, "Valid", metav1.ConditionTrue, "ValidationSucceeded")
}
func TestContentConfigurationSuite(t *testing.T) {
	suite.Run(t, new(ContentConfigurationTestSuite))
}

func (suite *ContentConfigurationTestSuite) verifyCondition(conditions []metav1.Condition, conditionType string, status metav1.ConditionStatus, reason string) {
	condition := getCondition(conditions, conditionType)
	suite.Require().NotNil(condition)
	suite.Equal(status, condition.Status)
	suite.Equal(reason, condition.Reason)
}

func getCondition(conditions []metav1.Condition, conditionType string) *metav1.Condition {
	for _, condition := range conditions {
		if condition.Type == conditionType {
			return &condition
		}
	}
	return nil
}
