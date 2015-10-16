package gmh

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Create a new instance of the VM
func New() *VM {
	return &VM{
		stack: newStack(),
		heap:  make([]int, 0),
		buf:   "",
	}
}

// The interrupter function will be called in the main loop
// If the function returns true, the execution will be terminated immediately.
func (this *VM) SetInterrupter(f func() bool) *VM {
	this.interrupter = f
	return this
}

// The input function will be called when input buffer runs out
// The function should return a string for input buffer.
func (this *VM) SetInput(f func(string) string) *VM {
	this.input = f
	return this
}

func (this *VM) readInput(res *string) (r byte) {
	if this.input == nil {
		r = ' '
		return
	}

	if len(this.buf) == 0 {
		this.buf = this.input(*res)
		(*res) = ""
	}

	if len(this.buf) > 0 {
		r = this.buf[0]
		this.buf = this.buf[1:]
	} else {
		r = ' '
	}

	return
}

// Execute commands
func (this *VM) Exec(cmds []string) (res string, err error) {
	labels := make(map[int]int)
	calls := newStack()

	// Scan for labels
	for i := 0; i < len(cmds); i++ {
		cmd := cmds[i]
		if strings.HasPrefix(cmd, FLOW) {
			c := cmd[len(FLOW):]

			if strings.HasPrefix(c, MARK) {
				var num int
				num, err = parseNum(c[len(MARK):], false)

				if err == nil {
					labels[num] = i
				} else {
					return
				}
			}
		}
	}

	for i := 0; i < len(cmds); i++ {
		if (this.interrupter != nil) && this.interrupter() {
			err = errors.New("Interrupted.")
			return
		}

		cmd := cmds[i]
		switch {
		case strings.HasPrefix(cmd, STACK):
			err = this.opStack(cmd[len(STACK):])
		case strings.HasPrefix(cmd, MATH):
			err = this.opMath(cmd[len(MATH):])
		case strings.HasPrefix(cmd, HEAP):
			err = this.opHeap(cmd[len(HEAP):])
		case strings.HasPrefix(cmd, IO):
			// I/O control, do it here
			c := cmd[len(IO):]

			switch {
			case strings.HasPrefix(c, PRINT_CHAR):
				m := this.stack.pop()

				if m != nil {
					res += string(m.(int))
				}
			case strings.HasPrefix(c, PRINT_NUM):
				m := this.stack.pop()

				if m != nil {
					res += fmt.Sprintf("%d", m.(int))
				}
			case strings.HasPrefix(c, SCAN_CHAR):
				m := this.stack.pop()

				if m != nil {
					r := this.readInput(&res)

					if this.checkHeap(m.(int), true) {
						this.heap[m.(int)] = int(r)
					}
				}
			case strings.HasPrefix(c, SCAN_NUM):
				m := this.stack.pop()

				if m != nil {
					r := this.readInput(&res)

					var n int64
					n, err = strconv.ParseInt(string(r), 10, 32)
					if this.checkHeap(m.(int), true) {
						this.heap[m.(int)] = int(n)
					}
				}
			}
		case strings.HasPrefix(cmd, FLOW):
			// Flow control, do it here
			c := cmd[len(FLOW):]

			switch {
			case strings.HasPrefix(c, END1):
				return
			//case strings.HasPrefix(c, MARK):
			// Mark a point
			//	var num int
			//	num, err = parseNum(c[len(MARK):], false)

			//	if err == nil {
			//		labels[num] = i
			//	}
			case strings.HasPrefix(c, CALL):
				// Call a label as a function
				var num int
				num, err = parseNum(c[len(CALL):], false)

				if err == nil {
					calls.push(i)

					if val, ok := labels[num]; ok {
						i = val
					} else {
						err = errors.New(fmt.Sprintf("Label %d not found.", num))
					}
				}
			case strings.HasPrefix(c, RETURN):
				m := calls.pop()

				if m == nil {
					err = errors.New("Cannot return to caller.")
					return
				}

				i = m.(int)
			case strings.HasPrefix(c, JUMP):
				// Jump to a label
				var num int
				num, err = parseNum(c[len(JUMP):], false)

				if err == nil {
					if val, ok := labels[num]; ok {
						i = val
					}
				}
			case strings.HasPrefix(c, JUMP_ZERO):
				// Jump to a label if top of stack is 0
				m := this.stack.pop()

				if m == nil {
					err = errors.New("Stack has no elements.")
					return
				}

				if m.(int) == 0 {
					var num int
					num, err = parseNum(c[len(JUMP_ZERO):], false)

					if err == nil {
						if val, ok := labels[num]; ok {
							i = val
						}
					}
				}
			case strings.HasPrefix(c, JUMP_SUBZERO):
				// Jump to a label if top of stack is negative
				m := this.stack.pop()

				if m == nil {
					err = errors.New("Stack has no elements.")
					return
				}

				if m.(int) < 0 {
					var num int
					num, err = parseNum(c[len(JUMP_SUBZERO):], false)

					if err == nil {
						if val, ok := labels[num]; ok {
							i = val
						}
					}
				}

			}
		case cmd == END2:
			return
		}

		//fmt.Println(i)

		if err != nil {
			break
		}
	}

	return
}

// Stack operation
func (this *VM) opStack(cmd string) (err error) {
	switch {
	case strings.HasPrefix(cmd, PUSH):
		// Push into stack
		var num int
		num, err = parseNum(cmd[len(PUSH):], true)

		if err == nil {
			this.stack.push(num)
		}
	case strings.HasPrefix(cmd, DUPLICATE):
		// Duplicate the top item
		m := this.stack.pop()

		// Double-push
		this.stack.push(m).push(m)
	case strings.HasPrefix(cmd, COPY):
		// Copy an item onto the top
		var num int
		num, err = parseNum(cmd[len(COPY):], true)

		if err == nil {
			tmp := newStack()
			for i := 0; i < num-1; i++ {
				tmp.push(this.stack.pop())
			}

			m := this.stack.pop()

			this.stack.push(m)

			for {
				t := tmp.pop()

				if t == nil {
					break
				}

				this.stack.push(t)
			}

			this.stack.push(m)
		}
	case strings.HasPrefix(cmd, SWAP):
		// Swap the top two items
		m := this.stack.pop()
		n := this.stack.pop()

		this.stack.push(m)
		this.stack.push(n)
	case strings.HasPrefix(cmd, DISCARD):
		// Discard the top
		this.stack.pop() // Throw it to GC
	case strings.HasPrefix(cmd, SLIDE):
		// Slide n items from the top but keep the top
		var num int
		num, err = parseNum(cmd[len(SLIDE):], true)

		if err == nil {
			m := this.stack.pop()

			for i := 0; i < num; i++ {
				this.stack.pop()
			}

			this.stack.push(m)
		}
	}

	return
}

// Math operation
func (this *VM) opMath(cmd string) (err error) {
	m := this.stack.pop()
	n := this.stack.pop()

	if (m == nil) || (n == nil) {
		err = errors.New("Can't evaluate when arguments are null")
		return
	}

	x := m.(int)
	y := n.(int)

	switch {
	case strings.HasPrefix(cmd, ADD):
		this.stack.push(y + x)
		return
	case strings.HasPrefix(cmd, SUB):
		this.stack.push(y - x)
		return
	case strings.HasPrefix(cmd, DIV):
		if x == 0 {
			err = errors.New("Cannot devide 0")
			return
		}
		this.stack.push(y / x)
		return
	case strings.HasPrefix(cmd, MUL):
		this.stack.push(y * x)
		return
	case strings.HasPrefix(cmd, MOD):
		this.stack.push(y % x)
		return
	}

	// If arrived here, we've got an invalid operation
	err = errors.New("Invalid operation" + cmd)
	return
}

// Heap
func (this *VM) opHeap(cmd string) (err error) {
	switch {
	case strings.HasPrefix(cmd, STORE):
		m := this.stack.pop()
		n := this.stack.pop()

		if (m == nil) || (n == nil) {
			err = errors.New("Missing argument for heap operation")
			return
		}

		x := m.(int)
		y := n.(int)

		if !this.checkHeap(y, true) {
			err = errors.New(fmt.Sprintf("Heap address %d inaccessible", y))
		}

		this.heap[y] = x
		return
	case strings.HasPrefix(cmd, RETRIEVE):
		m := this.stack.pop()

		if m == nil {
			err = errors.New("Missing argument for heap operation")
			return
		}

		x := m.(int)

		if !this.checkHeap(x, false) {
			err = errors.New(fmt.Sprintf("Invalid heap address %d", x))
			return
		}

		this.stack.push(this.heap[x])
		return
	}

	// Errored
	err = errors.New("Illegal operation " + cmd)
	return
}

func (this *VM) checkHeap(i int, alloc bool) bool {
	if (i >= 65535) || (i < 0) {
		// Overflown
		return false
	}

	if alloc && (i >= len(this.heap)) {
		for {
			this.heap = append(this.heap, 0)

			if i < len(this.heap) {
				break
			}
		}
	}

	return true
}

func parseNum(str string, signed bool) (res int, err error) {
	str = strings.Trim(str, " ")

	if len(str) == 0 {
		err = errors.New("Can't parse empty number.")
		return
	}

	i := strings.LastIndex(str, NUM_END)

	if i > 0 {
		str = str[:i]
	}

	if signed {
		if strings.HasPrefix(str, ZERO) {
			str = str[len(ZERO):]
		} else if strings.HasPrefix(str, ONE) {
			str = "-" + str[len(ONE):]
		}
	}

	str = strings.Replace(str, ZERO, "0", -1)
	str = strings.Replace(str, ONE, "1", -1)

	r, e := strconv.ParseInt(str, 2, 32)

	res = int(r)
	err = e
	return
}
