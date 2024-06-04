package subroutines

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/openmfp/extension-content-operator/api/v1alpha1"
	"github.com/openmfp/extension-content-operator/pkg/validation"
	"github.com/openmfp/golang-commons/controller/lifecycle"
	"github.com/openmfp/golang-commons/errors"
	"github.com/openmfp/golang-commons/logger"
)

const (
	ContentConfigurationSubroutineName = "ContentConfigurationSubroutine"
)

type ContentConfigurationSubroutine struct {
	validator validation.ContentConfigurationInterface
}

func NewContentConfigurationSubroutine() *ContentConfigurationSubroutine {
	return &ContentConfigurationSubroutine{
		validator: validation.NewContentConfiguration(),
	}
}

func (r *ContentConfigurationSubroutine) WithValidator(validator validation.ContentConfigurationInterface) {
	r.validator = validator
}

func (r *ContentConfigurationSubroutine) GetName() string {
	return ContentConfigurationSubroutineName
}

func (r *ContentConfigurationSubroutine) Finalize(
	ctx context.Context,
	runtimeObj lifecycle.RuntimeObject) (ctrl.Result, errors.OperatorError) {
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

	var rawConfig []byte
	var contentType string
	// InlineConfiguration has higher priority than RemoteConfiguration
	switch {
	case instance.Spec.InlineConfiguration.Content != "":
		rawConfig = []byte(instance.Spec.InlineConfiguration.Content)
		contentType = instance.Spec.InlineConfiguration.ContentType
	case instance.Spec.RemoteConfiguration.URL != "":
		bytes, err, retry := getRemoteConfig(instance.Spec.RemoteConfiguration.URL)
		if err != nil {
			log.Err(err).Msg("failed to fetch remote configuration")

			return ctrl.Result{}, errors.NewOperatorError(err, retry, true)
		}
		rawConfig = bytes
		contentType = instance.Spec.RemoteConfiguration.ContentType
	default:
		return ctrl.Result{}, errors.NewOperatorError(errors.New("no configuration provided"), false, true)
	}

	schema, err := getSchema()
	if err != nil {
		log.Err(err).Msg("failed to read schema")

		return ctrl.Result{}, errors.NewOperatorError(err, false, true)
	}

	configResult, err := r.validator.Validate(schema, rawConfig, contentType)
	if err != nil {
		log.Err(err).Msg("failed to validate configuration")

		return ctrl.Result{}, errors.NewOperatorError(err, false, true)
	}

	instance.Status.ConfigurationResult = configResult

	return ctrl.Result{}, nil
}

// Do makes an HTTP request to the specified URL.
func getRemoteConfig(url string) (res []byte, err error, retry bool) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err), false
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err), false
	}
	defer resp.Body.Close() // nolint: errcheck

	if resp.StatusCode != http.StatusOK {
		// Give the caller signal to retry if we have 5xx status codes
		if resp.StatusCode >= http.StatusInternalServerError && resp.StatusCode < 600 {
			return nil, fmt.Errorf("received non-200 status code: %d", resp.StatusCode), true
		}

		return nil, fmt.Errorf("received non-200 status code: %d", resp.StatusCode), false
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err), false
	}

	return body, nil, false
}

func getSchema() ([]byte, error) {
	fileContent, err := os.ReadFile("../validation/example_schema.json")
	if err != nil {
		return nil, err
	}

	return fileContent, nil
}
