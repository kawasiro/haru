package main

import "os"

func main() {
	if len(os.Getenv("ID")) > 0 {
		mainCli()
	} else {
		mainSvr()
	}
}
