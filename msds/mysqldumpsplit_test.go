package msds

import (
	"bufio"
	"os"
	"reflect"
	"testing"
)

func Test_isGzip(t *testing.T) {
	type args struct {
		b *bufio.Reader
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isGzip(tt.args.b); got != tt.want {
				t.Errorf("isGzip() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_openReader(t *testing.T) {
	type args struct {
		f *os.File
	}
	tests := []struct {
		name string
		args args
		want *bufio.Reader
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := openReader(tt.args.f); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("openReader() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLineReader(t *testing.T) {
	type args struct {
		file *os.File
		bus  ChannelBus
	}
	tests := []struct {
		name string
		args args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			LineReader(tt.args.file, tt.args.bus)
		})
	}
}

func TestLineParser(t *testing.T) {
	type args struct {
		bus          ChannelBus
		combineFiles bool
	}
	tests := []struct {
		name string
		args args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			LineParser(tt.args.bus, tt.args.combineFiles)
		})
	}
}

func TestWriter(t *testing.T) {
	type args struct {
		outputDir  string
		skipTables []string
		bus        ChannelBus
	}
	tests := []struct {
		name string
		args args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Writer(tt.args.outputDir, tt.args.skipTables, tt.args.bus)
		})
	}
}

func TestCombineFiles(t *testing.T) {
	type args struct {
		filePath  string
		outputDir string
		bus       ChannelBus
	}
	tests := []struct {
		name string
		args args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CombineFiles(tt.args.filePath, tt.args.outputDir, tt.args.bus)
		})
	}
}

func TestLogger(t *testing.T) {
	type args struct {
		bus ChannelBus
	}
	tests := []struct {
		name string
		args args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Logger(tt.args.bus)
		})
	}
}
