package main

import (
	"testing"
)

type TestCleanInputs struct {
	text string
	want string
}

func verifyResults(t testing.TB,input,want,got string){
	t.Helper()
	if got != want {
		t.Errorf("testinput: %v\n want: %v\n got: %v",input,want,got)
	}
}

func TestCleanInput(t *testing.T) {
	t.Run("clean singe input",func(t *testing.T) {
		testInput := " Test This  "
		
		got :=cleanInput(testInput)
		want := "test this"

		verifyResults(t,testInput,want,got)
	})

	t.Run("Clean Multiple inputs",func(t *testing.T) {
		testInputs := []TestCleanInputs{{text: " AnothER One",want: "another one"},{text: " a SiLLY OnE",want: "a silly one"},{text: "    ",want: ""},{text: " vVvVvV   ",want: "vvvvvv"}}

		for _,input := range testInputs {
			verifyResults(t,input.text,input.want,cleanInput(input.text))
		}
	})
}