package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/afrase/mysqldumpsplit/msds"
)

var version = "master"

type config struct {
	InputFile       string
	OutputPath      string
	CombineFilePath string
	Combine         bool
	Version         bool
	SkipTable       msds.CsvFlagType
	SkipData        msds.CsvFlagType
}

func parseFlags() *config {
	conf := new(config)

	flag.StringVar(&conf.InputFile, "i", "", "The file to read from, can be a gzip file")
	flag.StringVar(&conf.OutputPath, "o", "output", "The output path ")

	flag.Var(&conf.SkipData, "skipData",
		"Comma separated list of tables you want to skip outputting the data for.\n\tUse '*' to skip all.")
	flag.Var(&conf.SkipTable, "skipTable",
		"Comma separated list of tables to skip.\n\tNames can contain '*' for wildcard values")

	flag.BoolVar(&conf.Combine, "combine", false,
		"Combine all tables into a single file, deletes individual table files")
	flag.StringVar(&conf.CombineFilePath, "combineFile", "dumpfile.sql",
		"The path to output a single SQL file\n\tOnly used if combine flag is set")
	flag.BoolVar(&conf.Version, "version", false, "Display the version and exit")

	flag.Parse()
	return conf
}

func main() {
	conf := parseFlags()
	if conf.Version {
		fmt.Println(version)
		os.Exit(0)
	}

	if conf.InputFile == "" {
		flag.PrintDefaults()
		os.Exit(0)
	}

	file, err := msds.OpenFile(conf.InputFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	bus := msds.ChannelBus{
		Finished:    make(chan bool),
		Log:         make(chan string),
		TableData:   make(chan string),
		TableScheme: make(chan string),
		TableName:   make(chan string),
		CurrentLine: make(chan string),
	}

	go msds.Logger(bus)

	bus.Log <- fmt.Sprintf("outputing all tables to %s\n", conf.OutputPath)
	if len(conf.SkipData) > 0 {
		bus.Log <- fmt.Sprintf("skiping data from tables %s\n", strings.Join(conf.SkipData, ", "))
	}
	if len(conf.SkipTable) > 0 {
		bus.Log <- fmt.Sprintf("skiping tables %s\n", strings.Join(conf.SkipTable, ", "))
	}

	start := time.Now()
	bus.Log <- fmt.Sprintf("begin processing %s\n", conf.InputFile)
	// create a pipeline of goroutines
	go msds.LineReader(file, bus)
	go msds.LineParser(bus, conf.Combine)
	go msds.Writer(conf.OutputPath, conf.SkipData, conf.SkipTable, bus)

	// wait for the writer to finish.
	<-bus.Finished

	if conf.Combine {
		msds.CombineFiles(conf.CombineFilePath, conf.OutputPath, bus)
	}

	bus.Log <- fmt.Sprintf("finished in %s", time.Now().Sub(start))
	bus.Log <- ""
	close(bus.Log)
	close(bus.Finished)
}
