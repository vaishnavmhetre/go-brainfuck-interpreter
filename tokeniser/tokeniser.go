package tokeniser

import (
	"bufio"
	"fmt"
	"os"
)

type Op uint

const (
	MoveRight Op = iota
	MoveLeft
	Increment
	Decrement
	Print
	Replace
	LParen
	RParen
)

type Loc struct {
	FileName string
	Line     uint
	Col      uint
}

func NewLoc(fileName string, line, col uint) Loc {
	return Loc{
		FileName: fileName,
		Line:     line,
		Col:      col,
	}
}

func (loc Loc) Disp() string {
	return fmt.Sprintf("%s:%d:%d", loc.FileName, loc.Line, loc.Col)
}

type Token struct {
	Op       Op
	Loc      Loc
	JumpAddr uint
}

func NewToken(op Op, fileName string, line, col uint, jumpAddr uint) (token *Token, err error) {
	return &Token{
		Op:       op,
		Loc:      NewLoc(fileName, line, col),
		JumpAddr: jumpAddr,
	}, nil
}

func Tokenise(file *os.File) (tokens []*Token, err error) { // NOSONAR - I wont fix the Cognitive issue.
	// tokens := []rune(program)

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve File Info to estimate buffer size to be allocated. Reported %w", err)
	}

	approximateBufferSize := float64(fileInfo.Size()) * 1.1

	tokens = make([]*Token, 0, int(approximateBufferSize))

	openLParensStack := make([]int, 0, int(float64(cap(tokens))*0.2)) // 20% could be [ barckets
	openLParensStackHead := -1

	scanner := bufio.NewScanner(file)

	var lineNum uint = 1
	var colNum uint = 1

	for scanner.Scan() {
		line := scanner.Text()

		// Loop through each character in the line
		for i, char := range line {
			_ = i
			var op Op
			skip := false
			switch char {
			case '>':
				op = MoveRight
			case '<':
				op = MoveLeft
			case '+':
				op = Increment
			case '-':
				op = Decrement
			case '.':
				op = Print
			case ',':
				op = Replace
			case '[':
				op = LParen
			case ']':
				op = RParen
			default:
				skip = true
			}

			if skip {
				continue
			}

			tokenPtr, err := NewToken(op, file.Name(), lineNum, colNum, 0)
			if err != nil {
				return nil, fmt.Errorf("%s:%d:%d: failed to parse Parenthsis Token `%c`. Reported %w", file.Name(), lineNum, colNum, char, err)
			}

			tokens = append(tokens, tokenPtr)

			if op == LParen {
				openLParensStackHead++
				if openLParensStackHead == len(openLParensStack) {
					openLParensStack = append(openLParensStack, len(tokens)-1)
				} else {
					openLParensStack[openLParensStackHead] = len(tokens) - 1
				}
			} else if op == RParen {
				if openLParensStackHead == -1 {
					return nil, fmt.Errorf("%s: no Left Paren `[` remained to match this Right Paren `]`", tokenPtr.Loc.Disp())
				}

				lParenIdx := openLParensStack[openLParensStackHead]
				openLParensStackHead--

				lTokenPtr := tokens[lParenIdx]

				if lTokenPtr.Op == LParen {
					// lParenTokenPtr := tmp.(*ParenToken)
					lTokenPtr.JumpAddr = uint(len(tokens) - 1)
					tokenPtr.JumpAddr = uint(lParenIdx)
				} else {
					return nil, fmt.Errorf("%s: found token %+v being pointed by index stored in Left Paren `[` Stack, while a search for a match for this Right Paren `]` was attempted. THIS IS DEFINITELY A BUG IN TOKENISER", tokenPtr.Loc.Disp(), *tokens[lParenIdx])
				}
			}
			colNum++ // Increment column number for each character
		}

		lineNum++  // Increment line number after processing the line
		colNum = 1 // Reset column number for the next line
	}

	if openLParensStackHead >= 0 {
		return nil, fmt.Errorf("some open Left Parens are still left in stack, no Right Parens found for them. Remaining: %v", openLParensStack[:openLParensStackHead+1])
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error scanning file: %w", err)
	}

	return
}
