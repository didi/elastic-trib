package elastic

import (
	"context"
	"fmt"
	"net/url"
)

// CatNodeAttrsService allows to get the node attrs of the cluster.
//
// See https://www.elastic.co/guide/en/elasticsearch/reference/master/cat-nodeattrs.html
// for details.
type CatNodeAttrsService struct {
	client        *Client
	format        string
	pretty        bool
	local         *bool
	masterTimeout string
	timeout       string
}

// NewCatNodeAttrsService creates a new CatNodeAttrsService.
func NewCatNodeAttrsService(client *Client) *CatNodeAttrsService {
	return &CatNodeAttrsService{
		client: client,
		format: "json",
	}
}

// Format indicates that the JSON response be indented and human readable.
func (s *CatNodeAttrsService) Format(format string) *CatNodeAttrsService {
	if format != "" {
		s.format = format
	} else {
		s.format = "json"
	}
	return s
}

// Pretty indicates that the JSON response be indented and human readable.
func (s *CatNodeAttrsService) Pretty(pretty bool) *CatNodeAttrsService {
	s.pretty = pretty
	return s
}

// Local indicates whether to return local information. If it is true,
// we do not retrieve the state from master node (default: false).
func (s *CatNodeAttrsService) Local(local bool) *CatNodeAttrsService {
	s.local = &local
	return s
}

// MasterTimeout specifies an explicit operation timeout for connection to master node.
func (s *CatNodeAttrsService) MasterTimeout(masterTimeout string) *CatNodeAttrsService {
	s.masterTimeout = masterTimeout
	return s
}

// Timeout specifies an explicit operation timeout.
func (s *CatNodeAttrsService) Timeout(timeout string) *CatNodeAttrsService {
	s.timeout = timeout
	return s
}

// buildURL builds the URL for the operation.
func (s *CatNodeAttrsService) buildURL() (string, url.Values, error) {
	// Build URL
	var path string

	path = "/_cat/nodeattrs"

	// Add query string parameters
	params := url.Values{}
	if s.pretty {
		params.Set("pretty", "true")
	}
	if s.local != nil {
		params.Set("local", fmt.Sprintf("%v", *s.local))
	}
	if s.masterTimeout != "" {
		params.Set("master_timeout", s.masterTimeout)
	}
	if s.timeout != "" {
		params.Set("timeout", s.timeout)
	}

	return path, params, nil
}

// Validate checks if the operation is valid.
func (s *CatNodeAttrsService) Validate() error {
	return nil
}

// Do executes the operation.
func (s *CatNodeAttrsService) Do(ctx context.Context) (*CatNodeAttrsResponse, error) {
	// Check pre-conditions
	if err := s.Validate(); err != nil {
		return nil, err
	}

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
	ret := new(CatNodeAttrsResponse)
	if err := s.client.decoder.Decode(res.Body, &ret.NodeAttrs); err != nil {
		return nil, err
	}
	return ret, nil
}

// CatNodeAttrsResponse is the response of CatNodeAttrsService.Do
type CatNodeAttrsResponse struct {
	NodeAttrs []*nodeAttrRecord
}

type nodeAttrRecord struct {
	Node  string `json:"node"`
	Host  string `json:"host"`
	IP    string `json:"ip"`
	Attr  string `json:"attr"`
	Value string `json:"value"`
}
