package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/gophercises/quiet_hn/hn"
)

var wg sync.WaitGroup

func main() {
	// parse flags

	var port, numStories int
	flag.IntVar(&port, "port", 3000, "the port to start the web server on")
	flag.IntVar(&numStories, "num_stories", 30, "the number of top stories to display")
	flag.Parse()
	tpl := template.Must(template.ParseFiles("./index.gohtml"))

	http.HandleFunc("/", handler(numStories, tpl))

	// Start the server
	fmt.Printf("listening on port %d", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func handler(numStories int, tpl *template.Template) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		var client hn.Client
		ids, err := client.TopItems()
		if err != nil {
			http.Error(w, "Failed to load top stories", http.StatusInternalServerError)
			return
		}
		var stories []item
		channelResults := []Results{}
		ch := make(chan Results)
		for idx, id := range ids[0:30] {
			fmt.Println(idx)
			go fetchStory(id, idx, &stories, &ch)
			time.Sleep(5 * time.Millisecond)
		}
		for i := 0; i < numStories; i++ {
			channelResults = append(channelResults, <-ch) //this blocks execution until channelResults receives something

			//fmt.Println(channelResults)
		}
		sort.Slice(channelResults, func(i, j int) bool {
			return channelResults[i].pos < channelResults[j].pos
		})
		for _, res := range channelResults {
			stories = append(stories, res.item)
		}

		data := templateData{
			Stories: stories,
			Time:    time.Since(start),
		}
		err = tpl.Execute(w, data)
		if err != nil {
			http.Error(w, "Failed to process the template", http.StatusInternalServerError)
			return
		}
	})
}

var count int

func fetchStory(id int, idx int, stories *[]item, ch *chan Results) {
	client := hn.Client{}

	hnItem, err := client.GetItem(id)
	if err != nil {
		//fmt.Println(err)
	}
	item := parseHNItem(hnItem)

	if err == nil {

		*ch <- Results{item: item, pos: idx}
		//*stories = append(*stories, item)
	}

}

func isStoryLink(item item) bool {
	return item.Type == "story" && item.URL != ""
}

func parseHNItem(hnItem hn.Item) item {
	ret := item{Item: hnItem}
	url, err := url.Parse(ret.URL)
	if err == nil {
		ret.Host = strings.TrimPrefix(url.Hostname(), "www.")
	}
	return ret
}

// item is the same as the hn.Item, but adds the Host field
type item struct {
	hn.Item
	Host string
}

type templateData struct {
	Stories []item
	Time    time.Duration
}

type Results struct {
	item item
	err  error
	pos  int
}
