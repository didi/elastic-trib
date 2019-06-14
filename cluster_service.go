package main

import (
	"fmt"
	"os"
	"os/user"
	"path"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/olivere/elastic"
	"github.com/spf13/viper"
)

// NewElasticClient use cluster name / host / localhost.
func NewElasticClient(context *cli.Context) (*elastic.Client, error) {
	var options []elastic.ClientOptionFunc
	var addr, esAddr string

	if context.GlobalString("host") != "" {
		addr = context.GlobalString("host")
	} else if context.GlobalString("cluster") != "" {
		cluster := context.GlobalString("cluster")
		addr = viper.GetString("clusters." + cluster)

		if addr == "" {
			return nil, fmt.Errorf("get error cluster name: %s in cfgFile:(%s)", cluster, viper.ConfigFileUsed())
		}
	} else {
		addr = "http://127.0.0.1:9200"
	}

	esAddr = checkURLScheme(addr, "http")
	if esAddr != "" {
		options = append(options, elastic.SetURL(esAddr))
	} else {
		return nil, fmt.Errorf("Es addr checkURLScheme failed: %s", addr)
	}

	basicAuth := context.GlobalString("http-auth")
	if basicAuth != "" {
		auths := strings.Split(basicAuth, ":")
		if len(auths) == 2 {
			options = append(options, elastic.SetBasicAuth(auths[0], auths[1]))
		} else {
			logrus.Warnf("Set basic auth failed: %s", auths)
		}
	}

	options = append(options, elastic.SetSniff(false))

	usr, err := user.Current()
	if err != nil {
		fmt.Println(err)
	}

	// log operation record
	now := time.Now()
	record := fmt.Sprintf("%s %s %s %s %s %s %s", now, usr.Gid, usr.HomeDir, usr.Name, usr.Uid, usr.Username, os.Args)
	logOprationRecord(record)

	// Create a client and connect to addr.
	return elastic.NewClient(options...)
}

// logOprationRecord
func logOprationRecord(content string) {
	// Create a new instance of the logger. You can have any number of instances.
	var log = logrus.New()

	// The API for setting attributes is a little different than the package level
	// exported logger. See Godoc.
	log.Out = os.Stdout

	filename := viper.GetString("authlog")
	if filename == "" {
		filename = "elastic-auth.log"
	}
	filepath := path.Join(GetCurrPath(), filename)
	fd, err := os.OpenFile(filepath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err == nil {
		log.Out = fd
	} else {
		log.Info("Failed to log to file, using default stderr")
	}

	log.Info(content)
}
