package parser

import (
	"errors"
	"fmt"
	"io"
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
func GetTextAndLinks(htmlStr string, domain string) ([]string, []string, string, error) {
	// obtains tree strucute of the html
	htmlStr = CleanText(htmlStr)
	doc, err := html.Parse(strings.NewReader(htmlStr))
	if err != nil {
		return nil, nil, "", errors.New("Unable to parse URL")
	}

	// perform DFS to obtain all text nodes
	text, links := []string{}, []string{}
	var title string
	dfs(doc, &text, &links, &title)
	return text, links, title, nil
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

func GetDomain(url string) string {
	// remove http:// or https://
	if strings.HasPrefix(url, "http://") {
		url = url[len("http://"):]
	} else if strings.HasPrefix(url, "https://") {
		url = url[len("https://"):]
	} else {
		return "" // not http
	}
	// remove path
	if idx := strings.Index(url, "/"); idx != -1 {
		url = url[:idx]
	}

	// Remove port if present
	if idx := strings.Index(url, ":"); idx != -1 {
		url = url[:idx]
	}

	return url
}

func NormalizeURL(raw string) string {
	u, err := url.Parse(raw)
	if err != nil {
		return raw
	}

	// append with index.html when the url doesnt change (domain only)
	if u.Path == "" || u.Path == "/" {
		u.Path = "/index.html"
		// append if trailing / (not empty path)
	} else if strings.HasSuffix(u.Path, "/") {
		u.Path = u.Path + "index.html"
		// add .html as we only care for this
	} else if !strings.HasSuffix(u.Path, ".html") {
		u.Path = u.Path + ".html"
	}

	// remove multiple slashes inside path
	u.Path = strings.ReplaceAll(u.Path, "//", "/")

	return u.String()
}

func RemoveFragment(raw string) string {
	u, err := url.Parse(raw)
	if err != nil {
		// return original if parsing fails
		return raw
	}

	u.Fragment = ""
	return u.String()
}

func ValidateLinks(links []string, curUrl string, domain string, allowExternal bool) []string {
	valids := []string{}

	for _, link := range links {
		var validLink string
		var ok error

		// normalize urls by removing index.html (as implied by browser already)
		link = strings.TrimSuffix(link, "index.html")

		if isValidURL(link) {
			validLink = link
		} else if validLink, ok = absolutePathToUrl(link, curUrl); ok == nil {
		} else if validLink, ok = relativePathToUrl(link, curUrl); ok == nil {
		} else {
			continue
		}

		// check domain after we have a valid URL
		curDomain := GetDomain(validLink)
		if !allowExternal && curDomain != domain {
			continue
		}

		valids = append(valids, NormalizeURL(validLink))
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
	if len(absPath) < 1 {
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
	addTrailingSlash := false
	if string(absPath[0]) != "/" && string(curPath[len(curPath)-1]) != "/" {
		addTrailingSlash = true
	}

	fullURL := fmt.Sprintf("%s://%s", parsed.Scheme, parsed.Host)
	if addTrailingSlash {
		fullURL += "/"
	}
	fullURL += absPath
	return fullURL, nil
}

func relativePathToUrl(relPath string, curPath string) (string, error) {
	base, err := url.Parse(curPath)
	if err != nil {
		return "", err
	}

	ref, err := url.Parse(relPath)
	if err != nil {
		return "", err
	}

	resolved := base.ResolveReference(ref)
	return resolved.String(), nil
}

// general purpose dfs that parses through html nodes looking for queried tag values
func dfs(head *html.Node, result *[]string, links *[]string, title *string) {
	if head == nil {
		return
	}

	// Skip <script> and <style>, cant have subtree within this
	if head.Type == html.ElementNode &&
		(head.Data == "script" || head.Data == "style") {
		return
	}

	// get title
	if head.Type == html.ElementNode && head.Data == "title" && head.FirstChild != nil {
		*title = head.FirstChild.Data
	}

	// get text nodes
	if head.Type == html.TextNode {
		*result = append(*result, strings.ToLower(head.Data))
	}

	// get links
	if isLinkNode(head) {
		href := getHref(head)
		*links = append(*links, href)
	}

	// check children DFS
	for child := head.FirstChild; child != nil; child = child.NextSibling {
		dfs(child, result, links, title)
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
