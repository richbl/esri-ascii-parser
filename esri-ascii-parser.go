/*-----------------------------------------------------------------------------
The MIT License (MIT)

 Copyright (C) Business Learning Incorporated (www.businesslearninginc.com)

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
-----------------------------------------------------------------------------*/

//
//  EsriASCIIParser - ESRI ASCII Raster Data File Parser
//
//  Validate and parse ESRI ASCII raster data file to various devices (i.e.,
//  console, file, database)
//
//  ESRI format specification reference:
//    http://resources.esri.com/help/9.3/ArcGISDesktop/com/Gp_ToolRef \
//      /Spatial_Analyst_Tools/esri_ascii_raster_format.htm
//

package main

import (
	"bufio"
	"flag"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"

	"errors"
	"os"
	"time"
)

// ================================================================================================
// local global constants/variables

// user-assigned values to override command-line interface flags (based on useCommandLineFlags)
const useCommandLineFlags bool = true

var outputToConsole = true
var outputToFile = true
var outputToDb = false

var inFilePath = "/tmp/cymuds00ag.asc"
var outFilePath = "/tmp/file.out"
var regionName = "esri_region_name"
var dbName = "root:root@/test_world"

// non-user-assigned constants
const appVersion = "0.10"

// esri ascii raster data format parameters
// http://resources.esri.com/help/9.3/ArcGISDesktop/com/Gp_ToolRef/Spatial_Analyst_Tools/esri_ascii_raster_format.htm
const (
	nCols = iota
	nRows
	xllCorner
	yllCorner
	cellSize
	noDataValue
)

type esriHeaderNameStruct struct {
	parameter string
	value     string
}

type esriHeaderDataStruct []struct {
	parameter string
	value     float64
}

type esriHeaderStruct struct {
	esriHeaderName esriHeaderNameStruct
	esriHeaderData esriHeaderDataStruct
}

// ================================================================================================
// functions

// checkErr is a simple wrapper to check err value passed to it
func checkErr(err error) {

	if err != nil {
		panic(err)
	}

}

// boolToInt converts a bool to an int (true = 1, else 0)
func boolToInt(b bool) int {

	if b {
		return 1
	}
	return 0
}

// showBanner displays the title banner to Stdout (configured for 80-columns)
func showBanner() {

	fmt.Fprintf(os.Stdout, "--------------------------------------------------------------------------------\n")
	fmt.Fprintf(os.Stdout, "  ESRI ASCII raster data file parser (v%s)\n", appVersion)
	fmt.Fprintf(os.Stdout, "  validate and parse file to console, file, or database\n")
	fmt.Fprintf(os.Stdout, "--------------------------------------------------------------------------------\n\n")

}

// parseFlags parses command-line flags and passes them on function return
//   flags used:
//     -dev=con, file, db (device output, can be cumulative)
//     -in=infile (fully qualified pathname to input esri ascii raster format file)
//     -out=outfile (relative path to output file to persist parse results)
//     -db=database (database name to persist parse results)
//	   -name=region name (when persisted to database)
func parseFlags() (inFile string, outputToConsole bool, outputToFile bool, outputToDb bool, outFile string, dbFile string, regionName string) {

	devFlag := flag.String("dev", "con", "output device(s): con, file, db, e.g,. -dev=con,file")
	inFileFlag := flag.String("in", "", "fully qualified pathname of esri ascii raster file to parse (required)")
	outFileFlag := flag.String("out", "/tmp/file.out", "file to persist parsed results (used when -dev=file)")
	dbFlag := flag.String("db", "", "database to persist parsed results (used when -dev=db)")
	nameFlag := flag.String("name", "region", "name to uniquely identify a region (used when -dev=db)")

	flag.Usage = func() {

		fmt.Fprintf(os.Stderr, "\nCurrent settings:\n")
		fmt.Fprintf(os.Stderr, "  -in=%s\n", *inFileFlag)
		fmt.Fprintf(os.Stderr, "  -out=%s\n", *outFileFlag)
		fmt.Fprintf(os.Stderr, "  -dev=%s\n", *devFlag)
		fmt.Fprintf(os.Stderr, "  -db=%s\n", *dbFlag)
		fmt.Fprintf(os.Stderr, "  -name=%s\n", *nameFlag)

		fmt.Fprintf(os.Stderr,
			"\nUsage:\n  %s (-in=\"</path/to/esri_file>\") "+
				"[-out=\"</path/to/file.out\">]\n"+
				"    [-dev=<\"con\">, <\"file\">, <\"db\">]\n"+
				"    [-db=\"[username[:password]@][protocol[(address)]]/dbname\"]\n"+
				"    [-name=<\"region\">]\n\n",
			filepath.Base(os.Args[0]))

		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()

		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  %s -in=/abc.asc\n", filepath.Base(os.Args[0]))
		fmt.Fprintf(os.Stderr, "  %s -in=/abc.asc -dev=con,file\n", filepath.Base(os.Args[0]))
		fmt.Fprintf(os.Stderr, "  %s -in=/abc.asc -dev=file -out=/tmp/results.out\n", filepath.Base(os.Args[0]))
		fmt.Fprintf(os.Stderr, "  %s -in=/abc.asc -dev=db -db=usr:pwd@tcp(10.10.10.1:3306)/db\n", filepath.Base(os.Args[0]))
		fmt.Fprintf(os.Stderr, "\n")
		os.Exit(0)
	}

	flag.Parse()
	flagErr := false

	if *inFileFlag == "" {
		fmt.Fprintf(os.Stderr, "-in must be set to a fully qualified pathname for input\n")
		flagErr = true
	}

	if strings.Contains(*devFlag, "con") {
		outputToConsole = true
	}

	if strings.Contains(*devFlag, "file") {
		outputToFile = true

		if *outFileFlag == "" {
			fmt.Fprintf(os.Stderr, "when -dev=file, -out must be set to a filename\n")
			flagErr = true
		}
	}

	if strings.Contains(*devFlag, "db") {
		outputToDb = true

		if *dbFlag == "" {
			fmt.Fprintf(os.Stderr, "when -dev=db, -db must be set to a database\n")
			flagErr = true
		}

		if *nameFlag == "" {
			fmt.Fprintf(os.Stderr, "when -dev=db, -name must be set to identify the region\n")
			flagErr = true
		}

	}

	if outputToDb == false && outputToFile == false && outputToConsole == false {
		fmt.Fprintf(os.Stderr, "-dev must be con, file, and/or db\n")
		flagErr = true
	}

	if flagErr {
		flag.Usage()
	}

	return *inFileFlag, outputToConsole, outputToFile, outputToDb, *outFileFlag, *dbFlag, *nameFlag

}

// esriHeaderProcess scans only the esri raster data header block and determines if the header is valid
// and returns whether the optional NoDataValue parameter is present
func esriHeaderProcess(esriHeader esriHeaderStruct) (tableID int64, hasNoDataValue bool) {

	// noDataValue is an optional parameter in esri ascii raster data, so parse for all parameters, but special case for noDataValue
	// and return bool on presence in header block
	hasNoDataValue = true

	tableID = 0
	inputFile, err := os.Open(inFilePath)
	checkErr(err)
	defer inputFile.Close()

	scanner := bufio.NewScanner(inputFile)

	for i := 0; i < 6; i++ {

		scanner.Scan()
		parameter := strings.ToLower(strings.Fields(scanner.Text())[0])

		// parse for esriHeaderData parameters
		if esriHeader.esriHeaderData[i].parameter == parameter {
			if len(strings.Fields(scanner.Text())) < 2 {
				checkErr(errors.New("invalid esri file format: header parameter value for " + esriHeader.esriHeaderData[i].parameter + " not found"))
			}
			value, err := strconv.ParseFloat(strings.Fields(scanner.Text())[1], 16)
			checkErr(err)
			esriHeader.esriHeaderData[i].value = value
		} else {

			if i == noDataValue {
				hasNoDataValue = false                     // optional noDataValue parameter not found in header
				esriHeader.esriHeaderData[i].value = -9999 // default parameter value
			} else {
				checkErr(errors.New("invalid esri file format: header parameter " + esriHeader.esriHeaderData[i].parameter + " not found"))
			}
		}

	}

	if outputToDb {
		tableID = esriHeaderPersist(esriHeader)
	}

	if outputToConsole {
		fmt.Println("esri header processed as:\n", esriHeader, "\n")
	}

	if outputToFile {
		file, err := os.Create(outFilePath)
		checkErr(err)
		defer file.Close()

		fmt.Fprintln(file, esriHeader)
	}

	return tableID, hasNoDataValue

}

// esriDataProcess scans only the esri raster data block and builds lat/lon values from header block
// configuration details, persisting them to the specified output device
func esriDataProcess(esriHeader esriHeaderStruct, tableID int64, hasNoDataValue bool) {

	var file *os.File

	inputFile, err := os.Open(inFilePath)
	checkErr(err)
	defer inputFile.Close()
	scanner := bufio.NewScanner(inputFile)

	// advance file pointer over esri file header block
	for i := 0; i < 5+boolToInt(hasNoDataValue); i++ {
		scanner.Scan()
	}

	if outputToFile {
		file, err = os.Create(outFilePath)
		checkErr(err)
		defer file.Close()
	}

	// calculate upper-left reference datum
	lat := esriHeader.esriHeaderData[yllCorner].value + (esriHeader.esriHeaderData[nRows].value * esriHeader.esriHeaderData[cellSize].value)

	// loop through esri data block (nCols x nRows)
	for i := 0; i < int(esriHeader.esriHeaderData[nRows].value); i++ {

		scanner.Scan()
		value := strings.Split(scanner.Text(), " ")

		if int(esriHeader.esriHeaderData[nCols].value) != len(value)-1 {
			checkErr(errors.New("nCols != line length"))
		}

		lon := esriHeader.esriHeaderData[xllCorner].value
		lat = lat - esriHeader.esriHeaderData[cellSize].value

		for j := 0; j < len(value)-1; j++ {

			value, err := strconv.ParseFloat(value[j], 16)
			checkErr(err)

			if outputToConsole {
				fmt.Println("lat:", float32(lat), ", lon:", float32(lon), ", value:", value)
			}

			if outputToFile {
				fmt.Fprintln(file, float32(lat), ",", float32(lon), ",", value)
			}

			if outputToDb {
				esriDataPersist(lat, lon, value, tableID)
			}

			lon = lon + esriHeader.esriHeaderData[cellSize].value
		}
	}

}

// esriHeaderPersist takes the esriHeaderStruct object and persists into mysql database as specified
// by user through command-line flags
func esriHeaderPersist(esriHeader esriHeaderStruct) (tableID int64) {

	db, err := sql.Open("mysql", dbName)
	err = db.Ping()
	checkErr(err)
	defer db.Close()

	result, err := db.Exec("INSERT INTO esriHeader(name, nCols, nRows, xllCorner, yllCorner, cellSize, noDataValue) VALUES(?, ?, ?, ?, ?, ?, ?)",
		esriHeader.esriHeaderName.value, esriHeader.esriHeaderData[nCols].value, int(esriHeader.esriHeaderData[nRows].value), esriHeader.esriHeaderData[xllCorner].value, esriHeader.esriHeaderData[yllCorner].value,
		esriHeader.esriHeaderData[cellSize].value, esriHeader.esriHeaderData[noDataValue].value)
	checkErr(err)

	tableID, err = result.LastInsertId()
	checkErr(err)

	return tableID

}

// esriDataPersist takes the data value objects (lat, lon, and value) and persists into mysql database
func esriDataPersist(lat float64, lon float64, value float64, tableID int64) {

	// TODO: non-optimal, db open/close should be moved into separate func, since this function gets called iteratively
	db, err := sql.Open("mysql", dbName)
	err = db.Ping()
	checkErr(err)
	defer db.Close()

	_, err = db.Exec("INSERT INTO esriData(lat, lon, value, esriHeader_id_esriHeader) VALUES(?, ?, ?, ?)", lat, lon, value, tableID)
	checkErr(err)

}

//
// ================================================================================================
//

func main() {

	if outputToConsole {
		showBanner()
	}

	if useCommandLineFlags {
		inFilePath, outputToConsole, outputToFile, outputToDb, outFilePath, dbName, regionName = parseFlags()
	}

	startTime := time.Now().Local()

	var esriHeader = esriHeaderStruct{
		esriHeaderNameStruct{
			"name", regionName,
		},
		esriHeaderDataStruct{
			{"ncols", 0},
			{"nrows", 0},
			{"xllcorner", 0},
			{"yllcorner", 0},
			{"cellsize", 0},
			{"nodata_value", 0},
		},
	}

	// get esri file HEADER block
	// need to know if NoDataValue parameter is present for esri file data block parsing
	tableID, hasNoDataValue := esriHeaderProcess(esriHeader)

	// get esri file DATA block
	esriDataProcess(esriHeader, tableID, hasNoDataValue)

	if outputToConsole {
		fmt.Println("\n===========\nTime elapsed:", time.Since(startTime))
	}

}
