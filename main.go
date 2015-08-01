package main

import (
	"GetWeb/funcs"
	"flag"
	"github.com/PuerkitoBio/goquery"
	"log"
	"os"
	"path"
	"strings"
)

func init() {
	log.SetFlags(log.Lshortfile)
	log.SetPrefix("GetWeb | ")
}

var outDir string

func main() {
	var url = flag.String("u", "", "抓取url地址")
	flag.StringVar(&outDir, "d", "outdir", "抓取结果存储目录")
	flag.Parse()

	if !flag.Parsed() || *url == "" {
		flag.PrintDefaults()
		os.Exit(-1)
	}

	doc, err := goquery.NewDocument(*url)
	if err != nil {
		log.Println("get url:", err)
		os.Exit(-1)
	}

	sel := doc.Find("head link[rel=stylesheet]")
	sel = sel.Add("script")
	sel = sel.Add("img")
	htmlFind := funcs.GetValidName("")
	cssFind := funcs.GetValidName("css")
	jsFind := funcs.GetValidName("js")
	imgFind := funcs.GetValidName("img")
	seg := strings.Split(*url, "/")
	if len(seg) < 3 {
		log.Println("url format invalid")
		os.Exit(-1)
	}

	outDir2 := htmlFind(seg[len(seg)-1])
	os.Remove(outDir)
	os.Remove(outDir2)
	if err = os.Mkdir(outDir2, 0744); err != nil {
		if err = os.Mkdir(outDir, 0744); err != nil {
			log.Println("mkdir output directory:", outDir, "failed, error:", err)
			os.Exit(-1)
		}
	} else {
		outDir = outDir2
	}

	if err = funcs.CrawlHtml(*url, outDir+"/index.html"); err != nil {
		log.Println("CrawHtml failed:", *url, ", error:", err)
		os.Exit(-1)
	}

	sel.Each(func(indx int, selection *goquery.Selection) {
		fileType := ""
		innerfile := ""
		var finder func(string) string
		switch selection.Nodes[0].Data {
		case "link":
			innerfile, _ = selection.Attr("href")
			fileType = "css"
			finder = cssFind
		case "script":
			innerfile, _ = selection.Attr("src")
			fileType = strings.ToLower(path.Ext(innerfile))
			if fileType == "js" {
				finder = jsFind
			} else {
				innerfile = ""
				fileType = "invalid"
			}
		case "img":
			innerfile, _ = selection.Attr("src")
			fileType = "img"
			finder = imgFind
		default:
			log.Println("not supported tag")
		}
		seg = strings.Split(innerfile, "/")
		innerFileName := innerfile
		if innerfile != "" {
			if len(seg) > 0 {
				lastpart := seg[len(seg)-1]
				innerFileName = finder(lastpart)
				if innerFileName == "" {
					log.Println("find valid", fileType, ":", innerfile, "failed")
					innerFileName = lastpart
				} else {
					log.Println("find", fileType, ":", innerFileName)
				}
			}
			if finalOut, err := funcs.CrawlInnerFile(*url, innerfile, outDir, innerFileName); err != nil {
				log.Println("get url to local failed:", err)
				os.Exit(-1)
			} else {
				if fileType == "css" {
					backimgs := funcs.CssFileBackgroundImages(finalOut)
					for _, img := range backimgs {
						seg := strings.Split(img, "/")
						lastpart := seg[len(seg)-1]
						imgName := imgFind(lastpart)
						if imgName == "" {
							imgName = lastpart
						}
						if img[0] != '/' {
							img = path.Dir(innerfile) + "/" + img
						}
						if imgOut, err := funcs.CrawlInnerFile(*url, img, outDir, imgName); err != nil {
							log.Println("download", img, "failed, error:", err)
						} else {
							log.Println("download", img, "at: [", imgOut, "] success!")
						}
					}
				}
			}

		}

	})
}
