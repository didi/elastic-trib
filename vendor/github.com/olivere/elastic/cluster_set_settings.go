// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"context"
	"fmt"
	"net/url"
)

// ClusterPutSettingsService allows to retrieve settings of cluster.
//
// See https://www.elastic.co/guide/en/elasticsearch/reference/6.0/cluster-update-settings.html
// for more details.
type ClusterPutSettingsService struct {
	client        *Client
	pretty        bool
	flatSettings  *bool
	masterTimeout string
	timeout       string
	bodyJson      interface{}
	bodyString    string
}

// NewClusterPutSettingsService creates a new ClusterPutSettingsService.
func NewClusterPutSettingsService(client *Client) *ClusterPutSettingsService {
	return &ClusterPutSettingsService{
		client: client,
	}
}

// Pretty cluster settings that the JSON response be indented and human readable.
func (s *ClusterPutSettingsService) Pretty(pretty bool) *ClusterPutSettingsService {
	s.pretty = pretty
	return s
}

// FlatSettings cluster whether to return settings in flat format (default: false).
func (s *ClusterPutSettingsService) FlatSettings(flatSettings bool) *ClusterPutSettingsService {
	s.flatSettings = &flatSettings
	return s
}

// MasterTimeout specifies an explicit operation timeout for connection to master node.
func (s *ClusterPutSettingsService) MasterTimeout(masterTimeout string) *ClusterPutSettingsService {
	s.masterTimeout = masterTimeout
	return s
}

// Timeout specifies an explicit operation timeout.
func (s *ClusterPutSettingsService) Timeout(timeout string) *ClusterPutSettingsService {
	s.timeout = timeout
	return s
}

// BodyJson is documented as: The cluster settings to be updated.
func (s *ClusterPutSettingsService) BodyJson(body interface{}) *ClusterPutSettingsService {
	s.bodyJson = body
	return s
}

// BodyString is documented as: The cluster settings to be updated.
func (s *ClusterPutSettingsService) BodyString(body string) *ClusterPutSettingsService {
	s.bodyString = body
	return s
}

// buildURL builds the URL for the operation.
func (s *ClusterPutSettingsService) buildURL() (string, url.Values, error) {
	var path string

	// Build URL
	path = "/_cluster/settings"

	// Add query string parameters
	params := url.Values{}
	if s.pretty {
		params.Set("pretty", "true")
	}
	if s.flatSettings != nil {
		params.Set("flat_settings", fmt.Sprintf("%v", *s.flatSettings))
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
func (s *ClusterPutSettingsService) Validate() error {
	return nil
}

// Do executes the operation.
func (s *ClusterPutSettingsService) Do(ctx context.Context) (*ClusterPutSettingsResponse, error) {
	// Check pre-conditions
	if err := s.Validate(); err != nil {
		return nil, err
	}

	// Get URL for request
	path, params, err := s.buildURL()
	if err != nil {
		return nil, err
	}

	// Setup HTTP request body
	var body interface{}
	if s.bodyJson != nil {
		body = s.bodyJson
	} else {
		body = s.bodyString
	}

	// Get HTTP response
	res, err := s.client.PerformRequest(ctx, PerformRequestOptions{
		Method: "PUT",
		Path:   path,
		Params: params,
		Body:   body,
	})
	if err != nil {
		return nil, err
	}

	// Return operation response
	ret := new(ClusterPutSettingsResponse)
	if err := s.client.decoder.Decode(res.Body, &ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// ClusterPutSettingsResponse is the response of ClusterPutSettingsService.Do.
type ClusterPutSettingsResponse struct {
	Acknowledged bool                   `json:"acknowledged"`
	Persistent   map[string]interface{} `json:"persistent"`
	Transient    map[string]interface{} `json:"transient"`
}
