package link

import (
	"database/sql"
	"fmt"
	"time"
)

type MySQLResponse struct {
	username    string
	discordID   string
	xuid        string
	linkedSince time.Time
}

func (m *MySQLResponse) Username() string       { return m.username }
func (m *MySQLResponse) DiscordID() string      { return m.discordID }
func (m *MySQLResponse) XUID() string           { return m.xuid }
func (m *MySQLResponse) LinkedSince() time.Time { return m.linkedSince }

type Linker struct {
	db *sql.DB
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

func (l *Linker) LinkedFromDiscordID(discordID string) (*MySQLResponse, bool, error) {
	var v string
	r := &MySQLResponse{}
	rows, err := l.db.Query(fmt.Sprintf("SELECT * FROM link WHERE discord_id='%s';", discordID))
	if err != nil {
		return r, false, err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&r.username, &r.discordID, &v, &r.xuid)
		if err != nil {
			return r, false, err
		}
	}
	if r.discordID == "" || r.username == "" || r.xuid == "" {
		return r, false, nil
	}
	r.linkedSince, _ = time.Parse(time.RFC3339, v)
	return r, r != nil, err
}
func (l *Linker) LinkedFromXUID(xuid string) (*MySQLResponse, bool, error) {
	var v string
	r := &MySQLResponse{}
	rows, err := l.db.Query(fmt.Sprintf("SELECT * FROM link WHERE xuid='%s';", xuid))
	if err != nil {
		return r, false, err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&r.username, &r.discordID, &v, &r.xuid)
		if err != nil {
			return r, false, err
		}
	}
	if r.discordID == "" || r.username == "" || r.xuid == "" {
		return r, false, nil
	}
	r.linkedSince, _ = time.Parse(time.RFC3339, v)
	return r, r != nil, err
}

func (l *Linker) LinkedFromGamerTag(gamertag string) (*MySQLResponse, bool, error) {
	var v string
	r := &MySQLResponse{}
	rows, err := l.db.Query(fmt.Sprintf("SELECT * FROM link WHERE username='%s';", gamertag))
	if err != nil {
		return r, false, err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&r.username, &r.discordID, &v, &r.xuid)
		if err != nil {
			return r, false, err
		}
	}
	if r.discordID == "" || r.username == "" || r.xuid == "" {
		return r, false, nil
	}
	r.linkedSince, _ = time.Parse(time.RFC3339, v)
	return r, r != nil, err
}

func (l *Linker) Link(username, code, discordID string) (err error) {
	if c, xuid, ok := l.LoadByUser(username); !ok {
		return err
	} else if c.Code == code {
		err = link(username, discordID, xuid, l.db)
		if err != nil {
			return err
		}
	}
	return
}

func link(username, discordID, xuid string, db *sql.DB) error {
	unLink(username, db)
	insert, err := db.Query("INSERT INTO link VALUES (?, ?, ?, ?);", username, discordID, xuid, time.Now().Format(time.RFC3339))
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
