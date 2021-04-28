package main

import (
	"fmt"
	"os"

	"github.com/fernandomalmeida/gohttpspeedtest"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s [fast|ookla]\n", os.Args[0])
		os.Exit(1)
	}

	var provider *gohttpspeedtest.Provider
	var err error

	switch os.Args[1] {
	case "fast":
		provider, err = gohttpspeedtest.FastProvider()
		if err != nil {
			fmt.Printf("error on fast provider: %s", err)
			os.Exit(2)
		}
	case "ookla":
		provider, err = gohttpspeedtest.OoklaProvider()
		if err != nil {
			fmt.Printf("error on ookla provider: %s", err)
			os.Exit(3)
		}
	default:
		fmt.Printf("Usage: %s [fast|ookla]\n", os.Args[0])
		os.Exit(1)
	}

	downloadSpeed, uploadSpeed, err := gohttpspeedtest.MeasureDownloadAndUpload(provider)
	if err != nil {
		fmt.Printf("error on measure: %s", err)
	}

	fmt.Printf("Download speed: %8.2fMbps\nUpload speed: %8.2fMbps\n", downloadSpeed, uploadSpeed)
}
