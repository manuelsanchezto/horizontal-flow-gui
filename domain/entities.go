package domain

import (
	"errors"
	"slices"
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

func (c *Context) AddStepAtTheEnd (stepName string) bool {
	c.Steps = append(c.Steps, Step {Name: stepName, Pos: len(c.Steps)})
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
	c.Steps[stepPos] = step 
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

func (c *Context) ChangeValueOfVariable(vIndex int, vValueIndex int, vValue string) bool {
	if  vIndex >= len(c.Variables) || vValueIndex >= len(c.Variables[vIndex].Values) {
		c.ErrorText = "Either the variable does not exists or the value is not yet set"
		return false
	}
	c.Variables[vIndex].Values[vValueIndex] = vValue
	return true
}

func (c *Context) AddNewValueOnVariable (sIndex int, vIndex int, vValue string) bool{
	// TODO: Is missing a check that for that step there is not already a value for that Variable
	if sIndex >= len(c.Steps) && vIndex >= len(c.Variables) {
		c.ErrorText = "The value set parameters are not valid"
		return false
	}
	c.Steps[sIndex].VariableTriads = append( c.Steps[sIndex].VariableTriads, 
		[3]int {
			vIndex, 
			c.Steps[sIndex].NOfRows, 
			len(c.Variables[vIndex].Values),
		})
	c.Variables[vIndex].Values = append(c.Variables[vIndex].Values, vValue)
	return true
}

func (c *Context) DeleteVariable(variableIndex int) bool {
	if variableIndex > len(c.Variables) {
		c.ErrorText = "Could not find the Variable to delete"
		return false
	}
	for i  := range c.Steps {
		for j, triad:= range c.Steps[i].VariableTriads {
			if triad[0] == variableIndex {
				c.Steps[i].VariableTriads = slices.Delete(c.Steps[i].VariableTriads, j, j+1)
			}
		}
	}
	c.Variables = slices.Delete(c.Variables, variableIndex, variableIndex+1)
	return true
}

