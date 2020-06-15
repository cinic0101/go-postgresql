package main

import (
	"fmt"
	"github.com/cinic0101/go-postgresql/pg"
	_ "github.com/joho/godotenv/autoload"
	"os"
)

var client *pg.Client

func main() {
	client = &pg.Client{
		MasterDbUrl: os.Getenv("MASTER_DB_URL"),
		SlaveDbUrl:  os.Getenv("SLAVE_DB_URL"),
	}

	if err := client.InitConnPools(); err != nil {
		panic(err)
	}

	insertNewTestAccount()
	listAllAccounts()
}

func insertNewTestAccount() {
	_, err := client.Exec("INSERT INTO account (name, email) VALUES ($1, $2)", "somebody", "somebody@example.com")
	if err != nil {
		panic(err)
	}
}

func listAllAccounts() {
	rows, err := client.Query("select id, name, email from account")
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var id int32
		var name string
		var email string

		err := rows.Scan(&id, &name, &email)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%d, %s, %s\n", id, name, email)
	}
}
