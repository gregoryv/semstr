package semstr

import (
	"fmt"
	"testing"
)

func ExampleMustParse() {
	valid := []string{
		"v1",
		"1",
		"1.8.0",
		"1.0",
		"2.93.144-beta",
		"2.3.0-alpha+SHA123",
	}
	for _, v := range valid {
		MustParse(v)
	}
	// output:
}

func ExampleCompare() {
	cmp := func(a, b string) {
		op := map[int]string{
			-1: "<",
			0:  "=",
			1:  ">",
		}
		got, _ := Compare(a, b)

		fmt.Printf("%s %s %s\n", a, op[got], b)
	}
	cmp("1.0", "0.9")
	cmp("0.1", "0.2")
	cmp("2.0", "2.0")
	cmp("1.0.0-alpha", "1.0.0-beta")
	cmp("1.0.0-rc1", "1.0.0-rc2")
	cmp("1.0.0-rc1", "1.0.0-rc11")
	// output:
	// 1.0 > 0.9
	// 0.1 < 0.2
	// 2.0 = 2.0
	// 1.0.0-alpha < 1.0.0-beta
	// 1.0.0-rc1 < 1.0.0-rc2
	// 1.0.0-rc1 < 1.0.0-rc11
}

func ExampleVersion_String() {
	fmt.Println(MustParse("v1.9").String())
	fmt.Println(MustParse("2.0.0-alpha").String())
	fmt.Println(MustParse("2.0.0-alpha+SHA1").String())
	// output:
	// 1.9.0
	// 2.0.0-alpha
	// 2.0.0-alpha+SHA1
}

func TestParse(t *testing.T) {
	bad := func(in string) {
		t.Helper()
		_, err := Parse(in)
		if err == nil {
			t.Errorf("%s: expect error", in)
		}
	}
	bad("")
	bad("1.x")
	bad("1..")
	bad("1.x.0")
	bad("1.0.x")
	bad("1.0.x-")
	bad("1.0.1-")
	bad("1.0.1-+SHA1")
	bad("1.0.1-alpha+")
	bad("x.0.0")
}

func TestMustParse_panics(t *testing.T) {
	defer catchPanic(t)
	MustParse("abc")
}

func TestMustCompare(t *testing.T) {
	ok := func(exp int, a, b string) {
		t.Helper()
		got := MustCompare(a, b)
		if got != exp {
			t.Error(got, a, b, "expected", exp)
		}
	}
	ok(0, "1.0.0", "1.0.0")
	ok(0, "1.0.0-dev+A", "1.0.0-dev+B")
	ok(1, "1.0", "0.0")
	ok(1, "1.1", "1.0")
	ok(1, "1.0.1", "1.0")
	ok(1, "1.0.1", "1.0.1-beta")
	ok(1, "1.0.1-beta", "1.0.1-alpha")
	ok(1, "1.0.1-rc2", "1.0.1-rc1")

	defer catchPanic(t)
	MustCompare("1.a.0", "1.0")
}

func TestCompare(t *testing.T) {
	bad := func(a, b string) {
		t.Helper()
		_, err := Compare(a, b)
		if err == nil {
			t.Errorf("Compare(%q, %q), expect error", a, b)
		}
	}
	bad("1.b.0", "1.0")
	bad("1.0", "1.c.0")
}

func BenchmarkParse(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Parse("1.342.12-dev")
	}
}

func BenchmarkCompare(b *testing.B) {
	v := MustParse("1.342.12-dev")
	o := MustParse("1.342.0")
	for i := 0; i < b.N; i++ {
		v.Compare(o)
	}
}

func catchPanic(t *testing.T) {
	if err := recover(); err == nil {
		t.Fail()
	}
}
