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
