package exec

import (
	"context"
	"math/rand"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSync(t *testing.T) {
	tests := []struct {
		it     string
		args   []string
		assert func(t *testing.T, out []string, err error)
	}{
		{
			it:   "returns an error if the command cannot be executed",
			args: []string{"sh", "-c", "exit 1"},
			assert: func(t *testing.T, out []string, err error) {
				require.Error(t, err)
				require.Nil(t, out)
			},
		},
		{
			it:   "returns command output as an array of strings",
			args: []string{"seq", "1", "3"},
			assert: func(t *testing.T, out []string, err error) {
				require.NoError(t, err)
				require.Equal(t, []string{"1", "2", "3"}, out)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.it, func(t *testing.T) {
			got, err := DefaultExecutor.Sync(context.Background(), tt.args...)
			tt.assert(t, got, err)
		})
	}
}

func TestStream(t *testing.T) {
	tests := []struct {
		it     string
		args   []string
		assert func(t *testing.T, ch <-chan string, err error, errChan chan error)
	}{
		{
			it:   "returns an error if the command cannot be executed",
			args: []string{"INVALID_" + strconv.Itoa(rand.Int())},
			assert: func(t *testing.T, _ <-chan string, err error, _ chan error) {
				require.Error(t, err)
			},
		},
		{
			it:   "sends errors to the error channel if the command does not succeed",
			args: []string{"sh", "-c", "exit 1"},
			assert: func(t *testing.T, ch <-chan string, err error, errChan chan error) {
				require.NoError(t, err)
				require.NoError(t, <-errChan, "no error should be encountered reading output")
				require.Error(t, <-errChan, "errors from non-zero exit codes should propagate")
			},
		},
		{
			it:   "send all lines of command output to a channel",
			args: []string{"seq", "1", "3"},
			assert: func(t *testing.T, ch <-chan string, err error, errChan chan error) {
				require.NoError(t, err)
				for i := 1; i <= 3; i++ {
					require.Equal(t, strconv.Itoa(i), <-ch)
				}
				require.NoError(t, <-errChan)
			},
		},
		{
			it:   "supports reading large lines of output",
			args: []string{"sh", "-c", "head -c 1000000 < /dev/zero"},
			assert: func(t *testing.T, ch <-chan string, err error, errChan chan error) {
				require.NoError(t, err)
				longLine := <-ch
				require.Equal(t, 1000000, len(longLine))
				for _, c := range longLine {
					require.Equal(t, 0, int(c))
				}
				require.NoError(t, <-errChan)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.it, func(t *testing.T) {
			errChan := make(chan error)
			ch, err := DefaultExecutor.Stream(context.Background(), errChan, tt.args...)
			tt.assert(t, ch, err, errChan)
		})
	}
}
