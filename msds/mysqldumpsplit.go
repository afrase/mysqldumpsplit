package msds

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const (
	sentinelString = "****SENTINEL-STRING****"
	headerData     = `/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
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
)

// Producer reads `file` line-by-line and adds it to the `readCh` channel.
// Note: This function closes `file`.
func Producer(file *os.File, readCh chan string) {
	r := bufio.NewReader(file)
	for line, err := r.ReadString('\n'); err == nil; line, err = r.ReadString('\n') {
		readCh <- line
	}
	file.Close()
	close(readCh)
}

// Consumer splits the file up and fills the different channels.
func Consumer(readCh, tableNameCh, tableSchemeCh, tableDataCh chan string) {
	onTableScheme, onTableData := false, false
	for line := range readCh {
		if strings.Contains(line, "Table structure for table") {
			onTableScheme, onTableData = true, false
			tableName := strings.Replace(line, "-- Table structure for table ", "", 1)
			tableNameCh <- strings.TrimSpace(strings.Replace(tableName, "`", "", -1))
			tableSchemeCh <- line
		} else if strings.Contains(line, "LOCK TABLES `") {
			onTableData, onTableScheme = true, false
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
				tableSchemeCh <- sentinelString
			} else if strings.Contains(line, "UNLOCK TABLES;") {
				onTableData = false
				tableDataCh <- sentinelString
			}
		}
	}

	close(tableNameCh)
	close(tableDataCh)
	close(tableSchemeCh)
}

// Writer writes the data from the different channels to different files.
func Writer(outputDir string, skipTables []string, tableNameCh, tableSchemeCh, tableDataCh chan string, doneCh chan bool) {
	os.Mkdir(outputDir, os.ModePerm)
	numTables := 0

	for tableName := range tableNameCh {
		fmt.Printf("extracting table: %s\n", tableName)
		numTables++
		tablePath := filepath.Join(outputDir, tableName+".sql")
		tableFile, _ := os.Create(tablePath)

		for tableData := range tableSchemeCh {
			if tableData == sentinelString {
				break
			}
			tableFile.WriteString(tableData)
		}

		for tableData := range tableDataCh {
			if tableData == sentinelString {
				break
			}

			if !StringInArray(tableName, &skipTables) {
				tableFile.WriteString(tableData)
			}
		}
		tableFile.Close()
	}
	fmt.Printf("\nExtracted %d tables\n", numTables)
	doneCh <- true
}

// CombineFiles combines all files ina directory into a single file
func CombineFiles(filePath, outputDir string) {
	combineFile, _ := os.Create(filePath)
	combineFile.WriteString(headerData)

	files, _ := ioutil.ReadDir(outputDir)
	fmt.Printf("Combining all %d files into %s\n", len(files), filePath)

	for _, file := range files {
		sqlFile, _ := os.Open(filepath.Join(outputDir, file.Name()))
		r := bufio.NewReader(sqlFile)

		for line, err := r.ReadString('\n'); err == nil; line, err = r.ReadString('\n') {
			combineFile.WriteString(line)
		}
		// write a newline between each file
		combineFile.WriteString("\n")
	}

	info, _ := combineFile.Stat()
	fmt.Printf("New file size %s\n", StringifyFileSize(info.Size()))
	combineFile.Close()

	fmt.Println("Deleting output directory")
	os.RemoveAll(outputDir)
}
