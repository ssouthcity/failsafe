package bungie

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const root = "https://bungie.net/platform"

type DestinyIconSequenceDefinition struct {
	Frames []string `json:"frames"`
}

type DestinyDisplayPropertiesDefinition struct {
	Name          string                        `json:"name"`
	Description   string                        `json:"description"`
	Icon          string                        `json:"icon"`
	IconSequences DestinyIconSequenceDefinition `json:"iconSequences"`
	HighResIcon   string                        `json:"highResIcon"`
	HasIcon       bool                          `json:"hasIcon"`
}

type DestinyEntitySearchResultItem struct {
	Hash              uint32                             `json:"hash"`
	EntityType        string                             `json:"entityType"`
	DisplayProperties DestinyDisplayPropertiesDefinition `json:"displayProperties"`
	Weight            float64                            `json:"weight"`
}

type PagedQuery struct {
	ItemsPerPage             int32  `json:"itemsPerPage"`
	CurrentPage              int32  `json:"currentPage"`
	RequestContinuationToken string `json:"requestContinuationToken"`
}

type SearchResultOfDestinyEntitySearchResultItem struct {
	Results                      []DestinyEntitySearchResultItem `json:"results"`
	TotalResults                 int32                           `json:"totalResults"`
	HasMore                      bool                            `json:"hasMore"`
	Query                        PagedQuery                      `json:"query"`
	ReplacementContinuationToken string                          `json:"replacementContinuationToken"`
	UseTotalResults              bool                            `json:"useTotalResults"`
}

type DestinyEntitySearchResult struct {
	SuggestedWords []string                                    `json:"suggestedWords"`
	Results        SearchResultOfDestinyEntitySearchResultItem `json:"results"`
}

type SearchEntityResult struct {
	Response           DestinyEntitySearchResult `json:"Response"`
	ErrorCode          int32                     `json:"ErrorCode"`
	ThrottleSeconds    int32                     `json:"ThrottleSeconds"`
	ErrorStatus        string                    `json:"ErrorStatus"`
	Message            string                    `json:"Message"`
	MessageData        map[string]string         `json:"MessageData"`
	DetailedErrorTrace string                    `json:"DetailedErrorTrace"`
}

func SearchEntity(def ManifestDefinition, term string, opts ...RequestOption) (SearchEntityResult, error) {
	var result SearchEntityResult

	if term == "" {
		return result, fmt.Errorf("cannot search with empty term")
	}

	endpoint := fmt.Sprintf("/Destiny2/Armory/Search/%s/%s", def, term)
	fullurl := root + endpoint

	client := http.Client{}

	req, err := http.NewRequest(http.MethodGet, fullurl, nil)
	if err != nil {
		return result, err
	}

	for _, opt := range opts {
		opt(req)
	}

	resp, err := client.Do(req)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return result, err
	}

	return result, nil
}
