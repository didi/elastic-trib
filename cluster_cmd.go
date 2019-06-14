package main

import (
	ctx "context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/olivere/elastic"
	"github.com/spf13/viper"
)

var clusterCommand = cli.Command{
	Name:    "cluster",
	Aliases: []string{"c"},
	Usage:   "Elastic cluster operation cmd.",
	Subcommands: []cli.Command{
		// cluster health
		clusterHealthCommand,
		// cluster master
		clusterMasterCommand,
		// cluster state
		// clusterStateCommand,
		// cluster settings
		clusterSettingsCommand,
		// cluster list
		clusterListCommand,
		// cluster stats
		clusterStatsCommand,
	},
}

// state
var clusterStateCommand = cli.Command{
	Name:        "state",
	Aliases:     []string{"s"},
	Usage:       "get state from elastic cluster.",
	Description: `Display the state of elastic cluster.`,
	Action: func(context *cli.Context) error {
		return clusterStateCmd(context)
	},
}

func clusterStateCmd(context *cli.Context) error {
	// Create a client and connect to addr.
	client, err := NewElasticClient(context)
	if err != nil {
		return err
	}
	defer client.Stop()

	// Starting with elastic.v5, you must pass a context to execute each service
	ctx := ctx.Background()

	res, err := client.ClusterState().Do(ctx)
	if err != nil {
		return err
	}

	jsonStr, err := json.Marshal(res)
	if err != nil {
		return err
	}
	fmt.Println(jsonPrettyPrint(string(jsonStr)))

	return nil
}

// health
var clusterHealthCommand = cli.Command{
	Name:        "health",
	Aliases:     []string{"h"},
	Usage:       "get health status from elastic cluster.",
	Description: `Display the health status of elastic cluster.`,
	Action: func(context *cli.Context) error {
		return clusterHealthCmd(context)
	},
}

func clusterHealthCmd(context *cli.Context) error {
	// Create a client and connect to addr.
	client, err := NewElasticClient(context)
	if err != nil {
		return err
	}
	defer client.Stop()

	// Starting with elastic.v5, you must pass a context to execute each service
	ctx := ctx.Background()

	res, err := client.ClusterHealth().Do(ctx)
	if err != nil {
		return err
	}

	jsonStr, err := json.Marshal(res)
	if err != nil {
		return err
	}
	fmt.Println(jsonPrettyPrint(string(jsonStr)))

	return nil
}

// master
var clusterMasterCommand = cli.Command{
	Name:        "master",
	Aliases:     []string{"m"},
	Usage:       "get the master node of elastic cluster.",
	Description: `Display the master node of elastic cluster.`,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "format",
			Value: "text",
			Usage: "set the format of output('text' (default), or 'json').",
		},
	},
	Action: func(context *cli.Context) error {
		return clusterMasterCmd(context)
	},
}

func clusterMasterCmd(context *cli.Context) error {
	// Create a client and connect to addr.
	client, err := NewElasticClient(context)
	if err != nil {
		return err
	}
	defer client.Stop()

	// Starting with elastic.v5, you must pass a context to execute each service
	ctx := ctx.Background()

	res, err := client.CatMasterService().Do(ctx)
	if err != nil {
		return err
	}

	format := context.String("format")
	switch format {
	case "text":
		printMasterList(res)
	case "json":
		jsonStr, err := json.Marshal(res)
		if err != nil {
			return err
		}
		fmt.Println(jsonPrettyPrint(string(jsonStr)))
	default:
		return fmt.Errorf("unknown format %q", context.String("format"))
	}

	return nil
}

// id          host          ip          node
func printMasterList(masterInfoResponse *elastic.CatMasterResponse) error {
	if masterInfoResponse == nil {
		return nil
	}

	display := NewTableDisplay()
	display.AddRow([]string{"id", "host", "ip", "node"})
	for _, masterInfo := range masterInfoResponse.Masters {
		display.AddRow([]string{masterInfo.ID, masterInfo.Host, masterInfo.IP, masterInfo.Node})
	}

	display.Flush()
	return nil
}

// settings
var clusterSettingsCommand = cli.Command{
	Name:        "settings",
	Usage:       "get or set the settings of the elasticsearch cluster.",
	Aliases:     []string{"set"},
	ArgsUsage:   `[-g or -s '{settings_json}']`,
	Description: `The command get or set the settings of the elasticsearch cluster.`,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "get, g",
			Usage: "get the settings of cluster.",
		},
		cli.StringFlag{
			Name:  "set, s",
			Value: "",
			Usage: "set the settings of cluster(-s '{settings_json}').",
		},
		cli.StringFlag{
			Name:  "file, f",
			Value: "",
			Usage: "set the settings of cluster from file content.",
		},
	},
	Action: func(context *cli.Context) error {
		return clusterSettingsCmd(context)
	},
}

func clusterSettingsCmd(context *cli.Context) error {
	// Create a client and connect to addr.
	client, err := NewElasticClient(context)
	if err != nil {
		return err
	}
	defer client.Stop()

	// Starting with elastic.v5, you must pass a context to execute each service
	ctx := ctx.Background()

	if context.Bool("get") {
		clusterGetSetting := client.ClusterGetSettings()

		res, err := clusterGetSetting.FlatSettings(true).Do(ctx)
		if err != nil {
			return err
		}
		jsonStr, err := json.Marshal(res)
		if err != nil {
			return err
		}
		fmt.Println(jsonPrettyPrint(string(jsonStr)))

	} else if str := context.String("set"); str != "" {
		str = strings.Trim(str, " ")
		if !isJSON(str) {
			return fmt.Errorf("'%s' is not a json string", str)
		}

		clusterPutSetting := client.ClusterPutSettings()
		res, err := clusterPutSetting.FlatSettings(true).BodyJson(str).Do(ctx)
		if err != nil {
			return err
		}

		jsonStr, err := json.Marshal(res)
		if err != nil {
			return err
		}
		fmt.Println(jsonPrettyPrint(string(jsonStr)))

	} else if fileName := context.String("file"); fileName != "" {
		var content string

		if fileName == "-" {
			data, err := ioutil.ReadAll(os.Stdin)
			if err != nil {
				return err
			}
			content = string(data)
		} else {
			data, err := ioutil.ReadFile(fileName)
			if err != nil {
				return err
			}
			content = string(data)
		}

		if !isJSON(string(content)) {
			return fmt.Errorf("'%s' is not a json string", string(content))
		}

		clusterPutSetting := client.ClusterPutSettings()
		res, err := clusterPutSetting.FlatSettings(true).BodyJson(string(content)).Do(ctx)
		if err != nil {
			return err
		}

		jsonStr, err := json.Marshal(res)
		if err != nil {
			return err
		}
		fmt.Println(jsonPrettyPrint(string(jsonStr)))

	} else {
		cli.ShowCommandHelp(context, "settings")
		return fmt.Errorf("cluster settings must provide one of (-g, -s str, -f file) parameters")
	}

	return nil
}

func putClusterSettings(setStr string) {

}

// list cluster from elastic.yaml
var clusterListCommand = cli.Command{
	Name:        "list",
	Usage:       "get name list of the elasticsearch cluster.",
	Aliases:     []string{"l"},
	Description: `The command get name list of the elasticsearch cluster.`,
	Action: func(context *cli.Context) error {
		return clusterListCmd(context)
	},
}

func clusterListCmd(context *cli.Context) error {
	// get cluster mappings
	if clusters := viper.GetStringMapString("clusters"); clusters != nil {
		var keys []string
		for key := range clusters {
			keys = append(keys, key)
		}

		sort.Strings(keys)
		fmt.Println("|---Clusters:")
		for _, key := range keys {
			fmt.Printf("|->  %-25s:\t%s\n", key, clusters[key])
		}
	} else {
		return fmt.Errorf("get cluster list from %s error", viper.ConfigFileUsed())
	}

	return nil
}

// stats
var clusterStatsCommand = cli.Command{
	Name:        "stats",
	Usage:       "get stats from elastic cluster.",
	Description: `Display the stats info of elastic cluster.`,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "indices, i",
			Usage: "get indices stats from elastic cluster.",
		},
		cli.BoolFlag{
			Name:  "nodes, n",
			Usage: "get nodes stats from elastic cluster.",
		},
	},
	Action: func(context *cli.Context) error {
		return clusterStatsCmd(context)
	},
}

func clusterStatsCmd(context *cli.Context) error {
	// Create a client and connect to addr.
	client, err := NewElasticClient(context)
	if err != nil {
		return err
	}
	defer client.Stop()

	ctx := ctx.Background()
	res, err := client.ClusterStats().Do(ctx)
	if err != nil {
		return err
	}

	jsonStr, err := json.Marshal(res)
	if err != nil {
		return err
	}

	if context.Bool("indices") {
		var clusterStatsIndices elastic.ClusterStatsResponse
		if err := json.Unmarshal([]byte(jsonStr), &clusterStatsIndices); err != nil {
			return err
		}

		indices, err := json.Marshal(clusterStatsIndices.Indices)
		if err != nil {
			return err
		}
		fmt.Println(jsonPrettyPrint(string(indices)))

	} else if context.Bool("nodes") {
		var clusterStatsIndices elastic.ClusterStatsResponse
		if err := json.Unmarshal([]byte(jsonStr), &clusterStatsIndices); err != nil {
			return err
		}

		nodes, err := json.Marshal(clusterStatsIndices.Nodes)
		if err != nil {
			return err
		}
		fmt.Println(jsonPrettyPrint(string(nodes)))
	} else {
		fmt.Println(jsonPrettyPrint(string(jsonStr)))
	}
	return nil
}
