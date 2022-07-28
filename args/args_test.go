package args

import (
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/ryantate13/hash-set"
	"github.com/stretchr/testify/require"

	"github.com/ryantate13/klogs/fn"
)

func Test_optFlags(t *testing.T) {
	a := &Args{}
	opts, flags := optFlags(a)
	require.Equal(t, hash_set.Of(
		"-h", "--help",
		"-v", "--version",
		"-a", "--all",
		"--all-namespaces",
		"--all-containers",
		"--list-themes",
		"-f", "--follow",
		"--timestamps",
		"--previous",
		"-p", "--prefix",
		"-j", "--json",
	), flags)
	require.Equal(t, hash_set.Of(
		"-l", "--label",
		"-s", "--since",
		"--since-time",
		"--tail",
		"-n", "--namespace",
		"-c", "--container",
		"--limit-bytes",
		"-k", "--kubeconfig",
		"-C", "--context",
		"-t", "--theme",
	), opts)
}

func TestUsage(t *testing.T) {
	a := &Args{}
	opts, flags := optFlags(a)
	usage := a.Usage()
	for _, kind := range []*hash_set.Set[string]{opts, flags} {
		kind.Foreach(func(arg string) {
			require.Contains(t, usage, arg)
		})
	}
	isArg := regexp.MustCompile("^-{1,2}[a-zA-Z-]+")
	helpWords := hash_set.Of(strings.Fields(usage)...).Filter(func(w string) bool {
		return isArg.MatchString(w)
	})
	require.Equal(t, opts.Len()+flags.Len(), helpWords.Len())
	isFullyDocumented := fn.Reduce(helpWords.Slice(), func(a bool, c string) bool {
		return a && (opts.Has(c) || flags.Has(c))
	}, true)
	require.True(t, isFullyDocumented)
}

func TestParse(t *testing.T) {
	for _, k := range []string{
		"KLOGS_ALL",
		"KUBECONFIG",
		"KLOGS_CONTEXT",
		"KLOGS_NAMESPACE",
		"KLOGS_PREFIX",
		"KLOGS_JSON",
		"KLOGS_THEME",
	} {
		require.NoError(t, os.Unsetenv(k))
	}
	tests := []struct {
		it   string
		args []string
		want *Args
		env  map[string]string
	}{
		{
			it:   "handles setting defaults",
			args: []string{"klogs"},
			want: defaults(),
		},
		{
			it: "parses short args",
			args: []string{
				"klogs",
				"-h",
				"-v",
				"-a",
				"-f",
				"-p",
				"-P",
				"-j",
				"-l", "test",
				"-s", "test",
				"-n", "test",
				"-c", "test",
				"-k", "test",
				"-C", "test",
				"-t", "test",
				"test",
			},
			want: &Args{
				Help:       true,
				Version:    true,
				Query:      []string{"test"},
				All:        true,
				Label:      []string{"test"},
				Since:      "test",
				Follow:     true,
				KubeConfig: "test",
				Container:  "test",
				Namespace:  "test",
				Prefix:     true,
				JSON:       true,
				Theme:      "test",
			},
		},
		{
			it: "parses long args",
			args: []string{
				"klogs",
				"--help",
				"--version",
				"--all",
				"--all-namespaces",
				"--all-containers",
				"--follow",
				"--timestamps",
				"--previous",
				"--prefix",
				"--json",
				"--label", "test",
				"--since", "test",
				"--since-time", "test",
				"--tail", "test",
				"--namespace", "test",
				"--container", "test",
				"--limit-bytes", "test",
				"--kubeconfig", "test",
				"--context", "test",
				"--theme", "test",
				"test",
			},
			want: &Args{
				Help:          true,
				Version:       true,
				Query:         []string{"test"},
				All:           true,
				AllNamespaces: true,
				AllContainers: true,
				Label:         []string{"test"},
				LimitBytes:    "test",
				Since:         "test",
				SinceTime:     "test",
				Tail:          "test",
				Follow:        true,
				Timestamps:    true,
				Previous:      true,
				KubeConfig:    "test",
				Context:       "test",
				Container:     "test",
				Namespace:     "test",
				Prefix:        true,
				JSON:          true,
				Theme:         "test",
			},
		},
		{
			it:   "reads defaults from the environment",
			args: []string{"klogs"},
			want: &Args{
				All:        true,
				KubeConfig: "test",
				Context:    "test",
				Namespace:  "test",
				Prefix:     true,
				JSON:       true,
				Theme:      "test",
			},
			env: map[string]string{
				"KLOGS_ALL":       "1",
				"KUBECONFIG":      "test",
				"KLOGS_CONTEXT":   "test",
				"KLOGS_NAMESPACE": "test",
				"KLOGS_PREFIX":    "1",
				"KLOGS_JSON":      "1",
				"KLOGS_THEME":     "test",
			},
		},
	}
	for _, tt := range tests {
		for k, v := range tt.env {
			require.NoError(t, os.Setenv(k, v))
		}
		t.Run(tt.it, func(t *testing.T) {
			require.Equal(t, tt.want, Parse(tt.args))
		})
	}
}
