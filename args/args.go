package args

import (
	"os"
	"reflect"
	"strings"

	"github.com/ryantate13/hash-set"

	"github.com/ryantate13/klogs/fn"
)

func defaults() *Args {
	return &Args{
		All:        os.Getenv("KLOGS_ALL") == "1",
		KubeConfig: os.Getenv("KUBECONFIG"),
		Context:    os.Getenv("KLOGS_CONTEXT"),
		Namespace:  os.Getenv("KLOGS_NAMESPACE"),
		Prefix:     os.Getenv("KLOGS_PREFIX") == "1",
		JSON:       os.Getenv("KLOGS_JSON") == "1",
		Theme:      fn.Coalesce(os.Getenv("KLOGS_THEME"), "nord"),
	}
}

func optFlags[T any](argStruct *T) (*hash_set.Set[string], *hash_set.Set[string]) {
	opts := hash_set.New[string]()
	flags := hash_set.New[string]()
	a := reflect.ValueOf(argStruct).Elem()
	for i := 0; i < a.NumField(); i++ {
		f := a.Type().Field(i)
		if positional, hasPositional := f.Tag.Lookup("positional"); hasPositional && positional == "true" {
			continue
		}
		name := strings.ToLower(f.Name)
		kind := opts
		if f.Type.Name() == "bool" {
			kind = flags
		}
		short, hasShort := f.Tag.Lookup("short")
		if hasShort {
			if short != "" {
				kind.Add("-" + short)
			}
		} else {
			kind.Add("-" + name[:1])
		}
		long, hasLong := f.Tag.Lookup("long")
		if hasLong {
			if long != "" {
				kind.Add("--" + long)
			}
		} else {
			kind.Add("--" + name)
		}
	}
	return opts, flags
}

// Args encapsulates all the various flags/options for klogs
type Args struct {
	Help          bool
	Version       bool
	Query         []string `positional:"true" description:""`
	All           bool
	AllNamespaces bool `short:"" long:"all-namespaces"`
	AllContainers bool `short:"" long:"all-containers"`
	Label         []string
	LimitBytes    string `short:"" long:"limit-bytes"`
	Since         string
	SinceTime     string `short:"" long:"since-time"`
	Tail          string `short:"" long:"tail"`
	Follow        bool
	Timestamps    bool `short:"" long:"timestamps"`
	Previous      bool `short:""`
	KubeConfig    string
	Context       string `short:"C"`
	Container     string
	Namespace     string
	Prefix        bool
	JSON          bool
	Theme         string
	ListThemes    bool `short:"" long:"list-themes"`
}

// Usage returns the documentation string for the command
// TODO - support missing args from https://kubernetes.io/docs/reference/generated/kubectl/kubectl-commands#logs
func (a *Args) Usage() string {
	return `klogs - Displays logs for kubernetes pods matching either a pod name query, a set of labels, or both

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
	theme:      KLOGS_THEME`
}

// Parse takes an array of string args and returns the parsed Args struct
func Parse(argv []string) *Args {
	a := &Args{}
	opts, flags := optFlags(a)
	for i, arg := range argv {
		switch {
		case arg == "-h" || arg == "--help":
			a.Help = true
		case arg == "-v" || arg == "--version":
			a.Version = true
		case arg == "-k" || arg == "--kubeconfig":
			a.KubeConfig = argv[i+1]
		case arg == "--context":
			a.Context = argv[i+1]
		case arg == "-c" || arg == "--container":
			a.Container = argv[i+1]
		case arg == "-n" || arg == "--namespace":
			a.Namespace = argv[i+1]
		case arg == "-l" || arg == "--label":
			a.Label = append(a.Label, argv[i+1])
		case arg == "--limit-bytes":
			a.LimitBytes = argv[i+1]
		case arg == "-s" || arg == "--since":
			a.Since = argv[i+1]
		case arg == "--since-time":
			a.SinceTime = argv[i+1]
		case arg == "--tail":
			a.Tail = argv[i+1]
		case arg == "--timestamps":
			a.Timestamps = true
		case arg == "-a" || arg == "--all":
			a.All = true
		case arg == "--all-namespaces":
			a.AllNamespaces = true
		case arg == "--all-containers":
			a.AllContainers = true
		case arg == "-f" || arg == "--follow":
			a.Follow = true
		case arg == "--previous":
			a.Previous = true
		case arg == "-p" || arg == "--prefix":
			a.Prefix = true
		case arg == "-j" || arg == "--json":
			a.JSON = true
		case arg == "-t" || arg == "--theme":
			a.Theme = argv[i+1]
		case arg == "--list-themes":
			a.ListThemes = true
		}
	}
	for i := len(argv) - 1; i >= 1; i-- {
		if flags.Has(argv[i]) || opts.Has(argv[i-1]) {
			break
		} else {
			a.Query = append(a.Query, argv[i])
		}
	}
	d := defaults()
	if a.All == false {
		a.All = d.All
	}
	if a.KubeConfig == "" {
		a.KubeConfig = d.KubeConfig
	}
	if a.Context == "" {
		a.Context = d.Context
	}
	if a.Namespace == "" {
		a.Namespace = d.Namespace
	}
	if a.Prefix == false {
		a.Prefix = d.Prefix
	}
	if a.JSON == false {
		a.JSON = d.JSON
	}
	if a.Theme == "" {
		a.Theme = d.Theme
	}
	return a
}
