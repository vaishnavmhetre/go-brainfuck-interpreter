package base

type Tape struct {
	cells []byte
	head  uint
}

func NewTape(initialCapacity uint) *Tape {
	return &Tape{
		cells: make([]byte, initialCapacity),
		head:  0,
	}
}

func (tape *Tape) MoveRight() {
	if tape.head == uint(len(tape.cells)-1) {
		tape.cells = append(tape.cells, 0)
	}
	tape.head++
}

func (tape *Tape) MoveLeft() {
	if tape.head == 0 {
		tape.cells = append(make([]byte, 1), tape.cells...)
	} else {
		tape.head--
	}
}

func (tape *Tape) Increment() {
	tape.cells[tape.head]++
}

func (tape *Tape) Decrement() {
	tape.cells[tape.head]--
}

func (tape *Tape) Get() byte {
	return tape.cells[tape.head]
}

func (tape *Tape) Set(value byte) {
	tape.cells[tape.head] = value
}
