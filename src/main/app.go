package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"sync"

	"../beans"
	"../config"
)

var conf *config.Config
var pageUrlChannel chan string
var poemUrlChannel chan string
var poems map[string]*beans.Poem
var wg1 sync.WaitGroup
var wg2 sync.WaitGroup

func init() {
	conf = config.GetConfig()
	pageUrlChannel = make(chan string, conf.Request_parallel_size)
	poemUrlChannel = make(chan string, conf.Request_parallel_size)
	poems = map[string]*beans.Poem{}
}

func main() {
	wg1.Add(10)
	//初始化pageUrls
	pages := getPageUrls()
	//启动pageUrls获取任务
	go producePagesUrls(pages)
	//启动poemUrl获取任务
	go producePoemUrls()
	//启动content获取任务
	go produceContents()
	wg1.Wait()
	wg2.Wait()
	jsonStr, _ := json.MarshalIndent(poems, "", "\t")
	ioutil.WriteFile(`C:\Users\fengl\Desktop\go\poemCrawler_Go\output.txt`, jsonStr, 0644)
}

func getPageUrls() []string {
	pages := []string{9: ""}
	for index := 0; index < 10; index++ {
		pages[index] = fmt.Sprintf("%s%d%s", "https://www.gushiwen.org/shiwen/default_4A111111111111A", index+1, ".aspx")
	}
	return pages
}

func producePagesUrls(pageUrls []string) {
	for _, url := range pageUrls {
		pageUrlChannel <- url
	}
}

func producePoemUrls() {
	for url := range pageUrlChannel {
		response, err := http.Get(url)
		if err != nil {
			log.Printf("%s>>>%s", url, err.Error())
			continue
		}
		defer response.Body.Close()
		content, _ := ioutil.ReadAll(response.Body)
		r := regexp.MustCompile(`<a style="font-size:18px; line-height:22px; height:22px;" href="(https://.*aspx)" target="_blank">`)
		ss := r.FindAllStringSubmatch(string(content), 20)
		for _, line := range ss {
			poemUrlChannel <- line[1]
			wg2.Add(1)
		}

		wg1.Done()
	}
}

// var yet bool

func produceContents() {
	for url := range poemUrlChannel {
		response, err := http.Get(url)
		if err != nil {
			log.Printf("%s===%s", url, err.Error())
		}
		defer response.Body.Close()
		poem := &beans.Poem{}
		content, _ := ioutil.ReadAll(response.Body)
		r := regexp.MustCompile(`<h1 style="font-size:20px; line-height:22px; height:22px; margin-bottom:10px;">(.+?)</h1>`)
		title := r.FindStringSubmatch(string(content))
		poem.Title = title[1]
		r1 := regexp.MustCompile(`<a href="/author.+?aspx">(.+?)</a>`)
		author := r1.FindStringSubmatch(string(content))
		if len(author) >= 2 {
			poem.Author = author[1]
		}
		r2 := regexp.MustCompile(`<textarea style=" background-color:#F0EFE2; border:0px;overflow:hidden;" cols="1" rows="1" id="txtare.+?">(.+?)——.*?https:.*?</textarea>`)
		con := r2.FindStringSubmatch(string(content))
		poem.Content = con[1]
		poems[url] = poem
		wg2.Done()
	}

}
