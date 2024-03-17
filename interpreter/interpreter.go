package interpreter

import (
	"brainfuck-interpreter/base"
	"brainfuck-interpreter/tokeniser"
	"fmt"
)

func Interpret(tokens []*tokeniser.Token) error { // NOSONAR - I wont fix the Cognitive issue.

	tokensLength := uint(len(tokens))

	tape := base.NewTape(tokensLength)

	var pc uint = 0
	for pc < tokensLength {
		tokenPtr := tokens[pc]
		op := tokenPtr.Op
		switch op {
		case tokeniser.MoveRight:
			tape.MoveRight()
		case tokeniser.MoveLeft:
			tape.MoveLeft()
		case tokeniser.Increment:
			tape.Increment()
		case tokeniser.Decrement:
			tape.Decrement()
		case tokeniser.Print:
			fmt.Printf("%c", tape.Get())
		case tokeniser.Replace:
			var input uint8
			if _, err := fmt.Scanf("%d", &input); err != nil {
				return fmt.Errorf("%s: failed to read input, halting execution", tokenPtr.Loc.Disp())
			}
			tape.Set(input)
		case tokeniser.LParen:
			if tape.Get() == 0 {
				pc = tokenPtr.JumpAddr
			}
		case tokeniser.RParen:
			if tape.Get() != 0 {
				pc = tokenPtr.JumpAddr
			}
		default:
			return fmt.Errorf("%s: unexpected token found: %v with operator %v. Halting Program Now. THIS IS DEFINITELY A BUG IN INTERPRETER", tokenPtr.Loc.Disp(), *tokenPtr, op)
		}

		if pc >= tokensLength {
			return fmt.Errorf("%s program counter has gone beyond actual count of tokens in program. Total tokens in Program: %d, PC is at %d. THIS IS DEFINTELY A BUG IN INTERPRETER", tokenPtr.Loc.Disp(), tokensLength, pc)
		}
		pc++
	}

	return nil

}
