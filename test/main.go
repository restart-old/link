package main

import (
	"github.com/SGPractice/link"
)

func main() {
	link.NewJSONStorer("./link/").Store("RestartFU", link.NewCode(7))
}
