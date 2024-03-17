package main

import (
	"brainfuck-interpreter/interpreter"
	"brainfuck-interpreter/tokeniser"
	"fmt"
	"os"
)

func main() {

	for _, fileName := range os.Args[1:] {

		file, err := os.Open(fileName)

		if err != nil {
			panic(fmt.Errorf("failed to open file %s. Reported %w", fileName, err))
		}

		tokens, err := tokeniser.Tokenise(file)
		_ = tokens

		if err != nil {
			panic(fmt.Errorf("failed to parse tokens. Reported %w", err))
		}

		if err := interpreter.Interpret(tokens); err != nil {
			panic(fmt.Errorf("failed to execute program: Reported %w", err))
		}

	}

}
