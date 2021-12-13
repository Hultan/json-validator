package main

import (
	"fmt"

	"github.com/hultan/json-validator/internal/json-parser"
)

func main() {
	p := json_parser.NewParser("test.json")
	result := p.Validate()
	fmt.Println("Validation : ", result)
	for _, e := range p.Errors {
		fmt.Printf("%s : Near %s (%v,%v)\n", e.Message, string(e.Token.Kind), e.Token.Line, e.Token.Column)
	}

}
