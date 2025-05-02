package api

import (
	"encoding/json"
	"fmt"
	"obsimcp/src/config"
	"obsimcp/src/utils"
)

type LocalRestApi interface {
	// Status

	// GetStatus return the status of the local REST API.
	GetStatus() (bool, error)

	// Vault Files

	// GetVaultFile  get file's details in vault
	GetVaultFile(filename string) (NoteDetails, error)
	// GetVaultFileFrontmatter  get file's frontmatter in vault
	GetVaultFileFrontmatter(filename string) (map[string]interface{}, error)
	// DeleteVaultFile  delete file in vault
	DeleteVaultFile(filename string) error
	// AppendVaultFile Insert content into an existing note relative to a heading within that document.
	// PatchVaultFile() error
	// AppendVaultFile Append content to a new or existing file.
	AppendVaultFile(filename string, content string) error
	// ListDirectory return the directory structure of your vault.
	ListDirectory(path string) ([]string, error)

	// CreateOrUpdateVaultFile Create a new file in your vault or update the content of an existing one.
	CreateOrUpdateVaultFile(filename string, content string) error

	// Search

	// SimpleSearch return notes' content around the search term within a specified context length.
	SimpleSearch(query string, contextLength int) ([]*SearchResult, error)
}

type localRestApi struct {
	// BaseUrl is the base URL of the local REST API of Obsidian.
	BaseUrl string
	// AuthToken is the authentication token for the local REST API of Obsidian.
	AuthToken string
}

var Client LocalRestApi

func NewLocalRestApi() LocalRestApi {
	return &localRestApi{
		BaseUrl:   config.Cfg.Plugins.Rest.BaseUrl,
		AuthToken: config.Cfg.Plugins.Rest.AuthToken,
	}
}

func InitLocalRestApi() error {
	Client = NewLocalRestApi()
	ok, err := Client.GetStatus()
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("local REST API is not running")
	}
	fmt.Println("local REST API is running")
	return nil
}

func (l localRestApi) GetStatus() (bool, error) {
	url := l.BaseUrl + "/"
	header := map[string]string{
		"accept": "application/json",
	}
	respBody, code, err := utils.Request(url, "GET", nil, l.AuthToken, header)
	if err != nil {
		return false, err
	}
	if code != 200 {
		apiError := Error{}
		if err = json.Unmarshal(respBody, &apiError); err == nil {
			return false, fmt.Errorf("ErrorCode: %d Error: %s", apiError.ErrorCode, apiError.Message)
		}
	}
	return true, nil
}
