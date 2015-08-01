package funcs

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"
)

var webClient http.Client

// 用闭包来隐藏局部变量在Go中变得不那么重要了，更多时候可以通过
// package实现
func GetValidName(suffix string) func(str string) string {
	pat := regexp.MustCompile(`(?i)^[\d\w_\-]+\.?` + suffix)
	if suffix == "img" {
		pat = regexp.MustCompile(`(?i)^[\d\w_\-]+\.(png|jpg|jpeg|gif|svg)`)
	}
	return func(str string) string {
		return pat.FindString(str)
	}
}

func CrawlHtml(url, outfile string) error {
	req, err := webClient.Get(url)
	defer req.Body.Close()
	if err != nil {
		return err
	}
	content := []byte{}
	if content, err = ioutil.ReadAll(req.Body); err == nil {
		err = ioutil.WriteFile(outfile, content, 0644)
	}
	return err

}
func CrawlInnerFile(url, innerfile, outDir, outfile string) (string, error) {
	getUrl := ""
	url = strings.ToLower(url)
	innerfile = strings.ToLower(innerfile)
	seg := strings.Split(url, "/")
	if len(seg) < 3 {
		return "", fmt.Errorf("cant get host of this url:", url)
	}
	host := strings.Join(seg[:3], "/")
	base := strings.Join(seg[:len(seg)-1], "/")
	if innerfile[:4] == "http" || innerfile[:2] == "//" {
		seg2 := strings.Split(innerfile, "/")
		if seg2[2] == seg[2] {
			innerfile = "/" + strings.Join(seg2[3:], "/")
			log.Println("trim innerfile:", innerfile)
		}
	}
	seg = strings.Split(innerfile, "/")
	outdir := ""
	if innerfile[0] == '/' {
		getUrl = host + innerfile
		outdir = outDir + "/" + strings.Join(seg[1:len(seg)-1], "/") + "/" + outfile
	} else {
		getUrl = base + innerfile
		outdir = outDir + "/" + strings.Join(seg[:len(seg)-1], "/") + "/" + outfile
	}
	log.Println("outdir:", outdir)
	dirName := path.Dir(outdir)
	if err := os.MkdirAll(dirName, 0744); err != nil {
		return outdir, err
	}

	return outdir, CrawlHtml(getUrl, outdir)
}

var imgBackgroundPat = regexp.MustCompile(`url\(([\d\w]+\.(png|jpg|jpeg|gif|svg))\)`)

func CssFileBackgroundImages(cssfile string) []string {
	reader, err := os.Open(cssfile)
	if err != nil {
		log.Println("when crawl background images of cssfile :", cssfile, "failed, error:", err)
		return []string{}
	}
	if cssbyte, err := ioutil.ReadAll(reader); err == nil {
		return CssByteBackgroundImages(cssbyte)
	} else {
		return []string{}
	}
}

func CssByteBackgroundImages(cssfile []byte) []string {
	ret := []string{}
	backgroundList := imgBackgroundPat.FindAllSubmatch(cssfile, -1)
	for _, m := range backgroundList {
		ret = append(ret, string(m[1]))
	}

	return ret
}
