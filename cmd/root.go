package cmd

import (
	"flag"
	"fmt"
	app "fs_scan/cli"
	"os"
)

func Execute() {
	rescanFlag := flag.Bool("rescan", false, "Автоматически открыть базу данных и перезагрузить данные")
	serveFlag := flag.Bool("serve", false, "Запустить сервер")
	flag.Parse()

	if *rescanFlag {
		performRescan()
	} else if *serveFlag {
		startServer()
	} else {
		launchCLI()
	}
}

func performRescan() {
	a := app.NewApp()
	err := a.OpenFileDb()
	if err != nil {
		fmt.Println("Ошибка при открытии базы данных:", err)
		os.Exit(1)
	}
	defer a.ExitApp()

	fmt.Println("Начало пересканирования файловой системы...")
	a.ReloadData()
	fmt.Println("Пересканирование завершено.")
}

func startServer() {
	a := app.NewApp()
	err := a.OpenFileDb()
	if err != nil {
		fmt.Println("Ошибка при открытии базы данных:", err)
		os.Exit(1)
	}
	defer a.ExitApp()

	fmt.Println("Запуск сервера...")
	err = a.Serve()
	if err != nil {
		fmt.Println("Ошибка при запуске сервера:", err)
		os.Exit(1)
	}
}

func launchCLI() {
	a := app.NewApp()
	a.LaunchCLI()
}
