package main

import (
	"fmt"

	webcrawler "github.com/firozt/crawler/internal"
)

func main() {
	body := webcrawler.ParseSite("https://en.wikipedia.org/wiki/Chair")
	text, links := webcrawler.GetTextAndLinks(body)
	fmt.Print(len(text))
	fmt.Print(len(links))
}
