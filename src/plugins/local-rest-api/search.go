package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"obsimcp/src/utils"
)

type SimpleSearchResponse struct {
	Filename string `json:"filename" mapstructure:"filename"`
	Matches  []struct {
		Context string `json:"context" mapstructure:"context"`
		Match   struct {
			End   int `json:"end" mapstructure:"end"`
			Start int `json:"start" mapstructure:"start"`
		} `json:"match" mapstructure:"match"`
	} `json:"matches" mapstructure:"matches"`
	Score float64 `json:"score" mapstructure:"score"`
}

const (
	SimpleSearchRoute = "/search/simple/?query=%s&contextLength=%d"
)

func (l localRestApi) SimpleSearch(query string, contextLength int) (searchResults []*SearchResult, err error) {
	url := l.BaseUrl + fmt.Sprintf(SimpleSearchRoute, query, contextLength)
	header := map[string]string{
		"accept": "application/json",
	}
	respBody, code, err := utils.Request(url, "POST", nil, l.AuthToken, header)
	if err != nil {
		return
	}

	if code != http.StatusOK {
		apiError := Error{}
		if err = json.Unmarshal(respBody, &apiError); err == nil {
			err = fmt.Errorf("ErrorCode: %d Error: %s", apiError.ErrorCode, apiError.Message)
		}
		return
	}

	var searchResponse []*SimpleSearchResponse
	err = json.Unmarshal(respBody, &searchResponse)
	if err != nil {
		return
	}
	for _, result := range searchResponse {
		searchResult := &SearchResult{
			Filename: result.Filename,
			Matches: []struct {
				Context string `json:"context,omitempty" mapstructure:"context"`
			}(make([]struct{ Context string }, len(result.Matches))),
		}
		for i, match := range result.Matches {
			searchResult.Matches[i].Context = match.Context
		}
		searchResults = append(searchResults, searchResult)
	}
	return
}
