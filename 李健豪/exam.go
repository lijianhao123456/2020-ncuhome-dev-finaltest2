package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	t1 := time.Now()
	const dataQuantity = 54
	const pageAmout = 7
	title := make(chan string, dataQuantity)
	date := make(chan string, dataQuantity)
	label := make(chan string, dataQuantity)
	excerpt := make(chan string, dataQuantity)
	for i := 1; i <= pageAmout; i++ {
		res, err := http.Get("https://blog.lenconda.top/page/" + strconv.Itoa(i) + "/")
		if err != nil {
			panic(err)
		}
		fmt.Println(res)
		defer res.Body.Close()
		doc, err := goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			panic(err)
		}
		doc.Find("h2.post-title>a").Each(func(i int, s *goquery.Selection) {
			title <- s.Text()
		})
		doc.Find("span.post-meta>time").Each(func(i int, s *goquery.Selection) {
			date <- s.Text()
		})
		doc.Find("span.post-meta:nth-child(3)").Each(func(i int, s *goquery.Selection) {
			result := strings.Replace(s.Contents().Text(), " ", "", -1)
			result = strings.Replace(result, "\n", "", -1)
			result = strings.Replace(result, "，", "", -1)
			label <- result
		})
		doc.Find("p.post-excerpt").Each(func(i int, s *goquery.Selection) {
			excerpt <- s.Text()
		})

	}

	close(title)
	close(date)
	close(label)
	close(excerpt)

	f, err := os.Create("./exam.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	for i := 1; i <= dataQuantity; i++ {
		f.WriteString(<-title + "\t" + <-date + "\t" + <-label + "\t" + <-excerpt + "\t" + "\r\n")
	}
	for i := 1; i <= pageAmout; i++ {
		url := ("https://blog.lenconda.top/page/" + strconv.Itoa(i) + "/")
		resp, _ := http.Get(url)
		defer resp.Body.Close()
		r, _ := (ioutil.ReadAll(resp.Body))
		result := string(r)
		fileName := "page" + strconv.Itoa(i) + ".html"
		f, err1 := os.Create(fileName)
		if err1 != nil {
			fmt.Println("os Create err1 = ", err1)
			continue
		}
		f.WriteString(result)
		f.Close()
	}
	elapsed := time.Since(t1)
	fmt.Println("爬虫结束,总共耗时: ", elapsed)
}
