package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/spf13/viper"
)

// version will be the hash that the binary was built from
// and will be populated by the Makefile
var version = ""

// gitCommit will be the hash that the binary was built from
// and will be populated by the Makefile
var gitCommit = ""

const (
	// name holds the name of this program
	name  = "elastic-trib"
	usage = `Elasticsearch Cluster command line utility.

For cluster/indices/nodes/tasks operation etc, you must specify the cluster name or host:port
of any working node in the cluster.`
)

// runtimeFlags is the list of supported global command-line flags
var runtimeFlags = []cli.Flag{
	cli.StringFlag{
		Name:  "host, H",
		Usage: "a host of elastic node.",
	},
	cli.StringFlag{
		Name:  "cluster, c",
		Usage: "appoint cluster name: elastic-log (host url in elastic-trib.yaml config).",
	},
	cli.StringFlag{
		Name:  "config",
		Usage: "appoint config file name.(default: ./elastic-trib.yaml)",
	},
	cli.StringFlag{
		Name:  "http-auth, A",
		Usage: "use basic authentication ex: user:pass.",
	},
	cli.StringFlag{
		Name:  "log",
		Value: "",
		Usage: "set the log file path where internal debug information is written.",
	},
	cli.StringFlag{
		Name:  "log-format",
		Value: "text",
		Usage: "set the format used by logs ('text' or 'json').",
	},
	cli.BoolFlag{
		Name:  "debug",
		Usage: "enable debug output for logging.",
	},
}

// runtimeBeforeSubcommands is the function to run before command-line
// parsing occurs.
var runtimeBeforeSubcommands = beforeSubcommands

// runtimeCommandNotFound is the function to handle an invalid sub-command.
var runtimeCommandNotFound = commandNotFound

// runtimeCommands is all sub-command
var runtimeCommands = []cli.Command{
	clusterCommand,
	indicesCommand,
	nodesCommand,
	tasksCommand,
}

func beforeSubcommands(context *cli.Context) error {
	if context.GlobalBool("debug") {
		logrus.SetLevel(logrus.DebugLevel)
	}

	config := context.GlobalString("config")
	initConfig(config)

	if path := context.GlobalString("log"); path != "" {
		f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND|os.O_SYNC, 0666)
		if err != nil {
			return err
		}
		logrus.SetOutput(f)
	}

	switch context.GlobalString("log-format") {
	case "text":
		// retain logrus's default.
	case "json":
		logrus.SetFormatter(new(logrus.JSONFormatter))
	default:
		logrus.Fatalf("unknown log-format %q", context.GlobalString("log-format"))
	}
	return nil
}

// function called when an invalid command is specified which causes the
// runtime to error.
func commandNotFound(c *cli.Context, command string) {
	err := fmt.Errorf("invalid command %q", command)
	fatal(err)
}

// makeVersionString returns a multi-line string describing the runtime version.
func makeVersionString() string {
	v := []string{
		version,
	}
	if gitCommit != "" {
		v = append(v, fmt.Sprintf("commit: %s", gitCommit))
	}

	return strings.Join(v, "\n")
}

// initConfig reads in config file and ENV variables if set.
func initConfig(cfgFile string) {
	if cfgFile != "" {
		// Use config file from the flag.
		path, _ := filepath.Abs(cfgFile)
		viper.AddConfigPath(path)
		viper.SetConfigFile(cfgFile)
	} else {
		// Search config in current directory with name "elastic-trib.yaml".
		viper.AddConfigPath(GetCurrPath())
		cfgFile = filepath.Join(GetCurrPath(), "elastic-trib.yaml")
		viper.SetConfigFile(cfgFile)
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		logrus.Warn("Read config file failed:", err.Error())
	} else {
		logrus.Debug("Using config file:", viper.ConfigFileUsed())
	}
}

func main() {
	app := cli.NewApp()

	app.Name = name
	app.Writer = os.Stdout
	app.Usage = usage
	app.Version = makeVersionString()
	app.Flags = runtimeFlags
	app.EnableBashCompletion = true
	app.CommandNotFound = runtimeCommandNotFound
	app.Before = runtimeBeforeSubcommands
	app.Commands = runtimeCommands

	if err := app.Run(os.Args); err != nil {
		fatal(err)
	}
}
