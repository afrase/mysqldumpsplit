package msds

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

const sentinelString = "****SENTINEL-STRING****"

// ChannelBus a struct to hold all channels used by the different go routines
type ChannelBus struct {
	Finished    chan bool
	Log         chan string
	CurrentLine chan string
	TableName   chan string
	TableScheme chan string
	TableData   chan string
}

func isGzip(b *bufio.Reader) bool {
	if m, err := b.Peek(2); err == nil {
		return m[0] == 0x1f && m[1] == 0x8b
	}
	return false
}

func openReader(f *os.File) *bufio.Reader {
	pageSize := os.Getpagesize() * 2
	buf := bufio.NewReaderSize(f, pageSize)
	if isGzip(buf) {
		gbuf, _ := gzip.NewReader(buf)
		return bufio.NewReaderSize(gbuf, pageSize)
	}
	return buf
}

// LineReader reads `file` line-by-line and adds it to the `bus.CurrentLine` channel.
// Note: This function closes `file`.
func LineReader(file *os.File, bus ChannelBus) {
	r := openReader(file)
	for line, err := r.ReadString('\n'); err == nil; line, err = r.ReadString('\n') {
		bus.CurrentLine <- line
	}
	file.Close()
	close(bus.CurrentLine)
}

// LineParser reads the CurrentLine and figures out which channel to put it in.
func LineParser(bus ChannelBus, combineFiles bool) {
	onTableScheme, onTableData, pastHeader := false, false, false
	headerMetaData := fmt.Sprintf("-- Generated with mysqldumpsplit on %s\n\n", time.Now())
	for line := range bus.CurrentLine {
		// The beginning of a mysqldump has some flags at the top of the file. Capture them into a variable.
		if !pastHeader && strings.Contains(line, "/*!40") {
			headerMetaData += line
		}

		if strings.Contains(line, "Table structure for table") {
			onTableScheme, onTableData = true, false
			tableName := strings.Replace(line, "-- Table structure for table ", "", 1)
			bus.TableName <- strings.TrimSpace(strings.Replace(tableName, "`", "", -1))
			// add headers to each file unless we are combining all of them into 1 file.
			if !combineFiles {
				bus.TableScheme <- headerMetaData
			} else if !pastHeader {
				// add the meta data to only the first table.
				bus.TableScheme <- headerMetaData
			}

			pastHeader = true
			bus.TableScheme <- "\n--\n" + line
		} else if strings.Contains(line, "LOCK TABLES `") {
			onTableData, onTableScheme = true, false
			bus.TableData <- line
		} else {
			if onTableScheme {
				bus.TableScheme <- line
			}
			if onTableData {
				bus.TableData <- line
			}

			if strings.Contains(line, "-- Dumping data for table") {
				onTableScheme = false
				bus.TableScheme <- "--\n"
				bus.TableScheme <- sentinelString
			} else if strings.Contains(line, "UNLOCK TABLES;") {
				onTableData = false
				bus.TableData <- sentinelString
			}
		}
	}

	close(bus.TableName)
	close(bus.TableData)
	close(bus.TableScheme)
}

// Writer writes the data from the different channels to different files.
func Writer(outputDir string, skipData []string, skipTables []string, bus ChannelBus) {
	os.Mkdir(outputDir, os.ModePerm)
	numTables := 0

	for tableName := range bus.TableName {
		var skipTableData bool
		skipTable := StringInArray(tableName, &skipTables)

		if skipTable {
			// also skip the data for the table.
			skipTableData = true
			bus.Log <- fmt.Sprintf("skipping table: %s\n", tableName)
		} else {
			bus.Log <- fmt.Sprintf("extracting table: %s\n", tableName)
			numTables++
			skipTableData = StringInArray(tableName, &skipData)
		}

		// to keep the code somewhat DRY, we create the sql file even if skipping the table.
		// we then delete it at the end of the loop.
		tablePath := filepath.Join(outputDir, tableName+".sql")
		tableFile, _ := os.Create(tablePath)

		for tableData := range bus.TableScheme {
			if tableData == sentinelString {
				break
			}

			if !skipTable {
				tableFile.WriteString(tableData)
			}
		}

		if skipTableData && !skipTable {
			// not skipping table but skipping data
			bus.Log <- fmt.Sprintf("skipping data for table: %s\n", tableName)
		}

		for tableData := range bus.TableData {
			if tableData == sentinelString {
				break
			}

			if !skipTableData {
				tableFile.WriteString(tableData)
			}
		}

		tableFile.Close()

		if skipTable {
			os.Remove(tablePath)
		}
	}
	bus.Log <- fmt.Sprintf("extracted %d tables\n", numTables)
	bus.Finished <- true
}

// CombineFiles combines all files ina directory into a single file
func CombineFiles(filePath, outputDir string, bus ChannelBus) {
	combineFile, _ := os.Create(filePath)
	cleanUpOutputDir := true

	files, _ := ioutil.ReadDir(outputDir)
	bus.Log <- fmt.Sprintf("Combining all %d files into %s\n", len(files), filePath)

	for _, file := range files {
		fullPath := path.Join(outputDir, file.Name())
		if combineFile.Name() == fullPath {
			cleanUpOutputDir = false
			continue
		}

		sqlFile, _ := OpenFile(filepath.Join(outputDir, file.Name()))
		r := bufio.NewReader(sqlFile)

		for line, err := r.ReadString('\n'); err == nil; line, err = r.ReadString('\n') {
			combineFile.WriteString(line)
		}
		// write a newline between each file
		combineFile.WriteString("\n")
		// close then delete the table file
		fileName := sqlFile.Name()
		sqlFile.Close()
		os.Remove(fileName)
	}

	info, _ := combineFile.Stat()
	bus.Log <- fmt.Sprintf("New file size %s\n", StringifyFileSize(info.Size()))
	combineFile.Close()

	if cleanUpOutputDir {
		bus.Log <- fmt.Sprintf("Deleting output directory")
		os.RemoveAll(outputDir)
	}
}

// Logger reads messages from `bus.Log` and outputs them to the logger.
func Logger(bus ChannelBus) {
	for msg := range bus.Log {
		if msg != "" {
			log.Output(3, msg)
		}
	}
}
