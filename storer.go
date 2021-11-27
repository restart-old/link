package link

import (
	"encoding/json"
	"os"
	"time"
)

type Storer interface {
	Store(username string, Code Code) error
	LoadByCode(code string) error
	LoadByUser(username string) error
}

type JSONStorer struct {
	Folder string
}

func NewJSONStorer(folder string) *JSONStorer {
	os.Mkdir(folder, 727)
	return &JSONStorer{Folder: folder}
}

// Store stores the code and username provided
func (s JSONStorer) Store(username string, code Code) error {
	f, err := os.OpenFile(s.codepath(), os.O_RDWR, 727)
	if os.IsNotExist(err) {
		os.WriteFile(s.codepath(), []byte("{}"), 727)
	}
	defer f.Close()
	return s.store(username, code)
}

// LoadByCode loads the username that currently has the code provided
func (s JSONStorer) LoadByCode(code string) (username string, ok bool) {
	if code == "" {
		return username, false
	}
	username, ok = loadbycode(code, s.codepath())
	if !ok {
		removeCode(s.codepath(), username)
		return "", false
	}
	return
}

// LoadByUser loads the code that the user provided currently has
func (s JSONStorer) LoadByUser(username string) (code Code, ok bool) {
	if username == "" {
		return code, false
	}
	code, ok = loadbyuser(username, s.codepath())
	if !ok {
		removeCode(s.codepath(), username)
		return Code{}, false
	}
	return
}

// codepath returns the path of the codes.json file
func (s JSONStorer) codepath() string { return s.Folder + "codes.json" }

// store...
func (s JSONStorer) store(username string, code Code) error {
	codes, err := collectCodesData(s.codepath())
	if err != nil {
		return err
	}
	codes[username] = code

	dataBuf, err := json.MarshalIndent(codes, "", "\t")
	if err != nil {
		return err
	}
	return os.WriteFile(s.codepath(), dataBuf, 727)
}

// loadbycode...
func loadbycode(code, file string) (username string, ok bool) {
	codes, err := collectCodesData(file)
	if err != nil {
		return username, false
	}
	for u, c := range codes {
		if c.Code == code {
			username = u
			return username, c.Expiration.After(time.Now())
		}
	}
	return username, false
}

//loadbyuser...
func loadbyuser(username, file string) (code Code, ok bool) {
	codes, err := collectCodesData(file)
	if err != nil {
		return code, false
	}
	code, ok = codes[username]
	return code, ok && code.Expiration.After(time.Now())
}
