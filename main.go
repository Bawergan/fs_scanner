package main

import (
	app "fs_scan/app"
)

func main() {
	app := app.NewApp()
	app.LaunchCLI()
}
