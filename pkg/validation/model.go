package validation

type ContentConfiguration struct {
	Name                string                `json:"name,omitempty" yaml:"name,omitempty"`
	LuigiConfigFragment []LuigiConfigFragment `json:"luigiConfigFragment,omitempty" yaml:"luigiConfigFragment,omitempty"`
}

type LuigiConfigFragment struct {
	Data LuigiConfigData `json:"data,omitempty" yaml:"data,omitempty"`
}

type LuigiConfigData struct {
	NodeDefaults NodeDefaults `json:"nodeDefaults,omitempty" yaml:"nodeDefaults,omitempty"`
	Nodes        []Node       `json:"nodes,omitempty" yaml:"nodes,omitempty"`
	Texts        []Text       `json:"texts,omitempty" yaml:"texts,omitempty"`
}

type NodeDefaults struct {
	EntityType  string `json:"entityType,omitempty" yaml:"entityType,omitempty"`
	IsolateView bool   `json:"isolateView,omitempty" yaml:"isolateView,omitempty"`
}

type Text struct {
	Locale         string            `json:"locale,omitempty" yaml:"locale,omitempty"`
	TextDictionary map[string]string `json:"textDictionary,omitempty" yaml:"textDictionary,omitempty"`
}

type Node struct {
	EntityType  string `json:"entityType,omitempty" yaml:"entityType,omitempty"`
	PathSegment string `json:"pathSegment,omitempty" yaml:"pathSegment,omitempty"`
	Label       string `json:"label,omitempty" yaml:"label,omitempty"`
	Icon        string `json:"icon,omitempty" yaml:"icon,omitempty"`
}
