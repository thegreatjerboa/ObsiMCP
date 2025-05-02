package local_rest_api

import "obsimcp/src/config"

type LocalRestApi interface {
	// Vault Files

	// GetVaultFile  get file's details in vault
	GetVaultFile(filename string) (NoteDetails, error)
	// DeleteVaultFile  delete file in vault
	DeleteVaultFile(filename string) error
	// AppendVaultFile Insert content into an existing note relative to a heading within that document.
	// PatchVaultFile() error
	// AppendVaultFile Append content to a new or existing file.
	AppendVaultFile(filename string, content string) error
	// GetDirectory return the directory structure of your vault.
	GetDirectory(path string) ([]string, error)

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

func NewLocalRestApi() LocalRestApi {
	return &localRestApi{
		BaseUrl:   config.Cfg.Plugins.Rest.BaseUrl,
		AuthToken: config.Cfg.Plugins.Rest.AuthToken,
	}
}
