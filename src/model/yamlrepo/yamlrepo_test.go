package yamlrepo

import (
	"io/ioutil"
	"os"
	"testing"
)

var yamlStr = `
NodeGroups:
    - GroupName: groupxxx
      Nodes:
          - Name: namexxx
            Host: ipxxx
            User: userxxx
            Password: passwordxxx
            KeyPath: kaypathxxx
          - Name: nameyyy
            Host: ipyyy
            User: useryyy
            Password: passwordyyy
            KeyPath: keypathyyy

    - GroupName: groupyyy
      Nodes:
          - Name: namexxx
            Host: ipxxx
            User: userxxx
            Password: passwordxxx
            KeyPath: keypathxxx
          - Name: namezzz
            Host: ipzzz
            User: userzzz
            Password: passwordzzz
            KeyPath: keypathzzz
`

func prepareYAML() (string, error) {
	tempfile, err := ioutil.TempFile("", "test.yaml")
	if err != nil {
		return "", err
	}

	tempfile.WriteString(yamlStr)
	if err = tempfile.Close(); err != nil {
		return "", err
	}

	return tempfile.Name(), nil
}

func TestYAMLRepo(t *testing.T) {
	yamlFilePath, err := prepareYAML()
	if err != nil {
		t.Errorf("prepare yaml file content failed: %v", err)
	}
	defer os.Remove(yamlFilePath)

	var yamlRepo *YAMLRepo
	if yamlRepo, err = New(yamlFilePath); err != nil {
		t.Errorf("new yaml repo failed: %v", err)
	}

	if len(yamlRepo.NodeGroups) != 2 {
		t.Errorf("HostGroups count is wrong")
	}

	for _, g := range yamlRepo.NodeGroups {
		if len(g.Nodes) != 2 {
			t.Errorf("Hosts count is wrong")
		}
	}
}

func TestFilterNodeGroups(t *testing.T) {
	yamlFilePath, err := prepareYAML()
	if err != nil {
		t.Errorf("prepare yaml file content failed: %v", err)
	}
	defer os.Remove(yamlFilePath)

	var yamlRepo *YAMLRepo
	groups, err := yamlRepo.FilterNodeGroups("*")
	if len(groups) != 0 {
		t.Errorf("group count must equal 0")
	}

	if yamlRepo, err = New(yamlFilePath); err != nil {
		t.Errorf("new yaml repo failed: %v", err)
	}

	groups, err = yamlRepo.FilterNodeGroups("*")
	if len(groups) != 2 {
		t.Errorf("group count must equal 2")
	}

	groups, err = yamlRepo.FilterNodeGroups("*xx*")
	if len(groups) != 1 {
		t.Errorf("group count must equal 1")
	}

	groups, err = yamlRepo.FilterNodeGroups("*grou*")
	if len(groups) != 2 {
		t.Errorf("group count must equal 2")
	}
}

func TestFilterNodeGroupsAndNodes(t *testing.T) {
	yamlFilePath, err := prepareYAML()
	if err != nil {
		t.Errorf("prepare yaml file content failed: %v", err)
	}
	defer os.Remove(yamlFilePath)

	var yamlRepo *YAMLRepo
	groups, err := yamlRepo.FilterNodeGroupsAndNodes("*", "*")
	if len(groups) != 0 {
		t.Errorf("group count must equal 0")
	}

	if yamlRepo, err = New(yamlFilePath); err != nil {
		t.Errorf("new yaml repo failed: %v", err)
	}

	groups, err = yamlRepo.FilterNodeGroupsAndNodes("*", "*")
	if (len(groups[0].Nodes) + len(groups[1].Nodes)) != 4 {
		t.Errorf("node count must equal 4")
	}

	groups, err = yamlRepo.FilterNodeGroupsAndNodes("*xx*", "*")
	if len(groups[0].Nodes) != 2 {
		t.Errorf("node count must equal 2")
	}

	groups, err = yamlRepo.FilterNodeGroupsAndNodes("grou*", "*")
	if (len(groups[0].Nodes) + len(groups[1].Nodes)) != 4 {
		t.Errorf("node count must equal 4")
	}

	groups, err = yamlRepo.FilterNodeGroupsAndNodes("*", "*xx")
	if (len(groups[0].Nodes) + len(groups[1].Nodes)) != 2 {
		t.Errorf("node count must equal 2")
	}

	groups, err = yamlRepo.FilterNodeGroupsAndNodes("*", "*zz")
	if len(groups[0].Nodes) != 1 {
		t.Errorf("node count must equal 1")
	}

	groups, err = yamlRepo.FilterNodeGroupsAndNodes("*", "*name*")
	if (len(groups[0].Nodes) + len(groups[1].Nodes)) != 4 {
		t.Errorf("node count must equal 4")
	}
}

func TestFilterNodes(t *testing.T) {
	yamlFilePath, err := prepareYAML()
	if err != nil {
		t.Errorf("prepare yaml file content failed: %v", err)
	}
	defer os.Remove(yamlFilePath)

	var yamlRepo *YAMLRepo
	nodes, err := yamlRepo.FilterNodes("*", "*")
	if len(nodes) != 0 {
		t.Errorf("node count must equal 0")
	}

	if yamlRepo, err = New(yamlFilePath); err != nil {
		t.Errorf("new yaml repo failed: %v", err)
	}

	nodes, err = yamlRepo.FilterNodes("*", "*")
	if len(nodes) != 4 {
		t.Errorf("node count must equal 4")
	}

	nodes, err = yamlRepo.FilterNodes("*xx*", "*")
	if len(nodes) != 2 {
		t.Errorf("node count must equal 2")
	}

	nodes, err = yamlRepo.FilterNodes("grou*", "*")
	if len(nodes) != 4 {
		t.Errorf("node count must equal 4")
	}

	nodes, err = yamlRepo.FilterNodes("*", "*xx")
	if len(nodes) != 2 {
		t.Errorf("node count must equal 2")
	}

	nodes, err = yamlRepo.FilterNodes("*", "*zz")
	if len(nodes) != 1 {
		t.Errorf("node count must equal 1")
	}

	nodes, err = yamlRepo.FilterNodes("*", "*name*")
	if len(nodes) != 4 {
		t.Errorf("node count must equal 4")
	}
}
