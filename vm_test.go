package gmh

import (
	"testing"
	"fmt"
	"strings"
	"time"
)

func TestParseNum(t *testing.T) {
	testNum("泥草草", false, 4, t)
	testNum("草泥草草", true, 4, t)
	testNum("泥泥泥草", true, -6, t)
}

func TestOneToTen(t *testing.T) {
	res, err := New().Exec(strings.Split(
		`草草草泥马 马草草草泥草草草草泥泥马 草马草 泥马草泥 草草草泥草泥草马 泥马草草 草草草泥马 泥草草草 草马草 草草草泥草泥泥马 泥草草泥 马泥草草泥草草草泥草泥马 马草马草泥草草草草泥泥马 马草草草泥草草草泥草泥马 草马马 马马马`,
		" "))

	fmt.Println(res, err)
}

func TestHeap(t *testing.T) {
	res, err := New().Exec(strings.Split(
		"草草草泥 草草草泥草草 泥泥草 草草草泥 泥泥泥 泥马草泥 河蟹",
		" "))

	fmt.Println(res, err)
}

func TestCall(t *testing.T) {
	res, err := New().Exec(strings.Split(
		"草草草泥草马 草草草泥草草马 草草草泥泥草马 马草泥泥 泥马草泥 河蟹 马草草泥马 泥草草草 草马泥 泥草泥草 马泥马",
		" "))

	fmt.Println(res, err)
}

func TestInterrupt(t *testing.T) {
	end := time.Now().Add(5 * time.Second)
	res, err := New().SetInterrupter(func() bool {
		return time.Now().After(end)
	}).Exec(strings.Split(
		"马草草泥马 马草马泥马",
		" "))

	fmt.Println(res, err)
}

func testNum(str string, signed bool, expect int, t *testing.T) {
	res, err := parseNum(str, signed)

	fmt.Printf("%s -> %d\n", str, res)

	if (res != expect) || (err != nil) {
		fmt.Println(err)
		t.Fail()
	}
}
