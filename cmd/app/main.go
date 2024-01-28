package main

import (
	"fmt"
	"github.com/m1ker1n/go-developer-test/internal/config"
)

func main() {
	cfg := config.MustLoad()
	//TODO: remove printf
	fmt.Printf("%+v\n", cfg)
}
