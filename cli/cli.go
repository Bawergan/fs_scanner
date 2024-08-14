package app

import (
	"bufio"
	"errors"
	"fmt"
	data "fs_scan/data"
	lt "fs_scan/tools"
	"net/http"
	"os"
	"strings"
)

type App struct {
	db *data.FileDb
}

func NewApp() *App {
	return &App{db: nil}
}

func (a *App) LaunchCLI() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Enter command: ")
		command, _ := reader.ReadString('\n')
		command = strings.TrimSpace(command)

		switch command {
		case "serve":
			fmt.Println("Statring server...")
			err := a.serve()
			if err != nil {
				fmt.Println(err)
			}
		case "open file db":
			fmt.Println("Opening file db...")

			err := a.openFileDb()
			if err != nil {
				fmt.Println(err)
			}
		case "reload data":
			fmt.Println("Reloading data...")
			a.reloadData()
		case "update data":
			fmt.Println("Updating data...")
		case "count":
			err := a.countEnteries()
			if err != nil {
				fmt.Println(err)
			}
		case "help":
			fmt.Println("Available commands:")
			fmt.Println("- help")
			fmt.Println("- exit")
		case "exit":
			fmt.Println("Exiting CLI...")

			err := a.exitApp()
			if err != nil {
				fmt.Println(err)
			}
			if err == nil {
				return
			}
		default:
			fmt.Println("Unknown command. Type 'help' to see available commands.")
		}
	}
}

func (a *App) openFileDb() error {
	if a.db != nil {
		return errors.New("db already open")
	}
	db, err := data.CreateFileDb()
	if err != nil {
		return err
	}
	a.db = db
	return nil
}

func (a *App) exitApp() error {
	if a.db != nil {
		err := a.db.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *App) reloadData() {
	lt.ReloadDb(a.db)
}

func (a *App) updateData() {

}

func (a *App) countEnteries() error {
	count, err := a.db.CountEnteries()
	if err != nil {
		return err
	}
	fmt.Printf("Count: %v\n", count)
	return nil
}

func (a *App) reloadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// Call your function here
		lt.ReloadDb(a.db)

		// Respond to the client
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Update performed successfully")
	} else {
		// If the method is not GET, return a 405 Method Not Allowed
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintln(w, "Method not allowed")
	}
}

func (a *App) serve() error {
	http.HandleFunc("/api/reload", a.reloadHandler)
	return http.ListenAndServe(":5000", nil)
}
