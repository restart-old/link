package link

import (
	"database/sql"
	"time"
)

type Linker struct {
	Token string
	db    *sql.DB
	users map[string]struct {
		LinkedAt time.Time
	}
	Storer
}

func NewLinker(token string, db *sql.DB) *Linker {
	return &Linker{
		Token: token,
		db:    db,
	}
}

func (l *Linker) Cache() {

}

func (l *Linker) LinkedFromDiscordID(discordID string) bool
func (l *Linker) LinkedFromGamerTag(gamertag string) bool

func (l *Linker) Link(username, code, discordID string) (err error) {
	if c, err := l.LoadByUser(username); err != nil {
		return err
	} else if c == code {
		err = link(username, discordID, l.db)
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

func (*Linker) UnLink(username, discordID string) error {
	return nil
}

func unLink(username, discordID string) error {
	return nil
}
