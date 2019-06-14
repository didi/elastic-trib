// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"context"
	"fmt"
	"net/url"
)

// ClusterGetSettingsService allows to retrieve settings of one
// or more indices.
//
// See https://www.elastic.co/guide/en/elasticsearch/reference/6.0/cluster-update-settings.html
// for more details.
type ClusterGetSettingsService struct {
	client        *Client
	pretty        bool
	flatSettings  *bool
	masterTimeout string
	timeout       string
}

// NewClusterGetSettingsService creates a new ClusterGetSettingsService.
func NewClusterGetSettingsService(client *Client) *ClusterGetSettingsService {
	return &ClusterGetSettingsService{
		client: client,
	}
}

// Pretty cluster settings that the JSON response be indented and human readable.
func (s *ClusterGetSettingsService) Pretty(pretty bool) *ClusterGetSettingsService {
	s.pretty = pretty
	return s
}

// FlatSettings cluster whether to return settings in flat format (default: false).
func (s *ClusterGetSettingsService) FlatSettings(flatSettings bool) *ClusterGetSettingsService {
	s.flatSettings = &flatSettings
	return s
}

// MasterTimeout specifies an explicit operation timeout for connection to master node.
func (s *ClusterGetSettingsService) MasterTimeout(masterTimeout string) *ClusterGetSettingsService {
	s.masterTimeout = masterTimeout
	return s
}

// Timeout specifies an explicit operation timeout.
func (s *ClusterGetSettingsService) Timeout(timeout string) *ClusterGetSettingsService {
	s.timeout = timeout
	return s
}

// buildURL builds the URL for the operation.
func (s *ClusterGetSettingsService) buildURL() (string, url.Values, error) {
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
func (s *ClusterGetSettingsService) Validate() error {
	return nil
}

// Do executes the operation.
func (s *ClusterGetSettingsService) Do(ctx context.Context) (*ClusterGetSettingsResponse, error) {
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
	ret := new(ClusterGetSettingsResponse)
	if err := s.client.decoder.Decode(res.Body, &ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// ClusterGetSettingsResponse is the response of ClusterGetSettingsService.Do.
type ClusterGetSettingsResponse struct {
	Persistent map[string]interface{} `json:"persistent"`
	Transient  map[string]interface{} `json:"transient"`
}
