package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type Page struct {
	Title   string
	Body    []byte
	Files   []string
	Dirs    []string
	CurrDir string
}

func serveImage(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Path[len("/image"):]
	//body := r.FormValue("fn")
	log.Printf("Filename: %s", filename)
	file, err := os.Open(filename)
	if err != nil {
		log.Printf("Error: %s", err)
	}
	file.Close()
	//image, err := ioutil.ReadFile(filename)
	//if err != nil {
	//fmt.Fprintf(w, "<h1>%s</h1>", err)
	//}
	http.ServeFile(w, r, filename)
}

func loadPage(title, directory string) (*Page, error) {
	searchDir := directory //"/home/syn/Dropbox/Dev/2017/ImageLight"
	//searchDir := "/home/syn/images/test11"
	fileList := []string{}
	dirList := []string{}
	dirList = append(dirList, filepath.Dir(directory))
	err := filepath.Walk(searchDir, func(path string, f os.FileInfo, err error) error {
		ext := filepath.Ext(path)
		if f.IsDir() {
			dirList = append(dirList, path)
			//log.Printf("Directory: %s", path)
			//log.Printf("Dirdir: %s", filepath.Dir(path))
		}
		if ext == ".png" || ext == ".jpg" || ext == ".jpeg" || ext == ".gif" || ext == ".mp4" {
			fileList = append(fileList, path)
		}
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}

	filename := title + ".html"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body, Files: fileList, Dirs: dirList}, nil
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	directory := r.URL.Path
	log.Printf(directory)
	p, err := loadPage("index", directory)
	if err != nil {
		p = &Page{Title: "index"}
	}
	t, _ := template.ParseFiles("index.html")
	t.Execute(w, p)
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello %s!", r.URL.Path[1:])
}

// http://www.alexedwards.net/blog/serving-static-sites-with-go
func main() {

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/image/", serveImage)

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	//http.HandleFunc()

	log.Println("Listening...")

	http.ListenAndServe(":8080", nil)
}
