package msds

import (
	"errors"
	"fmt"
	"math"
	"os"
	"strings"
)

// SkipTables converts a comma separated list into an array of strings.
type SkipTables []string

func (s *SkipTables) String() string {
	return fmt.Sprint(*s)
}

// Set sets the value
func (s *SkipTables) Set(value string) error {
	if len(*s) > 0 {
		return errors.New("skip tables flag already set")
	}
	for _, v := range strings.Split(value, ",") {
		*s = append(*s, strings.TrimSpace(v))
	}
	return nil
}

// StringInArray loops over `arrayOfStrings` and returns `true` if `str` is in the array.
func StringInArray(str string, arrayOfStrings *[]string) bool {
	for _, a := range *arrayOfStrings {
		if WildcardMatch(a, str) {
			return true
		}
	}
	return false
}

// StringifyFileSize converts bytes to something more readable.
func StringifyFileSize(size int64) string {
	suffix := []string{"B", "KB", "MB", "GB", "TB", "PB", "EB"}
	if size == 0 {
		return "0 B"
	}
	e := math.Floor(math.Log(float64(size)) / math.Log(1000))
	val := math.Floor(float64(size)/math.Pow(1000, e)*10+0.5) / 10
	return fmt.Sprintf("%.1f %s", val, suffix[int(e)])
}

// OpenFile tries to open the file at `path`.
func OpenFile(path string) (*os.File, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("the file '%s' does not exist", path)
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	return file, nil
}

// WildcardMatch matches strings using "*" and "?".
func WildcardMatch(str, pattern string) bool {
	s := []rune(str)
	p := []rune(pattern)

	// empty pattern only matches an empty string
	if len(p) == 0 {
		return len(s) == 0
	}

	lookup := initLookupTable(len(s)+1, len(p)+1)

	// empty pattern and empty string match
	lookup[0][0] = true

	// only '*' can match an empty string.
	for i := 1; i < len(p)+1; i++ {
		if p[i-1] == '*' {
			lookup[0][i] = lookup[0][i-1]
		}
	}

	for i := 1; i < len(s)+1; i++ {
		for j := 1; j < len(p)+1; j++ {
			if p[j-1] == '*' {
				// '*' matches any character
				lookup[i][j] = lookup[i][j-1] || lookup[i-1][j]
			} else if p[j-1] == '?' || s[i-1] == p[j-1] {
				// if pattern character is '?' or if the pattern and string actually match.
				lookup[i][j] = lookup[i-1][j-1]
			} else {
				// characters don't match.
				lookup[i][j] = false
			}
		}
	}

	return lookup[len(s)][len(p)]
}

func initLookupTable(rows, columns int) [][]bool {
	lookup := make([][]bool, rows)
	for i := range lookup {
		lookup[i] = make([]bool, columns)
	}
	return lookup
}
