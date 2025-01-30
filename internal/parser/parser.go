package parser

import (
	"bufio"
	"bytes"
	"os"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

/*
	ParseMarkdownFile takes a filepath and delimeter and returns
	markdown content as sections divided by the delimiter

	Paramaters:
		- filePath: A string that contains the path of the markdown file
		- delimimter: A []byte that contains the delimiter to decide where the markdown file will be split

	Returns:
		- [][]byte: A slice of byte slice containg all the markdown sections
		- error: An error for when the parsing fails, nil when it succeeds
*/
func ParseMarkdownFile(filePath string, delimiter []byte) ([][]byte, error) {
	// open the file
	file, err := os.Open(filePath)
	if err != nil {
		return [][]byte{}, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var buffer []byte
	var slideSections [][]byte

	// scan the markdown file
	for scanner.Scan() {
		line := scanner.Bytes()
		buffer = append(buffer, line...)
		buffer = append(buffer, '\n')

		// check for page end delimiter
		if bytes.Contains(buffer, delimiter) {
			parts := bytes.Split(buffer, delimiter)
			slideSections = append(slideSections, parts[0])
			buffer = parts[1]
		}
	}

	if len(buffer) > 0 {
		slideSections = append(slideSections, buffer)
	}

	return slideSections, nil
}

/*
markdownToHTML takes a markdown section and converts
it into an HTML section.

Parametes:
  - mdBuffer: A []byte containing the markdown content to be converted
  - mdConvertor: A goldmark.Markdown instance used to perform the conversion

Returns:
  - bytes.Buffer: A buffer containing the resulting HTML content
  - error: An error if the conversion fails, or nil if its successful
*/
func markdownToHTML(mdBuffer []byte, mdConvertor goldmark.Markdown) (bytes.Buffer, error) {
	var htmlBuffer bytes.Buffer
	if err := mdConvertor.Convert(mdBuffer, &htmlBuffer); err != nil {
		return bytes.Buffer{}, err
	}
	return htmlBuffer, nil
}

/*
ConvertsToHTML takes the markdown sections of a slide and converts
each markdown section to an HTML section

Parameters:
  - slideSections: [][]byte containg the slide as markdown sections

Returns:
  - []bytes.Buffer: slice of bytes.Buffer containing the HTML sections
  - error: returns an error when conversion to HTML failed, nil if succeeds
*/
func ConvertToHTML(slideSections [][]byte) ([]bytes.Buffer, error) {
	var htmlSections []bytes.Buffer
	// define the goldmark instance with parser/renderer options
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.Table,
			extension.Strikethrough,
			extension.TaskList,
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithXHTML(),
		),
	)

	for _, section := range slideSections {
		html, err := markdownToHTML(section, md)
		if err != nil {
			return []bytes.Buffer{}, err
		}
		htmlSections = append(htmlSections, html)
	}
	return htmlSections, nil
}
