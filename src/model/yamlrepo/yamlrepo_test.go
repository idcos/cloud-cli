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
            IP: ipxxx
            User: userxxx
            Password: passwordxxx
            KeyPath: kaypathxxx
          - Name: nameyyy
            IP: ipyyy
            User: useryyy
            Password: passwordyyy
            KeyPath: keypathyyy

    - GroupName: groupyyy
      Nodes:
          - Name: namexxx
            IP: ipxxx
            User: userxxx
            Password: passwordxxx
            KeyPath: keypathxxx
          - Name: namezzz
            IP: ipzzz
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

func TestGetNodeGroups(t *testing.T) {
	yamlFilePath, err := prepareYAML()
	if err != nil {
		t.Errorf("prepare yaml file content failed: %v", err)
	}
	defer os.Remove(yamlFilePath)

	var yamlRepo *YAMLRepo
	groups, err := yamlRepo.GetNodeGroups("")
	if len(groups) != 0 {
		t.Errorf("group count must equal 0")
	}

	if yamlRepo, err = New(yamlFilePath); err != nil {
		t.Errorf("new yaml repo failed: %v", err)
	}

	groups, err = yamlRepo.GetNodeGroups("")
	if len(groups) != 2 {
		t.Errorf("group count must equal 2")
	}

	groups, err = yamlRepo.GetNodeGroups("xx")
	if len(groups) != 1 {
		t.Errorf("group count must equal 1")
	}

	groups, err = yamlRepo.GetNodeGroups("grou")
	if len(groups) != 2 {
		t.Errorf("group count must equal 2")
	}
}

func TestGetNodesByGroupName(t *testing.T) {
	yamlFilePath, err := prepareYAML()
	if err != nil {
		t.Errorf("prepare yaml file content failed: %v", err)
	}
	defer os.Remove(yamlFilePath)

	var yamlRepo *YAMLRepo
	groups, err := yamlRepo.GetNodesByGroupName("", "")
	if len(groups) != 0 {
		t.Errorf("group count must equal 0")
	}

	if yamlRepo, err = New(yamlFilePath); err != nil {
		t.Errorf("new yaml repo failed: %v", err)
	}

	groups, err = yamlRepo.GetNodesByGroupName("", "")
	if len(groups) != 4 {
		t.Errorf("group count must equal 4")
	}

	groups, err = yamlRepo.GetNodesByGroupName("xx", "")
	if len(groups) != 2 {
		t.Errorf("group count must equal 2")
	}

	groups, err = yamlRepo.GetNodesByGroupName("grou", "")
	if len(groups) != 4 {
		t.Errorf("group count must equal 4")
	}

	groups, err = yamlRepo.GetNodesByGroupName("", "xx")
	if len(groups) != 2 {
		t.Errorf("group count must equal 2")
	}

	groups, err = yamlRepo.GetNodesByGroupName("", "zz")
	if len(groups) != 1 {
		t.Errorf("group count must equal 1")
	}

	groups, err = yamlRepo.GetNodesByGroupName("", "name")
	if len(groups) != 4 {
		t.Errorf("group count must equal 4")
	}
}
