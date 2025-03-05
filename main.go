package main

import (
	"fmt"
	"os"

	"tiktok-live-logger/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
