package elastic

import (
	"context"
	"fmt"
	"net/url"
)

// CatRecoveryService allows to get the master node of the cluster.
//
// See http://www.elastic.co/guide/en/elasticsearch/reference/master/cat-recovery.html
// for details.
type CatRecoveryService struct {
	client        *Client
	format        string
	pretty        bool
	local         *bool
	masterTimeout string
	timeout       string
}

// NewCatRecoveryService creates a new CatRecoveryService.
func NewCatRecoveryService(client *Client) *CatRecoveryService {
	return &CatRecoveryService{
		client: client,
	}
}

// Format indicates that the JSON response be indented and human readable.
func (s *CatRecoveryService) Format(format string) *CatRecoveryService {
	if format != "" {
		s.format = format
	} else {
		s.format = "json"
	}
	return s
}

// Pretty indicates that the JSON response be indented and human readable.
func (s *CatRecoveryService) Pretty(pretty bool) *CatRecoveryService {
	s.pretty = pretty
	return s
}

// Local indicates whether to return local information. If it is true,
// we do not retrieve the state from master node (default: false).
func (s *CatRecoveryService) Local(local bool) *CatRecoveryService {
	s.local = &local
	return s
}

// MasterTimeout specifies an explicit operation timeout for connection to master node.
func (s *CatRecoveryService) MasterTimeout(masterTimeout string) *CatRecoveryService {
	s.masterTimeout = masterTimeout
	return s
}

// Timeout specifies an explicit operation timeout.
func (s *CatRecoveryService) Timeout(timeout string) *CatRecoveryService {
	s.timeout = timeout
	return s
}

// buildURL builds the URL for the operation.
func (s *CatRecoveryService) buildURL() (string, url.Values, error) {
	// Build URL
	var path string

	path = "/_cat/recovery"

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
func (s *CatRecoveryService) Validate() error {
	return nil
}

// Do executes the operation.
func (s *CatRecoveryService) Do(ctx context.Context) (*CatRecoveryResponse, error) {
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
	ret := new(CatRecoveryResponse)
	if err := s.client.decoder.Decode(res.Body, &ret.Recoverys); err != nil {
		return nil, err
	}
	return ret, nil
}

// CatRecoveryResponse is the response of CatRecoveryService.Do
type CatRecoveryResponse struct {
	Recoverys []*recoveryRecord
}

//index                      shard time    type       stage source_host   target_host   repository snapshot files files_percent bytes        bytes_percent total_files total_bytes  translog translog_percent total_translog
//credit_20180416            0     13      replica    done  10.89.83.34   10.89.83.34   n/a        n/a      1     100.0%        130          100.0%        1           130          0        100.0%           0
type recoveryRecord struct {
	Index           string `json:"index"`
	Shard           string `json:"shard"`
	Time            string `json:"time"`
	Type            string `json:"type"`
	Stage           string `json:"stage"`
	SourceHost      string `json:"source_host"`
	TargetHost      string `json:"target_host"`
	Repository      string `json:"repository"`
	Snapshot        string `json:"snapshot"`
	Files           string `json:"files"`
	FilesPercent    string `json:"files_percent"`
	Bytes           string `json:"bytes"`
	BytesPercent    string `json:"bytes_percent"`
	TotalFiles      string `json:"total_files"`
	TotalBytes      string `json:"total_bytes"`
	Translog        string `json:"translog"`
	TranslogPercent string `json:"translog_percent"`
	TotalTranslog   string `json:"total_translog"`
}
