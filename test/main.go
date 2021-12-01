package main

import (
	"database/sql"
	"fmt"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	config := mysql.NewConfig()
	config.DBName = "GlowHCF"
	config.User = "root"
	config.Addr = "51.222.111.58:3306"
	config.Passwd = "f37JZEUm2QFexguhRuyscW{AdrKr86KajFGf%VT2h6BJUUF"
	config.Net = "tcp"

	connector, _ := mysql.NewConnector(config)
	db := sql.OpenDB(connector)
	defer db.Close()
	fmt.Println(db.Ping())
}
