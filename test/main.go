package main

import (
	"database/sql"
	"fmt"

	"github.com/SGPractice/link"
	"github.com/go-sql-driver/mysql"
)

func main() {
	config := mysql.NewConfig()
	config.DBName = "GlowHCF"
	config.User = "root"
	config.Addr = ":3306"
	config.Passwd = "f37JZEUm2QFexguhRuyscW{AdrKr86KajFGf%VT2h6BJUUF"
	config.Net = "tcp"

	connector, _ := mysql.NewConnector(config)
	db := sql.OpenDB(connector)
	defer db.Close()
	storer := link.NewJSONStorer("./link/")
	storer.Store("RestartFU", link.NewCode(7))
	linker := link.NewLinker(db, storer)
	if code, ok := storer.LoadByUser("RestartFU"); ok {
		if err := linker.Link("RestartFU", code.Code, "12412413453"); err != nil {
			fmt.Println(err)
		}
	}
	r, ok := linker.LinkedFromGamerTag("RestartFU")
	if ok {
		fmt.Println(r.LinkedSince().Date())
	}
}
