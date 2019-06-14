package main

import (
	"bytes"
	ctx "context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"text/template"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/olivere/elastic"
)

var nodesCommand = cli.Command{
	Name:    "nodes",
	Aliases: []string{"n"},
	Usage:   "Elastic nodes operation cmd.",
	Subcommands: []cli.Command{
		// nodes attrs
		nodesAttrsCommand,
		// nodes cat
		nodesCatNodesCommand,
		// nodes alloc
		nodesCatAllocCommand,
		// nodes exclude
		nodesExcludeCommand,
		// nodes include
		nodesIncludeCommand,
		// nodes info
		nodesInfoCommand,
	},
}

// name    pid   attr     value
// node-0 19566 testattr test
// cat nodeattrs
var nodesAttrsCommand = cli.Command{
	Name:        "attrs",
	Aliases:     []string{"attr"},
	Usage:       "Display the nodes attrs of elastic cluster.",
	Description: `get nodes attrs from elastic cluster.`,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "format",
			Value: "text",
			Usage: "set the format of output('text' (default), or 'json').",
		},
	},
	Action: func(context *cli.Context) error {
		return nodesCatAttrsCmd(context)
	},
}

func nodesCatAttrsCmd(context *cli.Context) error {
	// Create a client and connect to addr.
	client, err := NewElasticClient(context)
	if err != nil {
		return err
	}
	defer client.Stop()

	// Starting with elastic.v5, you must pass a context to execute each service
	ctx := ctx.Background()

	res, err := client.CatNodeAttrsService().Do(ctx)
	if err != nil {
		return err
	}

	format := context.String("format")
	switch format {
	case "text":
		printNodeAttrsList(res)
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

func printNodeAttrsList(catNodeAttrsResp *elastic.CatNodeAttrsResponse) error {
	if catNodeAttrsResp == nil {
		return nil
	}

	display := NewTableDisplay()
	display.AddRow([]string{"node", "host", "ip", "attr", "value"})
	for _, nodeAttrInfo := range catNodeAttrsResp.NodeAttrs {
		display.AddRow([]string{
			nodeAttrInfo.Node,
			nodeAttrInfo.Host,
			nodeAttrInfo.IP,
			nodeAttrInfo.Attr,
			nodeAttrInfo.Value})
	}

	display.Flush()
	return nil
}

// cat nodes
var nodesCatNodesCommand = cli.Command{
	Name:        "cat",
	Aliases:     []string{"c"},
	Usage:       "Display the nodes of elastic cluster.",
	Description: `get nodes list from elastic cluster.`,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "format",
			Value: "text",
			Usage: "set the format of output('text' (default), or 'json').",
		},
	},
	Action: func(context *cli.Context) error {
		return nodesCatNodesCmd(context)
	},
}

func nodesCatNodesCmd(context *cli.Context) error {
	// Create a client and connect to addr.
	client, err := NewElasticClient(context)
	if err != nil {
		return err
	}
	defer client.Stop()

	// Starting with elastic.v5, you must pass a context to execute each service
	ctx := ctx.Background()

	res, err := client.CatNodesService().Do(ctx)
	if err != nil {
		return err
	}

	format := context.String("format")
	switch format {
	case "text":
		printNodesList(res)
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

// host          ip            heap.percent ram.percent  load node.role master name
// 100.69.145.39 100.69.145.39           49         100  1.91 -         -      bigdata-ser543
func printNodesList(catNodesResp *elastic.CatNodesResponse) error {
	if catNodesResp == nil {
		return nil
	}

	display := NewTableDisplay()
	display.AddRow([]string{"host", "ip", "heap.percent", "ram.percent", "load", "role", "master", "name"})
	for _, nodeInfo := range catNodesResp.Nodes {
		display.AddRow([]string{
			nodeInfo.Host,
			nodeInfo.IP,
			nodeInfo.HeapPercent,
			nodeInfo.RAMPercent,
			nodeInfo.Load,
			nodeInfo.NodeRole,
			nodeInfo.Master,
			nodeInfo.Name})
	}

	display.Flush()
	return nil
}

// cat alloc
var nodesCatAllocCommand = cli.Command{
	Name:        "allocation",
	Aliases:     []string{"alloc"},
	Usage:       "Display the nodes alloc of elastic cluster.",
	Description: `get nodes alloc from elastic cluster.`,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "format",
			Value: "text",
			Usage: "set the format of output('text' (default), or 'json').",
		},
	},
	Action: func(context *cli.Context) error {
		return nodesCatAllocCmd(context)
	},
}

func nodesCatAllocCmd(context *cli.Context) error {
	// Create a client and connect to addr.
	client, err := NewElasticClient(context)
	if err != nil {
		return err
	}
	defer client.Stop()

	// Starting with elastic.v5, you must pass a context to execute each service
	ctx := ctx.Background()

	res, err := client.CatAllocService().Do(ctx)
	if err != nil {
		return err
	}

	format := context.String("format")
	switch format {
	case "text":
		printNodeAllocList(res)
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

// host          ip            heap.percent ram.percent  load node.role master name
func printNodeAllocList(catAllocResp *elastic.CatAllocResponse) error {
	if catAllocResp == nil {
		return nil
	}

	display := NewTableDisplay()
	display.AddRow([]string{"shards", "indices", "used", "avail", "total", "percent", "host", "ip", "node"})
	for _, allocInfo := range catAllocResp.Allocs {
		display.AddRow([]string{
			allocInfo.Shards,
			allocInfo.Indices,
			allocInfo.Used,
			allocInfo.Avail,
			allocInfo.Total,
			allocInfo.Percent,
			allocInfo.Host,
			allocInfo.Ip,
			allocInfo.Node})
	}

	display.Flush()
	return nil
}

/////////////////////
const (
	excludeStr  = "cluster.routing.allocation.exclude._ip"
	settingTemp = `{
			"persistent": {
				"cluster.routing.rebalance.enable": "all",
				"cluster.routing.allocation.exclude._ip": "{{.}}" }
			}`
)

// get "cluster.routing.allocation.exclude._ip" from cluster settings.
func getSrcIPFromCluster(client *elastic.Client, ctx ctx.Context) (ipArray []string, err error) {
	clusterGetSetting := client.ClusterGetSettings()
	res, err := clusterGetSetting.FlatSettings(true).Do(ctx)
	if err != nil {
		return nil, err
	}

	resVal := reflect.ValueOf(*res)
	persis := resVal.FieldByName("Persistent")
	imap := persis.Interface()
	a := imap.(map[string]interface{})

	if value, ok := a[excludeStr]; ok {
		valueStr := strings.TrimSpace(value.(string))
		ipArray = strings.Split(valueStr, ",")
	}

	return ipArray, nil
}

// node exclude
var nodesExcludeCommand = cli.Command{
	Name:        "exclude",
	Aliases:     []string{"e"},
	Usage:       "Exclude hosts from elastic cluster.",
	ArgsUsage:   `ip1,ip2`,
	Description: `exclude hosts from elastic cluster.`,
	Action: func(context *cli.Context) error {
		if context.NArg() != 1 {
			fmt.Printf("Incorrect Usage.\n\n")
			cli.ShowCommandHelp(context, "exclude")
			logrus.Fatalf("Must provide ip addrs for exclude command. (Example: ip1,ip2)")
		}

		return nodesExcludeCmd(context)
	},
}

func nodesExcludeCmd(context *cli.Context) error {
	var ipList string
	var settingStr string
	var buffer bytes.Buffer

	if ipList = context.Args().Get(0); ipList == "" {
		return errors.New("please check ip addr for exclude command")
	}

	ipArray := strings.Split(ipList, ",")
	if err := checkIPAddr(ipArray); err != nil {
		return err
	}

	// Create a client and connect to addr.
	client, err := NewElasticClient(context)
	if err != nil {
		return err
	}
	defer client.Stop()

	// Starting with elastic.v5, you must pass a context to execute each service
	ctx := ctx.Background()

	srcIPArray, err := getSrcIPFromCluster(client, ctx)
	if err != nil {
		return err
	}

	ipArray = append(ipArray, srcIPArray...)
	ipArray = DeDuplicate(ipArray)
	ipStr := strings.TrimRight(strings.Join(ipArray, ","), ",")
	t := template.Must(template.New("settingTemp").Parse(settingTemp))
	if err := t.Execute(&buffer, &ipStr); err != nil {
		return err
	}

	settingStr = buffer.String()
	clusterPutSetting := client.ClusterPutSettings()
	ret, err := clusterPutSetting.FlatSettings(true).BodyJson(settingStr).Do(ctx)
	if err != nil {
		return err
	}

	jsonStr, err := json.Marshal(ret)
	if err != nil {
		return err
	}
	fmt.Println(jsonPrettyPrint(string(jsonStr)))

	return nil
}

// node include
var nodesIncludeCommand = cli.Command{
	Name:        "include",
	Aliases:     []string{"i"},
	Usage:       "Add hosts to elastic cluster.",
	ArgsUsage:   `ip1,ip2`,
	Description: `add hosts to elastic cluster.`,
	Action: func(context *cli.Context) error {
		if context.NArg() != 1 {
			fmt.Printf("Incorrect Usage.\n\n")
			cli.ShowCommandHelp(context, "include")
			logrus.Fatalf("Must provide ip addrs for include command. (Example: ip1,ip2)")
		}

		return nodesIncludeCmd(context)
	},
}

func nodesIncludeCmd(context *cli.Context) error {
	var ipList string
	var settingStr string
	var buffer bytes.Buffer
	var srcIPArray []string

	if ipList = context.Args().Get(0); ipList == "" {
		return errors.New("please check ip addr for include command")
	}

	rmArray := strings.Split(ipList, ",")
	if err := checkIPAddr(rmArray); err != nil {
		return err
	}

	// Create a client and connect to addr.
	client, err := NewElasticClient(context)
	if err != nil {
		return err
	}
	defer client.Stop()

	// Starting with elastic.v5, you must pass a context to execute each service
	ctx := ctx.Background()

	if srcIPArray, err = getSrcIPFromCluster(client, ctx); err != nil {
		return err
	}

	for _, ip := range rmArray {
		for i, elem := range srcIPArray {
			if ip == elem {
				srcIPArray = append(srcIPArray[:i], srcIPArray[i+1:]...)
			}
		}
	}

	srcIPArray = DeDuplicate(srcIPArray)
	ipStr := strings.TrimRight(strings.Join(srcIPArray, ","), ",")
	t := template.Must(template.New("settingTemp").Parse(settingTemp))
	if err := t.Execute(&buffer, &ipStr); err != nil {
		return err
	}

	settingStr = buffer.String()

	clusterPutSetting := client.ClusterPutSettings()
	ret, err := clusterPutSetting.FlatSettings(true).BodyJson(settingStr).Do(ctx)
	if err != nil {
		return err
	}
	jsonRes, err := json.Marshal(ret)
	if err != nil {
		return err
	}
	fmt.Println(jsonPrettyPrint(string(jsonRes)))

	return nil
}

// nodesInfo
var nodesInfoCommand = cli.Command{
	Name:        "info",
	Usage:       "Display the cluster nodeinfo of elastic cluster",
	Description: `get cluster nodeinfo from elastic cluster`,
	Action: func(context *cli.Context) error {
		return nodesInfoCmd(context)
	},
}

func nodesInfoCmd(context *cli.Context) error {
	// Create a client and connect to addr.
	client, err := NewElasticClient(context)

	defer client.Stop()

	ctx := ctx.Background()
	res, err := client.NodesInfo().Do(ctx)
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
