// Grass-Mud-Horse VM in Golang
// https://code.google.com/p/grass-mud-horse/wiki/A_Brife_To_GrassMudHorse_Language
// An accent of Whitespace language
// [Space] -> 草
// [Tab] -> 泥
// [LF] -> 马
package gmh

// The Grass-Mud-Horse Virtual Machine
type VM struct {
	stack *stack
	heap []int
	interrupter func() bool
}

// Grass-Mud-Horse Instruction Modification Parameter
const (
	STACK = "草"
	MATH = "泥草"
	HEAP = "泥泥"
	FLOW = "马"
	IO = "泥马"
)

// Stack Manipulation
const (
	PUSH = "草"
	DUPLICATE = "马草"
	COPY = "泥草"
	SWAP = "马泥"
	DISCARD = "马马"
	SLIDE = "泥马"
)

// Arithmetic
const (
	ADD = "草草"
	SUB = "草泥"
	MUL = "草马"
	DIV = "泥草"
	MOD = "泥泥"
)

// Heap
const (
	STORE = "草"
	RETRIEVE = "泥"
)

// Flow control
const (
	MARK = "草草"
	CALL = "草泥"
	JUMP = "草马"
	JUMP_ZERO = "泥草"
	JUMP_SUBZERO = "泥泥"
	RETURN = "泥马"
	END1 = "马马"
	END2 = "河蟹"
)

// I/O
const (
	PRINT_CHAR = "草草"
	PRINT_NUM = "草泥"
	SCAN_CHAR = "泥草"
	SCAN_NUM = "泥泥"
)

// Numbers
const (
	ZERO = "草"
	ONE = "泥"
	NUM_END = "马"
)
