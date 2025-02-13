package subroutines

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/openmfp/golang-commons/controller/lifecycle"
	"github.com/openmfp/golang-commons/errors"
	"github.com/openmfp/golang-commons/logger"
	"k8s.io/apimachinery/pkg/api/meta"
	apimachinery "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/openmfp/extension-manager-operator/api/v1alpha1"
	"github.com/openmfp/extension-manager-operator/pkg/transformer"
	"github.com/openmfp/extension-manager-operator/pkg/validation"
)

const (
	ContentConfigurationSubroutineName = "ContentConfigurationSubroutine"
	ValidationConditionType            = "Valid"
	ValidationConditionReasonSuccess   = "ValidationSucceeded"
	ValidationConditionReasonFailed    = "ValidationFailed"
	ConditionStatusTrue                = "True"
	ConditionStatusFalse               = "False"
)

type ContentConfigurationSubroutine struct {
	client      *http.Client
	validator   validation.ExtensionConfiguration
	transformer []transformer.ContentConfigurationTransformer
}

func NewContentConfigurationSubroutine(validator validation.ExtensionConfiguration,
	client *http.Client) *ContentConfigurationSubroutine {

	transformers := []transformer.ContentConfigurationTransformer{
		&transformer.UrlSuffixTransformer{},
	}
	return &ContentConfigurationSubroutine{
		client:      client,
		validator:   validator,
		transformer: transformers,
	}
}

func (r *ContentConfigurationSubroutine) GetName() string {
	return ContentConfigurationSubroutineName
}

func (r *ContentConfigurationSubroutine) Finalize(_ context.Context, _ lifecycle.RuntimeObject) (ctrl.Result, errors.OperatorError) {
	return ctrl.Result{}, nil
}

func (r *ContentConfigurationSubroutine) Finalizers() []string {
	return []string{}
}

func (r *ContentConfigurationSubroutine) Process(
	ctx context.Context, runtimeObj lifecycle.RuntimeObject,
) (ctrl.Result, errors.OperatorError) {
	log := logger.LoadLoggerFromContext(ctx)

	instance := runtimeObj.(*v1alpha1.ContentConfiguration)

	log.Debug().Str("name", instance.Name).Msg("processing content configuration")

	// Download or Retrieve ContentConfiguration Json
	contentType, rawConfig, err := r.retrieveContentConfigurationData(instance, log)
	if err != nil {
		return ctrl.Result{}, errors.NewOperatorError(err, true, false)
	}

	// Validate ContentConfiguration Json
	validatedConfig, valErr := r.validator.Validate(rawConfig, contentType)
	if valErr != nil && valErr.Len() > 0 {
		log.Err(valErr).Msg("failed to validate configuration")
		condition := apimachinery.Condition{
			Type:    ValidationConditionType,
			Status:  ConditionStatusFalse,
			Reason:  ValidationConditionReasonFailed,
			Message: valErr.Error(),
		}
		meta.SetStatusCondition(&instance.Status.Conditions, condition)
		return ctrl.Result{}, nil
	}

	condition := apimachinery.Condition{
		Type:    ValidationConditionType,
		Status:  ConditionStatusTrue,
		Reason:  ValidationConditionReasonSuccess,
		Message: "OK",
	}
	meta.SetStatusCondition(&instance.Status.Conditions, condition)

	// Transform ContentConfiguration Json
	contentConfiguration := &validation.ContentConfiguration{}
	err = json.Unmarshal([]byte(validatedConfig), contentConfiguration)
	if err != nil {
		return ctrl.Result{}, errors.NewOperatorError(errors.Wrap(err, "failed to unmarshal contentConfiguration"), false, false)
	}
	for _, configurationTransformer := range r.transformer {
		err := configurationTransformer.Transform(contentConfiguration, instance)
		if err != nil {
			return ctrl.Result{}, errors.NewOperatorError(errors.Wrap(err, "failed to transform contentConfiguration"), false, true)
		}
	}

	validatedConfigBytes, err := json.Marshal(contentConfiguration)
	if err != nil {
		return ctrl.Result{}, errors.NewOperatorError(errors.Wrap(err, "failed to marshal contentConfiguration"), false, false)
	}
	validatedConfig = string(validatedConfigBytes)

	// Store resulting configuration in the status
	instance.Status.ConfigurationResult = validatedConfig
	return ctrl.Result{}, nil
}

func (r *ContentConfigurationSubroutine) retrieveContentConfigurationData(instance *v1alpha1.ContentConfiguration, log *logger.Logger) (string, []byte, error) {
	var contentType string
	var rawConfig []byte
	// InlineConfiguration has higher priority than RemoteConfiguration
	switch {
	case instance.Spec.InlineConfiguration != nil && instance.Spec.InlineConfiguration.Content != "":
		contentType = instance.Spec.InlineConfiguration.ContentType
		rawConfig = []byte(instance.Spec.InlineConfiguration.Content)
	case instance.Spec.RemoteConfiguration != nil && instance.Spec.RemoteConfiguration.URL != "":
		url := instance.Spec.RemoteConfiguration.URL
		if instance.Spec.RemoteConfiguration.InternalUrl != "" {
			url = instance.Spec.RemoteConfiguration.InternalUrl
		}
		bytes, err := r.getRemoteConfig(url, log)
		if err != nil {
			log.Err(err).Msg("failed to fetch remote configuration")

			return "", nil, err
		}
		log.Info().Msg("fetched remote configuration")
		contentType = instance.Spec.RemoteConfiguration.ContentType
		rawConfig = bytes
	default:
		return "", nil, errors.New("no configuration provided")
	}
	return contentType, rawConfig, nil
}

// Do makes an HTTP request to the specified URL.
func (r *ContentConfigurationSubroutine) getRemoteConfig(url string, log *logger.Logger) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			log.Err(closeErr).Msg("failed to close response body")
		}
	}()

	if resp.StatusCode != http.StatusOK {
		// Give the caller signal to retry if we have 5xx status codes
		if resp.StatusCode >= http.StatusInternalServerError && resp.StatusCode < 600 {
			return nil, fmt.Errorf("received non-200 status code: %d", resp.StatusCode)
		}

		return nil, fmt.Errorf("received non-200 status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}
	// TODO
	// we need to check the size of the received body before loading it to memory.
	// In case it exceeds a certain size we should reject it.
	// https://github.com/openmfp/extension-manager-operator/pull/23#discussion_r1622598363

	return body, nil
}
