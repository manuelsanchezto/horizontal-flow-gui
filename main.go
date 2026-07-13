package main

import (
	"errors"
	"html/template"
	"net/http"
	"strconv"
)

type Context struct {
	TableName string
	Steps     []Step
}

type Step struct {
	Name      string
	Pos       int
	Variables []Variable
}

type Variable struct {
	Name           string
	PresentedValue string
	NumericValue   float32
}

func (c *Context) findStepByPos(pos int) (*Step, error) {
	for i := range c.Steps {
		if c.Steps[i].Pos == pos {
			return &c.Steps[i], nil
		}
	}
	return &Step{}, errors.New("Step not found")
}

func handler(w http.ResponseWriter, r *http.Request) {
	fillIndex(w)
}

func stepAddHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	stepName := r.Form.Get("stepName")
	if stepName == "" {
		stepName = "Default Step"
	}
	context.Steps = append(context.Steps, Step{Name: stepName, Pos: len(context.Steps) + 1})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func valFromStepAddHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	stepPosition := r.Form.Get("step")
	stepPositionParsed, err := strconv.Atoi(stepPosition)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	step, err := context.findStepByPos(stepPositionParsed)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	variable := Variable{
		Name:           "a",
		PresentedValue: "b",
		NumericValue:   0}
	step.Variables = append(step.Variables, variable)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func fillIndex(w http.ResponseWriter) {
	template, error := template.ParseFiles("index.html")
	if error != nil {
		http.Error(w, error.Error(), http.StatusInternalServerError)
		return
	}
	template.Execute(w, context)
}

var context Context

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/add", stepAddHandler)
	http.HandleFunc("/add-val", valFromStepAddHandler)
	http.ListenAndServe(":8090", nil)
}
