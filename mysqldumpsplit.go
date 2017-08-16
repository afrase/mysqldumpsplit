package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const headerData = `
/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;
`

func producer(filePath string, readCh chan string) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		readCh <- line
	}
	close(readCh)
}

func consumer(readCh, tableNameCh, tableSchemeCh, tableDataCh chan string) {
	onTableScheme := false
	onTableData := false
	for line := range readCh {
		if strings.Contains(line, "Table structure for table") {
			onTableScheme = true
			onTableData = false
			tableName := strings.Replace(line, "-- Table structure for table ", "", 1)
			tableNameCh <- strings.TrimSpace(strings.Replace(tableName, "`", "", -1))
			tableSchemeCh <- line
		} else if strings.Contains(line, "LOCK TABLES `") {
			onTableData = true
			onTableScheme = false
			tableDataCh <- line
		} else {
			if onTableScheme {
				tableSchemeCh <- line
			}
			if onTableData {
				tableDataCh <- line
			}

			if strings.Contains(line, "-- Dumping data for table") {
				onTableScheme = false
				tableSchemeCh <- "SENTINEL_STRING"
			} else if strings.Contains(line, "UNLOCK TABLES;") {
				onTableData = false
				tableDataCh <- "SENTINEL_STRING"
			}
		}
	}
	close(tableNameCh)
	close(tableDataCh)
	close(tableSchemeCh)
}

func writer(outputDir string, tableNameCh, tableSchemeCh, tableDataCh chan string, doneCh chan bool) {
	os.Mkdir(outputDir, os.ModePerm)

	for tableName := range tableNameCh {
		fmt.Printf("%s\n", tableName)
		tablePath := filepath.Join(outputDir, tableName+".sql")
		tableFile, _ := os.Create(tablePath)

		for tableData := range tableSchemeCh {
			if tableData == "SENTINEL_STRING" {
				break
			}
			tableFile.WriteString(tableData)
		}

		for tableData := range tableDataCh {
			if tableData == "SENTINEL_STRING" {
				break
			}
			if tableName != "default_api_logs" && tableName != "default_ci_sessions" {
				tableFile.WriteString(tableData)
			}
		}
		tableFile.Close()
	}
	doneCh <- true
}

func combineFiles(filePath, outputDir string) {
	combineFile, _ := os.Create(filePath)
	combineFile.WriteString(headerData)

	defer combineFile.Close()

	files, _ := ioutil.ReadDir(outputDir)
	for _, file := range files {
		sqlFile, _ := os.Open(filepath.Join(outputDir, file.Name()))
		sqlFileReader := bufio.NewReader(sqlFile)

		for line, err := sqlFileReader.ReadString('\n'); err == nil; line, err = sqlFileReader.ReadString('\n') {
			combineFile.WriteString(line)
		}
		combineFile.WriteString("\n")
	}
	os.RemoveAll(outputDir)
}

func main() {
	var inputFile, outputPath, combineFilePath string
	var combine bool

	flag.StringVar(&inputFile, "input", "", "The file to read from")
	flag.StringVar(&inputFile, "i", "", "The file to read from")
	flag.StringVar(&outputPath, "output", "output", "The output path")
	flag.StringVar(&outputPath, "o", "output", "The output path")

	flag.StringVar(&combineFilePath, "combinefile", "dumpfile.sql", "The path to output single SQL file")
	flag.BoolVar(&combine, "combine", false, "Combine all tables into a single file")

	flag.Parse()

	readCh := make(chan string)
	tableNameCh := make(chan string)
	tableDataCh := make(chan string)
	tableSchemeCh := make(chan string)
	doneCh := make(chan bool)

	go producer(inputFile, readCh)
	go consumer(readCh, tableNameCh, tableSchemeCh, tableDataCh)
	go writer(outputPath, tableNameCh, tableSchemeCh, tableDataCh, doneCh)

	<-doneCh

	if combine {
		combineFiles(combineFilePath, outputPath)
	}
}
