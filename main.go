package main

import (
	"errors"
	"html/template"
	"net/http"
	"strconv"
)

type Variable struct {
	Name           string
	Values         []string
	IsVisible      bool
}

type Step struct {
	Name            string
	Pos             int
	NumberOfColumns int
	VariableTriads  [][3]int
	/*
The block of triads is as follows:
- [0] represents the index of the variable on the Context slice
- [1] represents the row that occupies on the table
- [2] represents the value that will be placed
*/
}

type Context struct {
	TableName string
	Steps     []Step
	Variables []Variable
}



func (c *Context) findStepByPos(pos int) (*Step, error) {
	for i := range c.Steps {
		if c.Steps[i].Pos == pos {
			return &c.Steps[i], nil
		}
	}
	return &Step{}, errors.New("Step not found")
}

func (c *Context) evaluate () {
	aliveVars := 0
	for i := range c.Steps {
		aliveVars = aliveVars + len(c.Steps[i].VariableTriads)
		c.Steps[i].NumberOfColumns = aliveVars
	}

}

func (s *Step) addVariableTriad (variable int, row int, value int) {
	s.VariableTriads = append(s.VariableTriads, [3]int{variable,row,value})
}

func (s *Step) addNewVariable (variable int) {
	s.addVariableTriad(variable, 1,0)
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
	context.evaluate()
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
	index := step.NumberOfColumns
	step.addNewVariable(index)
	variableName := r.Form.Get("variableName")
	variableValue := r.Form.Get("variableValue")
	variable := Variable{Name:variableName, Values: []string{}, IsVisible: true}
	variable.Values = append(variable.Values, variableValue)
	context.Variables = append(context.Variables, variable)

	context.evaluate()
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func rangeFunc(n int) []int {
	r:= make([]int ,n)
	for i := range n{
		r[i] = i
	}
	return r
}

func fillIndex(w http.ResponseWriter) {
	funcMap := template.FuncMap{"ranger": rangeFunc,}
	template, error := template.New("index.html").Funcs(funcMap).ParseFiles("index.html")
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
