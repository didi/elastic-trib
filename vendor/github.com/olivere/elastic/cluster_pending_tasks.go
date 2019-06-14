// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"context"
	"fmt"
	"net/url"

	"github.com/olivere/elastic/uritemplates"
)

// ClusterPendingTasksService is documented at
// https://www.elastic.co/guide/en/elasticsearch/reference/6.0/cluster-pending.html.
type ClusterPendingTasksService struct {
	client       *Client
	pretty       bool
	flatSettings *bool
	human        *bool
}

// NewClusterPendingTasksService creates a new ClusterPendingTasksService.
func NewClusterPendingTasksService(client *Client) *ClusterPendingTasksService {
	return &ClusterPendingTasksService{
		client: client,
	}
}

// FlatSettings is documented as: Return settings in flat format (default: false).
func (s *ClusterPendingTasksService) FlatSettings(flatSettings bool) *ClusterPendingTasksService {
	s.flatSettings = &flatSettings
	return s
}

// Human is documented as: Whether to return time and byte values in human-readable format..
func (s *ClusterPendingTasksService) Human(human bool) *ClusterPendingTasksService {
	s.human = &human
	return s
}

// Pretty indicates that the JSON response be indented and human readable.
func (s *ClusterPendingTasksService) Pretty(pretty bool) *ClusterPendingTasksService {
	s.pretty = pretty
	return s
}

// buildURL builds the URL for the operation.
func (s *ClusterPendingTasksService) buildURL() (string, url.Values, error) {
	// Build URL
	var err error
	var path string

	path, err = uritemplates.Expand("/_cluster/pending_tasks", map[string]string{})
	if err != nil {
		return "", url.Values{}, err
	}

	// Add query string parameters
	params := url.Values{}
	if s.pretty {
		params.Set("pretty", "true")
	}
	if s.flatSettings != nil {
		params.Set("flat_settings", fmt.Sprintf("%v", *s.flatSettings))
	}
	if s.human != nil {
		params.Set("human", fmt.Sprintf("%v", *s.human))
	}
	return path, params, nil
}

// Validate checks if the operation is valid.
func (s *ClusterPendingTasksService) Validate() error {
	return nil
}

// Do executes the operation.
func (s *ClusterPendingTasksService) Do(ctx context.Context) (*ClusterPendingTasksResponse, error) {
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
	ret := new(ClusterPendingTasksResponse)
	if err := s.client.decoder.Decode(res.Body, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// ClusterPendingTasksResponse is the response of ClusterPendingTasksService.Do.
type ClusterPendingTasksResponse struct {
	Tasks []*ClusterPendingTask `json:"tasks"`
}

// ClusterPendingTask is the struct of ClusterPendingTask
type ClusterPendingTask struct {
	InsertOrder       int    `json:"insert_order"`
	Priority          string `json:"priority"`
	Source            string `json:"source"`
	TimeInQueueMillis int    `json:"time_in_queue_millis"`
	TimeInQueue       string `json:"time_in_queue"`
}
