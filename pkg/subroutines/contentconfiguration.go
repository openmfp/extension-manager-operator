package subroutines

import (
	"context"
	"fmt"
	"io"
	"net/http"

	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/openmfp/extension-content-operator/api/v1alpha1"
	"github.com/openmfp/golang-commons/controller/lifecycle"
	"github.com/openmfp/golang-commons/errors"
	"github.com/openmfp/golang-commons/logger"
)

const (
	ContentConfigurationSubroutineName = "ContentConfigurationSubroutine"
)

type ContentConfigurationSubroutine struct {
	client *http.Client
}

func NewContentConfigurationSubroutine() *ContentConfigurationSubroutine {
	return &ContentConfigurationSubroutine{
		client: http.DefaultClient,
	}
}

func (r *ContentConfigurationSubroutine) WithClient(client *http.Client) {
	r.client = client
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
	// InlineConfiguration has higher priority than RemoteConfiguration
	switch {
	case instance.Spec.InlineConfiguration.Content != "":
		rawConfig = []byte(instance.Spec.InlineConfiguration.Content)
	case instance.Spec.RemoteConfiguration.URL != "":
		bytes, err, retry := r.getRemoteConfig(instance.Spec.RemoteConfiguration.URL)
		if err != nil {
			log.Err(err).Msg("failed to fetch remote configuration")

			return ctrl.Result{}, errors.NewOperatorError(err, retry, true)
		}
		rawConfig = bytes
	default:
		return ctrl.Result{}, errors.NewOperatorError(errors.New("no configuration provided"), false, true)
	}

	// TODO replace it with validation function
	validatedConfig := string(rawConfig)

	instance.Status.ConfigurationResult = validatedConfig

	return ctrl.Result{}, nil
}

// Do makes an HTTP request to the specified URL.
func (r *ContentConfigurationSubroutine) getRemoteConfig(url string) (res []byte, err error, retry bool) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err), false
	}

	resp, err := r.client.Do(req)
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

	// TODO
	// we need to check the size of the received body before loading it to memory.
	// In case it exceeds a certain size we should reject it.
	// https://github.com/openmfp/extension-content-operator/pull/23#discussion_r1622598363

	return body, nil, false
}
