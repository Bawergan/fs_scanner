package app

import (
	"bufio"
	"errors"
	"fmt"
	. "fs_scan/db"
	lt "fs_scan/tools"
	"os"
	"strings"
)

const file_db_name = "files"
const store_path = "./store/"

type App struct {
	db *Database
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
		case "create file db":
			fmt.Println("Creating file db...")

			err := a.createNewFileDb()
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
            if err != nil{
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

func (a *App) createNewFileDb() error {
	if a.db != nil {
		return errors.New("db already open!")
	}
	db, err := OpenDatabase(store_path + file_db_name)
	if err != nil {
		return err
	}

	sqlStmt := `
    CREATE TABLE IF NOT EXISTS files (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL,
        created_at DATETIME NOT NULL,
        tags TEXT NOT NULL
    );
    `
	err = db.CreateTable(sqlStmt)
	if err != nil {
		return err
	}
	a.db = db
	return nil
}
func (a *App) openFileDb() error {
	if a.db != nil {
		return errors.New("db already open!")
	}
	db, err := OpenDatabase(store_path + file_db_name)
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

	rows, err := a.db.Query("SELECT count(*) FROM files")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next(){
        count := 0
		err = rows.Scan(&count)
		if err != nil {
			return err
		}
		fmt.Printf("Count: %v\n", count)
	}

	err = rows.Err()
	if err != nil {
		return err
	}
    return nil
}
