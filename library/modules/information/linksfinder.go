package information

import (
	"bufio"
	"bytes"
	"compress/zlib"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"gowebfuzz/library/modules"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

type LinksFinder struct {
	modules.ModuleInfo
}

func (finder *LinksFinder)MatchRequest(req *http.Request) error {
	return nil
}

func (finder *LinksFinder)MatchResponse(res *http.Response) error {
	if strings.Contains(res.Header.Get("Content-Type"),"image") {
		return nil
	}

	if strings.Contains(res.Header.Get("Content-Type"),"font") {
		return nil
	}

	resBody,err := ioutil.ReadAll(res.Body)
	if err != nil {

	}
	res.Body = ioutil.NopCloser(bytes.NewReader(resBody))
	//currPath := req.RequestURI


	if strings.Contains(res.Header.Get("Content-Encoding"),"gzip") {
		r,err := zlib.NewReader(bytes.NewReader(resBody))
		if err != nil {
			return err
		}
		resBody,err = ioutil.ReadAll(r)
		if err != nil {
			return err
		}
	}

	if strings.Contains(res.Header.Get("Content-Type"),"html") {

		doc, err := goquery.NewDocumentFromReader(bufio.NewReader(bytes.NewReader(resBody)))
		if err != nil {
			return err
		}

		doc.Find("a[href]").Each(func(index int, item *goquery.Selection) {
			href, _ := item.Attr("href")
			fmt.Printf("link: %s - anchor text: %s\n", href, item.Text())

		})
	}

	if strings.Contains(res.Header.Get("Content-Type"),"javascript") {
		for {
			urlPattern, err := regexp.Compile("([a-zA-Z]{1,10}://|//)[^\"'/]{1,}\\.[a-zA-Z]{2,}[^\"' ]{0,}")
			if err != nil {
				break
			}
			pathPattern, err := regexp.Compile("(/|\\.\\./|\\./)[^\"'><,;| *()(%%$^/\\\\\\[\\]][^\"'><,;|()]{1,}")
			if err != nil {
				break
			}
			filePattern, err := regexp.Compile("[a-zA-Z0-9_\\-]{1,}\\.(php|asp|aspx|jsp|json|action|html|js|txt|xml)(\\?[^\"|']{0,}|)")
			if err != nil {
				break
			}

			urls := urlPattern.FindAll(resBody,-1)
			for _, url := range urls {
				fmt.Println(string(url))
			}

			paths := pathPattern.FindAll(resBody,-1)
			for _,path := range paths {
				fmt.Println(string(path))
			}

			files := filePattern.FindAll(resBody,-1)
			for _,file := range files {
				fmt.Println(string(file))
			}
			break
		}
	}
	return nil
}