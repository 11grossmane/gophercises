package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gophercises/urlshort"
)

func main() {
	yamlFileNamePointer := flag.String("path", "none", "this is the file path")
	flag.Parse()
	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	// Build the YAMLHandler using the mapHandler as the
	// fallback
	yml := `
- path: /urlshort
  url: https://github.com/gophercises/urlshort
- path: /urlshort-final
  url: https://github.com/gophercises/urlshort/tree/solution
`
	var yamlHandler http.HandlerFunc
	var err error
	if *yamlFileNamePointer == "none" {
		yamlHandler, err = urlshort.YAMLHandler([]byte(yml), mapHandler)
	} else {
		ymlFile, ymlErr := ioutil.ReadFile(*yamlFileNamePointer)
		fmt.Println("both yaml slices:", []byte(yml), "hhhhhhhhhhhhhhhh", ymlFile)
		if err != nil {
			fmt.Printf("yamlFile.Get err   #%v ", ymlErr)
		}
		yamlHandler, err = urlshort.YAMLHandler(ymlFile, mapHandler)
	}

	if err != nil {
		panic(err)
	}
	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", yamlHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
