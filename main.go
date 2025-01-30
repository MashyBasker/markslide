package main

import (
	"fmt"

	"github.com/MashyBasker/markslide/internal/parser"
)

func main() {
	path := "./data/basic.md"
	delim := "[end page]"
	sections, err := parser.ParseMarkdownFile(path, []byte(delim))
	if err != nil {
		panic(err)
	}
	htmlsections, err := parser.ConvertToHTML(sections)

	if err != nil {
		panic(err)
	}
	for _, html := range htmlsections {
		fmt.Println(html.String())
		fmt.Printf("==================\n")
	}
}
