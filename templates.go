package kinli

import (
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"strings"
)

type simpleTemplate struct {
	prefix      string
	partialsDir string
	t           *template.Template
}

var (
	// PathTemplate is the template/ dir . default is tmpl/
	PathTemplate = "tmpl/"
	// PathPartialTemplate is the relative path for partials/ directory. default is tmpl/partials/
	PathPartialTemplate = "tmpl/partials/"
	// CacheMode is true to cache templates. set to false during development
	CacheMode = true
	// ViewFuncs supports extra functions which user can use in views
	ViewFuncs template.FuncMap
	// templates is the internal reference to keep track of parsed files
	templates *simpleTemplate
)

func init() {
	// the default directories are expected to be present
	// InitTmpl()
}

// InitTmpl should be called after any of the Template variables are modified
func InitTmpl() {
	loadTemplates()
}

// InitTmpl should be called to for set templates context wide
// loadTemplates is initialised once on init
// when CacheMode is set to false, it is called on every page load
func loadTemplates() {
	// Templates with functions available to them
	templates = &simpleTemplate{
		PathTemplate,
		PathPartialTemplate,
		template.New("").Funcs(internalFuncMap).Funcs(ViewFuncs),
	}
	load()
}

// DisplayText takes in the context for creating a NewPage and displays the text as page
// Uses "_text.html" file in your template directory
func DisplayText(hc *HttpContext, w io.Writer, text string) {
	page := NewPage(hc, text, "Action Required!", text, nil)
	DisplayPage(w, "_text", page)
}

// DisplayPage exposes writer
// path is the name of the file without ".html"
// page can be the structure that you can access in the view.
// Page struct is recommended or inherit the Page struct
func DisplayPage(w io.Writer, path string, page interface{}) {
	if !CacheMode {
		loadTemplates()
	}

	tv := templates.t.Lookup(templates.prefix + path)
	tv.Execute(w, page)
}

// GetPage gets the raw page content
func GetPage(path string, page interface{}) io.Reader {
	pr, pw := io.Pipe()
	go func() {
		defer pw.Close()
		DisplayPage(pw, path, page)
	}()
	return pr
}

// GetPageContent gets the page content as string
func GetPageContent(path string, page interface{}) (string, error) {
	t, err := ioutil.ReadAll(GetPage(path, page))
	if err != nil {
		return "", err
	}
	return string(t), nil
}

func load() {
	// load templates
	loadTmplFromDir(templates.prefix)
	// load partials
	loadTmplFromDir(templates.partialsDir)
}

func loadTmplFromDir(dirpath string) {
	files, err := ioutil.ReadDir(dirpath)
	if err != nil {
		log.Print(err)
		return
	}

	for _, file := range files {
		// ignore directories
		if file.IsDir() {
			continue
		}
		name := dirpath + file.Name()
		tmplname := strings.Replace(file.Name(), ".html", "", 1)

		b, err := ioutil.ReadFile(name)
		_, err = templates.t.New(dirpath + tmplname).Parse(string(b))

		if err != nil {
			log.Print(err)
		}
	}
}
