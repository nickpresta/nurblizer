package main

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"runtime"
	"strings"
)

var templates = template.Must(template.ParseFiles("index.html", "nurble.html"))

func mainRequestHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "index.html", "")
}

func nurbleRequestHandler(w http.ResponseWriter, r *http.Request) {
	text := r.FormValue("text")
	text = nurble(text)
	safe := template.HTMLEscapeString(text)
	safe = strings.Replace(safe, "\n", "<br>", -1)
	safe = strings.Replace(safe, "nurble", " <span class='nurble'>nurble</span> ", -1)
	templates.ExecuteTemplate(w, "nurble.html", template.HTML(safe))
}

func nurble(text string) string {
	upper := strings.ToUpper(text)
	words := wordRegexp.ReplaceAllString(upper, "")
	splitWords := strings.Fields(words)

	for _, word := range splitWords {
		if NOUNS[word] == false {
			pattern := regexp.MustCompile(" " + word + " ")
			upper = pattern.ReplaceAllString(upper, " nurble ")
		}
	}

	return upper
}

func readNouns() {
	file, err := ioutil.ReadFile("./nouns.txt")
	if err != nil {
		log.Fatal(err)
	}

	parts := bytes.Split(file, []byte{'\n'})
	NOUNS = make(map[string]bool, len(parts))
	for _, line := range parts {
		NOUNS[strings.ToUpper(string(line))] = true
	}
}

var (
	NOUNS      map[string]bool
	wordRegexp = regexp.MustCompile("[^A-Z ]")
)

func main() {
	numCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPU)

	readNouns()

	http.HandleFunc("/", mainRequestHandler)
	http.HandleFunc("/nurble", nurbleRequestHandler)
	err := http.ListenAndServe(":17562", http.DefaultServeMux)
	if err != nil {
		log.Fatal(err)
	}
}
