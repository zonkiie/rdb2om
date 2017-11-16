/****** MIT License **********
Copyright (c) 2017 Datzer Rainer
Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
***************************/

package main

import (
	"fmt"
	"os"
	"github.com/jmoiron/sqlx"
)

// use special syntay for mysql connection
// @see https://github.com/gogits/gogs/issues/39
// user:pass@tcp(localhost:3306)/database
func connect() {
	if *debugOutput { fmt.Fprintf(os.Stderr, "Driver: %s, DSN: %s\n", *dbDriver, *dbDSN) }
	if *dbDriver == "" { panic("Driver is not defined!") }
	if *dbDSN == "" { panic("DSN is not defined!") }
	db, err = sqlx.Open(*dbDriver, *dbDSN)
	//db, err = sql.Open(*dbDriver, *dbDSN)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		//panic(err)
	}
}

// https://stackoverflow.com/questions/35362459/golang-create-a-slice-of-maps
// https://stackoverflow.com/questions/6372474/how-to-determine-an-interface-values-real-type
func executeQuery(db *sqlx.DB, query string)(resultmap []map[string]interface{}) {
	resultmap = make([]map[string]interface{}, 0, 0)
	rows, err := db.Queryx(query)
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		singularmap := make(map[string]interface{})
		err = rows.MapScan(singularmap)
		if *convertToString { convertToStringFunc(singularmap) }
		resultmap = append(resultmap, singularmap)
	}
	return resultmap
}

func convertToStringFunc(singularmap map[string]interface{}) {
	for key, value := range singularmap {
		switch value.(type) {
			case int, int32, int64, uint32, uint64:
				singularmap[key] = fmt.Sprintf("%v", value)
			case float32, float64:
				singularmap[key] = fmt.Sprintf("%lf", value)
			case nil:
				singularmap[key] = fmt.Sprintf("")
			case bool:
				singularmap[key] = fmt.Sprintf("%t", value)
			default:
				singularmap[key] = fmt.Sprintf("%s", value)
		}
	}
}
