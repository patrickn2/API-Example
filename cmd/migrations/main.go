package main

import (
	"flag"
	"fmt"

	"github.com/patrickn2/api-challenge/infra/migrations"
)

func main() {
	upFlag := flag.Bool("up", false, "Run migrations UP")
	downFlag := flag.Bool("down", false, "Run migrations DOWN")
	createFlag := flag.String("create", "", "Create a new migration\nUse -create migrationName")
	flag.Parse()

	if *upFlag {
		fmt.Println("Migrations up")
		migrations.Up()
		return
	}
	if *downFlag {
		fmt.Println("Migrations up")
		migrations.Down()
		return
	}

	if createFlag != nil {
		if *createFlag == "" {
			fmt.Println("Migration name can't be empty")
			return
		}
		fmt.Println("Create Migration", *createFlag)
		migrations.Create(*createFlag)
		return
	}
}
