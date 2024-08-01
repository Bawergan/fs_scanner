package app

import (
	"bufio"
	"fmt"
	"os"
	"strings"
    . "fs_scan/db"
)

type App struct{
    db *Database
}

func NewApp() *App{
    return &App{db: nil}
}

func (a *App) LaunchCLI(){
    reader := bufio.NewReader(os.Stdin)
    for {
        fmt.Print("Enter command: ")
        command, _ := reader.ReadString('\n')
        command = strings.TrimSpace(command)

        switch command {
        case "open db":
        case "reload data":
        case "update data":
        case "help":
            fmt.Println("Available commands:")
            fmt.Println("- help")
            fmt.Println("- exit")
        case "exit":
            fmt.Println("Exiting CLI...")
            return
        default:
            fmt.Println("Unknown command. Type 'help' to see available commands.")
        }
    }
}

func (a *App) openDb() {

}

func(a *App) reloadData(){

}

func (a *App) updateData(){

}
