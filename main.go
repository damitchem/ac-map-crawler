package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"github.com/pkg/errors"
)

func main() {

	file, _ := os.Open("dungeons.txt")

	defer file.Close()
	format := "https://asheron.fandom.com/wiki/File:%v"
	scanner := bufio.NewScanner(file)
	lineNo := 1
	if _, err := os.Stat("maps"); os.IsNotExist(err) {
		os.Mkdir("maps", 644)
	}
	for scanner.Scan() {
		fmt.Printf("Starting line %v\r\n", lineNo)
		line := scanner.Text()
		pieces := strings.Split(line, ";")
		if len(pieces) == 0 {
			panic(fmt.Sprint("invalid structure on line", lineNo))
		}
		file := fmt.Sprintf("%v.gif", pieces[0])
		url := fmt.Sprintf(format, file)
		err := handleDownload(file, url)
		if err != nil {
			fmt.Println("Failed to download", err)
		}

		lineNo++
	}

}

func handleDownload(file, url string) error {
	req, _ := http.NewRequest("GET", url, nil)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "error doing request")
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return errors.Wrap(err, "error creating document from response")
	}

	doc.Find("a.internal").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if !exists {
			return
		}
		req, _ := http.NewRequest("GET", href, nil)
		res, err := http.DefaultClient.Do(req)

		if err != nil {
			fmt.Println(err)
			return
		}

		defer res.Body.Close()
		b, _ := ioutil.ReadAll(res.Body)
		fmt.Printf("Writing out file %v\r\n", file)
		err = ioutil.WriteFile(fmt.Sprintf("maps/%v", file), b, 644)
		if err != nil {
			fmt.Println(err)
		}
	})

	return nil
}
