package parser

import (
	"errors"
	"fmt"
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

// removes the html tags such as <div> <h1> etc, returns clean text and a list of links fround within href's
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

func absolutePathToUrl(absPath string, curPath string) (string, error) {
	if len(absPath) < 1 || !strings.HasPrefix(absPath, "/") {
		return "", errors.New("absPath is not a valid absolute path (must start with /)")
	}

	parsed, err := url.Parse(curPath)
	if err != nil {
		return "", errors.New("curPath URL is not valid")
	}

	scheme := parsed.Scheme
	host := parsed.Host

	if scheme == "" || host == "" {
		return "", errors.New("curPath must include scheme and host")
	}

	fullURL := fmt.Sprintf("%s://%s%s", parsed.Scheme, parsed.Host, absPath)
	return fullURL, nil
}
func relativePathToUrl(relPath string, curPath string) (string, error) {
	// ./../path/to/file
	var err error

	if len(relPath) < 1 && !strings.HasPrefix(curPath, "https://") {
		return "", errors.New("invalid relative path / url")
	}

	var pathArr []string = strings.Split(curPath, "/")
	var relPathArr []string = strings.Split(relPath, "/")
	pathArr = pathArr[:len(pathArr)-1]

	for _, action := range relPathArr {
		fmt.Println(action)
		fmt.Println(strings.Join(pathArr, "/"))

		if action == "." {
			// skip
			continue
		} else if action == ".." {
			// invalid current path, doesnt make sense
			if len(pathArr)-1 < 0 {
				return "", errors.New("invalid curpath")
			}
			pathArr = pathArr[:len(pathArr)-1]
		} else {
			pathArr = append(pathArr, action)
		}
		fmt.Println(strings.Join(pathArr, "/"))
	}
	res := strings.Join(pathArr, "/")
	if !strings.Contains(pathArr[len(pathArr)-1], ".") {
		res += "/" // trailing / for folder
	}
	return res, err

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
		*links = append(*links, href)
	}

	// iterate over all children nodes
	for child := head.FirstChild; child != nil; child = child.NextSibling {
		dfs(child, result, links)
	}
}

// gets link from a href node (no url validation)
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
	req.Header.Set("User-Agent", "gowebcrawler/1.0 firozt03@gmail.com")

	// make request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}
