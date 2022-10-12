package logs

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/alecthomas/chroma/quick"
	"github.com/fatih/color"

	"github.com/ryantate13/klogs/args"
	"github.com/ryantate13/klogs/exec"
	"github.com/ryantate13/klogs/fn"
)

type colorFunc func(format string, a ...interface{}) string

var (
	colors = []colorFunc{
		color.CyanString,
		color.YellowString,
		color.GreenString,
		color.MagentaString,
		color.BlueString,
		color.HiCyanString,
		color.HiYellowString,
		color.HiGreenString,
		color.HiMagentaString,
		color.HiBlueString,
	}
	noColor colorFunc = fmt.Sprintf
)

func mkError(err map[string]interface{}) error {
	j, _ := json.MarshalIndent(err, "", "  ")
	return errors.New(string(j))
}

type pod struct {
	Name, Namespace string
}

func Read(ctx context.Context, opts *args.Args, ex exec.Executor, tty string) (<-chan string, <-chan error, error) {
	logChan := make(chan string)
	kubectl := []string{"kubectl"}
	for k, v := range map[string]string{"--kubeconfig": opts.KubeConfig, "--context": opts.Context} {
		if v != "" {
			kubectl = append(kubectl, k, v)
		}
	}
	getPods := append(kubectl, "get", "pods", "-o", "custom-columns=:metadata.name,:metadata.namespace")
	if opts.AllNamespaces {
		getPods = append(getPods, "--all-namespaces")
	} else if opts.Namespace != "" {
		getPods = append(getPods, "--namespace", opts.Namespace)
	}
	for _, l := range opts.Label {
		getPods = append(getPods, "-l", l)
	}
	nsPods, err := ex.Sync(ctx, getPods...)
	if err != nil {
		return nil, nil, mkError(map[string]interface{}{
			"code":    "get_pods_error",
			"command": getPods,
			"error":   err.Error(),
		})
	}
	pods := fn.Filter(
		fn.Map(nsPods, func(s string) *pod {
			f := strings.Fields(s)
			return &pod{f[0], f[1]}
		}), func(p *pod) bool {
			// filter by label only
			if len(opts.Query) == 0 {
				return true
			}
			// all search terms must match
			if opts.All {
				return fn.Reduce(opts.Query, func(a bool, c string) bool {
					return a && strings.Index(p.Name, c) != -1
				}, true)
			}
			// default behavior - one or more search terms must match
			return fn.Reduce(opts.Query, func(a bool, c string) bool {
				return a || strings.Index(p.Name, c) != -1
			}, false)
		})
	if len(pods) == 0 {
		return nil, nil, mkError(map[string]interface{}{
			"code":  "no_pods_found",
			"error": "no available pods match query terms",
			"opts":  opts,
		})
	}
	wg := &sync.WaitGroup{}
	wg.Add(len(pods))
	errChan := make(chan error)
	logCmd := append(kubectl, "logs")
	for k, v := range map[string]bool{
		"--all-containers": opts.AllContainers,
		"--follow":         opts.Follow,
		"--prefix":         opts.Prefix,
		"--previous":       opts.Previous,
		"--timestamps":     opts.Timestamps,
	} {
		if v {
			logCmd = append(logCmd, k)
		}
	}
	for k, v := range map[string]string{
		"--container":   opts.Container,
		"--limit-bytes": opts.LimitBytes,
		"--since":       opts.Since,
		"--since-time":  opts.SinceTime,
		"--tail":        opts.Tail,
	} {
		if v != "" {
			kubectl = append(kubectl, k, v)
		}
	}
	for i, p := range pods {
		podLogCmd := append(logCmd, "-n", p.Namespace, p.Name)
		colorize := noColor
		if tty != "" {
			colorize = colors[i%len(colors)]
		}
		c, err := ex.Stream(ctx, errChan, podLogCmd...)
		if err != nil {
			return nil, nil, mkError(map[string]interface{}{
				"code":    "logs_error",
				"command": podLogCmd,
				"error":   err.Error(),
				"pod":     p,
			})
		}
		prettyJSON := opts.JSON && tty != ""
		numParts := 1
		for _, b := range []bool{opts.Prefix, opts.Timestamps} {
			if b {
				numParts++
			}
		}
		go func(ch <-chan string) {
			for {
				select {
				case line, ok := <-ch:
					if !ok {
						wg.Done()
						return
					}
					parts := strings.SplitN(line, " ", numParts)
					prefix := ""
					if opts.Prefix {
						prefix = colorize(parts[0]) + " "
					}
					timestamp := ""
					if opts.Timestamps {
						i := 0
						if !opts.Prefix {
							i = 1
						}
						timestamp = parts[i] + " "
					}
					logEntry := parts[len(parts)-1]
					if prettyJSON {
						b := bytes.NewBuffer(nil)
						err := quick.Highlight(b, logEntry, "json", tty, opts.Theme)
						if err == nil {
							logEntry = b.String()
						}
					}
					logChan <- prefix + timestamp + logEntry
				case <-ctx.Done():
					wg.Done()
					return
				}
			}
		}(c)
	}
	go func() {
		wg.Wait()
		close(logChan)
	}()
	return logChan, errChan, nil
}
