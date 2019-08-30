package datastore

import (
	"log"
	"os"
	"strings"
)

func init() {
	dbname := os.Getenv("PGDATABASE")
	if dbname == "" {
		dbname = "bogthesrctest"
	}

	if !strings.HasSuffix(dbname, "test") {
		dbname += "test"
	}

	if err := os.Setenv("PGDATABASE", dbname); err != nil {
		log.Fatalf("error setting env %s", err)
	}

	Connect()
	Drop()
	Create()
}
