package main

import (
	"fmt"
	"os"

	"github.com/Masterminds/semver"
)

func main() {

	// Program Name is always the first (implicit) argument
	cmd := os.Args[0]
	fmt.Printf("Program Name: %s\n", cmd)

	// number of arguments being passed in.
	// os.Args[1:] simply says: “give me a new subslice starting with index 1 (not 0) to the end of the slice.”
	argCount := len(os.Args[1:])
	fmt.Printf("Total Arguments (excluding program name): %d\n", argCount)

	// loop through all the arguments being passed in
	for i := 1; i <= argCount; i++ {
		fmt.Printf("%s\n", os.Args[i])
	}

	v, err := semver.NewVersion("1.2.3-beta.1+build345")
	if err != nil {
		_ = fmt.Errorf("Error parsing version: %s", err)
	}

	fmt.Printf("%s\n", v.String())
}
