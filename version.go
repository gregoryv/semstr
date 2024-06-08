/*
Package semver provides semantic version parse and compare funcs.
*/
package semver

import (
	"fmt"
	"strconv"
	"strings"
)

func MustCompare(a, b string) int {
	v, err := Compare(a, b)
	if err != nil {
		panic(err.Error())
	}
	return v
}

// Compare returns -1, 0 or 1. If a parsing error occurs it returns 0
// and a non nil error.
func Compare(a, b string) (int, error) {
	v, err := Parse(a)
	if err != nil {
		return 0, err
	}

	o, err := Parse(b)
	if err != nil {
		return 0, err
	}

	return v.Compare(o), nil
}

// MustParse returns a valid version or panics.
func MustParse(in string) *Version {
	v, err := Parse(in)
	if err != nil {
		panic(err.Error())
	}
	return v
}

// Parse returns a valid sematic version or an error.
func Parse(in string) (*Version, error) {
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

// numEqual returns true if major, minor and patch fields are
// equal.
func numEqual(v, o *Version) bool {
	return v.Major == o.Major && v.Minor == o.Minor && v.Patch == o.Patch
}

type Version struct {
	Major int
	Minor int
	Patch int

	Text string
}

func (v *Version) String() string {
	var res []byte
	res = fmt.Append(res, v.Major, ".", v.Minor, ".", v.Patch)
	if len(v.Text) > 0 {
		res = fmt.Append(res, "-", v.Text)
	}
	return string(res)
}

// Compare returns
//
//	 1, v > o
//	 0, v == o
//	-1, v < o
func (v *Version) Compare(o *Version) int {
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
