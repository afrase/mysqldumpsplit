package msds

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestCsvFlagType_String(t *testing.T) {
	tests := []struct {
		name string
		s    *CsvFlagType
		want string
	}{
		{"Single table", &CsvFlagType{"table1"}, "[table1]"},
		{"Multiple tables", &CsvFlagType{"table1", "table2", "table3"}, "[table1 table2 table3]"},
		{"Same table twice", &CsvFlagType{"table1", "table2", "table2"}, "[table1 table2 table2]"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.String(); got != tt.want {
				t.Errorf("CsvFlagType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCsvFlagType_Set(t *testing.T) {
	type args struct {
		value string
	}

	tests := []struct {
		name    string
		s       *CsvFlagType
		args    args
		wantErr bool
	}{
		{"No error single table", &CsvFlagType{}, args{"table1"}, false},
		{"No error multiple tables", &CsvFlagType{}, args{"table1,table2, table3"}, false},
		{"Returns error single table", &CsvFlagType{"table1"}, args{"table1"}, true},
		{"Returns error multiple tables", &CsvFlagType{"table1"}, args{"table1,table2"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.Set(tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("CsvFlagType.Set() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStringInArray(t *testing.T) {
	type args struct {
		str            string
		arrayOfStrings *[]string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"String in array", args{"bar", &[]string{"foo", "bar", "baz"}}, true},
		{"String not in array", args{"bar", &[]string{"foo", "baz"}}, false},
		{"Wildcard in array", args{"foo", &[]string{"f*", "bar", "baz"}}, true},
		{"Single character in array", args{"bar", &[]string{"foo", "b?r", "baz"}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StringInArray(tt.args.str, tt.args.arrayOfStrings); got != tt.want {
				t.Errorf("StringInArray() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringifyFileSize(t *testing.T) {
	type args struct {
		size int64
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"bytes", args{275}, "275.0 B"},
		{"megabytes", args{42000000}, "42.0 MB"},
		{"gigabytes", args{42500000000}, "42.5 GB"},
		{"terabytes", args{42500000000000}, "42.5 TB"},
		{"petabytes", args{42500000000000000}, "42.5 PB"},
		{"exabytes", args{4250000000000000000}, "4.3 EB"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StringifyFileSize(tt.args.size); got != tt.want {
				t.Errorf("StringifyFileSize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOpenFile(t *testing.T) {
	type args struct {
		path string
	}
	tempFile1, _ := ioutil.TempFile("", "file1")

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"no error opening file", args{tempFile1.Name()}, false},
		{"error opening non existing file", args{"foo.txt"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := OpenFile(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("OpenFile() error = %v, wantErr %v", err, tt.wantErr)
			}
			f.Close()
		})
	}

	tempFile1.Close()
	os.Remove(tempFile1.Name())
}

func TestWildcardMatch(t *testing.T) {
	type args struct {
		str     string
		pattern string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"Empty pattern and string", args{"", ""}, true},
		{"Empty pattern", args{"foo", ""}, false},
		{"Empty string", args{"", "bar"}, false},
		{"Identical strings", args{"foo", "foo"}, true},
		{"Single character match start", args{"foo", "?oo"}, true},
		{"Single character wildcard", args{"foobar", "?*"}, true},
		{"Single character matches wildcard", args{"*", "?"}, true},
		{"Single character match end", args{"foo", "fo?"}, true},
		{"Single character match middle", args{"foobar", "foo?ar"}, true},
		{"Multiple single character matches", args{"foobar", "f?o?a?"}, true},
		{"Wildcard matches wildcard", args{"*", "*"}, true},
		{"Wildcard match", args{"foo", "f*"}, true},
		{"Wildcard with unmatched ending", args{"foo", "f*bar"}, false},
		{"Multiple wildcard matches", args{"foobar", "f*b*r"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WildcardMatch(tt.args.str, tt.args.pattern); got != tt.want {
				t.Errorf("WildcardMatch() = %v, want %v", got, tt.want)
			}
		})
	}
}
