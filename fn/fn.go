package fn

// Reduce reduces a collection to a single value given a reducer func
func Reduce[A, B any](elems []A, fn func(B, A) B, init B) B {
	for _, el := range elems {
		init = fn(init, el)
	}
	return init
}

// Filter filters a collection of elements given a filter func
func Filter[T any](elems []T, f func(T) bool) []T {
	return Reduce[T, []T](elems, func(a []T, c T) []T {
		if f(c) {
			a = append(a, c)
		}
		return a
	}, make([]T, 0))
}

// Map transforms a collection of elements given a mapping func
func Map[A, B any](as []A, fn func(A) B) []B {
	bs := make([]B, len(as))
	for i, a := range as {
		bs[i] = fn(a)
	}
	return bs
}

// Coalesce returns the first value in a collection that is not the zero value for its type
func Coalesce[T comparable](elems ...T) T {
	var zero T
	for _, el := range elems {
		if el != zero {
			return el
		}
	}
	return zero
}
