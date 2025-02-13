//go:build !solution

package main

import (
	"fmt"
	"strconv"
	"strings"
)

type Evaluator struct {
	stack       []int
	defenitions map[string]string
}

// NewEvaluator creates evaluator.
func NewEvaluator() *Evaluator {
	return &Evaluator{
		stack:       []int{},
		defenitions: map[string]string{},
	}
}

// Process evaluates sequence of words or definition.
//
// Returns resulting stack state and an error.
func (e *Evaluator) Process(row string) ([]int, error) {
	args := strings.Split(row, " ")

	if args[0] == ":" && args[len(args)-1] == ";" {
		err := e.defineWord(args)

		if err != nil {
			return e.stack, err
		}

		return e.stack, nil
	}

	for i := 0; i < len(args); i++ {
		var action string
		if valueInDefenition, ok := e.defenitions[strings.ToLower(args[i])]; ok {
			action = valueInDefenition
		} else {
			action = strings.ToLower(args[i])
		}

		actionArgs := strings.Split(action, " ")

		for j := 0; j < len(actionArgs); j++ {
			if number, err := strconv.Atoi(actionArgs[j]); err == nil {
				e.stack = append(e.stack, number)
				continue
			}

			if len(e.stack) == 0 {
				return e.stack, fmt.Errorf("not enough args")
			}

			if actionArgs[j] == "dup" {
				e.stack = append(e.stack, e.stack[len(e.stack)-1])
				continue
			}

			if actionArgs[j] == "drop" {
				e.stack = e.stack[:len(e.stack)-1]
				continue
			}

			if len(e.stack) == 1 {
				return e.stack, fmt.Errorf("not enough args")
			}

			firstNumber := e.stack[len(e.stack)-1]
			secondNumber := e.stack[len(e.stack)-2]
			e.stack = e.stack[:len(e.stack)-2]

			if actionArgs[j] == "-" {
				e.stack = append(e.stack, secondNumber-firstNumber)
				continue
			}

			if actionArgs[j] == "+" {
				e.stack = append(e.stack, secondNumber+firstNumber)
				continue
			}

			if actionArgs[j] == "/" {

				if firstNumber == 0 {
					return e.stack, fmt.Errorf("division by zero")
				}

				e.stack = append(e.stack, secondNumber/firstNumber)
				continue
			}

			if actionArgs[j] == "*" {
				e.stack = append(e.stack, secondNumber*firstNumber)
				continue
			}

			if actionArgs[j] == "over" {
				e.stack = append(e.stack, secondNumber)
				e.stack = append(e.stack, firstNumber)
				e.stack = append(e.stack, secondNumber)
				continue
			}

			if actionArgs[j] == "swap" {
				e.stack = append(e.stack, firstNumber)
				e.stack = append(e.stack, secondNumber)
				continue
			}

			return e.stack, fmt.Errorf("non-existent word")
		}

	}

	return e.stack, nil
}

func (e *Evaluator) defineWord(args []string) error {
	if len(args) < 4 {
		return fmt.Errorf("invalid definition")
	}

	if _, err := strconv.Atoi(args[1]); err == nil {
		return fmt.Errorf("number definition")
	}

	newDefenitions := []string{}

	for j := 2; j < len(args)-1; j++ {
		loweredArg := strings.ToLower(args[j])

		if valueInDefenition, ok := e.defenitions[loweredArg]; ok {
			newDefenitions = append(newDefenitions, valueInDefenition)
		} else {
			newDefenitions = append(newDefenitions, loweredArg)
		}
	}

	e.defenitions[strings.ToLower(args[1])] = strings.Join(newDefenitions, " ")

	return nil
}
