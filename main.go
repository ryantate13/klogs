package main

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/alecthomas/chroma/quick"
	"github.com/alecthomas/chroma/styles"
	"github.com/fatih/color"
	term "github.com/jwalton/go-supportscolor"
	"github.com/mattn/go-isatty"

	"klogs/args"
	"klogs/exec"
	"klogs/fn"
	"klogs/logs"
)

var (
	//go:embed VERSION
	v         string
	version   = strings.TrimSpace(v)
	ttyFormat string
)

func fatal(s string) {
	fmt.Fprintln(os.Stderr, color.RedString(s))
	os.Exit(1)
}

func init() {
	isTTY := isatty.IsTerminal(os.Stdout.Fd())
	color.NoColor = !isTTY
	if isTTY {
		out := term.Stdout()
		if out.Has16m {
			ttyFormat = "terminal16m"
		} else if out.Has256 {
			ttyFormat = "terminal256"
		} else if out.SupportsColor {
			ttyFormat = "terminal"
		}
	}
}

func main() {
	opts := args.Parse(os.Args)
	if opts.Help {
		fmt.Println(opts.Usage())
		os.Exit(0)
	}
	if opts.Version {
		fmt.Println("klogs " + version)
		os.Exit(0)
	}
	if opts.ListThemes {
		if ttyFormat == "" {
			for _, theme := range styles.Names() {
				fmt.Println(theme)
			}
		} else {
			example := `{"string":"test","number":123,"array":["1",2],"obj":{"foo":"bar"},"null":null,"bool":true}`
			colALen := fn.Reduce(styles.Names(), func(a int, c string) int {
				if len(c) > a {
					return len(c)
				}
				return a
			}, 0)
			colBLen := len(example)
			const (
				header = 'h'
				row    = 'r'
				footer = 'f'
			)
			printSep := func(t rune) {
				var l, c, r string
				switch t {
				case header:
					l, c, r = "┌", "┬", "┐"
				case row:
					l, c, r = "├", "┼", "┤"
				case footer:
					l, c, r = "└", "┴", "┘"
				}
				fmt.Println(l + strings.Repeat("─", colALen+2) + c + strings.Repeat("─", colBLen+2) + r)
			}
			printRow := func(colA, colB string) {
				fmt.Printf("│ %-"+strconv.Itoa(colALen)+"s │ %-"+strconv.Itoa(colBLen)+"s │\n", colA, colB)
			}
			printSep(header)
			printRow("Name", "Example")
			printSep(row)
			for _, theme := range styles.Names() {
				b := bytes.NewBuffer(nil)
				row := example
				if err := quick.Highlight(b, example, "json", ttyFormat, theme); err == nil {
					row = b.String()
				}
				printRow(theme, row)
			}
			printSep(footer)
		}
		os.Exit(0)
	}
	if len(opts.Query) == 1 && opts.Query[0] == "-" {
		stdin, err := io.ReadAll(os.Stdin)
		if err == nil {
			opts.Query = strings.Fields(string(stdin))
		} else {
			opts.Query = []string{}
		}
	}
	if len(opts.Query) == 0 && len(opts.Label) == 0 {
		fatal("Error: either pod name query or pod labels must be supplied\n\n" + opts.Usage())
	}

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		shutdown := make(chan os.Signal)
		signal.Notify(shutdown, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
		<-shutdown
		cancel()
	}()

	logChan, errChan, err := logs.Read(ctx, opts, exec.DefaultExecutor, ttyFormat)
	if err != nil {
		fatal(err.Error())
	}
	for {
		select {
		case err = <-errChan:
			if err != nil {
				fatal(err.Error())
			}
		case log, ok := <-logChan:
			if !ok {
				return
			}
			fmt.Println(log)
		}
	}
}
