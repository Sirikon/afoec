package structs

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"os"

	"strconv"

	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type Item struct {
	Model  Model
	ID     int
	Fields []ItemField
}

type ItemField struct {
	ModelField ModelField
	Value      []byte
}

// Model defines the structure of a model
type Model struct {
	Name   string
	Fields []ModelField
}

// ModelField is a field structure, Model have a list of this
type ModelField struct {
	Name string
	Type string
}

// GetModels returns the models defined as .JSON files
func GetModels() ([]Model, error) {
	var files []os.FileInfo
	if f, err := ioutil.ReadDir("./models"); err == nil {
		files = f
	} else {
		return nil, err
	}

	models := make([]Model, 0, len(files))

	for _, f := range files {
		var data []byte
		if d, err := ioutil.ReadFile("./models/" + f.Name()); err == nil {
			data = d
		} else {
			return nil, err
		}
		var model Model
		err := json.Unmarshal(data, &model)
		if err != nil {
			return nil, err
		}
		models = append(models, model)
	}

	return models, nil
}

// CreateSchema will create the schema in the DB
func (m Model) CreateSchema() {
	db, err := sql.Open("sqlite3", "./data.db")
	checkErr(err)

	sql := "CREATE TABLE `" + m.Name + "` ("
	sql += "`id` INTEGER PRIMARY KEY AUTOINCREMENT, "
	for i, field := range m.Fields {
		sql += "`" + field.Name + "` VARCHAR(200) NULL"
		if i < len(m.Fields)-1 {
			sql += ","
		}
	}
	sql += ");"

	createStatement, err := db.Prepare(sql)
	checkErr(err)

	_, err = createStatement.Exec()
	checkErr(err)
}

// Save will save the given item using it's associated model
func (item Item) Save() {
	db, err := sql.Open("sqlite3", "./data.db")
	checkErr(err)

	sql := "INSERT INTO " + item.Model.Name + "("
	for i, field := range item.Fields {
		sql += field.ModelField.Name
		if i < len(item.Fields)-1 {
			sql += ","
		}
	}
	sql += ") values ("
	for i, field := range item.Fields {
		sql += "\"" + string(field.Value) + "\""
		if i < len(item.Fields)-1 {
			sql += ","
		}
	}
	sql += ");"

	fmt.Println(sql)

	insertStatement, err := db.Prepare(sql)
	checkErr(err)

	_, err = insertStatement.Exec()
	checkErr(err)
}

// GetAll returns all the items
func (m Model) GetAll() []Item {
	db, err := sql.Open("sqlite3", "./data.db")
	checkErr(err)

	sql := "SELECT * FROM " + m.Name + ";"

	rows, err := db.Query(sql)
	checkErr(err)

	var result []Item

	for rows.Next() {
		values := make([]interface{}, 0, len(m.Fields)+1)
		values = append(values, &[]byte{})
		for _ = range m.Fields {
			values = append(values, &[]byte{})
		}
		rows.Scan(values...)
		item := Item{Model: m, Fields: []ItemField{}}
		for i, value := range values {
			valueBytes := unReferenceByteSlice(value.(*[]byte))
			if i == 0 {
				item.ID, _ = strconv.Atoi(string(valueBytes))
			} else {
				item.Fields = append(item.Fields, ItemField{ModelField: m.Fields[i-1], Value: valueBytes})
			}
		}

		result = append(result, item)
	}

	return result
}

func unReferenceByteSlice(src *[]byte) []byte {
	var result []byte
	for _, b := range *src {
		result = append(result, b)
	}
	return result
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
