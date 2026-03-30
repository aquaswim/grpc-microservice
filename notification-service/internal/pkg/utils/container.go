package utils

import "github.com/golobby/container/v3"

func Resolve[T any](c container.Container) T {
	var t T
	container.MustResolve(c, &t)
	return t
}

func ResolveNamed[T any](c container.Container, name string) T {
	var t T
	container.MustNamedResolve(c, &t, name)
	return t
}
