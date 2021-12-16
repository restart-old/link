package link

import (
	"encoding/json"
	"os"
	"time"
)

type Storer interface {
	Store(username string, Code Code) error

	LoadByCode(code string) (string, string, bool)
	LoadByUser(username string) (Code, string, bool)

	RemoveByCode(code string)
	RemoveByUser(username string)
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
	f, err := os.OpenFile(s.codepath(), os.O_RDWR, 0777)
	if os.IsNotExist(err) {
		os.WriteFile(s.codepath(), []byte("{}"), 0777)
	}
	defer f.Close()
	return s.store(username, code)
}

// LoadByCode loads the username that currently has the code provided
func (s JSONStorer) LoadByCode(code string) (username string, xuid string, ok bool) {
	if code == "" {
		return username, "", false
	}
	username, xuid, ok = loadbycode(code, s.codepath())
	if !ok {
		s.RemoveByUser(username)
		return "", "", false
	}
	return
}

// LoadByUser loads the code that the user provided currently has
func (s JSONStorer) LoadByUser(username string) (code Code, xuid string, ok bool) {
	if username == "" {
		return code, "", false
	}
	code, ok = loadbyuser(username, s.codepath())
	if !ok {
		s.RemoveByUser(username)
		return code, "", false
	}
	return code, code.XUID, true
}

func (s JSONStorer) RemoveByCode(code string) error {
	codes, err := collectCodesData(s.codepath())
	if err != nil {
		return err
	}
	for username, c := range codes {
		if c.Code == code {
			delete(codes, username)
		}
	}
	databuf, _ := json.MarshalIndent(codes, "", "\t")
	return os.WriteFile(s.codepath(), databuf, 0777)
}

func (s JSONStorer) RemoveByUser(username string) error {
	codes, err := collectCodesData(s.codepath())
	if err != nil {
		return err
	}
	delete(codes, username)
	databuf, _ := json.MarshalIndent(codes, "", "\t")
	return os.WriteFile(s.codepath(), databuf, 0777)
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
	return os.WriteFile(s.codepath(), dataBuf, 0777)
}

// loadbycode...
func loadbycode(code, file string) (username string, xuid string, ok bool) {
	codes, err := collectCodesData(file)
	if err != nil {
		return username, "", false
	}
	for u, c := range codes {
		if c.Code == code {
			username = u
			xuid = c.XUID
			return username, c.XUID, c.Expiration.After(time.Now())
		}
	}
	return username, xuid, false
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

func collectCodesData(file string) (v map[string]Code, err error) {
	b, err := os.ReadFile(file)
	if err != nil {
		return v, err
	}
	json.Unmarshal(b, &v)
	return v, err
}
