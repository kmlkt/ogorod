package utils

import (
	"iter"
)

func RunAll(funcs ...func() error) error {
	return RunAllIter(FuncsToIter(funcs...))
}

func RunAllIter(funcs iter.Seq[error]) error {
	for err := range funcs {
		if err != nil {
			return err
		}
	}
	return nil
}

func Then[T any](f func() (T, error), g func(T) error) func() error {
	return func() error {
		t, err := f()
		if err != nil {
			return err
		}
		return g(t)
	}
}

func ThenA[A any, T any](f func(A) (T, error), g func(T) error) func(A) error {
	return func(a A) error {
		t, err := f(a)
		if err != nil {
			return err
		}
		return g(t)
	}
}

func ThenAA[A any, B any, T any](
	f func(A) (T, error),
	g func(T, B) error,
) func(A, B) error {
	return func(a A, b B) error {
		t, err := f(a)
		if err != nil {
			return err
		}
		return g(t, b)
	}
}

func SafeAR[A any, R any](f func(A) R) func(A) (R, error) {
	return func(a A) (R, error) {
		return f(a), nil
	}
}
