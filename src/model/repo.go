package model

// Node store info for host
type Node struct {
	Name     string `yaml:"Name"`
	Host     string `yaml:"Host"`
	User     string `yaml:"User"`
	Password string `yaml:"Password"`
	KeyPath  string `yaml:"KeyPath"`
}

// NodeGroup store info for node group
type NodeGroup struct {
	Name  string `yaml:"GroupName"`
	Nodes []Node `yaml:"Nodes"`
}

// IRepo repo interface
type IRepo interface {
	FilterNodeGroups(gName string) ([]NodeGroup, error)
	FilterNodeGroupsAndNodes(gName string, nName string) ([]NodeGroup, error)
}
