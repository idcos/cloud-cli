package yamlrepo

import (
	"io/ioutil"
	"model"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

// YAMLRepo yaml format repo
type YAMLRepo struct {
	// YAMLFilePath yaml file path for repo
	YAMLFilePath string
	// NodeGroups node groups info from yaml
	NodeGroups []model.NodeGroup `yaml:"NodeGroups"`
}

// New create yaml repo
func New(yamlFilePath string) (*YAMLRepo, error) {
	var yamlRepo YAMLRepo
	yamlRepo.YAMLFilePath = yamlFilePath

	buf, err := ioutil.ReadFile(yamlFilePath)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(buf, &yamlRepo)
	return &yamlRepo, err
}

// FilterNodeGroups find groups info from yaml repo
func (yp *YAMLRepo) FilterNodeGroups(gName string) ([]model.NodeGroup, error) {
	var filterNodeGroups = make([]model.NodeGroup, 0)

	if yp == nil {
		return filterNodeGroups, nil
	}

	for _, g := range yp.NodeGroups {
		if strings.Contains(g.Name, gName) {
			filterNodeGroups = append(filterNodeGroups, g)
		}
	}

	return filterNodeGroups, nil
}

// FilterNodeGroupsAndNodes find nodes and groups info from yaml repo
func (yp *YAMLRepo) FilterNodeGroupsAndNodes(gName, nName string) ([]model.NodeGroup, error) {
	var groups, _ = yp.FilterNodeGroups(gName)
	var filterGroups = make([]model.NodeGroup, 0)

	for _, g := range groups {
		if g.Nodes == nil {
			continue
		}

		var filterNodes = make([]model.Node, 0)
		for _, n := range g.Nodes {
			if strings.Contains(n.Name, nName) {
				filterNodes = append(filterNodes, n)
			}
		}

		// only return groups which has nodes
		if len(filterNodes) > 0 {
			g.Nodes = filterNodes
			filterGroups = append(filterGroups, g)
		}
	}

	return filterGroups, nil
}
