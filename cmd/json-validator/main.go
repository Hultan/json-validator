package main

import (
	"fmt"

	"github.com/hultan/json-validator/internal/parser"
)

func main() {
	p := parser.NewParser("test3.json")
	result := p.Parse()
	fmt.Println("Validation : ", result)
	for _, e := range p.Errors {
		fmt.Printf("%s : Near %s (%v,%v)\n", e.Message, string(e.Token.Kind), e.Token.Line, e.Token.Column)
	}

}
