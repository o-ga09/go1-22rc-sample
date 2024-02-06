package main

import "github.com/o-ga09/go122rcsample/api/internal/presenter"

func main() {
	s := presenter.NewServer("8080")
	s.Run()
}
