package validation

type Node struct {
	EntityType  string `json:"entityType,omitempty" yaml:"entityType,omitempty"`
	PathSegment string `json:"pathSegment,omitempty" yaml:"pathSegment,omitempty"`
	Label       string `json:"label,omitempty" yaml:"label,omitempty"`
	Icon        string `json:"icon,omitempty" yaml:"icon,omitempty"`
}

type LuigiConfigData struct {
	Nodes []Node `json:"nodes,omitempty" yaml:"nodes,omitempty"`
}

type LuigiConfigFragment struct {
	Data LuigiConfigData `json:"data,omitempty" yaml:"data,omitempty"`
}

type ContentConfiguration struct {
	Name                string                `json:"name,omitempty" yaml:"name,omitempty"`
	LuigiConfigFragment []LuigiConfigFragment `json:"luigiConfigFragment,omitempty" yaml:"luigiConfigFragment,omitempty"`
}
