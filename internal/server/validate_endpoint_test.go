package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/openmfp/extension-content-operator/pkg/validation"
	"github.com/openmfp/golang-commons/logger"
	"github.com/stretchr/testify/assert"
)

type responseError struct {
	ValidationErrors []validationError `json:"validationErrors,omitempty"`
}

type responseSuccess struct {
	ParsedConfiguration string `json:"parsedConfiguration,omitempty"`
}

func TestHandlerValidate_Error(t *testing.T) {

	logcfg := logger.DefaultConfig()
	log, _ := logger.New(logcfg)

	handler := NewHttpValidateHandler(log, validation.NewContentConfiguration())

	reqBody := `{
	"contentType": "json",
	"contentConfiguration":"{\"luigiConfigFragment2\": {\"data\": {\"nodeDefaults\": {\"entityType\": \"global\",\"isolateView\": true},\"nodes\": [{\"entityType\": \"global\",\"icon\": \"home\",\"label\": \"Overview\",\"pathSegment\": \"home\"}],\"texts\": [{\"locale\": \"de\",\"textDictionary\": {\"hello\": \"Hallo\"}}]}},\"name\": \"overview\"}"}"
	}`
	req := httptest.NewRequest(http.MethodPost, "/validate", bytes.NewBufferString(reqBody))
	w := httptest.NewRecorder()

	handler.HandlerValidate(w, req)

	resp := w.Result()

	r := &responseError{}
	decoder := json.NewDecoder(resp.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(r)
	defer resp.Body.Close()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.GreaterOrEqual(t, len(r.ValidationErrors), 1)
}

func TestHandlerValidate_Success(t *testing.T) {

	logcfg := logger.DefaultConfig()
	log, _ := logger.New(logcfg)

	handler := NewHttpValidateHandler(log, validation.NewContentConfiguration())

	reqBody := `{
            "contentType": "json",
			"contentConfiguration":"{\"luigiConfigFragment\": {\"data\": {\"nodeDefaults\": {\"entityType\": \"global\",\"isolateView\": true},\"nodes\": [{\"entityType\": \"global\",\"icon\": \"home\",\"label\": \"Overview\",\"pathSegment\": \"home\"}],\"texts\": [{\"locale\": \"de\",\"textDictionary\": {\"hello\": \"Hallo\"}}]}},\"name\": \"overview\"}"}"
}`
	req := httptest.NewRequest(http.MethodPost, "/validate", bytes.NewBufferString(reqBody))
	w := httptest.NewRecorder()

	handler.HandlerValidate(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	decoder.DisallowUnknownFields()

	r := &responseSuccess{}
	decoder.DisallowUnknownFields()
	err := decoder.Decode(r)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.GreaterOrEqual(t, len(r.ParsedConfiguration), 0)
}

func TestYAML_Success(t *testing.T) {

	logcfg := logger.DefaultConfig()
	log, _ := logger.New(logcfg)

	handler := NewHttpValidateHandler(log, validation.NewContentConfiguration())

	reqBody := `{
            "contentType": "yaml",
            "contentConfiguration": "contentType: json\nluigiConfigFragment:\n  data:\n    nodes:\n    - dxpOrder: 6\n      entityType: global\n      hideSideNav: true\n      icon: business-one\n      label: '{{catalog}}'\n      order: 6\n      pathSegment: catalog\n      showBreadcrumbs: false\n      tabNav: true\n      urlSuffix: /#/global-catalog\n      visibleForFeatureToggles:\n      - '!global-catalog'\n    - dxpOrder: 6\n      entityType: global\n      hideSideNav: true\n      icon: business-one\n      label: '{{catalog}}'\n      order: 6\n      pathSegment: catalog\n      showBreadcrumbs: false\n      tabNav: true\n      urlSuffix: /#/new-global-catalog\n      visibleForFeatureToggles:\n      - global-catalog\n    - children:\n      - context:\n          extClassName: :extClassName\n        hideFromNav: true\n        pathSegment: :extClassName\n        urlSuffix: /#/extensions/:extClassName\n      entityType: global\n      hideFromNav: true\n      label: '{{extensions}}'\n      pathSegment: extensions\n    targetAppConfig:\n      _version: 1.13.0\n      sap.integration:\n        navMode: inplace\n        urlTemplateId: urltemplate.url\n        urlTemplateParams:\n          query: {}\n    texts:\n    - locale: \"\"\n      textDictionary:\n        catalog: Catalog\n        extensions: Extensions\n    - locale: en\n      textDictionary:\n        catalog: Catalog\n        extensions: Extensions\n    - locale: de\n      textDictionary:\n        catalog: Katalog\n        extensions: Erweiterungen\n    viewGroup:\n      preloadSuffix: /#/preload\nname: extension-manager\n"
        }`
	req := httptest.NewRequest(http.MethodPost, "/validate", bytes.NewBufferString(reqBody))
	w := httptest.NewRecorder()

	handler.HandlerValidate(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	decoder.DisallowUnknownFields()
	r := &responseSuccess{}
	err := decoder.Decode(r)
	assert.Nil(t, err)
	assert.Greater(t, len(r.ParsedConfiguration), 0)

	decoder = json.NewDecoder(resp.Body)
	decoder.DisallowUnknownFields()
	re := &responseError{}
	err = decoder.Decode(re)
	assert.NotNil(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, len(re.ValidationErrors), 0)
}

func TestYAML_FailureWrongType(t *testing.T) {

	logcfg := logger.DefaultConfig()
	log, _ := logger.New(logcfg)

	handler := NewHttpValidateHandler(log, validation.NewContentConfiguration())

	reqBody := `{
            "contentType": "json",
            "contentConfiguration": "contentType: json\nluigiConfigFragment:\n  data:\n    nodes:\n    - dxpOrder: 6\n      entityType: global\n      hideSideNav: true\n      icon: business-one\n      label: '{{catalog}}'\n      order: 6\n      pathSegment: catalog\n      showBreadcrumbs: false\n      tabNav: true\n      urlSuffix: /#/global-catalog\n      visibleForFeatureToggles:\n      - '!global-catalog'\n    - dxpOrder: 6\n      entityType: global\n      hideSideNav: true\n      icon: business-one\n      label: '{{catalog}}'\n      order: 6\n      pathSegment: catalog\n      showBreadcrumbs: false\n      tabNav: true\n      urlSuffix: /#/new-global-catalog\n      visibleForFeatureToggles:\n      - global-catalog\n    - children:\n      - context:\n          extClassName: :extClassName\n        hideFromNav: true\n        pathSegment: :extClassName\n        urlSuffix: /#/extensions/:extClassName\n      entityType: global\n      hideFromNav: true\n      label: '{{extensions}}'\n      pathSegment: extensions\n    targetAppConfig:\n      _version: 1.13.0\n      sap.integration:\n        navMode: inplace\n        urlTemplateId: urltemplate.url\n        urlTemplateParams:\n          query: {}\n    texts:\n    - locale: \"\"\n      textDictionary:\n        catalog: Catalog\n        extensions: Extensions\n    - locale: en\n      textDictionary:\n        catalog: Catalog\n        extensions: Extensions\n    - locale: de\n      textDictionary:\n        catalog: Katalog\n        extensions: Erweiterungen\n    viewGroup:\n      preloadSuffix: /#/preload\nname: extension-manager\n"
        }`
	req := httptest.NewRequest(http.MethodPost, "/validate", bytes.NewBufferString(reqBody))
	w := httptest.NewRecorder()

	handler.HandlerValidate(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	decoder.DisallowUnknownFields()
	re := &responseError{}
	err := decoder.Decode(re)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.GreaterOrEqual(t, len(re.ValidationErrors), 1)
}
