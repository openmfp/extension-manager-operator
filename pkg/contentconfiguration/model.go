package contentconfiguration

type ContentConfiguration struct {
	Name                string                `json:"name" yaml:"name"`
	LuigiConfigFragment []LuigiConfigFragment `json:"luigiConfigFragment" yaml:"luigiConfigFragment"`
}

type LuigiConfigFragment struct {
	Data LuigiConfigData `json:"data" yaml:"data"`
}

type LuigiConfigData struct {
	Nodes []Node `json:"nodes" yaml:"nodes"`
}

type Node struct {
	EntityType  string `json:"entityType" yaml:"entityType"`
	PathSegment string `json:"pathSegment" yaml:"pathSegment"`
	Label       string `json:"label" yaml:"label"`
	Icon        string `json:"icon" yaml:"icon"`
}
