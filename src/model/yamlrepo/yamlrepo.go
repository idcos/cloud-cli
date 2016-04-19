package yamlrepo

import (
	"io/ioutil"
	"model"

	yaml "gopkg.in/yaml.v2"
)

type YAMLRepo struct {
	YAMLFilePath string
	NodeGroups   []model.NodeGroup `yaml:"NodeGroups"`
}

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

func (yp *YAMLRepo) GetNodeGroups(gName string) ([]model.NodeGroup, error) {

	return nil, nil
}

func (yp *YAMLRepo) GetNodesByGroupName(gName, nName string) ([]model.Node, error) {

	return nil, nil
}
