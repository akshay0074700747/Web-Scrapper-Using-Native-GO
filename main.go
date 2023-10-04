package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func Gethref(t html.Token) (ok bool, url string) {

	for _, v := range t.Attr {
		if v.Key == "href" {
			return true, v.Val
		}
	}

	return false, ""
}

func Crawl(url string, ch chan string, finurl chan bool) {

	// the get request to a url return an html respoce which consists of the whole html of the url

	resp, err := http.Get(url)

	defer func() {
		finurl <- true
	}()

	if err != nil {
		fmt.Println(url, "---- the given url cannot be crawled...")
		finurl <- false
		return
	}

	b := resp.Body

	defer b.Close()

	tokens := html.NewTokenizer(b)

	for {

		token := tokens.Next()

		switch {
		//here ErrorToken means thet is is the end of the html document
		case token == html.ErrorToken:
			return
			//checks if the token is a starting token eg : <a>,<p>,<h1>, etc ...
		case token == html.StartTagToken:
			t := tokens.Token()
			//assigns the result of this condition to isanchor
			//and checks if the tag is an <a> tag since we want links and it will be in <a></a> these tag
			isanchor := t.Data == "a"

			if !isanchor {
				continue
			}

			ok, uurl := Gethref(t)

			if !ok {
				continue
			}

			//checks if its a valid url if its a valid url then it should start with http, if the returned index of the substring http in the uurl is zero
			hashttp := strings.Index(uurl, "http") == 0

			if hashttp {
				ch <- uurl
			}
		}

	}

}

func main() {

	foundlinks := make(map[string]bool)
	seedurls := os.Args[1:]

	chanlinks := make(chan string)
	seedfinished := make(chan bool)

	for _, url := range seedurls {

		go Crawl(url, chanlinks, seedfinished)

	}

	for i := 0; i < len(seedurls); {

		select {
		case url := <-chanlinks:
			foundlinks[url] = true
		case <-seedfinished:
			i++
		}

	}

	fmt.Println("Found.... ", len(foundlinks), " urls")

	for url := range foundlinks {

		fmt.Println("--" + url)

	}

	close(chanlinks)

}
