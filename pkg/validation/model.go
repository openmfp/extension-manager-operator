package validation

//go:generate go run schema/genschema.go
type ContentConfiguration struct {
	Name                string              `json:"name,omitempty" yaml:"name,omitempty" jsonschema:"oneof_required=string"`
	ContentType         string              `json:"contentType,omitempty" yaml:"contentType,omitempty"`
	LuigiConfigFragment LuigiConfigFragment `json:"luigiConfigFragment" yaml:"luigiConfigFragment"`
}

type LuigiConfigFragment struct {
	Data      LuigiConfigData `json:"data,omitempty" yaml:"data,omitempty" jsonschema:"oneof_required=object"`
	ViewGroup ViewGroup       `json:"viewGroup,omitempty" yaml:"viewGroup,omitempty"`
}

type ViewGroup struct {
	PreloadSuffix string `json:"preloadSuffix,omitempty" yaml:"preloadSuffix,omitempty"`
}

type LuigiConfigData struct {
	NodeDefaults    NodeDefaults    `json:"nodeDefaults,omitempty" yaml:"nodeDefaults,omitempty"`
	Nodes           []Node          `json:"nodes,omitempty" yaml:"nodes,omitempty" jsonschema:"oneof_required=array"`
	Texts           []Text          `json:"texts,omitempty" yaml:"texts,omitempty"`
	TargetAppConfig TargetAppConfig `json:"targetAppConfig,omitempty" yaml:"targetAppConfig,omitempty"`
}

type TargetAppConfig struct {
	Version        string         `json:"_version,omitempty" yaml:"_version,omitempty"`
	SapIntegration SapIntegration `json:"sapIntegration,omitempty" yaml:"sapIntegration,omitempty"`
}

type SapIntegration struct {
	NavMode           string            `json:"navMode,omitempty" yaml:"navMode,omitempty"`
	UrlTemplateId     string            `json:"urlTemplateId,omitempty" yaml:"urlTemplateId,omitempty"`
	UrlTemplateParams UrlTemplateParams `json:"urlTemplateParams,omitempty" yaml:"urlTemplateParams,omitempty"`
}

type UrlTemplateParams struct {
	Query interface{} `json:"query,omitempty" yaml:"query,omitempty"`
}

type NodeDefaults struct {
	EntityType  string `json:"entityType,omitempty" yaml:"entityType,omitempty"`
	IsolateView bool   `json:"isolateView,omitempty" yaml:"isolateView,omitempty"`
}

type Text struct {
	Locale         string            `json:"locale,omitempty" yaml:"locale,omitempty"`
	TextDictionary map[string]string `json:"textDictionary" yaml:"textDictionary"`
}

type Node struct {
	EntityType                string       `json:"entityType,omitempty" yaml:"entityType,omitempty"`
	PathSegment               string       `json:"pathSegment,omitempty" yaml:"pathSegment,omitempty"`
	Label                     string       `json:"label,omitempty" yaml:"label,omitempty"`
	Icon                      string       `json:"icon,omitempty" yaml:"icon,omitempty"`
	Category                  interface{}  `json:"category,omitempty" yaml:"category,omitempty" jsonschema:"anyof_type=string;object"`
	Url                       string       `json:"url,omitempty" yaml:"url,omitempty"`
	HideFromNav               bool         `json:"hideFromNav,omitempty" yaml:"hideFromNav,omitempty"`
	VisibleForFeatureToggles  []string     `json:"visibleForFeatureToggles,omitempty" yaml:"visibleForFeatureToggles,omitempty"`
	VirtualTree               bool         `json:"virtualTree,omitempty" yaml:"virtualTree,omitempty"`
	RequiredIFramePermissions interface{}  `json:"requiredIFramePermissions,omitempty" yaml:"requiredIFramePermissions,omitempty" jsonschema:"anyof_type=object"`
	Compound                  interface{}  `json:"compound,omitempty" yaml:"compound,omitempty" jsonschema:"anyof_type=object"`
	InitialRoute              string       `json:"initialRoute,omitempty" yaml:"initialRoute,omitempty"`
	LayoutConfig              interface{}  `json:"layoutConfig,omitempty" yaml:"layoutConfig,omitempty" jsonschema:"anyof_type=object"`
	Context                   interface{}  `json:"context,omitempty" yaml:"context,omitempty" jsonschema:"anyof_type=object"`
	Webcomponent              Webcomponent `json:"webcomponent,omitempty" yaml:"webcomponent,omitempty"`
	LoadingIndicator          interface{}  `json:"loadingIndicator,omitempty" yaml:"loadingIndicator,omitempty" jsonschema:"anyof_type=object"`
	DefineEntity              DefineEntity `json:"defineEntity,omitempty" yaml:"defineEntity,omitempty"`
	KeepSelectedForChildren   bool         `json:"keepSelectedForChildren,omitempty" yaml:"keepSelectedForChildren,omitempty"`
	Children                  []Node       `json:"children,omitempty" yaml:"children,omitempty"`
	UrlSuffix                 string       `json:"urlSuffix,omitempty" yaml:"urlSuffix,omitempty"`
	HideSideNav               bool         `json:"hideSideNav,omitempty" yaml:"hideSideNav,omitempty"`
	TabNav                    bool         `json:"tabNav,omitempty" yaml:"tabNav,omitempty"`
}

type DefineEntity struct {
	Id string `json:"id,omitempty" yaml:"id,omitempty"`
}

type Webcomponent struct {
	SelfRegistered bool `json:"selfRegistered,omitempty" yaml:"selfRegistered,omitempty"`
}

type Category struct {
	Label       string `json:"label,omitempty" yaml:"label,omitempty"`
	Icon        string `json:"icon,omitempty" yaml:"icon,omitempty"`
	Collapsible bool   `json:"collapsible,omitempty" yaml:"collapsible,omitempty"`
}
