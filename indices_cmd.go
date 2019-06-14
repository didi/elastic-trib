package main

import (
	ctx "context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/olivere/elastic"
)

var indicesCommand = cli.Command{
	Name:    "indices",
	Aliases: []string{"i"},
	Usage:   "Elastic indices operation cmd.",
	Subcommands: []cli.Command{
		// indices cat
		indicesCatCommand,
		// indices list
		indicesListCommand,
		// indices cat shards
		indicesCatShardsCommand,
		// indices open
		indicesOpenCommand,
		// indices close
		indicesCloseCommand,
		// indices delete
		indicesDeleteCommand,
		// indices settings
		indicesSettingsCommand,
		// indices template
		indicesTemplateCommand,
		// indices cat aliases
		indicesCatAliasesCommand,
	},
}

//cat alias
var indicesCatAliasesCommand = cli.Command{
	Name:        "alias",
	Usage:       "cat indices alias list from elastic cluster.",
	ArgsUsage:   `[-i "alias* or alias1,alias2"]`,
	Description: `Display the cat indices of elastic cluster.`,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "format",
			Value: "text",
			Usage: "set the format of output('text' (default),or 'json')",
		},
		cli.StringFlag{
			Name:  "alias",
			Value: "",
			Usage: "set alias for query(alias1,alias2).",
		},
	},
	Action: func(context *cli.Context) error {
		return indicesCatAliasCmd(context)
	},
}

func indicesCatAliasCmd(context *cli.Context) error {
	client, err := NewElasticClient(context)
	if err != nil {
		return err
	}
	defer client.Stop()

	ctx := ctx.Background()
	aliasService := client.CatAliasService()
	aliases := context.String("alias")
	if aliases != "" {
		aliaesarray := strings.Split(aliases, ",")
		if len(aliaesarray) > 0 {
			aliasService.Alias(aliaesarray...)
		}
	}

	res, err := aliasService.Do(ctx)
	if err != nil {
		return err
	}
	format := context.String("format")
	switch format {
	case "text":
		printAliasesList(res)
	case "json":
		jsonStr, err := json.Marshal(res)
		if err != nil {
			return err
		}
		fmt.Println(jsonPrettyPrint(string(jsonStr)))
	default:
		return fmt.Errorf("unknows format %q", format)
	}

	return nil
}

// aliases alias index filter routing.index routing.search
func printAliasesList(CatAliasResponse *elastic.CatAliasResponse) error {
	if CatAliasResponse == nil {
		return nil
	}

	display := NewTableDisplay()
	display.AddRow([]string{"alias", "index", "filter", "routingIndex", "routingSearch"})
	for _, indiceInfo := range CatAliasResponse.Aliases {
		display.AddRow([]string{
			indiceInfo.Alias,
			indiceInfo.Index,
			indiceInfo.Filter,
			indiceInfo.Routingindex,
			indiceInfo.Routingsearch})
	}

	display.Flush()
	return nil
}

// cat
var indicesCatCommand = cli.Command{
	Name:        "cat",
	Usage:       "cat indices list from elastic cluster.",
	ArgsUsage:   `[-i "indices* or index1,index2"]`,
	Description: `Display the cat indices of elastic cluster.`,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "format",
			Value: "text",
			Usage: "set the format of output('text' (default), or 'json').",
		},
		cli.StringFlag{
			Name:  "indices, i",
			Value: "",
			Usage: "set indices for query (index1,index2).",
		},
	},
	Action: func(context *cli.Context) error {
		return indicesCatCmd(context)
	},
}

func indicesCatCmd(context *cli.Context) error {
	// Create a client and connect to addr.
	client, err := NewElasticClient(context)
	if err != nil {
		return err
	}
	defer client.Stop()

	// Starting with elastic.v5, you must pass a context to execute each service
	ctx := ctx.Background()

	catService := client.CatIndicesService()

	indices := context.String("indices")
	if indices != "" {
		iarray := strings.Split(indices, ",")
		if len(iarray) > 0 {
			catService.Index(iarray...)
		}
	}

	res, err := catService.Do(ctx)
	if err != nil {
		return err
	}

	format := context.String("format")
	switch format {
	case "text":
		printIndicesList(res)
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

// health status index  pri rep docs.count docs.deleted store.size pri.store.size
func printIndicesList(indicesInfoResponse *elastic.CatIndicesResponse) error {
	if indicesInfoResponse == nil {
		return nil
	}

	display := NewTableDisplay()
	display.AddRow([]string{"health", "status", "index", "uuid", "pri", "rep", " count", "deleted", "size", "storeSize"})
	for _, indice := range indicesInfoResponse.Indices {
		display.AddRow([]string{
			indice.Health,
			indice.Status,
			indice.Index,
			indice.UUID,
			indice.Pri,
			indice.Rep,
			indice.Count,
			indice.Deleted,
			indice.Size,
			indice.StoreSize})
	}

	display.Flush()
	return nil
}

// shards
var indicesCatShardsCommand = cli.Command{
	Name:        "shards",
	Aliases:     []string{"s"},
	Usage:       "Display the cat shards of elastic cluster.",
	ArgsUsage:   `[-i "indices* or index1,index2"]`,
	Description: `get cat shards from elastic cluster.`,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "format",
			Value: "text",
			Usage: "set the format of output('text' (default), or 'json').",
		},
		cli.StringFlag{
			Name:  "indices, i",
			Value: "",
			Usage: "set indices for query (index1,index2).",
		},
	},
	Action: func(context *cli.Context) error {
		return indicesCatShardsCmd(context)
	},
}

func indicesCatShardsCmd(context *cli.Context) error {
	// Create a client and connect to addr.
	client, err := NewElasticClient(context)
	if err != nil {
		return err
	}
	defer client.Stop()

	// Starting with elastic.v5, you must pass a context to execute each service
	ctx := ctx.Background()

	catService := client.CatShardsService()

	indices := context.String("indices")
	if indices != "" {
		iarray := strings.Split(indices, ",")
		if len(iarray) > 0 {
			catService.Index(iarray...)
		}
	}

	res, err := catService.Do(ctx)
	if err != nil {
		return err
	}

	format := context.String("format")
	switch format {
	case "text":
		printShardsList(res)
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

// index shard prirep state docs store ip node
func printShardsList(shardsInfoResponse *elastic.CatShardsResponse) error {
	if shardsInfoResponse == nil {
		return nil
	}

	display := NewTableDisplay()
	display.AddRow([]string{"index", "shard", "prirep", "state", "docs", "store", "ip", "node"})
	for _, shard := range shardsInfoResponse.Shards {
		display.AddRow([]string{
			shard.Index,
			shard.Shard,
			shard.Prirep,
			shard.State,
			shard.Docs,
			shard.Store,
			shard.Ip,
			shard.Node})
	}

	display.Flush()
	return nil
}

// list
var indicesListCommand = cli.Command{
	Name:        "list",
	Usage:       "Display the indices list of elastic cluster.",
	Description: `get indices list from elastic cluster.`,
	Action: func(context *cli.Context) error {
		return indicesListCmd(context)
	},
}

func indicesListCmd(context *cli.Context) error {
	// Create a client and connect to addr.
	client, err := NewElasticClient(context)
	if err != nil {
		return err
	}
	defer client.Stop()

	// Starting with elastic.v5, you must pass a context to execute each service
	//ctx := ctx.Background()

	res, err := client.IndexNames()
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

// open            indicesName
var indicesOpenCommand = cli.Command{
	Name:        "open",
	Usage:       "The command open the elasticsearch indices.",
	ArgsUsage:   `indicesName`,
	Description: `open the elasticsearch indices.`,
	Action: func(context *cli.Context) error {
		if context.NArg() != 1 {
			fmt.Printf("Incorrect Usage.\n\n")
			cli.ShowCommandHelp(context, "open")
			logrus.Fatalf("Must provide indicesName for open command!")
		}

		return indicesOpenCmd(context)
	},
}

func indicesOpenCmd(context *cli.Context) error {
	var indicesName string

	if indicesName = context.Args().Get(0); indicesName == "" {
		return errors.New("please check indicesName for open command")
	}

	// Create a client and connect to addr.
	client, err := NewElasticClient(context)
	if err != nil {
		return err
	}
	defer client.Stop()

	// Starting with elastic.v5, you must pass a context to execute each service
	ctx := ctx.Background()

	res, err := client.OpenIndex(indicesName).Do(ctx)
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

// close            indicesName
var indicesCloseCommand = cli.Command{
	Name:        "close",
	Usage:       "Close the elasticsearch indices.",
	ArgsUsage:   `indicesName`,
	Description: `The command close the elasticsearch indices.`,
	Action: func(context *cli.Context) error {
		if context.NArg() != 1 {
			fmt.Printf("Incorrect Usage.\n\n")
			cli.ShowCommandHelp(context, "close")
			logrus.Fatalf("Must provide indicesName for close command!")
		}

		return indicesCloseCmd(context)
	},
}

func indicesCloseCmd(context *cli.Context) error {
	var indicesName string

	if indicesName = context.Args().Get(0); indicesName == "" {
		return errors.New("please check indicesName for close command")
	}

	// Create a client and connect to addr.
	client, err := NewElasticClient(context)
	if err != nil {
		return err
	}
	defer client.Stop()

	// Starting with elastic.v5, you must pass a context to execute each service
	ctx := ctx.Background()

	res, err := client.CloseIndex(indicesName).Do(ctx)
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

// delete            indicesName
var indicesDeleteCommand = cli.Command{
	Name:        "delete",
	Usage:       "Delete the elasticsearch indices.",
	Aliases:     []string{"del"},
	ArgsUsage:   `index1,index2`,
	Description: `The command delete the elasticsearch indices.`,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "yes, y",
			Usage: "Answer delete indices conform.",
		},
	},
	Action: func(context *cli.Context) error {
		if context.NArg() != 1 {
			fmt.Printf("Incorrect Usage.\n\n")
			cli.ShowCommandHelp(context, "delete")
			logrus.Fatalf("Must provide indicesName for delete command!")
		}

		return indicesDeleteCmd(context)
	},
}

func indicesDeleteCmd(context *cli.Context) error {
	var indicesName string
	if indicesName = context.Args().Get(0); indicesName == "" {
		return errors.New("please check indicesName for delete command")
	}

	indicesList := strings.Split(indicesName, ",")

	// Create a client and connect to addr.
	client, err := NewElasticClient(context)
	if err != nil {
		return err
	}
	defer client.Stop()

	// Starting with elastic.v5, you must pass a context to execute each service
	ctx := ctx.Background()
	fmt.Println(sgrBoldBlue("[Attention] Delete below indices? type (yes) to conform delete."))
	if !context.Bool("yes") {
		YesOrDie(strings.Join(indicesList, " "))
	}
	res, err := client.DeleteIndex(indicesList...).Do(ctx)
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

// settings            indicesName
var indicesSettingsCommand = cli.Command{
	Name:        "settings",
	Usage:       "Get settings of the elasticsearch indices.",
	Aliases:     []string{"set"},
	ArgsUsage:   `index1,index2`,
	Description: `The command get settings of the elasticsearch indices.`,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "get, g",
			Usage: "get the settings of indices(index1,index2).",
		},
		cli.StringFlag{
			Name:  "set, s",
			Value: "",
			Usage: "set the settings of indices(index1,index2): -s '{settings_json}'.",
		},
		cli.StringFlag{
			Name:  "replicas, r",
			Value: "",
			Usage: "set the number_of_replicas of indices(index1,index2): -r num.",
		},
	},
	Action: func(context *cli.Context) error {
		if context.NArg() != 1 {
			fmt.Printf("Incorrect Usage.\n\n")
			cli.ShowCommandHelp(context, "settings")
			logrus.Fatalf("Must provide indicesName for settings command!")
		}

		return indicesSettingsCmd(context)
	},
}

func indicesSettingsCmd(context *cli.Context) error {
	var indicesName string
	if indicesName = context.Args().Get(0); indicesName == "" {
		return errors.New("please check indicesName for settings command")
	}

	indicesList := strings.Split(indicesName, ",")

	// Create a client and connect to addr.
	client, err := NewElasticClient(context)
	if err != nil {
		return err
	}
	defer client.Stop()

	// Starting with elastic.v5, you must pass a context to execute each service
	ctx := ctx.Background()

	if context.Bool("get") {
		indexGetSetting := client.IndexGetSettings(indicesList...)

		res, err := indexGetSetting.FlatSettings(true).Do(ctx)
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

		indexPutSetting := client.IndexPutSettings(indicesList...)
		res, err := indexPutSetting.BodyJson(str).Do(ctx)
		if err != nil {
			return err
		}
		jsonStr, err := json.Marshal(res)
		if err != nil {
			return err
		}
		fmt.Println(jsonPrettyPrint(string(jsonStr)))

	} else if str := context.String("replicas"); str != "" {
		if num, err := strconv.Atoi(str); err != nil || num < 0 {
			return fmt.Errorf("Invalid replicas num: %s", str)
		}
		jsonStr := fmt.Sprintf("{\"index.number_of_replicas\": \"%s\"}", str)

		if !isJSON(jsonStr) {
			return fmt.Errorf("'%s' is not a json string", jsonStr)
		}

		indexPutSetting := client.IndexPutSettings(indicesList...)
		res, err := indexPutSetting.BodyJson(jsonStr).Do(ctx)
		if err != nil {
			return err
		}

		jsonRes, err := json.Marshal(res)
		if err != nil {
			return err
		}
		fmt.Println(jsonPrettyPrint(string(jsonRes)))

	} else {
		cli.ShowCommandHelp(context, "settings")
		return fmt.Errorf("indices settings must provide -g or -s parameters")
	}

	return nil
}

// template
var indicesTemplateCommand = cli.Command{
	Name:        "template",
	Usage:       "Get template of the elasticsearch indices.",
	Aliases:     []string{"tpl"},
	ArgsUsage:   `tpl1,tpl2`,
	Description: `The command get template of the elasticsearch indices.`,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "get, g",
			Usage: "get the template of templates(tpl1,tpl2).",
		},
		cli.StringFlag{
			Name:  "set, s",
			Value: "",
			Usage: "set the template of templates(tpl1,tpl2): -s '{settings_json}'.",
		},
		cli.StringFlag{
			Name:  "replicas, r",
			Value: "",
			Usage: "set the number_of_replicas of template(tpl1,tpl2): -r num.",
		},
	},
	Action: func(context *cli.Context) error {
		if context.NArg() != 1 {
			fmt.Printf("Incorrect Usage.\n\n")
			cli.ShowCommandHelp(context, "template")
			logrus.Fatalf("Must provide templateName for template command!")
		}

		return indicesTemplateCmd(context)
	},
}

func indicesTemplateCmd(context *cli.Context) error {
	var templatesName string
	if templatesName = context.Args().Get(0); templatesName == "" {
		return errors.New("please check templatesName for template command")
	}

	templatesList := strings.Split(templatesName, ",")

	// Create a client and connect to addr.
	client, err := NewElasticClient(context)
	if err != nil {
		return err
	}
	defer client.Stop()

	// Starting with elastic.v5, you must pass a context to execute each service
	ctx := ctx.Background()

	if context.Bool("get") {
		indexGetTemplate := client.IndexGetTemplate(templatesList...)

		res, err := indexGetTemplate.FlatSettings(true).Do(ctx)
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

		indexPutTemplate := client.IndexPutTemplate(templatesList[0])
		res, err := indexPutTemplate.BodyJson(str).Do(ctx)
		if err != nil {
			return err
		}
		jsonStr, err := json.Marshal(res)
		if err != nil {
			return err
		}
		fmt.Println(jsonPrettyPrint(string(jsonStr)))

	} else if str := context.String("replicas"); str != "" {
		if num, err := strconv.Atoi(str); err != nil || num < 0 {
			return fmt.Errorf("Invalid replicas num: %s", str)
		}

		indexGetTemplate := client.IndexGetTemplate(templatesList[0])

		res, err := indexGetTemplate.FlatSettings(true).Do(ctx)
		if err != nil {
			return err
		}
		jsonStr, err := json.Marshal(res)
		if err != nil {
			return err
		}

		out := map[string]interface{}{}
		json.Unmarshal([]byte(jsonStr), &out)

		fmt.Println(out)
		// jsonStr := fmt.Sprintf("{\"index.number_of_replicas\": \"%s\"}", str)

		// if !isJSON(jsonStr) {
		// 	return fmt.Errorf("'%s' is not a json string", jsonStr)
		// }

		// indexPutTemplate := client.IndexPutTemplate(templatesList[0])
		// res, err := indexPutTemplate.BodyJson(jsonStr).Do(ctx)
		// if err != nil {
		// 	return err
		// }

		// jsonRes, err := json.Marshal(res)
		// if err != nil {
		// 	return err
		// }
		// fmt.Println(jsonPrettyPrint(string(jsonRes)))

	} else {
		cli.ShowCommandHelp(context, "template")
		return fmt.Errorf("indices template must provide -g or -s parameters")
	}

	return nil
}
