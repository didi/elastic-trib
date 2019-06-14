package elastic

import (
	"context"
	"fmt"
	"net/url"
)

// CatAllocService allows to get node allocation of the cluster.
//
// See http://www.elastic.co/guide/en/elasticsearch/reference/master/cat-master.html
// for details.
type CatAllocService struct {
	client        *Client
	format        string
	local         *bool
	masterTimeout string
}

// NewCatAllocService creates a new CatAllocService.
func NewCatAllocService(client *Client) *CatAllocService {
	return &CatAllocService{
		client: client,
		format: "json",
	}
}

// Format indicates that the JSON response be indented and human readable.
//func (s *CatAllocService) Format(format string) *CatAllocService {
//	if format != "" {
//		s.format = format
//	} else {
//		s.format = "json"
//	}
//	return s
//}

// Local indicates whether to return local information. If it is true,
// we do not retrieve the state from master node (default: false).
func (s *CatAllocService) Local(local bool) *CatAllocService {
	s.local = &local
	return s
}

// MasterTimeout specifies an explicit operation timeout for connection to master node.
func (s *CatAllocService) MasterTimeout(masterTimeout string) *CatAllocService {
	s.masterTimeout = masterTimeout
	return s
}

// buildURL builds the URL for the operation.
func (s *CatAllocService) buildURL() (string, url.Values, error) {
	// Build URL
	var path string
	path = "/_cat/allocation"

	// Add query string parameters
	params := url.Values{}
	if s.local != nil {
		params.Set("local", fmt.Sprintf("%v", *s.local))
	}
	if s.masterTimeout != "" {
		params.Set("master_timeout", s.masterTimeout)
	}

	return path, params, nil
}

// Validate checks if the operation is valid.
func (s *CatAllocService) Validate() error {
	return nil
}

// Do executes the operation.
func (s *CatAllocService) Do(ctx context.Context) (*CatAllocResponse, error) {
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
	ret := new(CatAllocResponse)
	if err := s.client.decoder.Decode(res.Body, &ret.Allocs); err != nil {
		return nil, err
	}
	return ret, nil
}

// CatAllocResponse is the response of CatAllocService.Do
type CatAllocResponse struct {
	Allocs []*allocRecord
}

type allocRecord struct {
	Shards  string `json:"shards"`
	Indices string `json:"disk.indices"`
	Used    string `json:"disk.used"`
	Avail   string `json:"disk.avail"`
	Total   string `json:"disk.total"`
	Percent string `json:"disk.percent"`
	Host    string `json:"host"`
	Ip      string `json:"ip"`
	Node    string `json:"node"`
}
