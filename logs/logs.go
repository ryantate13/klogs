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

	"klogs/args"
	"klogs/exec"
	"klogs/fn"
)

type colorFunc = func(format string, a ...interface{}) string

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
	j, _ := json.Marshal(err)
	return errors.New(string(j))
}

func Read(ctx context.Context, opts *args.Args, ex exec.Executor, tty string) (<-chan string, <-chan error, error) {
	logChan := make(chan string)
	kubectl := []string{"kubectl"}
	if opts.KubeConfig != "" {
		kubectl = append(kubectl, "--kubeconfig", opts.KubeConfig)
	}
	if opts.Context != "" {
		kubectl = append(kubectl, "--context", opts.Context)
	}
	if opts.AllNamespaces {
		kubectl = append(kubectl, "--all-namespaces")
	} else if opts.Namespace != "" {
		kubectl = append(kubectl, "--namespace", opts.Namespace)
	}
	getPods := append(kubectl, "get", "pods", "-o", "custom-columns=:metadata.name")
	for _, l := range opts.Label {
		getPods = append(getPods, "-l", l)
	}
	pods, err := ex.Sync(ctx, getPods...)
	if err != nil {
		return nil, nil, mkError(map[string]interface{}{
			"code":    "get_pods_error",
			"command": getPods,
			"error":   err.Error(),
		})
	}
	pods = fn.Filter(pods, func(p string) bool {
		// filter by label only
		if len(opts.Query) == 0 {
			return true
		}
		// all search terms must match
		if opts.All {
			return fn.Reduce(opts.Query, func(a bool, c string) bool {
				return a && strings.Index(p, c) != -1
			}, true)
		}
		// default behavior - one or more search terms must match
		return fn.Reduce(opts.Query, func(a bool, c string) bool {
			return a || strings.Index(p, c) != -1
		}, false)
	})
	if len(pods) == 0 {
		return nil, nil, mkError(map[string]interface{}{
			"code":                 "no_pods_found",
			"error":                "no available pods match query terms",
			"labels":               opts.Label,
			"name_query":           opts.Query,
			"all_terms_must_match": opts.All,
		})
	}
	wg := &sync.WaitGroup{}
	wg.Add(len(pods))
	errChan := make(chan error)
	logCmd := append(kubectl, "logs")
	if opts.Follow {
		logCmd = append(logCmd, "-f")
	}
	if opts.AllContainers {
		logCmd = append(logCmd, "--all-containers")
	}
	if opts.Timestamps {
		logCmd = append(logCmd, "--timestamps")
	}
	if opts.Prefix {
		logCmd = append(logCmd, "--prefix")
	}
	if opts.Previous {
		logCmd = append(logCmd, "--previous")
	}
	if opts.Since != "" {
		logCmd = append(logCmd, "--since", opts.Since)
	}
	if opts.SinceTime != "" {
		logCmd = append(logCmd, "--since-time", opts.SinceTime)
	}
	if opts.Tail != "" {
		logCmd = append(logCmd, "--tail", opts.Tail)
	}
	if opts.Container != "" {
		logCmd = append(logCmd, "-c", opts.Container)
	}
	if opts.LimitBytes != "" {
		logCmd = append(logCmd, "--limit-bytes", opts.LimitBytes)
	}
	for i, pod := range pods {
		podLogCmd := append(logCmd, pod)
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
				"pod":     pod,
			})
		}
		prettyJSON := opts.JSON && tty != ""
		numParts := 1
		if opts.Prefix {
			numParts++
		}
		if opts.Timestamps {
			numParts++
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
