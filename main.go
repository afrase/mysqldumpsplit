package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/afrase/mysqldumpsplit/msds"
)

type config struct {
	InputFile       string
	OutputPath      string
	CombineFilePath string
	Combine         bool
	Skip            msds.SkipTables
}

func openFile(path string) *os.File {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "The file '%s' does not exist.\n", path)
		os.Exit(1)
	}

	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	info, _ := file.Stat()
	fmt.Printf("Original file size %s\n", msds.StringifyFileSize(info.Size()))

	return file
}

func parseFlags() *config {
	conf := new(config)

	flag.StringVar(&conf.InputFile, "i", "", "The file to read from")
	flag.StringVar(&conf.OutputPath, "o", "output", "The output path ")

	flag.Var(&conf.Skip, "skipData", "Comma separated list of tables you don't want the data for")

	flag.BoolVar(&conf.Combine, "combine", false, "Combine all tables into a single file")
	flag.StringVar(&conf.CombineFilePath, "combineFile", "dumpfile.sql",
		"The path to output a single SQL file\n\tOnly used if combine flag is set")

	flag.Parse()
	return conf
}

func main() {
	conf := parseFlags()

	readCh := make(chan string)
	tableNameCh := make(chan string)
	tableDataCh := make(chan string)
	tableSchemeCh := make(chan string)
	doneCh := make(chan bool)

	file := openFile(conf.InputFile)

	fmt.Printf("Outputing all tables to %s\n", conf.OutputPath)
	if len(conf.Skip) > 0 {
		fmt.Printf("Skiping data from tables %s\n", strings.Join(conf.Skip, ", "))
	}

	fmt.Printf("Begin processing %s\n\n", conf.InputFile)
	// create a pipeline of goroutines
	go msds.Producer(file, readCh)
	go msds.Consumer(readCh, tableNameCh, tableSchemeCh, tableDataCh)
	go msds.Writer(conf.OutputPath, conf.Skip, tableNameCh, tableSchemeCh, tableDataCh, doneCh)

	// wait for the writer to finish.
	<-doneCh

	if conf.Combine {
		msds.CombineFiles(conf.CombineFilePath, conf.OutputPath)
	}
}
