package indexer

import (
	"crypto/sha256"
	"fmt"

	"github.com/PuerkitoBio/goquery"
	"github.com/geziyor/geziyor"
)

func Index(Message message) {
	geziyor.NewGeziyor(&geziyor.Options{
		StartURLs: []string{Message.URL},
		ParseFunc: func(g *geziyor.Geziyor, r *client.Response) {
			fmt.Println(string(r.Body))
		},
	}).Start()

	doc, err := goquery.NewDocument(contents)
	doc.Find(message.Selector).Each(func(i int, s *goquery.Selection) {
		// Extract the text from the selected element
		text := s.Text()

		// Print the extracted text
		fmt.Println(text)
	})

	checksum := sha256.New().Sum([]byte(parsedContent))

}
