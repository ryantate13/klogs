package exec

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"io"
	"os/exec"
	"strings"
)

// Executor executes OS commands and returns the results, either synchronously, or asynchronously on a channel
type Executor interface {
	// Sync executes a command synchronously and returns the output as an array of strings
	Sync(ctx context.Context, cmdAndArgs ...string) ([]string, error)
	// Stream executes a command asynchronously and sends lines of output to a channel
	Stream(ctx context.Context, errChan chan<- error, cmdAndArgs ...string) (<-chan string, error)
}

type executor struct{}

// DefaultExecutor can be used to execute OS commands or overridden to change the default behavior
var DefaultExecutor Executor = &executor{}

func (e *executor) Sync(ctx context.Context, cmdAndArgs ...string) ([]string, error) {
	c := exec.CommandContext(ctx, cmdAndArgs[0], cmdAndArgs[1:]...)
	out, err := c.Output()
	if err != nil {
		return nil, err
	}
	return strings.Split(strings.TrimSpace(string(out)), "\n"), nil
}

func (e *executor) Stream(ctx context.Context, errChan chan<- error, cmdAndArgs ...string) (<-chan string, error) {
	c := exec.CommandContext(ctx, cmdAndArgs[0], cmdAndArgs[1:]...)
	ch := make(chan string)
	out, err := c.StdoutPipe()
	if err != nil {
		return nil, err
	}
	errBuf := bytes.NewBuffer(nil)
	c.Stderr = errBuf
	if err = c.Start(); err != nil {
		return nil, err
	}
	rd := bufio.NewReader(out)
	go func() {
		var (
			chunk    []byte
			isPrefix bool
			err      error
			buf      strings.Builder
		)
		for err == nil {
			chunk, isPrefix, err = rd.ReadLine()
			if err != nil {
				break
			}
			buf.Write(chunk)
			if !isPrefix {
				ch <- buf.String()
				buf = strings.Builder{}
			}
		}
		if err == io.EOF {
			errChan <- nil
		} else {
			errChan <- err
		}
		if err = c.Wait(); err != nil {
			errChan <- errors.New(errBuf.String())
		}
		close(ch)
	}()
	return ch, nil
}
