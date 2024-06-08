package sem

import (
	"fmt"
	"strconv"
	"strings"
)

// MustParseVersion returns a valid version or panics.
func MustParseVersion(in string) *Version {
	v, err := ParseVersion(in)
	if err != nil {
		panic(err.Error())
	}
	return v
}

// ParseVersion returns a valid sematic version or an error.
func ParseVersion(in string) (*Version, error) {
	if len(in) == 0 {
		return nil, fmt.Errorf("empty")
	}
	var v Version
	// major
	i := strings.Index(in, ".")
	if i == -1 {
		return nil, fmt.Errorf("missing dot")
	}
	var err error
	v.Major, err = strconv.Atoi(in[:i])
	if err != nil {
		return nil, err
	}

	// optional minor
	in = in[i+1:]
	i = strings.Index(in, ".")
	if i == -1 {
		v.Minor, err = strconv.Atoi(in)
		if err != nil {
			return nil, err
		}
		return &v, nil
	}
	v.Minor, err = strconv.Atoi(in[:i])
	if err != nil {
		return nil, err
	}

	// optional patch
	in = in[i+1:]
	i = strings.Index(in, "-")
	if i == -1 {
		v.Patch, err = strconv.Atoi(in)
		if err != nil {
			return nil, err
		}
		return &v, nil
	}
	v.Patch, err = strconv.Atoi(in[:i])
	if err != nil {
		return nil, err
	}
	v.Text = in[i+1:]
	return &v, nil
}

// Sless returns true if version a is considered to be less than b.
func Sless(a, b string) bool {
	A := MustParseVersion(a)
	B := MustParseVersion(b)
	return Less(A, B)
}

// Less returns true if version v is considered to come before version
// o.
func Less(v, o *Version) bool {
	return Compare(v, o) < 0
}

// numEqual returns true if major, minor and patch fields are
// equal.
func numEqual(v, o *Version) bool {
	return v.Major == o.Major && v.Minor == o.Minor && v.Patch == o.Patch
}

func Scompare(a, b string) int {
	A := MustParseVersion(a)
	B := MustParseVersion(b)
	return Compare(A, B)
}

func Compare(v, o *Version) int {
	// equal
	if numEqual(v, o) && v.Text == o.Text {
		return 0
	}
	if v.Major > o.Major {
		return 1
	}
	if v.Major == o.Major && v.Minor > o.Minor {
		return 1
	}
	if v.Major == o.Major && v.Minor == o.Minor && v.Patch > o.Patch {
		return 1
	}
	if numEqual(v, o) && len(v.Text) == 0 && len(o.Text) > 0 {
		return 1
	}
	if numEqual(v, o) && len(v.Text) > 0 && len(o.Text) > 0 && v.Text > o.Text {
		return 1
	}
	return -1
}

type Version struct {
	Major int
	Minor int
	Patch int

	Text string
}

func (v *Version) String() string {
	if len(v.Text) > 0 {
		return fmt.Sprintf("%v.%v.%v-%s", v.Major, v.Minor, v.Patch, v.Text)
	}
	return fmt.Sprintf("%v.%v.%v", v.Major, v.Minor, v.Patch)
}
