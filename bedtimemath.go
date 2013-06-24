package main

// Download Bedtime math and store it as PDF
// Require's API key for convert API
// Depends on github.com/moovweb/gokogiri for xpath
import (
	"fmt"
	"github.com/moovweb/gokogiri"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// Alias for Println
var print = fmt.Println

type Url struct {
	value    string
	filename string
}
type Data struct {
	contents []byte
	filename string
}

// channel buffer count
const (
	queuecount int = 5
)

// xpath to get the actual webpage for each post
var xpath = `//*[@id="post-5198"]/div/div[1]/div[1]/div/div[1]/div[2]/div[1]/div[1]/div/div[1]/div/h2/a/@href`

func main() {

	urlchannel := make(chan Url, queuecount)
	grepchannel := make(chan Data, queuecount)
	downloadchannel := make(chan Url, queuecount)
	writechannel := make(chan Data, queuecount)

	// function that gets the environment variables as a map.
	getenvironment := func(data []string, getkeyval func(item string) (key, val string)) map[string]string {
		items := make(map[string]string)
		for _, item := range data {
			key, val := getkeyval(item)
			items[key] = val
		}
		return items
	}
	environment := getenvironment(os.Environ(), func(item string) (key, val string) {
		splits := strings.Split(item, "=")
		key = splits[0]
		val = splits[1]
		return
	})

	// url to convert the webpage to pdf
	var convertPDF = "http://do.convertapi.com/Web2Pdf?ApiKey=" + environment["KEY"] + "&PageSize=a5&CUrl="

	// downloads the page like http://bedtimemath.org/page/2/ and http://bedtimemath.org/page/3 to extract the actual
	// post the url using the xpath.
	download := func() {
		for url := range urlchannel {
			response, err := http.Get(url.value)
			defer response.Body.Close()
			contents, err := ioutil.ReadAll(response.Body)
			if err != nil {
				fmt.Printf("%s", err)
			}
			grepchannel <- Data{contents, url.filename}
		}
	}

	// parses the above downloaded page to get the actual url to the webpage
	parsehtml := func() {
		for html := range grepchannel {
			doc, _ := gokogiri.ParseHtml(html.contents)
			defer doc.Free()
			n, _ := doc.Root().Search(xpath)
			if len(n) < 1 {
				print("Could not find the element in " + html.filename)
			} else {
				downloadchannel <- Url{"http://viewtext.org/?url=" + n[0].String(), html.filename}
			}
		}
	}

	// downloads the pdf from the above url
	downloadPDF := func() {
		for html := range downloadchannel {
			response, err := http.Get(convertPDF + html.value)
			defer response.Body.Close()
			contents, err := ioutil.ReadAll(response.Body)
			if err != nil {
				fmt.Printf("%s", err)
			}
			writechannel <- Data{contents, html.filename}
		}
	}

	// write's the above pdf to the disk
	writefile := func() {
		for data := range writechannel {
			ioutil.WriteFile(data.filename+".pdf", data.contents, 0x777)
			print("Finished writing " + data.filename + ".pdf pdf.")
		}
	}
	for i := 0; i < queuecount; i++ {
		go writefile()
		go download()
		go parsehtml()
		go downloadPDF()
	}
	for i := 1; i < 100; i++ {
		url := "http://bedtimemath.org/page/" + strconv.Itoa(i)
		urlchannel <- Url{url, strconv.Itoa(i)}
	}
	time.Sleep(time.Minute * 5)
}
