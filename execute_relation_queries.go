/****** MIT License **********
Copyright (c) 2017 Datzer Rainer
Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
***************************/

package main

import (
	"os"
	"fmt"
	"strings"
	"regexp"
	"github.com/jmoiron/sqlx"
)

type fk_kv_pair struct {
	key string
	value interface{}
}

type fk_table struct {
	kv []fk_kv_pair
	foreign_schema_name string
	foreign_table_name string
}

func get_relations(db *sqlx.DB, table string, schema string)(resultmap []map[string]interface{}) {
	resultmap = make([]map[string]interface{}, 0, 0)
	var query string
	var found bool
	query, found = querymap[*dbDriver]["relationquery"]
	if found {
		//var re *Regexp
		re := regexp.MustCompile("\\barg_table_name\\b")
		query = re.ReplaceAllString(query, "'" + table + "'")
		re = regexp.MustCompile("\\barg_schema_name\\b")
		query = re.ReplaceAllString(query, "'public'")
		//fmt.Printf("Query: %s\n", query)
		rows, err := db.Queryx(query)
		if err != nil {
			panic(err)
		}
		for rows.Next() {
			singularmap := make(map[string]interface{})
			err = rows.MapScan(singularmap)
			resultmap = append(resultmap, singularmap)
		}
	} else {
		resultmap = relationquery_func[*dbDriver](db, schema, table)
	}
	//fmt.Printf("get_relations Map: %#v\n", resultmap)
	return
}

func get_pk_fk_mapping(db *sqlx.DB, constraintname string, mapping []map[string]interface{})(fkmap map[string]string) {
	fkmap = make(map[string]string)
	for _, row := range mapping {
		//if fmt.Sprint(row["constraint_name"]) == constraintname {
			fkmap[fmt.Sprint(row["foreign_column_name"])] = fmt.Sprint(row["column_name"])
		//}
	}
	//fmt.Printf("get_pk_fk_mapping: %v\n", fkmap)
	return
}

func get_pk(schema string, table string)(result []string) {
	var query string
	var found bool
	result = make([]string, 0)
	query, found = querymap[*dbDriver]["primary_keys_query"]
	if found {
		re := regexp.MustCompile("\\barg_table_name\\b")
		query = re.ReplaceAllString(query, "'" + table + "'")
		re = regexp.MustCompile("\\barg_schema_name\\b")
		query = re.ReplaceAllString(query, "'" + schema + "'")
		resultmap := executeQuery(db, query)
		for _, row := range resultmap {
			result = append(result, fmt.Sprint(row["attname"]))
		}
	} else {
		result = primary_keys_query_func[*dbDriver](db, schema, table)
	}
	return
}

func fetch_recursive(db *sqlx.DB, schema string, table string, where string, rec_level int, detect_cycles bool)(resultmap []map[string]interface{}) {
	resultmap = make([]map[string]interface{}, 0, 0)
	
	pklist := get_pk(schema, table)
	
	relations := get_relations(db, table, schema)
	fkmap1 := get_pk_fk_mapping(db, "", relations)
	
	if *debugOutput { fmt.Fprintf(os.Stderr, "File/Line: %s, Primary keys: %v, fkmap1: %v\n", file_line(), pklist, fkmap1) }

	query := "SELECT * FROM " + table + " WHERE " + where
	resultmap = executeQuery(db, query)
	
	for resultkey, result := range resultmap {
	
		if *debugOutput { fmt.Fprintf(os.Stderr, "File/Line: %s, Result: %v\n", file_line(), result) }
		
		pk_key_value_list := make(map[string]string, 0)
		for _, keyname := range pklist {
			//if result[fkmap1[keyname]] == nil { continue }
			pk_key_value_list[keyname] = fmt.Sprint(result[keyname])
		}
		
		if *debugOutput { fmt.Fprintf(os.Stderr, "File/Line: %s, pk_key_value_list: %v, Resultkey: %s\n", file_line(), pklist, resultkey) }
		
		for _, el := range load_table_by_pkey(db, schema, table, pk_key_value_list, rec_level - 1, detect_cycles) {
			resultmap[resultkey] = el
		}
	}
	
	return resultmap
}

// @see https://stackoverflow.com/questions/32751537/why-do-i-get-a-cannot-assign-error-when-setting-value-to-a-struct-as-a-value-i
func load_table_by_pkey(db *sqlx.DB, schema string, table string, fk map[string]string, rec_level int, detect_cycles bool)(resultmap []map[string]interface{}) {
	resultmap = make([]map[string]interface{}, 0, 0)
	
	if len(fk) == 0 {
		return
		//panic("No data in fk in load_table_by_pkey!")
	}
	
	els := make([]string, 0)
	for mkey, mvalue := range fk {
		els = append(els, mkey + "='" + mvalue + "'")
	}
	where := strings.Join(els[:]," AND ")
	
	query := "SELECT * FROM " + table + " WHERE " + where
	
	if *debugOutput { fmt.Fprintf(os.Stderr, "File/Line: %s, Query: %s\n", file_line(), query) }
	
	resultmap = executeQuery(db, query)
	
	if len(resultmap) < 1 { return nil }
	
	if *debugOutput { fmt.Fprintf(os.Stderr, "File/Line: %s, Result: %s\n", file_line(), dump_r_v(resultmap)) }
	
	for resultkey, result := range resultmap {
	
		relations := get_relations(db, table, schema)
		
		if *debugOutput { fmt.Fprintf(os.Stderr, "File/Line: %s, Relations: %v\n", file_line(), relations) }
		
		fkmap := make(map[string]*fk_table)
		
		for _, avalue := range relations {
			constraintname := fmt.Sprint(avalue["constraint_name"])
			if *debugOutput { fmt.Fprintf(os.Stderr, "File/Line: %s, Type: %T\n", file_line(), fkmap[constraintname]) }
			
			if _, ok := fkmap[constraintname];!ok {
				fkmap[constraintname] = &fk_table{kv: nil, foreign_schema_name: fmt.Sprint(avalue["foreign_schema_name"]), foreign_table_name: fmt.Sprint(avalue["foreign_table_name"])}
			}
			if fkmap[constraintname].kv == nil {
				fkmap[constraintname].kv = make([]fk_kv_pair, 0)
			}
			
			fkmap[constraintname].kv = append(fkmap[constraintname].kv, fk_kv_pair{key: fmt.Sprint(avalue["column_name"]), value: result[fmt.Sprint(avalue["foreign_column_name"])]})
		}
		
		if *debugOutput { fmt.Fprintf(os.Stderr, "File/Line: %s, Fkmap: %s\n", file_line(), dump_r_v(fkmap)) }
		
		for constraintname, conditions := range fkmap {
			conditionmap := make(map[string]string)
			ok := true
			loadtable := conditions.foreign_table_name
			loadschema := conditions.foreign_schema_name
			for _, condition := range conditions.kv {
				if condition.value == nil {
					ok = false
					break
				}
				conditionmap[condition.key] = fmt.Sprint(condition.value)
			}
			if ok == false {
				conditionmap = make(map[string]string)
			} else {
				if *debugOutput { fmt.Fprintf(os.Stderr, "File/Line: %s, Entering recursion. constraintname: %s, conditionmap: %v\n", file_line(), constraintname, conditionmap) }
				if rec_level <= 0 {
					return resultmap
				}
				resultmap[resultkey]["table[" + loadtable + "][" + constraintname + "]"] = load_table_by_pkey(db, loadschema, loadtable, conditionmap, rec_level - 1, detect_cycles)
				if *debugOutput { fmt.Fprintf(os.Stderr, "File/Line: %s, Resultmap from load call: %v\n", file_line(), resultmap[resultkey]["table[" + loadtable + "][" + constraintname + "]"]) }
			}
		}
	}
	
	if *debugOutput { fmt.Fprintf(os.Stderr, "File/Line: %s, Resultmap to return: %v\n", file_line(), resultmap) }
	
	return resultmap
}

