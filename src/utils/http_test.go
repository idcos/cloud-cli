package utils

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestPostFile(t *testing.T) {
	var filePath = "testfile"
	f, err := os.Create(filePath)
	if err != nil {
		t.Errorf("prepare test failed: %v", err)
	}

	defer os.Remove(filePath)
	f.WriteString("testtesttset")
	f.Close()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "test post file")
	}))
	defer ts.Close()

	if err := PostFile("file", filePath, nil, ts.URL); err != nil {
		t.Errorf("post file error: %v", err)
	}
}
