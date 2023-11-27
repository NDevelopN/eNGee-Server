package utils

import "fmt"

func ValidateInputRefuseEmpty(input string, allowed map[string]struct{}) error {
	if input == "" {
		return fmt.Errorf("input is empty")
	}

	return ValidateInput(input, allowed)
}

func ValidateInput(input string, allowed map[string]struct{}) error {
	if len(allowed) == 0 {
		return nil
	}

	_, contains := allowed[input]
	if contains {
		return nil
	}

	return fmt.Errorf("%q is not a valid input", input)
}
