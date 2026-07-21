package main

import (
	"testing"
	d "github.com/manuelsanchezto/horizontal-flow-gui/domain"
	"math/rand/v2"
)

func TestEvaluateDoNotFail (t *testing.T){
	context = d.Context{}
	context.Evaluate()
}


func TestAddAStep (t *testing.T){
	context = d.Context{}
	context.AddStep("step", 0)
	if len(context.Steps) != 1 {
		t.Errorf("Instead of 1, %d steps were found", len(context.Steps))
	}
	if context.ErrorText != "" {
		t.Errorf("The following error was found on the context %s", context.ErrorText)
	}
}

func TestAddManyStep (t *testing.T){
	context = d.Context{}
	numberOfIterations := rand.IntN(15)
	for i :=range numberOfIterations {
		context.AddStep("stepName", i)
	}
	if len(context.Steps) !=  numberOfIterations{
		t.Errorf("Instead of %d, %d steps were found", numberOfIterations, len(context.Steps))
	}
	if context.ErrorText != "" {
		t.Errorf("The following error was found on the context %s", context.ErrorText)
	}
}

func TestAddAVariable (t *testing.T){
	context = d.Context{}
	context.AddStep("step", 0)
	context.AddNewVariable(0,"test","test")
	if len(context.Variables) != 1 {
		t.Errorf("Instead of 1, %d steps were found", len(context.Variables))
	}
	if context.ErrorText != "" {
		t.Errorf("The following error was found on the context %s", context.ErrorText)
	}
}

func TestExpectFailureIfTheVariablePointsToInvalidStep (t *testing.T){
	context = d.Context{}
	context.AddNewVariable(0,"test","test")
	if len(context.Variables) != 0 {
		t.Error("The variable has been added without a valid Step")
	}
	if context.ErrorText != "Variable creation failed, review the input" {
		t.Errorf("The following error was found on the context %s instead of the expected %s", context.ErrorText, "Variable creation failed, review the input")
	}
}

func TestVariableSetNewValueOfStartingValue (t *testing.T){
	testedVal := "test value 2"
	context = d.Context{}
	context.AddStepAtTheEnd("test step")
	context.AddNewVariable(0, "test var", "Basic Value")
	if context.ErrorText != "" {
		t.Errorf("Expected no error before changing the starting value, found %s", context.ErrorText)
	}
	changedVal := context.ChangeValueOfVariable(0, 0, testedVal)
	if !changedVal {
		t.Error("The value change failed")	
	}
	if context.ErrorText != "" {
		t.Errorf("Expected no error, found %s", context.ErrorText)
	}
	if context.Variables[0].Values[0] != testedVal {
		t.Error("The new value has not been set")
	}
}


func TestVariableAddNewVariable (t *testing.T){
	testedVal := "test value 2"
	context = d.Context{}
	context.AddStepAtTheEnd("test step")
	context.AddStepAtTheEnd("test step 2")
	context.AddNewVariable(0, "test var", "BasicValue")
	changedVal := context.AddNewValueOnVariable(1, 0, testedVal)
	if !changedVal {
		t.Error("The value change failed")	
	}
	if context.ErrorText != "" {
		t.Errorf("Expected no error, found %s", context.ErrorText)
	}
	if context.Variables[0].Values[1] != testedVal {
		t.Error("The new variable value has not been set")
	}
	if context.Steps[1].VariableTriads[0] != [3]int {0,0,1} {
		t.Error("The new value of the variable has not been correclty linked to the Step")
	}

}

