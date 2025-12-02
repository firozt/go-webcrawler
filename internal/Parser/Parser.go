package parser

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

// -------------------- PUBLIC -------------------- //

// parses the url, and returns the body of the http response
func ParseSite(url string) string {
	body, err := getBody(url)
	if err != nil {
		log.Fatal("Error trying to obtain body of the url: ", err)
	}
	return string(body)
}

// removes the html tags such as <div> <h1> etc, returns clean text
func GetTextAndLinks(htmlStr string) ([]string, []string) {
	// obtains tree strucute of the html
	doc, err := html.Parse(strings.NewReader(htmlStr))
	if err != nil {
		log.Fatal("Error parsing html: ", err)
	}

	// perform DFS to obtain all text nodes
	text, links := []string{}, []string{}
	dfs(doc, &text, &links)
	return text, links
}

// -------------------- PRIVATE -------------------- //

// checks if a url is valid
func isValidURL(u string) bool {
	parsed, err := url.Parse(u)
	if err != nil {
		return false // cannot parse
	}
	if (parsed.Scheme != "http" && parsed.Scheme != "https") || parsed.Host == "" {
		return false
	}
	return true
}

// general purpose dfs that parses through html nodes looking for queried tag values
func dfs(head *html.Node, result *[]string, links *[]string) {
	if head == nil {
		return
	}
	// check node type
	if head.Type == html.TextNode {
		*result = append(*result, head.Data)
	}
	if isLinkNode(head) {
		href := getHref(head)
		if isValidURL(href) {
			*links = append(*links, href)
		}
	}

	// iterate over all children nodes
	for child := head.FirstChild; child != nil; child = child.NextSibling {
		dfs(child, result, links)
	}
}

// gets href from a link node
func getHref(n *html.Node) string {
	for _, attr := range n.Attr {
		if attr.Key == "href" {
			return attr.Val
		}
	}

	return ""
}

func isLinkNode(n *html.Node) bool {
	return n.Type == html.ElementNode && n.Data == "a"
}

// runs a get request for a given url and returns its body
// may return errors
func getBody(url string) ([]byte, error) {
	// setup request
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	// set crawling headers for request
	req.Header.Set("User-Agent", "MyWebCrawler/1.0 firozt03@gmail.com")

	// make request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close() // we need to close the TCP connection after were done
	return io.ReadAll(resp.Body)
}
