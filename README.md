## Esri-ASCII-Parser
A CLI utility written in [Go](https://golang.org/ "Go Language") to validate and parse ESRI ASCII raster data file to various devices (*i.e.*, console, file, database).

An example of the how the output of this parser can be used is noted here, as a population density plot in Google Maps:
![Google Maps Population Density Plots](https://cloud.githubusercontent.com/assets/10182110/10521011/35a57632-7322-11e5-8b4b-780dd723a55a.png "Google Maps Population Density Plots")

For specific file format details, see the [ESRI ASCII raster data format specification](http://resources.esri.com/help/9.3/ArcGISDesktop/com/Gp_ToolRef/Spatial_Analyst_Tools/esri_ascii_raster_format.htm "ESRI ASCII Raster Data Format") reference.

## Requirements

 - [Go](https://golang.org/ "Go Language") language distribution
	 - Install through your OS package system (*e.g.*, [APT](http://en.wikipedia.org/wiki/Advanced_Packaging_Tool)) or download directly [here](https://golang.org/dl/).
 -  Supplemental [Go](https://golang.org/ "Go Language") packages:
    +  Go-SQL-Driver package (http://github.com/go-sql-driver/mysql)
 -  Operational database (if persisting to a database, `dev=db`)
    + See MySQL DDL project sample (test_world_create.sql)
    
While this package was written and tested under Linux (Ubuntu 14.04 LTS), there should be no reason why this won't work under other operating systems. 

## Installation
To meet the Go-SQL-Driver package dependency, install the package into your [Go](https://golang.org/ "Go Language") environment ([`$GOPATH`](http://code.google.com/p/go-wiki/wiki/GOPATH "GOPATH")) using the [`go`](http://golang.org/cmd/go/ "go command") command tool:

	$ /go/src$ go get github.com/github.com/go-sql-driver/mysql

Similarly, install this esri-ascii-parser package into the Go environment:

	$ /go/src$ go get github.com/richbl/esri-ascii-parser
	
Both package sources should be viewable under the `$GOPATH/src` folder (*e.g.*, `/go/src/github.com/richbl/esri-ascii-parser`).
	
## Basic Usage
Esri-ASCII-Parser is run through a command-line interface (CLI), so all of the command options are made available there.

Here's the default response when running `esri-ascii-parser` with no parameters:

    $ esri-ascii-parser

	--------------------------------------------------------------------------------
	  esri-ascii-parser: ESRI ASCII raster data file parser (v0.10)
	  validate and parse file to console, file, or database
	--------------------------------------------------------------------------------
	-in must be set to a fully qualified pathname for input
	
	Current settings:
	  -in=
	  -out=/tmp/file.out
	  -dev=con
	  -db=
	  -name=region
	
	Usage:
	  esri-ascii-parser (-in="</path/to/esri_file>") [-out="</path/to/file.out">]
	    [-dev=<"con">, <"file">, <"db">]
	    [-db="[username[:password]@][protocol[(address)]]/dbname"]
	    [-name=<"region">]
	
	Options:
	  -db="": database to persist parsed results (used when -dev=db)
	  -dev="con": output device(s): con, file, db, e.g,. -dev=con,file
	  -in="": fully qualified pathname of esri ascii raster file to parse (required)
	  -name="region": name to uniquely identify a region (used when -dev=db)
	  -out="/tmp/file.out": file to persist parsed results (used when -dev=file)
	
	Examples:
	  esri-ascii-parser -in=/abc.asc
	  esri-ascii-parser -in=/abc.asc -dev=con,file
	  esri-ascii-parser -in=/abc.asc -dev=file -out=/tmp/results.out
	  esri-ascii-parser -in=/abc.asc -dev=db -db=usr:pwd@tcp(10.10.10.1:3306)/db

In this example, the utility responds by indicating that the `-in` parameter must be set to an input file to be parsed (otherwise, what's the point if the input file is null, really).

The output generated from the parse of the utility is managed with the `-dev` parameter. By default, it's set to console (*i.e.*, `-dev=con`), but it can be set to multiple devices at once (*e.g.*, `-dev=con,file`). Note that depending on how this parameter is set, additional parameters will need to be set.

## License

The MIT License (MIT)

Copyright (c) 2015 Business Learning Incorporated

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
