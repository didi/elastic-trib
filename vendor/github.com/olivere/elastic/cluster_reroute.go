// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

// ClusterRerouteService execute a cluster reroute command.
//
// See https://www.elastic.co/guide/en/elasticsearch/reference/6.0/cluster-update-settings.html
// for more details.
type ClusterRerouteService struct {
	client        *Client
	pretty        bool
	flatSettings  *bool
	masterTimeout string
	timeout       string
	bodyJSON      interface{}
	bodyString    string
}
