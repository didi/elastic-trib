package elastic

import (
	"context"
	"net/url"
	"strings"

	"github.com/olivere/elastic/uritemplates"
)

// CatAliasService allows to get the master alias
type CatAliasService struct {
	client *Client
	format string
	pretty bool
	alias  []string
}

// NewCatAliasService creates a new CatAliasService.
func NewCatAliasService(client *Client) *CatAliasService {
	return &CatAliasService{
		client: client,
		format: "json",
		alias:  make([]string, 0),
	}
}

// Alias limits the information returned to specific indices.
func (s *CatAliasService) Alias(alias ...string) *CatAliasService {
	s.alias = append(s.alias, alias...)
	return s
}

// Pretty indicates that the JSON response be indented and human readable.
func (s *CatAliasService) Pretty(pretty bool) *CatAliasService {
	s.pretty = pretty
	return s
}

func (s *CatAliasService) buildURL() (string, url.Values, error) {
	var err error
	var path string

	if len(s.alias) > 0 {
		path, err = uritemplates.Expand("/_cat/aliases/{index}", map[string]string{
			"index": strings.Join(s.alias, ","),
		})
	} else {
		path = "/_cat/aliases"
	}

	if err != nil {
		return "", url.Values{}, err
	}

	// Add query string parameters
	params := url.Values{}
	if s.pretty {
		params.Set("pretty", "true")
	}
	return path, params, nil
}

// Do executes the operation.
func (s *CatAliasService) Do(ctx context.Context) (*CatAliasResponse, error) {

	// Get URL for request
	path, params, err := s.buildURL()
	if err != nil {
		return nil, err
	}

	// Get HTTP response
	res, err := s.client.PerformRequest(ctx, PerformRequestOptions{
		Method: "GET",
		Path:   path,
		Params: params,
	})
	if err != nil {
		return nil, err
	}

	// Return operation response
	ret := new(CatAliasResponse)
	if err := s.client.decoder.Decode(res.Body, &ret.Aliases); err != nil {
		return nil, err
	}
	return ret, nil
}

// CatAliasResponse is the response of CataliasService.Do
type CatAliasResponse struct {
	Aliases []*indicesAlias
}

type indicesAlias struct {
	Alias         string `json:"alias"`
	Index         string `json:"index"`
	Filter        string `json:"filter"`
	Routingindex  string `json:"routing.index"`
	Routingsearch string `json:"routing.search"`
}
