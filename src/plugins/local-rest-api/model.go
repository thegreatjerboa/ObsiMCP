package local_rest_api

type APIError struct {
	ErrorCode int    `json:"errorCode"`
	Message   string `json:"message"`
}

type NoteDetails struct {
	Tags        []interface{}          `json:"tags"`
	Frontmatter map[string]interface{} `json:"frontmatter"`
	Stat        struct {
		Ctime int64 `json:"ctime"`
		Mtime int64 `json:"mtime"`
		Size  int   `json:"size"`
	} `json:"stat"`
	Path    string `json:"path"`
	Content string `json:"content"`
}

type ListDirResult struct {
	Files []string `json:"files"`
}

type SearchResult struct {
	Filename string  `json:"filename,omitempty" mapstructure:"filename"`
	Score    float64 `json:"score,omitempty" mapstructure:"score"`
	Matches  []struct {
		Context string `json:"context,omitempty" mapstructure:"context"`
	} `json:"matches,omitempty" mapstructure:"matches"`
}
