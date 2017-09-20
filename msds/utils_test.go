package msds

import (
	"testing"
	"io/ioutil"
	"os"
)

func TestSkipTables_String(t *testing.T) {
	tests := []struct {
		name string
		s    *SkipTables
		want string
	}{
		{"Single table", &SkipTables{"table1"}, "[table1]"},
		{"Multiple tables", &SkipTables{"table1", "table2", "table3"}, "[table1 table2 table3]"},
		{"Same table twice", &SkipTables{"table1", "table2", "table2"}, "[table1 table2 table2]"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.String(); got != tt.want {
				t.Errorf("SkipTables.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSkipTables_Set(t *testing.T) {
	type args struct {
		value string
	}

	tests := []struct {
		name    string
		s       *SkipTables
		args    args
		wantErr bool
	}{
		{"No error single table", &SkipTables{}, args{"table1"}, false},
		{"No error multiple tables", &SkipTables{}, args{"table1,table2, table3"}, false},
		{"Returns error single table", &SkipTables{"table1"}, args{"table1"}, true},
		{"Returns error multiple tables", &SkipTables{"table1"}, args{"table1,table2"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.Set(tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("SkipTables.Set() error = %v, wantErr %v", err, tt.wantErr)
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

	tests := []struct{
		name string
		args args
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
