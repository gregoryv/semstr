package sem

import (
	"fmt"
	"testing"
)

func ExampleVersion_String() {
	v := Version{1, 9, 11, ""}
	fmt.Println(v.String())
	v.Text = "dev"
	fmt.Println(v.String())
	// output:
	// 1.9.11
	// 1.9.11-dev
}

func ExampleParseVersion() {
	fmt.Println(MustParseVersion("1.8.0"))
	fmt.Println(MustParseVersion("1.0"))
	fmt.Println(MustParseVersion("2.93.144-beta"))
	// output:
	// 1.8.0
	// 1.0.0
	// 2.93.144-beta
}

func TestParseVersion(t *testing.T) {
	bad := func(in string) {
		t.Helper()
		_, err := ParseVersion(in)
		if err == nil {
			t.Errorf("%s: expect error", in)
		}
	}
	bad("")
	bad("1")
	bad("1.x")
	bad("1.x.0")
	bad("1.0.x")
	bad("1.0.x-")
	bad("x.0.0")
}

func TestVersion_Less(t *testing.T) {
	ok := func(a, b string) {
		t.Helper()
		if !Sless(a, b) {
			t.Errorf("%s < %s", a, b)
		}
	}
	ok("0.3.5", "1.0.0")
	ok("0.3.5-dev", "0.3.5")
}

func TestMustParseVersion_panics(t *testing.T) {
	defer func() {
		e := recover()
		if e == nil {
			t.Fail()
		}
	}()
	MustParseVersion("abc")
}

func TestVersion_Compare(t *testing.T) {
	ok := func(exp int, a, b string) {
		t.Helper()
		got := Scompare(a, b)
		if got != exp {
			t.Error(got, a, b, "expected", exp)
		}
	}
	ok(0, "1.0.0", "1.0.0")
	ok(1, "1.0", "0.0")
	ok(1, "1.1", "1.0")
	ok(1, "1.0.1", "1.0")
	ok(1, "1.0.1", "1.0.1-beta")
	ok(1, "1.0.1-beta", "1.0.1-alpha")
	ok(1, "1.0.1-rc2", "1.0.1-rc1")
}

func BenchmarkParseVersion(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ParseVersion("1.342.12-dev")
	}
}

func BenchmarkCompare(b *testing.B) {
	v := MustParseVersion("1.342.12-dev")
	o := MustParseVersion("1.342.0")
	for i := 0; i < b.N; i++ {
		Compare(v, o)
	}
}
