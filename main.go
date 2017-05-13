package main

import (
	"fmt"

	"strings"

	"net/http"

	"io/ioutil"

	"github.com/Jeffail/gabs"
	"github.com/Sirikon/afoec/structs"
	"github.com/Sirikon/afoec/templates"
	"github.com/gin-gonic/gin"
)

// GinHandler .
func createItemHandler(model structs.Model) gin.HandlerFunc {
	return func(c *gin.Context) {
		var item structs.Item
		body, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			c.Error(err)
		}
		jsonParsed, _ := gabs.ParseJSON(body)
		children, _ := jsonParsed.ChildrenMap()
		for key, child := range children {
			item.Fields = append(item.Fields, structs.ItemField{ModelField: structs.ModelField{Name: key}, Value: []byte(child.Data().(string))})
		}
		item.Model = model
		// for i := range item.Fields {
		// 	item.Fields[i].ModelField = model.Fields[i]
		// }
		item.Save()
		c.String(http.StatusOK, "Yay")
	}
}

func getItemsHandler(model structs.Model) gin.HandlerFunc {
	return func(c *gin.Context) {
		// items := model.GetAll()
		// jsonObj := gabs.New()
		// jsonObj.Array()
		// for _, item := range items {
		// 	itemJSON := gabs.New()
		// 	for _, field := range item.Fields {
		// 		itemJSON.Set(string(field.Value), field.ModelField.Name)
		// 	}
		// 	jsonObj.ArrayAppend(itemJSON.Data())
		// }
		// c.String(200, jsonObj.String())
		items := model.GetAll()
		c.Data(200, "text/html", []byte(templates.Hello(model, items)))
	}
}

func main() {
	models, err := structs.GetModels()
	if err != nil {
		fmt.Println(err)
		return
	}

	r := gin.Default()

	for _, m := range models {
		fmt.Println("Processing model: " + m.Name)
		m.CreateSchema()

		r.POST("/"+strings.ToLower(m.Name), createItemHandler(m))
		r.GET("/"+strings.ToLower(m.Name), getItemsHandler(m))

		// item := Item{Model: m, Fields: []ItemField{ItemField{ModelField: m.Fields[0], Value: []byte("Pepito")}}}
		// item.Save()

		// items := m.GetAll()
		// fmt.Println("=== Items ===")
		// for _, item := range items {
		// 	fmt.Println("ID: " + strconv.Itoa(item.ID))
		// 	for _, field := range item.Fields {
		// 		fmt.Println(field.ModelField.Name + ": " + string(field.Value))
		// 	}
		// }
	}

	r.Run() // listen and serve on 0.0.0.0:8080

}
