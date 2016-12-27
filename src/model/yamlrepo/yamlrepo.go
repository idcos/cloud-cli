package yamlrepo

import (
	"io/ioutil"
	"model"

	"utils"

	"regexp"

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
	var err error
	var buf []byte
	var yamlRepo YAMLRepo

	yamlRepo.YAMLFilePath, err = utils.ConvertHomeDir(yamlFilePath)
	if err != nil {
		return nil, err
	}

	buf, err = ioutil.ReadFile(yamlRepo.YAMLFilePath)
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

	gNamePattern := utils.WildCharToRegexp(gName)
	for _, g := range yp.NodeGroups {
		matched, _ := regexp.MatchString(gNamePattern, g.Name)
		if matched {
			filterNodeGroups = append(filterNodeGroups, g)
		}
	}

	return filterNodeGroups, nil
}

// FilterNodeGroupsAndNodes find nodes and groups info from yaml repo
func (yp *YAMLRepo) FilterNodeGroupsAndNodes(gName string, nNames ...string) ([]model.NodeGroup, error) {
	var groups, _ = yp.FilterNodeGroups(gName)
	var filterGroups = make([]model.NodeGroup, 0)

	for _, g := range groups {
		if g.Nodes == nil {
			continue
		}

		var filterNodes = make([]model.Node, 0)

		for _, n := range g.Nodes {
			if utils.IsWildCharMatch(n.Name, nNames...) {
				filterNodes = append(filterNodes, n)
			}
		}

		// only return groups which has nodes
		if len(filterNodes) > 0 {
			g.Nodes = filterNodes
			g = initNodesByGroup(g)
			filterGroups = append(filterGroups, g)
		}
	}

	return filterGroups, nil
}

// FilterNodes find nodes info from yaml repo
func (yp *YAMLRepo) FilterNodes(gName string, nNames ...string) ([]model.Node, error) {
	var groups, _ = yp.FilterNodeGroups(gName)
	var filterNodes = make([]model.Node, 0)

	for _, g := range groups {
		if g.Nodes == nil {
			continue
		}

		g = initNodesByGroup(g)
		// without node filter
		if nNames == nil || len(nNames) == 0 {
			return g.Nodes, nil
		}

		for _, n := range g.Nodes {
			if utils.IsWildCharMatch(n.Name, nNames...) {
				filterNodes = append(filterNodes, n)
			}
		}
	}

	return filterNodes, nil
}

func initNodesByGroup(group model.NodeGroup) model.NodeGroup {
	if group.Nodes == nil {
		return group
	}

	for index, n := range group.Nodes {
		if n.User == "" {
			group.Nodes[index].User = group.User
		}

		if n.Password == "" {
			group.Nodes[index].Password = group.Password
		}

		if n.KeyPath == "" {
			group.Nodes[index].KeyPath = group.KeyPath
		}

		if n.Port == 0 {
			group.Nodes[index].Port = group.Port
		}
	}

	return group
}
