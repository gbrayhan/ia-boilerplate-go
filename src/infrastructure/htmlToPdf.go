package infrastructure

import (
	"bytes"
	"errors"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"text/template"
)

var FuncMapTemplates = template.FuncMap{
	"derefInt": DerefInt,
}

func RouteTemplateToPDF(routeTemplate string, data interface{}) (pdfContent string, err error) {
	wkhtmltopdfBin := os.Getenv("WKHTMLTOPDF_BIN")
	if wkhtmltopdfBin == "" {
		err = errors.New("environment variable WKHTMLTOPDF_BIN is not defined")
		return
	}
	filenamePDF := "archives/tmp/" + strconv.Itoa(rand.Int()) + "_file.pdf"
	filenameHTML := "archives/tmp/" + strconv.Itoa(rand.Int()) + "_file.html"
	file, err := os.Create(filenameHTML)
	if err != nil {
		return
	}
	htmlTemplate := processFile(routeTemplate, data)
	if _, err = file.WriteString(htmlTemplate); err != nil {
		return
	}
	if err = file.Close(); err != nil {
		return
	}
	args := []string{"-T", "0", "-B", "0", "-L", "0", "-R", "0", "-s", "Letter", "-O", "Portrait",
		"--javascript-delay", "5000",
		"--no-stop-slow-scripts",
		filenameHTML, filenamePDF}
	cmd := exec.Command(wkhtmltopdfBin, args...)
	_, err = cmd.CombinedOutput()
	if err != nil {
		return
	}
	content, err := ioutil.ReadFile(filenamePDF)
	if err != nil {
		return
	}
	if err = os.Remove(filenamePDF); err != nil {
		return
	}
	if err = os.Remove(filenameHTML); err != nil {
		return
	}
	pdfContent = string(content)
	return
}

func process(t *template.Template, vars interface{}) string {
	var tmplBytes bytes.Buffer
	err := t.Execute(&tmplBytes, vars)
	if err != nil {
		panic(err)
	}
	return tmplBytes.String()
}

func processFile(filename string, vars interface{}) string {

	tmpl, err := template.
		New(filename).
		Funcs(FuncMapTemplates).ParseFiles("templates/" + filename)
	if err != nil {
		panic(err)
	}
	return process(tmpl, vars)
}

func DerefInt(i *int) int {
	if i == nil {
		return 0
	}
	return *i
}
