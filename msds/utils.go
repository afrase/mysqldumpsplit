package msds

import (
	"errors"
	"fmt"
	"math"
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
		if str == a {
			return true
		}
	}
	return false
}

// StringifyFileSize converts bytes to something more readable.
func StringifyFileSize(size int64) string {
	suffix := []string{"B", "KB", "MB", "GB", "TB", "PB", "EB"}
	e := math.Floor(math.Log(float64(size)) / math.Log(1000))
	val := math.Floor(float64(size)/math.Pow(1000, e)*10+0.5) / 10
	return fmt.Sprintf("%.1f %s", val, suffix[int(e)])
}
