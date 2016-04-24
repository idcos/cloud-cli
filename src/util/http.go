package util

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

// ErrUploadFile error message when upload file failed
var ErrUploadFile = "upload file failed: URL(%s)\nfilename(%s)\nmessage(%s)"

// PostFile post file to targetUrl
func PostFile(fileParam, filePath string, extraParams map[string]string, targetURL string) error {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	var (
		fileWriter io.Writer
		fh         *os.File
		resp       *http.Response
		output     []byte
		err        error
	)

	// this step is very important
	fileWriter, err = bodyWriter.CreateFormFile(fileParam, filepath.Base(filePath))
	if err != nil {
		return err
	}

	// open file handle
	fh, err = os.Open(filePath)
	if err != nil {
		return err
	}

	//iocopy
	_, err = io.Copy(fileWriter, fh)
	if err != nil {
		return err
	}

	contentType := bodyWriter.FormDataContentType()
	// add params
	if extraParams != nil {
		for k, v := range extraParams {
			bodyWriter.WriteField(k, v)
		}
	}

	bodyWriter.Close()

	resp, err = http.Post(targetURL, contentType, bodyBuf)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	output, err = ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf(ErrUploadFile, targetURL, filePath, string(output))
	}

	return err
}
