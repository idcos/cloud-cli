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

func TestYAMLRepo(t *testing.T) {
	tempfile, err := ioutil.TempFile("", "test.yaml")
	if err != nil {
		t.Errorf("prepare yaml file failed: %v", err)
	}

	defer os.Remove(tempfile.Name())
	tempfile.WriteString(yamlStr)
	if err = tempfile.Close(); err != nil {
		t.Errorf("prepare yaml file content failed: %v", err)
	}

	var yamlRepo *YAMLRepo
	if yamlRepo, err = New(tempfile.Name()); err != nil {
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
