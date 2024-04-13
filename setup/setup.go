package setup

import (
	"fmt"
	"my-ether-tool/database"
	"os"
	"path"
)

func defaultDbPath() string {
	home := os.Getenv("HOME")
	dir := path.Join(home, ".met")
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		panic(fmt.Sprintf("create default db path: %s error: %v", dir, err))
	}

	return path.Join(dir, "met.db")
}

func SetupDb() {

	dbPath := os.Getenv("met_db")
	if dbPath == "" {
		dbPath = defaultDbPath()
	}
	database.InitDB("error", dbPath)
}
