package utils

import "iter"

func NextIter[T any](next func() (T, error)) iter.Seq[T] {
	return func(yield func(T) bool) {
		for v, err := next(); err == nil; v, err = next() {
			if !yield(v) {
				return
			}
		}
	}
}

func Apply[A any, R any](s iter.Seq[A], f func(A) R) iter.Seq[R] {
	return func(yield func(R) bool) {
		for t := range s {
			if !yield(f(t)) {
				return
			}
		}
	}
}

func FuncsToIter[R any](funcs ...func() R) iter.Seq[R] {
	return func(yield func(R) bool) {
		for _, f := range funcs {
			if !yield(f()) {
				return
			}
		}
	}
}
