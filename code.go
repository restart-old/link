package link

import (
	"crypto/rand"
	"encoding/json"
	"math/big"
	"os"
	"strings"
	"time"
)

type Code struct {
	Code       string
	XUID       string
	Expiration time.Time
}

func NewCode(length int, xuid string) Code {
	str := "ABCDEFGHIJKLMNOPQRSTUVWXYZ123456789"
	list := strings.Split(str, "")
	return Code{Code: newcode(list, length), Expiration: time.Now().Add(15 * time.Minute), XUID: xuid}
}
func newcode(list []string, length int) (s string) {
	for i := 0; i != length; i++ {
		n := randIntn(len(list) - 1)
		s += lowerUpper(list[n])
	}
	return s
}

func lowerUpper(str string) string {
	if randBool() {
		return strings.ToLower(str)
	}
	return str
}

func randIntn(n int) int {
	c, _ := rand.Int(rand.Reader, big.NewInt(int64(n)))
	return int(c.Int64())
}
func randBool() bool {
	return randIntn(2) == 1
}

func RemoveCode(file, username string) error {
	codes, err := collectCodesData(file)
	if err != nil {
		return err
	}
	delete(codes, username)
	databuf, _ := json.MarshalIndent(codes, "", "\t")
	return os.WriteFile(file, databuf, 0777)
}

func collectCodesData(file string) (v map[string]Code, err error) {
	b, err := os.ReadFile(file)
	if err != nil {
		return v, err
	}
	json.Unmarshal(b, &v)
	return v, err
}
