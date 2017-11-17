/****** MIT License **********
Copyright (c) 2017 Datzer Rainer
Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
***************************/

package main

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

var (
 
// http://www.ryanday.net/2009/11/12/postgres-table-relations/
/**
 *
<code language="sql">SELECT pc1.relname as ltable, pga2.attname as lcolumn, pc2.relname as rtable, pga1.attname as rcolumn
FROM pg_class pc1, pg_class pc2, pg_constraint, pg_attribute pga1, pg_attribute pga2
WHERE pc1.relname = '{table}' and pg_constraint.conrelid = pc1.oid
AND pc2.relkind = 'r' AND pc2.oid = pg_constraint.confrelid
AND pga1.attnum = pg_constraint.confkey[1]
AND pga1.attrelid = pc2.oid
AND pga2.attnum = pg_constraint.conkey[1]
AND pga2.attrelid = pc1.oid
</code>

@see https://stackoverflow.com/questions/1152260/postgres-sql-to-list-table-foreign-keys  , post answered Jun 18 '13 at 8:56 oscavi

<code>
select c.constraint_name
    , x.table_schema as schema_name
    , x.table_name
    , x.column_name
    , y.table_schema as foreign_schema_name
    , y.table_name as foreign_table_name
    , y.column_name as foreign_column_name
from information_schema.referential_constraints c
join information_schema.key_column_usage x
    on x.constraint_name = c.constraint_name
join information_schema.key_column_usage y
    on y.ordinal_position = x.position_in_unique_constraint
    and y.constraint_name = c.unique_constraint_name
order by c.constraint_name, x.ordinal_position
</code>
*/

/// https://github.com/StefanSchroeder/Golang-NestedDatastructures-Tutorial/blob/master/chapter04.markdown
querymap = map[string]map[string]string{
	"pgx":map[string]string{
		"relationquery":
		`
select c.constraint_name::name
    , x.table_schema::name as schema_name
    , x.table_name::name
    , x.column_name::name
    , y.table_schema::name as foreign_schema_name
    , y.table_name::name as foreign_table_name
    , y.column_name::name as foreign_column_name
from information_schema.referential_constraints c
join information_schema.key_column_usage x
    on x.constraint_name = c.constraint_name
join information_schema.key_column_usage y
    on y.ordinal_position = x.position_in_unique_constraint
    and y.constraint_name = c.unique_constraint_name
where y.table_name=arg_table_name and y.table_schema=arg_schema_name
order by c.constraint_name, x.ordinal_position
`,
/**
 * @see https://wiki.postgresql.org/wiki/Retrieve_primary_key_columns
 * @see https://stackoverflow.com/questions/1214576/how-do-i-get-the-primary-keys-of-a-table-from-postgres-via-plpgsql
 */
		"primary_keys_query":`
SELECT               
  pg_attribute.attname as attname, 
  format_type(pg_attribute.atttypid, pg_attribute.atttypmod),
  nspname,
  pg_class.relname
FROM pg_index, pg_class, pg_attribute, pg_namespace 
WHERE 
  pg_class.oid = arg_table_name::regclass AND 
  indrelid = pg_class.oid AND 
  nspname = arg_schema_name AND 
  pg_class.relnamespace = pg_namespace.oid AND 
  pg_attribute.attrelid = pg_class.oid AND 
  pg_attribute.attnum = any(pg_index.indkey)
 AND indisprimary
`,
// https://stackoverflow.com/questions/769683/show-tables-in-postgresql
"show_table_query":`
SELECT * FROM pg_catalog.pg_tables WHERE schemaname != 'pg_catalog' AND schemaname != 'information_schema'
`,
	},
	"mysql":map[string]string{


/**
@see https://stackoverflow.com/questions/20855065/how-to-find-all-the-relations-between-all-mysql-tables
<code>
SELECT 
  `TABLE_SCHEMA`,                          -- Foreign key schema
  `TABLE_NAME`,                            -- Foreign key table
  `COLUMN_NAME`,                           -- Foreign key column
  `REFERENCED_TABLE_SCHEMA`,               -- Origin key schema
  `REFERENCED_TABLE_NAME`,                 -- Origin key table
  `REFERENCED_COLUMN_NAME`                 -- Origin key column
FROM
  `INFORMATION_SCHEMA`.`KEY_COLUMN_USAGE`  -- Will fail if user don't have privilege
WHERE
  `TABLE_SCHEMA` = SCHEMA()                -- Detect current schema in USE 
  AND `REFERENCED_TABLE_NAME` IS NOT NULL; -- Only tables with foreign keys
</code>
 */
"relationquery":`
SELECT 
  CONSTRAINT_NAME AS constraint_name,                        -- Constraint Name
  TABLE_SCHEMA AS schema_name,                               -- Foreign key schema
  TABLE_NAME AS table_name,                                  -- Foreign key table
  COLUMN_NAME AS column_name,                                -- Foreign key column
  REFERENCED_TABLE_SCHEMA AS foreign_schema_name,            -- Origin key schema
  REFERENCED_TABLE_NAME AS foreign_table_name,               -- Origin key table
  REFERENCED_COLUMN_NAME AS foreign_column_name              -- Origin key column
FROM
  INFORMATION_SCHEMA.KEY_COLUMN_USAGE                        -- Will fail if user don't have privilege
WHERE
  REFERENCED_TABLE_NAME = arg_table_name
  AND TABLE_SCHEMA = SCHEMA()                                -- Detect current schema in USE 
  AND REFERENCED_TABLE_NAME IS NOT NULL;                     -- Only tables with foreign keys
`,
/// @see https://stackoverflow.com/questions/201621/how-do-i-see-all-foreign-keys-to-a-table-or-column
// SELECT i.TABLE_NAME, i.CONSTRAINT_TYPE, i.CONSTRAINT_NAME, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME 
"primary_keys_query":`
SELECT k.REFERENCED_COLUMN_NAME AS attname, k.REFERENCED_TABLE_SCHEMA AS nspname, i.CONSTRAINT_NAME AS relname
FROM information_schema.TABLE_CONSTRAINTS i 
LEFT JOIN information_schema.KEY_COLUMN_USAGE k ON i.CONSTRAINT_NAME = k.CONSTRAINT_NAME 
WHERE i.CONSTRAINT_TYPE = 'FOREIGN KEY' 
AND i.TABLE_SCHEMA = DATABASE()
AND i.TABLE_NAME = arg_table_name;
`,
	},
	"sqlite":map[string]string{

	},
}

/**
* sqlite:
* @see https://stackoverflow.com/questions/5499003/sqlite-list-all-foreign-keys-in-a-database
* PRAGMA foreign_key_list(table)
* PRAGMA table_info(table) 
*/

queryfuncs = map[string]map[string]func(*sqlx.DB, string, string) func() {
	/*"postgres":map[string]func(string, string) {
	},
	"mysql":map[string]func(string, string) string {
	},*/
	
	/*"sqlite":map[string]func(*sqlx.DB, string, string) func() {
		"relationquery":func(db *sqlx.DB, schema string, table string) func() {
			return nil
			//return "Schema:" + schema + ", Table:" + table
		},
		"primary_keys_query":func(db *sqlx.DB, schema string, table string) func() {
			return nil
			//return "Schema:" + schema + ", Table:" + table
		},
	},*/
}

relationquery_func = map[string]func(*sqlx.DB, string, string) []map[string]interface{} {
	"sqlite3":func(db *sqlx.DB, schema string, table string) []map[string]interface{} {
		resultmap := make([]map[string]interface{}, 0, 0)
		relquery := "PRAGMA foreign_key_list(" + table + ");"
		relmap := executeQuery(db, relquery)
		for _, row := range relmap {
			new_entry := make(map[string]interface{})
			new_entry["constraint_name"] = row["id"]
			new_entry["foreign_schema_name"] = schema
			new_entry["foreign_table_name"] = row["table"]
			new_entry["column_name"] = row["from"]
			new_entry["foreign_column_name"] = row["to"]
			resultmap = append(resultmap, new_entry)
		}
// 		fmt.Printf("relationquery_func: %v\n", resultmap)
		return resultmap
	},
}

primary_keys_query_func = map[string]func(db *sqlx.DB, schema string, table string) []string {
	"sqlite3":func(db *sqlx.DB, schema string, table string) []string {
		result := make([]string, 0)
		pkquery := "PRAGMA table_info(" + table + ");"
		pkmap := executeQuery(db, pkquery)
		for _, row := range pkmap {
			//fmt.Printf("primary_keys_query_func for loop, row: %v\n", row)
			if fmt.Sprint(row["pk"]) == "1" { result = append(result, fmt.Sprintf("%s", row["name"])) }
		}
// 		fmt.Printf("primary_keys_query_func: %v\n", result)
		return result
	},
}

)

