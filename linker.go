package link

import (
	"database/sql"
	"fmt"
	"time"
)

type Linker struct {
	db    *sql.DB
	users map[string]struct {
		LinkedAt time.Time
	}
	Storer
}

func NewLinker(db *sql.DB, s Storer) *Linker {
	return &Linker{
		db:     db,
		Storer: s,
	}
}

func (l *Linker) Cache() {

}

func (l *Linker) LinkedFromDiscordID(discordID string) bool {
	var v interface{}
	rows, err := l.db.Query(fmt.Sprintf("SELECT * FROM link WHERE discord_id='%s';", discordID))
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&v, &v)
		if err != nil {
			return false
		}
	}
	return err == nil && v != nil
}
func (l *Linker) LinkedFromGamerTag(gamertag string) bool {
	var v interface{}
	rows, err := l.db.Query(fmt.Sprintf("SELECT * FROM link WHERE username='%s';", gamertag))
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&v, &v)
		if err != nil {
			return false
		}
	}
	return err == nil && v != nil
}

func (l *Linker) Link(username, code, discordID string) (err error) {
	if c, ok := l.LoadByUser(username); !ok {
		return err
	} else if c.Code == code {
		err = link(username, discordID, l.db)
		if err != nil {
			return err
		}
	}
	return
}

func link(username, discordID string, db *sql.DB) error {
	insert, err := db.Query("INSERT INTO link VALUES (?, ?);", username, discordID)
	if err != nil {
		return err
	}
	defer insert.Close()
	return nil
}

func (l *Linker) UnLink(username string) error {
	return unLink(username, l.db)
}

func unLink(username string, db *sql.DB) error {
	delete, err := db.Query(fmt.Sprintf("DELETE FROM link WHERE username='%s';", username))
	if err != nil {
		return err
	}
	defer delete.Close()
	return nil
}
