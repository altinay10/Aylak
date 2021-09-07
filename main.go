package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/chromedp"
)

type Page struct {
	// page url withouht page number
	Url string `json:"url"`

	// Number of sub pages
	Count int `json:"count"`

	// if page has footer tag, it can be taken. Otherwise class or id that on the bottom of page can be given.
	FooterSelector string `json:"footerSelector"`

	// DOM object that enveloping items
	WrapperSelector string `json:"wrapperSelector"`

	// DOM object which has the desired value
	ItemSelector []string `json:"itemSelector"`
}

type Parser struct {
	// DOM object which has the desired value
	ItemSelector []string `json:"itemSelector"`

	// Content of wrapper element
	Content []byte
}

func newPages() []Page {
	var pages []Page
	jsonFile, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatal("err : ", err)
	}
	if err = json.Unmarshal(jsonFile, &pages); err != nil {
		log.Fatal("err : ", err)
	}
	return pages
}

func fetchPage(wrapperContent chan<- Parser, wg *sync.WaitGroup, page *Page, pageNumber int) {
	var wrapper string
	ctx, close := chromedp.NewContext(context.Background())
	defer close()

	ctx, close2 := context.WithTimeout(ctx, 40*time.Second)
	defer close2()

	defer wg.Done()
	fmt.Printf("%d. sayfa çekiliyor page = %+v\n", pageNumber, *page)
	if err := chromedp.Run(
		ctx,
		emulation.SetUserAgentOverride("aylak v0.1"),
		chromedp.Navigate(page.Url+strconv.Itoa(pageNumber)),
		chromedp.ScrollIntoView(page.FooterSelector),
		chromedp.WaitVisible("footer"),
		chromedp.OuterHTML(page.WrapperSelector, &wrapper)); err != nil {
		log.Println("Could not fetch website err: ", err)
		// will add rescue function
	}
	fmt.Println("url = ", page.Url+strconv.Itoa(pageNumber))
	wrapperContent <- Parser{Content: []byte(wrapper), ItemSelector: page.ItemSelector}
}

func Scrape(content chan<- Parser) {
	pages := newPages()
	var wg sync.WaitGroup
	wg.Add(len(pages))
	for _, page := range pages {
		go func(page Page) {
			for i := 1; i <= page.Count; i++ {
				wg.Add(1)
				go fetchPage(content, &wg, &page, i)
			}
			wg.Done()
		}(page)
	}
	wg.Wait()
	close(content)
}

func Parse(content <-chan Parser) {
	for parser := range content {
		doc, err := goquery.NewDocumentFromReader(bytes.NewReader(parser.Content))
		if err != nil {
			log.Println("Could not create document, err: ", err)
		}
		for _, selector := range parser.ItemSelector {
			f, err := os.OpenFile("data.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				log.Println("err: ", err)
			}
			doc.Find(selector).Each(func(i int, item *goquery.Selection) {
				if _, err := f.WriteString(item.Text()); err != nil {
					log.Println("err: ", err)
				}
			})
			if err := f.Close(); err != nil {
				log.Println("err: ", err)
			}
		}
	}
}

func main() {
	ch := make(chan Parser)

	go Scrape(ch)
	Parse(ch)
}
