package fn

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFilter(t *testing.T) {
	t.Run("filters a collection of elements given a filter func", func(t *testing.T) {
		require.Equal(t, []string{"foo", "bar", "baz"}, Filter([]string{
			"foo",
			"bar",
			"baz",
			"1234",
			"abcd",
		}, func(s string) bool {
			return len(s) == 3
		}))
	})
}

func TestMap(t *testing.T) {
	t.Run("transforms a collection of elements given a mapping func", func(t *testing.T) {
		require.Equal(t, []int{1, 4, 9, 16, 25}, Map[int]([]int{1, 2, 3, 4, 5}, func(i int) int {
			return i * i
		}))
	})
}

func TestReduce(t *testing.T) {
	t.Run("reduces a collection to a single value given a reducer func", func(t *testing.T) {
		require.Equal(t, 10, Reduce([]int{1, 2, 3, 4}, func(a, b int) int {
			return a + b
		}, 0))
		require.Equal(t, true, Reduce([]bool{false, false, false, true}, func(a, b bool) bool {
			return a || b
		}, false))
		require.Equal(t, false, Reduce([]bool{true, true, true, false}, func(a, b bool) bool {
			return a && b
		}, true))
		require.Equal(
			t,
			map[string]int{
				"foo": 2,
				"bar": 1,
			},
			Reduce(
				[]string{"foo", "foo", "bar"},
				func(a map[string]int, c string) map[string]int {
					a[c] = a[c] + 1
					return a
				},
				map[string]int{},
			),
		)
	})
}

func TestCoalesce(t *testing.T) {
	type test struct {
		foo string
	}
	t.Run("returns the first value in a collection that is not the zero value for its type", func(t *testing.T) {
		require.Equal(t, 123, Coalesce(0, 0, 0, 0, 0, 0, 123, 0, 0, 0))
		require.Equal(t, "test", Coalesce("", "", "", "", "", "", "test", "", ""))
		require.Equal(t, &test{}, Coalesce([]*test{nil, nil, nil, {}, nil, nil}...))
		require.Equal(t, test{foo: "bar"}, Coalesce([]test{{}, {}, {}, {}, {foo: "bar"}, {}}...))
	})
}
