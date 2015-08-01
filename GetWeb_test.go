package main

import (
	"GetWeb/funcs"
	"log"
	"strings"
	"testing"
)

func TestGetValidName(t *testing.T) {
	css := "214/214.css"
	css1 := "/214/214.css?v=8may2013"
	css2 := "214/214.css?v=8may2013"
	css3 := "214/214.css_v_8may2013"

	finder := funcs.GetValidName("css")
	imgFinder := funcs.GetValidName("img")
	expect := "214.css"

	seg := strings.Split(css, "/")
	n := len(seg)
	ret := finder(seg[n-1])
	if ret != expect {
		t.Fatal(css, "failed, retured:", ret)
	}
	seg = strings.Split(css1, "/")
	n = len(seg)
	ret = finder(seg[n-1])
	if ret != expect {
		t.Fatal(css1, "failed, retured:", ret)
	}
	seg = strings.Split(css2, "/")
	n = len(seg)
	ret = finder(seg[n-1])
	if ret != expect {
		t.Fatal(css2, "failed, retured:", ret)
	}
	seg = strings.Split(css3, "/")
	n = len(seg)
	ret = finder(seg[n-1])
	if ret != expect {
		t.Fatal(css3, "failed, retured:", ret)
	}

	url := "http://www.csszengarden.com/atlantis"
	htmlFinder := funcs.GetValidName("")
	seg = strings.Split(url, "/")
	n = len(seg)
	ret = htmlFinder(seg[n-1])
	if ret != "atlantis" {
		t.Fatal(url, "failed, retured:", ret)
	}
	url1 := "http://www.csszengarden.com/atlantis?a=123"
	seg = strings.Split(url1, "/")
	n = len(seg)
	ret = htmlFinder(seg[n-1])
	if ret != "atlantis" {
		t.Fatal(url1, "failed, retured:", ret)
	}

	img := "/img/img.jpg"
	seg = strings.Split(img, "/")
	n = len(seg)
	expect = "img.jpg"
	ret = imgFinder(seg[n-1])
	if ret != expect {
		t.Fatal(img, "failed, retured:", ret)
	}
	img1 := "/img/img.png"
	seg = strings.Split(img1, "/")
	n = len(seg)
	expect = "img.png"
	ret = imgFinder(seg[n-1])
	if ret != expect {
		t.Fatal(img1, "failed, retured:", ret)
	}
	img1 = "a/b/img/img.svg?v=123"
	seg = strings.Split(img1, "/")
	n = len(seg)
	expect = "img.svg"
	ret = imgFinder(seg[n-1])
	if ret != expect {
		t.Fatal(img1, "failed, retured:", ret)
	}
}

func TestCrawUrl(t *testing.T) {
	url := "http://www.csszengarden.com/"
	css1 := "/214/214.css?v=8may2013"
	outDir := "outdir"
	finder := funcs.GetValidName("css")

	seg := strings.Split(css1, "/")
	n := len(seg)
	name := finder(seg[n-1])
	log.Println("name =", name)
	finalOut, err := funcs.CrawlInnerFile(url, css1, outDir, name)
	if err != nil || finalOut != outDir+"/214/214.css" {
		t.Fatal("test CrawlFile failed, error:", err, "output at:", finalOut)
	}
}

func TestBackgroundImages(t *testing.T) {
	cssfile := []byte(`background: #2d6360 50% 0 url(huntington.jpg) no-repeat; /* old IE fallback */  background-attachment: fixed, fixed, fixed, scroll;  background-image: url(contours.png), url(noise.png), url(gridlines.png), url(huntington.jpg);  background-position: 0 0, 0 0, -5px -25px, 0 50%;`)

	l := funcs.CssByteBackgroundImages(cssfile)

	if len(l) > 0 {
		log.Println("find css background images:", l)
	} else {
		t.Fatal("TestBackgroundImages failed")
	}
}
