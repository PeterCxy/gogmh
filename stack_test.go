package gmh

import (
	"fmt"
	"math/rand"
	"testing"
)

func TestPushPop(*testing.T) {
	stack := newStack()

	for i := 0; i < 10; i++ {
		r := rand.Int31()
		stack.push(r)
		fmt.Printf("Pushing %d\n", r)
	}

	for {
		m := stack.pop()
		fmt.Printf("Poped %d\n", m)

		if m == nil {
			break
		}
	}
}
