package main

import (
	d "github.com/manuelsanchezto/horizontal-flow-gui/domain"
	"html/template"
	"net/http"
	"strconv"
)

func fillIndex(w http.ResponseWriter) {
	funcMap := template.FuncMap{"iterate": rangeFunc, "findVar":findVariableByColumn}
	template, error := template.New("index.html").Funcs(funcMap).ParseFiles("index.html")
	if error != nil {
		http.Error(w, error.Error(), http.StatusInternalServerError)
		return
	}
	template.Execute(w, context)
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
	Ok := context.AddStepAtTheEnd(stepName)
	if Ok {
		context.Evaluate()
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func valueAddHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	stepIndex, err := strconv.Atoi(r.Form.Get("step"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	variableName := r.Form.Get("variableName")
	variableValue := r.Form.Get("variableValue")

	Ok := context.AddNewVariable(stepIndex, variableName, variableValue)
	if Ok {
		context.Evaluate()
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func changeValueHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	variableIndex, err := strconv.Atoi(r.Form.Get("variableIndex"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	valueIndex, errVI := strconv.Atoi(r.Form.Get("valueIndex"))
	if errVI != nil {
		http.Error(w, errVI.Error(), http.StatusInternalServerError)
	}
	newValue := r.Form.Get("valueIndex")
	Ok := context.ChangeValueOfVariable(variableIndex, valueIndex, newValue)
	if Ok {
		context.Evaluate()
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}


func addValueHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	stepIndex, err := strconv.Atoi(r.Form.Get("step"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	variableIndex, errVI := strconv.Atoi(r.Form.Get("variableIndex"))
	if errVI != nil {
		http.Error(w, errVI.Error(), http.StatusInternalServerError)
	}
	variableValue := r.Form.Get("variableValue")
	Ok := context.AddNewValueOnVariable(stepIndex, variableIndex, variableValue)
	if Ok {
		context.Evaluate()
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func rangeFunc(n int) []int {
	r:= make([]int ,n)
	for i := range n{
		r[i] = i
	}
	return r
}


func findVariableByColumn(s d.Step, i int) (d.Result){
	for _, triad:= range s.VariableTriads{
		if i == triad[1]{
			return d.Result{
				Variable:context.Variables[triad[0]],
				Value:context.Variables[triad[0]].Values[triad[2]],
				Ok:true,
			}
		}
	}
	return d.Result{}
}

var context d.Context

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/add-step", stepAddHandler)
	http.HandleFunc("/add-variable", valueAddHandler)
	http.HandleFunc("/change-value", changeValueHandler)
	http.HandleFunc("/add-value", addValueHandler)
	http.ListenAndServe(":8090", nil)
}
