EsriASCIIParser
==

CLI utility written in Go to validate and parse ESRI ASCII raster data file to various devices (i.e., console, file, database).

For format details, see the ESRI ASCII raster data format specification reference:
http://resources.esri.com/help/9.3/ArcGISDesktop/com/Gp_ToolRef/Spatial_Analyst_Tools/esri_ascii_raster_format.htm

Requirements
--
+ Go language distribution
    + Install through your OS package system or go here: https://golang.org/dl/
+  Supplemental Go packages:
    +  Go-SQL-Driver (http://github.com/go-sql-driver/mysql)
+  Operational database (if persisting to a database, dev=db)
    + See MySQL DDL project sample (test_world_create.sql)
    
While this package was written and tested under Linux (Ubuntu 14.04 LTS), there should be no reason why this won't work just fine under other operating systems. 

Basic Usage
--
EsriASCIIParser is run through a command-line interface (CLI), so all of the command options are made available there.

Here's the default response when running esri-ascii-parser with no parameters:

    richbl@main:~$ esri-ascii-parser

```
--------------------------------------------------------------------------------
  ESRI ASCII raster data file parser (v0.10)
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
  esri_ascii_parser (-in="</path/to/esri_file>") [-out="</path/to/file.out">]
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
  esri_ascii_parser -in=/abc.asc
  esri_ascii_parser -in=/abc.asc -dev=con,file
  esri_ascii_parser -in=/abc.asc -dev=file -out=/tmp/results.out
  esri_ascii_parser -in=/abc.asc -dev=db -db=usr:pwd@tcp(10.10.10.1:3306)/db
```
The utility responses by indicating that the -in parameter must be set to a file to be parsed (otherwise, what's the point if the input file is null, really).

The output generated from the parse of the utility is managed with the -dev parameter. By default, it's set to console (i.e., -dev=con), but it can be set to multiple devices at once (e.g., -dev=con,file). Note that depending on how this parameter is set, additional parameters will need to be set.
License
--
The MIT License (MIT)

Copyright (c) 2014 Business Learning Incorporated

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
