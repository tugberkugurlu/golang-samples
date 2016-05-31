package main

import (
	"fmt"
	"os"

	"github.com/Masterminds/semver"
)

func main() {
	v, err := semver.NewVersion("1.2.3-beta.1+build345")

	for i := 0; i < len(os.Args); i++ {
		fmt.Printf("%s\n", os.Args[i])
	}

	if err != nil {
		_ = fmt.Errorf("Error parsing version: %s", err)
	}

	fmt.Printf(v.String())
}
