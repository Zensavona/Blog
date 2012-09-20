package main

import (
	"fmt"
	"path"
	"github.com/russross/blackfriday"
	"github.com/hoisie/mustache"
	"io/ioutil"
	"strings"
	"net/http"
)

type Post struct {
	Title string
	Date string
	Body string
}

const sep = "---"

func getContent(title string) (content string) {
	filename := title + ".md"
    file, _ := ioutil.ReadFile(path.Join("notes", filename))
    return string(file)
}

func loadPost(title string) *Post {
	content := getContent(title)

	sepLength := len(sep)
    i := strings.LastIndex(content, sep)
    headers := content[sepLength:i]
    body := content[i+sepLength+1:]
    html := blackfriday.MarkdownCommon([]byte(body))
    meta := strings.Split(headers, "\n")
    return &Post{meta[1], meta[2], string(html)}
}

func loadTemplate(name string) string {
    file, _ := ioutil.ReadFile(name + ".html.moustache")
    return string(file)
}

func noteHandler(w http.ResponseWriter, r *http.Request) {
    //request := "why-you-should-never-use-godaddy-ever-again"
	note := loadPost(r.URL.Path[1:])
	rendered := mustache.RenderInLayout(note.Body, loadTemplate("note"), nil)
    fmt.Fprintf(w, rendered)
}

func main() {
	http.HandleFunc("/", noteHandler)
    http.ListenAndServe(":8080", nil)
}