package msds

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"os"
	"reflect"
	"testing"
)

func Test_isGzip(t *testing.T) {
	createGzipData := func(data []byte) []byte {
		var buf bytes.Buffer
		gz := gzip.NewWriter(&buf)
		_, _ = gz.Write(data)
		_ = gz.Close()
		return buf.Bytes()
	}

	tests := []struct {
		name    string
		content []byte
		want    bool
	}{
		{
			name:    "gzip file",
			content: createGzipData([]byte("test data")),
			want:    true,
		},
		{
			name:    "non-gzip file",
			content: []byte("test data"),
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := bytes.NewReader(tt.content)
			bufReader := bufio.NewReader(b)
			if got := isGzip(bufReader); got != tt.want {
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
		outputDir string
		skipTable []string
		skipData  []string
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
			Writer(tt.args.outputDir, tt.args.skipData, tt.args.skipTable, tt.args.bus)
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
