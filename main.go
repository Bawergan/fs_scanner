package main

import (
	cli "fs_scan/cli"
)

func main() {
	app := cli.NewApp()
	app.LaunchCLI()
}
