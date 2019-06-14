package main

import (
	ctx "context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/codegangsta/cli"
	"github.com/olivere/elastic"
)

var tasksCommand = cli.Command{
	Name:    "tasks",
	Aliases: []string{"t"},
	Usage:   "Elastic tasks operation cmd.",
	Subcommands: []cli.Command{
		// tasks cancel
		tasksCancelCommand,
		// tasks get
		tasksGetCommand,
		// tasks list
		tasksListCommand,
		// tasks recovery
		tasksRecoveryCommand,
		// tasks pending
		tasksPendingCommand,
	},
}

// Tasks List Command
var tasksListCommand = cli.Command{
	Name:        "list",
	Aliases:     []string{"l"},
	Usage:       "Get tasks list from elastic cluster.",
	Description: `get tasks list from elastic cluster.`,
	Action: func(context *cli.Context) error {
		return tasksListCmd(context)
	},
}

func tasksListCmd(context *cli.Context) error {
	// Create a client and connect to addr.
	client, err := NewElasticClient(context)
	if err != nil {
		return err
	}
	defer client.Stop()

	// Starting with elastic.v5, you must pass a context to execute each service
	ctx := ctx.Background()

	res, err := client.TasksList().Do(ctx)
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

// Tasks Get Command
var tasksGetCommand = cli.Command{
	Name:        "get",
	Aliases:     []string{"g"},
	Usage:       "Get tasks info from elastic cluster.",
	Description: `get tasks info from elastic cluster.`,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "id",
			Usage: "tasks get id.",
		},
	},
	Action: func(context *cli.Context) error {
		return tasksGetCmd(context)
	},
}

func tasksGetCmd(context *cli.Context) error {
	// Create a client and connect to addr.
	client, err := NewElasticClient(context)
	if err != nil {
		return err
	}
	defer client.Stop()

	// Starting with elastic.v5, you must pass a context to execute each service
	ctx := ctx.Background()

	tget := client.TasksGetTask()
	id := context.String("id")
	if id != "" {
		if _, err := strconv.ParseInt(id, 10, 32); err != nil {
			tget.TaskId(id)
		}
	}

	res, err := tget.Do(ctx)
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

// Tasks Cancel Command
var tasksCancelCommand = cli.Command{
	Name:        "cancel",
	Aliases:     []string{"c"},
	Usage:       "Cancel running task in elastic cluster.",
	Description: `cancel running task in elastic cluster.`,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "id",
			Usage: "tasks list ids.",
		},
	},
	Action: func(context *cli.Context) error {
		return tasksCancelCmd(context)
	},
}

func tasksCancelCmd(context *cli.Context) error {
	// Create a client and connect to addr.
	client, err := NewElasticClient(context)
	if err != nil {
		return err
	}
	defer client.Stop()

	// Starting with elastic.v5, you must pass a context to execute each service
	ctx := ctx.Background()

	tcancel := client.TasksCancel()

	id := context.String("id")
	if id != "" {
		if idnum, err := strconv.ParseInt(id, 10, 32); err == nil {
			tcancel.TaskId(idnum)
		}
	}

	res, err := tcancel.Do(ctx)
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

// Tasks Recovery Command
var tasksRecoveryCommand = cli.Command{
	Name:        "recovery",
	Aliases:     []string{"r"},
	Usage:       "Get recovery task list from elastic cluster.",
	Description: `get recovery tasks list from elastic cluster.`,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "format",
			Value: "text",
			Usage: "set the format of output('text' (default), or 'json').",
		},
		cli.BoolFlag{
			Name:  "all, a",
			Usage: "deplay the all of recovery tasks.",
		},
	},
	Action: func(context *cli.Context) error {
		return tasksRecoveryCmd(context)
	},
}

func tasksRecoveryCmd(context *cli.Context) error {
	// Create a client and connect to addr.
	client, err := NewElasticClient(context)
	if err != nil {
		return err
	}
	defer client.Stop()

	// Starting with elastic.v5, you must pass a context to execute each service
	ctx := ctx.Background()

	res, err := client.CatRecoveryService().Do(ctx)
	if err != nil {
		return err
	}

	format := context.String("format")
	switch format {
	case "text":
		printRecoveryRecordList(res, context.Bool("all"))
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

//index shard time type stage source_host target_host repository snapshot files files_percent bytes bytes_percent total_files total_bytes  translog translog_percent total_translog
func printRecoveryRecordList(catRecoveryRecordResp *elastic.CatRecoveryResponse, all bool) error {
	if catRecoveryRecordResp == nil {
		return nil
	}

	display := NewTableDisplay()
	display.AddRow([]string{"index", "shard", "time", "type", "stage", "sourcehost", "targethost", "repository", "snapshot", "files", "filespercent", "bytes", "bytespercent", "totalfiles", "totalbytes", "translog", "translogpercent", "totaltranslog"})
	for _, recovery := range catRecoveryRecordResp.Recoverys {
		if !all && recovery.Stage == "done" {
			continue
		}
		display.AddRow([]string{
			recovery.Index,
			recovery.Shard,
			recovery.Time,
			recovery.Type,
			recovery.Stage,
			recovery.SourceHost,
			recovery.TargetHost,
			recovery.Repository,
			recovery.Snapshot,
			recovery.Files,
			recovery.FilesPercent,
			recovery.Bytes,
			recovery.BytesPercent,
			recovery.TotalFiles,
			recovery.TotalBytes,
			recovery.Translog,
			recovery.TranslogPercent,
			recovery.TotalTranslog})
	}

	display.Flush()
	return nil
}

// Tasks Pending Command
var tasksPendingCommand = cli.Command{
	Name:        "pending",
	Aliases:     []string{"p"},
	Usage:       "Get pending task list from elastic cluster.",
	Description: `get pending tasks list from elastic cluster.`,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "format",
			Value: "text",
			Usage: "set the format of output('text' (default), or 'json').",
		},
	},
	Action: func(context *cli.Context) error {
		return tasksPendingCmd(context)
	},
}

func tasksPendingCmd(context *cli.Context) error {
	// Create a client and connect to addr.
	client, err := NewElasticClient(context)
	if err != nil {
		return err
	}
	defer client.Stop()

	// Starting with elastic.v5, you must pass a context to execute each service
	ctx := ctx.Background()

	res, err := client.ClusterPendingTasksService().Do(ctx)
	if err != nil {
		return err
	}

	format := context.String("format")
	switch format {
	case "text":
		printPendingRecordList(res)
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

// printPendingRecordList print pending record list.
func printPendingRecordList(clusterPendingTasksResp *elastic.ClusterPendingTasksResponse) error {
	if clusterPendingTasksResp == nil {
		return nil
	}

	display := NewTableDisplay()
	display.AddRow([]string{"insertorder", "priority", "source", "timeInQueueMillis", "timeInQueue"})
	for _, task := range clusterPendingTasksResp.Tasks {
		display.AddRow([]string{
			strconv.Itoa(task.InsertOrder),
			task.Priority,
			task.Source,
			strconv.Itoa(task.TimeInQueueMillis),
			task.TimeInQueue})
	}

	display.Flush()
	return nil
}
