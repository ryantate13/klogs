![klogs](./klogs.svg)

# `klogs`

Tired of reading logs from kubernetes feeling flavorless? You and everyone else. Well, now there's `klogs` to sprinkle
your log entries with a healthy dose of sweetness. Klogs provides optional syntax highlighting for any pods that log
their output as JSON, with 46 available themes! Klogs allows you to stream logs from pods where:

* Labels match one or more label queries (`k=v`, `k!=v`), and/or
* Names match one or more pod name queries (configurable to match any search term vs. all, see [usage](#usage))

## Installation

```console
$ go install github.com/ryantate13/klogs@latest
```

## Usage

```console
$ klogs --help

klogs - Displays logs for kubernetes pods matching either a pod name query, a set of labels, or both

Usage: klogs [flags] [options] <search terms>...

Example: klogs -f service-one service-two # follow logs of all pods for service-one and service-two

Flags:
	-h | --help           Show this help message and quit
	-v | --version        Show the application version
	-a | --all            All pod name queries must match. Default is to show logs for pods where any name query matches
	   | --all-namespaces Query for pods in all namespaces
	   | --all-containers Get all containers' logs in the pod(s)
	-f | --follow         Follow log output
	   | --timestamps     Include timestamps on each line in the log output. Defaults to false
	   | --previous       If true, print the logs for the previous instance of the container in a pod if it exists. Defaults to false
	-p | --prefix         Prefix each pod's logs entries with [pod name]
	-j | --json           Add syntax highlighting for JSON log entries. Only available if outputting to a TTY that supports color
	   | --list-themes    List all available JSON highlighting theme names and exit

Options:
	<search terms>...  One or more case-sensitive search terms for pod names. Pass "-" to read search terms from stdin. Default is to show logs for a pod if any term is a match 
	-l | --label       Filter pods by one or more labels, pass additional -l arguments to add labels. Filtering is performed prior to name matching
	-s | --since       Show logs only since this timestamp
	   | --since-time  Only return logs after a specific date (RFC3339). Defaults to all logs. Only one of since-time / since may be used.
	   | --tail        Lines of recent log file to display. Defaults to -1, showing all log lines.
	-n | --namespace   Namespace pods must be in. Default is the default namespace for the cluster
	-c | --container   Print the logs of this container
	   | --limit-bytes Maximum bytes of logs to return. Defaults to no limit.
	-k | --kubeconfig  Path to kube config file. Defaults to value of env var KUBECONFIG or ~/.kube/config if not present
	-C | --context     The name of the kubeconfig context to use
	-t | --theme       Theme to use for JSON syntax highlighting. Default is "nord". See "--list-themes"

Environment Variables:
	Example: KLOGS_NAMESPACE=foo KLOGS_CONTEXT=bar KLOGS_PREFIX=1 KLOGS_JSON=1 KLOGS_THEME=monokai klogs search-terms...
	The following options/flags can be overridden via environment variables. Set value to "1" to enable a flag.
	context:    KLOGS_CONTEXT
	namespace:  KLOGS_NAMESPACE
	prefix:     KLOGS_PREFIX
	json:       KLOGS_JSON
	theme:      KLOGS_THEME

```