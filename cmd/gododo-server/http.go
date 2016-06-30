package main

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/russross/blackfriday"
)

type DocHandler struct {
	cfg Config
	tpl *template.Template
}

func mkHandler(cfg Config) *DocHandler {
	tplfiles := []string{"page.html", "error.html"}
	tplpaths := make([]string, len(tplfiles), len(tplfiles))
	for i := range tplpaths {
		tplpaths[i] = filepath.Join(cfg.TemplateRoot, tplfiles[i])
	}
	tpl, err := template.New("template").ParseFiles(tplpaths...)
	assert(err)

	return &DocHandler{
		cfg: cfg,
		tpl: tpl,
	}
}

type PageData struct {
	Title string
	HTML  template.HTML
}

func (d *DocHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// Convert from HTTP path to filesystem path
	cleanpath := filepath.Clean(req.URL.Path)
	fspath := filepath.Join(d.cfg.DocumentRoot, cleanpath)

	if filepath.Base(fspath) == "" {
		http.Redirect(w, req, req.RequestURI+"home", http.StatusMovedPermanently)
		return
	}

	stat, err := os.Stat(fspath)
	if err != nil {
		if os.IsNotExist(err) {
			d.Error(w, http.StatusNotFound, "File not found")
			return
		}
		log.Println(err.Error())
		d.Error(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	if stat.IsDir() {
		http.Redirect(w, req, req.RequestURI+"/home", http.StatusMovedPermanently)
		return
	}

	data, err := ioutil.ReadFile(fspath)
	if err != nil {
		log.Println(err.Error())
		d.Error(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	html := blackfriday.MarkdownCommon(data)
	err = d.tpl.ExecuteTemplate(w, "page.html", PageData{
		Title: getTitle(req.URL.Path),
		HTML:  template.HTML(html),
	})
	if err != nil {
		log.Println(err.Error())
		d.Error(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
}

func (d *DocHandler) Error(w http.ResponseWriter, errcode int, errmsg string) {
	errno := strconv.Itoa(errcode)
	errfancy := strings.Replace(errno, "0", "o", -1)

	w.WriteHeader(errcode)
	d.tpl.ExecuteTemplate(w, "error.html", struct {
		ErrorNo      string
		ErrorFancyNo string
		ErrorMsg     string
	}{
		errno,
		errfancy,
		strings.ToLower(errmsg),
	})
}

func getTitle(path string) string {
	// Extract basename
	parts := strings.Split(path, "/")
	pagename := parts[len(parts)-1]

	// Decode URL encoded glyphs
	if decoded, err := url.QueryUnescape(pagename); err == nil {
		pagename = decoded
	}

	// Replace _ with whitespace
	pagename = strings.Replace(pagename, "_", " ", -1)

	// Capitalize
	if len(pagename) > 1 {
		pagename = strings.ToUpper(pagename[0:1]) + pagename[1:]
	}

	return pagename
}
