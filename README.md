# **elastic-trib** <sup><sub>_elastic cluster command line tool._</sub></sup>

Create and administrate your [Elastic Cluster][cluster-tutorial] from the Command Line.

## Dependencies

* [elastic][]
* [cli][]
* [logrus][]

Dependencies are handled by [godep][], simple install it and type `godep restore` to fetch them.

## Install

#### Restore project env in first build
```console
$ git clone https://github.com/soarpenguin/elastic-trib.git
$ cd elastic-trib
$ make godep
$ make bin
$ PROG=./elastic-trib source ./autocomplete/bash_autocomplete
```

#### Build the code
```console
$ cd elastic-trib
$ make bin
```

## Usage

```console
NAME:
   elastic-trib - Elasticsearch Cluster command line utility.

For cluster/indices/nodes/tasks operation etc, you must specify the cluster name or host:port
of any working node in the cluster.

USAGE:
   elastic-trib [global options] command [command options] [arguments...]

VERSION:
   v0.1.d6ea303

COMMANDS:
     cluster, c  Elastic cluster operation cmd.
     indices, i  Elastic indices operation cmd.
     nodes, n    Elastic nodes operation cmd.
     tasks, t    Elastic tasks operation cmd.
     help, h     Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --host value, -H value       a host of elastic node.
   --cluster value, -c value    appoint cluster name: elastic-log (host url in elastic-trib.yaml config).
   --config value               appoint config file name.(default: ./elastic-trib.yaml)
   --http-auth value, -A value  use basic authentication ex: user:pass.
   --log value                  set the log file path where internal debug information is written.
   --log-format value           set the format used by logs ('text' or 'json'). (default: "text")
   --debug                      enable debug output for logging.
   --help, -h                   show help
   --version, -v                print the version
```

[elastic]: https://github.com/olivere/elastic
[cli]: https://github.com/codegangsta/cli
[logrus]: https://github.com/Sirupsen/logrus
[godep]: https://github.com/tools/godep 
