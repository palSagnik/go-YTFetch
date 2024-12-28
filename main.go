package main

import (
	"fmt"
	"time"

	"github.com/palSagnik/go-YTFetch.git/database"
)

func main() {
	
	// loop till database is initialised
	for {
		if err := database.ConnectDB(); err != nil {
			fmt.Println(err)
			fmt.Println("waiting for 30 seconds before trying again")
			time.Sleep(time.Second * 30)
			continue
		}
		break
	}

	err := database.MigrateUp()
	if err != nil {
		panic(err)
	}
}