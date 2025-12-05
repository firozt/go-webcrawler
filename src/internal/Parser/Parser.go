package parser

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

// -------------------- PUBLIC -------------------- //

// parses the url, and returns the body of the http response
func ParseSite(url string) (string, error) {
	var err error
	body, err := getBody(url)
	if err != nil {
		return "", err
	}
	return string(body), err
}

// removes the html tags such as <div> <h1> etc, returns clean text and a list of links fround within href's
// texts, links
func GetTextAndLinks(htmlStr string) ([]string, []string) {
	// obtains tree strucute of the html
	htmlStr = CleanText(htmlStr)
	doc, err := html.Parse(strings.NewReader(htmlStr))
	if err != nil {
		log.Fatal("Error parsing html: ", err)
	}

	// perform DFS to obtain all text nodes
	text, links := []string{}, []string{}
	dfs(doc, &text, &links)
	return text, links
}

// removes whitespaces in htmls
func CleanText(raw string) string {
	// Remove leading/trailing whitespace
	text := strings.TrimSpace(raw)

	// Collapse multiple spaces/newlines/tabs into a single space
	re := regexp.MustCompile(`\s+`)
	text = re.ReplaceAllString(text, " ")

	return text

}

func ValidateLinks(links []string, curUrl string) []string {
	valids := []string{}
	for _, link := range links {

		// see if its a https link
		if isValidURL(link) {
			valids = append(valids, link)
			continue
		}

		// see if its an absolute path
		valid_link, ok := absolutePathToUrl(link, curUrl)

		if ok == nil {
			valids = append(valids, valid_link)
			continue
		}

		// see if its a relative path
		valid_link, ok = relativePathToUrl(link, curUrl)
		if ok == nil {
			valids = append(valids, valid_link)
			continue
		}

	}
	return valids
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
