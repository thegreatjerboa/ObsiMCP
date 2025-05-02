package local_rest_api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"obsimcp/src/utils"
	"strings"
)

const (
	FileRoute          = "/vault/%s"
	DirectoryRoute     = "/vault/%s/"
	RootDirectoryRoute = "/vault/"
)

// GetVaultFile Return the content of a single file in your vault.
// The filename should be the relative path to the file in your vault.
func (l localRestApi) GetVaultFile(filename string) (details NoteDetails, err error) {
	url := l.BaseUrl + fmt.Sprintf(FileRoute, filename)
	header := map[string]string{
		"accept": "application/vnd.olrapi.note+json",
	}
	respBody, code, err := utils.Request(url, "GET", nil, l.AuthToken, header)
	if err != nil {
		return
	}

	if code != http.StatusOK {
		apiError := APIError{}
		if err = json.Unmarshal(respBody, &apiError); err == nil {
			err = fmt.Errorf("ErrorCode: %d Error: %s", apiError.ErrorCode, apiError.Message)
		}
		return
	}
	err = json.Unmarshal(respBody, &details)
	return
}

// DeleteVaultFile Delete a file in your vault.
func (l localRestApi) DeleteVaultFile(filename string) (err error) {
	url := l.BaseUrl + fmt.Sprintf(FileRoute, filename)
	header := map[string]string{
		"accept": "*/*",
	}
	respBody, code, err := utils.Request(url, "DELETE", nil, l.AuthToken, header)
	if err != nil {
		return
	}

	if code != http.StatusNoContent {
		apiError := APIError{}
		if err = json.Unmarshal(respBody, &apiError); err == nil {
			return fmt.Errorf("ErrorCode: %d Error: %s", apiError.ErrorCode, apiError.Message)
		}
	}
	return
}

// AppendVaultFile Append content to a new or existing file.
// Appends content to the end of an existing note. If the specified file does not yet exist, it will be created as an empty file.
func (l localRestApi) AppendVaultFile(filename string, content string) (err error) {
	url := l.BaseUrl + fmt.Sprintf(FileRoute, filename)
	reader := strings.NewReader(content)
	header := map[string]string{
		"Content-Type": "text/markdown",
	}
	respBody, code, err := utils.Request(url, "POST", reader, l.AuthToken, header)
	if err != nil {
		return
	}

	if code != http.StatusNoContent {
		apiError := APIError{}
		if err = json.Unmarshal(respBody, &apiError); err == nil {
			return fmt.Errorf("ErrorCode: %d Error: %s", apiError.ErrorCode, apiError.Message)
		}
	}
	return
}

// CreateOrUpdateVaultFile
// Creates a new file in your vault or updates the content of an existing one if the specified file already exists.
// if the file already exists, it will be overwritten.
func (l localRestApi) CreateOrUpdateVaultFile(filename string, content string) (err error) {
	url := l.BaseUrl + fmt.Sprintf(FileRoute, filename)
	reader := strings.NewReader(content)
	header := map[string]string{
		"Content-Type": "text/markdown",
	}
	respBody, code, err := utils.Request(url, "PUT", reader, l.AuthToken, header)
	if err != nil {
		return
	}

	if code != http.StatusNoContent {
		apiError := APIError{}
		if err = json.Unmarshal(respBody, &apiError); err == nil {
			return fmt.Errorf("ErrorCode: %d Error: %s", apiError.ErrorCode, apiError.Message)
		}
	}
	return
}

func (l localRestApi) GetDirectory(path string) (res []string, err error) {
	var url string
	if path == "" || path == "/" {
		url = l.BaseUrl + RootDirectoryRoute
	} else {
		url = l.BaseUrl + fmt.Sprintf(DirectoryRoute, path)
	}
	header := map[string]string{
		"accept": "application/json",
	}
	respBody, code, err := utils.Request(url, "GET", nil, l.AuthToken, header)
	if err != nil {
		return
	}
	if code != http.StatusOK {
		apiError := APIError{}
		if err = json.Unmarshal(respBody, &apiError); err == nil {
			err = fmt.Errorf("ErrorCode: %d Error: %s", apiError.ErrorCode, apiError.Message)
		}
		return
	}
	var resp ListDirResult
	err = json.Unmarshal(respBody, &resp)
	if err != nil {
		return
	}
	if len(resp.Files) == 0 {
		err = fmt.Errorf("no files found in directory: %s", path)
		return
	}
	res = resp.Files
	return
}
