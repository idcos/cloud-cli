package model

// Node store info for host
type Node struct {
	Name     string `yaml:"Name"`
	Host     string `yaml:"Host"`
	Port     int    `yaml:"Port"`
	User     string `yaml:"User"`
	Password string `yaml:"Password"`
	KeyPath  string `yaml:"KeyPath"`
}

// NodeGroup store info for node group
type NodeGroup struct {
	Name     string `yaml:"GroupName"`
	User     string `yaml:"User"`
	Password string `yaml:"Password"`
	KeyPath  string `yaml:"KeyPath"`
	Port     int    `yaml:"Port"`
	Nodes    []Node `yaml:"Nodes"`
}

// IRepo repo interface
type IRepo interface {
	FilterNodeGroups(gName string) ([]NodeGroup, error)
	FilterNodeGroupsAndNodes(gName string, nNames ...string) ([]NodeGroup, error)
	FilterNodes(gName string, nNames ...string) ([]Node, error)
}
