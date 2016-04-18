package model

// Node store info for host
type Node struct {
	Name     string
	IP       string
	User     string
	Password string
	KeyPath  string
}

// Group store info for node group
type Group struct {
	Name  string
	Nodes []Node
}

type IRepo interface {
	GetAllGroups() ([]Group, error)
	GetNodesByGroupName(groupName string) ([]Node, error)
}
