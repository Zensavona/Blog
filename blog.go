package main

import (
	"fmt"
	"path"
	"github.com/russross/blackfriday"
	"github.com/hoisie/mustache"
	"code.google.com/p/gorilla/mux"
	"io/ioutil"
	"strings"
	"net/http"
)

type Note struct {
	Uglyname string
	Title string
	Date string
	Body string
}

const sep = "---"
const notePath = "/note/"
const lenPath = len(notePath)

var notes = loadNotes()

// less than ideal
var home, _ = ioutil.ReadFile("home.md")
var homeMarkup = string(blackfriday.MarkdownCommon(home))

func loadNotes() []Note {
	files, _ := ioutil.ReadDir("notes/")
	var notes []Note
	for _, file := range files {
		title := strings.Replace(file.Name(), ".md", "", -1)
		notes = append(notes, *loadPost(title))
	}
	return notes
}

func getContent(title string) (content string) {
	filename := title + ".md"
    file, _ := ioutil.ReadFile(path.Join("notes", filename))
    return string(file)
}

func loadPost(title string) *Note {
	content := getContent(title)
	sepLength := len(sep)
    i := strings.LastIndex(content, sep)
    headers := content[sepLength:i]
    body := content[i+sepLength+1:]
    html := blackfriday.MarkdownCommon([]byte(body))
    meta := strings.Split(headers, "\n")
    return &Note{title, meta[1], meta[2], string(html)}
}

func loadTemplate(name string) string {
    file, _ := ioutil.ReadFile(name + ".html")
    return string(file)
}

func noteHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    title := vars["note"]
    title = strings.Replace(title, ".html", "", -1)
    var rendered string
	for _, note := range notes {
		if note.Uglyname == title {
			rendered = mustache.RenderInLayout(note.Body, loadTemplate("note"), nil)
		}
	}
	if len(rendered) == 0 {
		fmt.Fprintf(w, "404 page not found")
	}
	
    fmt.Fprintf(w, rendered)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	content := make(map[string][]map[string]string)
	for _, note := range notes {
		thisNote := map[string]string{
			"Title":note.Title,
			"url":note.Uglyname,
		}

		content["notes"] = append(content["notes"], thisNote)
	}
	rendered := mustache.RenderInLayout(homeMarkup, loadTemplate("home"), content)
	fmt.Fprintf(w, rendered)
}

func main() {
	r := mux.NewRouter()
    r.HandleFunc("/", indexHandler)
    r.HandleFunc(notePath+"{note}", noteHandler)
    http.Handle("/", r)
    http.ListenAndServe(":8080", nil)
}
