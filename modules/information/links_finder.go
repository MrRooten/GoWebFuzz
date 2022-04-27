package information

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"gowebfuzz/modules"
	"net/http"
	"regexp"
	"strings"
)

type LinksFinder struct {
	modules.ModuleInfo
}

func (finder LinksFinder) MatchRequest(req *http.Request) error {
	return nil
}

func (finder LinksFinder) MatchResponse(req *http.Request, res *http.Response, body *[]byte) error {
	if strings.Contains(res.Header.Get("Content-Type"), "image") {
		return nil
	}

	if strings.Contains(res.Header.Get("Content-Type"), "font") {
		return nil
	}

	fmt.Println(req.RequestURI)
	if strings.Contains(res.Header.Get("Content-Type"), "html") {

		doc, err := goquery.NewDocumentFromReader(bufio.NewReader(bytes.NewReader(*body)))
		if err != nil {
			return err
		}

		doc.Find("a[href]").Each(func(index int, item *goquery.Selection) {
			href, _ := item.Attr("href")
			fmt.Printf("\tlink: %s - anchor text: %s\n", href, item.Text())
		})
	}

	if strings.Contains(res.Header.Get("Content-Type"), "xml") {
		fmt.Println(*body)
	}

	if strings.Contains(res.Header.Get("Content-Type"), "javascript") {
		for {
			urlPattern, _ := regexp.Compile("([a-zA-Z]{1,10}://|//)[^\"'/]{1,}\\.[a-zA-Z]{2,}[^\"' ]{0,}")
			pathPattern, _ := regexp.Compile("\\.?(/[a-zA-Z0-9]+)+/?")
			filePattern, _ := regexp.Compile("[a-zA-Z0-9_\\-]{1,}\\.(php|asp|aspx|jsp|json|action|html|js|txt|xml)(\\?[^\"|']{0,}|)")

			urls := urlPattern.FindAll(*body, -1)
			for _, url := range urls {
				fmt.Println("\t", string(url))
			}

			paths := pathPattern.FindAll(*body, -1)
			for _, path := range paths {
				fmt.Println("\t", string(path))
			}

			files := filePattern.FindAll(*body, -1)
			for _, file := range files {
				fmt.Println("\t", string(file))
			}
			break
		}
	}
	return nil
}
