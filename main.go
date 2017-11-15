/****** MIT License **********
Copyright (c) 2017 Datzer Rainer
Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
***************************/

package main

/* @see https://stackoverflow.com/questions/16465705/how-to-handle-configuration-in-go
 * Used Packages: iniflags, URL: https://github.com/vharitonsky/iniflags
 * 
 */

import (
	"fmt"
	/*"reflect"
	"regexp"
	"encoding/json"
	"encoding/xml"*/
	"flag"
	//"strconv"
	"os"
	"os/user"
	/*
	"strings"
	"database/sql"
	"database/sql/driver"
	"log"
	"testing"
	"time"*/
	
	"github.com/vharitonsky/iniflags"
	"github.com/jmoiron/sqlx"
	_ "github.com/go-sql-driver/mysql"
	//"github.com/jmoiron/sqlx/reflectx"
	//_ "github.com/lib/pq"
	_ "github.com/jackc/pgx"
	_ "github.com/jackc/pgx/pgtype"
	_ "github.com/jackc/pgx/stdlib"
	_ "github.com/mattn/go-sqlite3"
	//_ "github.com/gwenn/gosqlite"
)

var (
	dbDriver = flag.String("DbDriver", "", "Database Driver")
	dbDSN = flag.String("DSN", "", "Database Data Source Name")
	dbSchema = flag.String("Schema", "public", "Default Database Schema")
	dbAnonUser = flag.String("AnonUser", "", "Anonymous User")
	outFormat = flag.String("Format", "json", "Output Format")
	cquery = flag.String("Query", "", "The Query to execute")
	ctable = flag.String("Table", "", "The Table to query")
	cid = flag.String("ID", "", "The Primary Key of the Table")
	cwhere = flag.String("Where", "", "The Where Condition")
	crecdeepth = *flag.Int("RecDeepth", 999, "The Recursion Deepth")
	convertToString = flag.Bool("ConvertToString", false, "Convert every row to String")
	db *sqlx.DB
	debugOutput = flag.Bool("DebugOutput", false, "Print debug values")
	testAnon = flag.Bool("TestAnon", false, "Test the anonymous functions")
	//db *sql.DB
	err interface{}
	defaultConfigFile string
)

func all_vars()(result string) {
	result = fmt.Sprintf("dbDriver: %s\ndbDSN: %s\ndbSchema: %s\ndbAnonUser: %s\noutFormat: %s\nquery: %s\n", *dbDriver, *dbDSN, *dbSchema, *dbAnonUser, *outFormat, *cquery)
	return
}

func dump_vars() {
	fmt.Print(all_vars())
}

func prog_run() {
	CUser, uerror := user.Current()
	if uerror != nil {
		panic("Could not determine User!")
	}
	defaultConfigFile = CUser.HomeDir + "/.rdb2om"
	iniflags.SetAllowMissingConfigFile(true)
	/// @see https://stackoverflow.com/questions/12518876/how-to-check-if-a-file-exists-in-go
	if _, err := os.Stat(defaultConfigFile); !os.IsNotExist(err) {
		iniflags.SetConfigFile(defaultConfigFile)
		if *debugOutput { fmt.Fprintf(os.Stderr, "Read from Default Configfile: %s\n", defaultConfigFile) }
	} else  {
		if *debugOutput { fmt.Fprintf(os.Stderr, "No Default Configfile %s found.\n", defaultConfigFile) }
	}
	
	iniflags.Parse()
	
	if *debugOutput { fmt.Fprintf(os.Stderr, "Flags: %v, OS Flags: %#v\n", flag.Args(), os.Args) }
	
	if *testAnon {
		fmt.Printf("queryfunc: %s\n", queryfuncs["sqlite"]["relationquery"](db, "schema_public", "table_tab"))
		os.Exit(0)
	}
	
	//dump_vars()
	connect()
	if *cquery != "" {
		if *debugOutput { fmt.Fprintf(os.Stderr, "Executing Query: %s\n", *cquery) }
		resultmap := executeQuery(db, *cquery)
		fmt.Println(Marshal(resultmap, *outFormat))
	}
	if *ctable != "" && *cid == "" && *debugOutput {
		/*var input string
		fmt.Printf("Table to inspect:")
		fmt.Scanln(&input)
		if input != "" {
			fmt.Printf(JsonMarshal(get_relations(db, *input)))
		}*/
		fmt.Println(Marshal(get_relations(db, *ctable, *dbSchema), *outFormat))
	}
	if *ctable != "" && *cid != "" {
		if *debugOutput { fmt.Fprintf(os.Stderr, "Load Table %s by id %s.\n", *ctable, *cid) }
		// attributes := map[string]string{"id":"c3989ad4-1cb6-49fb-aada-b7fb89a1d46f", "class_id":"1"}
		mylist := load_table_by_pkey(db, *dbSchema, *ctable, map[string]string{"id":*cid}, crecdeepth, false)
		fmt.Println(Marshal(mylist, *outFormat))
	}
	if *ctable != "" && *cwhere != "" {
		if *debugOutput { fmt.Fprintf(os.Stderr, "Load Table entries from %s recursive by where clause \"%s\".\n", *ctable, *cwhere) }
		mylist := fetch_recursive(db, *dbSchema, *ctable, *cwhere, crecdeepth, false)
		if *debugOutput { fmt.Fprintf(os.Stderr, "Result: %v\n", mylist) }
		fmt.Println(Marshal(mylist, *outFormat))
	}
	
}

func main() {
	prog_run()
	
}

