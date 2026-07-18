package domain

import (
	"errors"
)

type Variable struct {
	Name           string
	Values         []string
	IsVisible      bool
}

type Step struct {
	Name            string
	Pos             int
	NOfRows int
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
	ErrorText string
}

type Result struct {
	Variable Variable
	Value string
	Ok bool
}

func (c *Context) AddStep(stepName string, position int) bool {
	if stepName == "" || position < len(c.Steps) {
		c.ErrorText = "Step creation failed, either no name or invalid position"
		return false
	}
	c.Steps = append(c.Steps, Step {Name: stepName, Pos: position})
	return true
}

func (c *Context) AddNewVariable(stepPos int, vName string, vValue string) bool{
	if stepPos >= len(c.Steps) || stepPos < 0 || vName == "" || vValue == "" {
		c.ErrorText = "Variable creation failed, review the input"
		return false
	}
	vExists := c.FindVariableByName(vName)
	if vExists {
		c.ErrorText = "Variable already exists"
		return false
	}
	step := c.Steps[stepPos]
	step.NOfRows = step.NOfRows + 1
	step.VariableTriads = append(step.VariableTriads, [3]int {len(c.Variables),step.NOfRows,0})
	c.Variables = append(c.Variables, Variable{Name: vName, Values: []string{vValue}, IsVisible: true})
	return true
}

func (c *Context) FindVariableByName (vName string) bool {
	for i := range c.Variables {
		if c.Variables[i].Name == vName {
			return true
		}
	}
	return false
}

func (c *Context) FindStepByPos(pos int) (*Step, error) {
	for i := range c.Steps {
		if c.Steps[i].Pos == pos {
			return &c.Steps[i], nil
		}
	}
	return &Step{}, errors.New("Step not found")
}

func (c *Context) Evaluate () {
	aliveVars := 0
	for i := range c.Steps {
		aliveVars = aliveVars + len(c.Steps[i].VariableTriads)
		c.Steps[i].NOfRows = aliveVars
	}
	c.ErrorText = ""

}



