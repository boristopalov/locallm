package main

import (
	"fmt"
	// "os"
	// "github.com/boristopalov/localsearch/utils"
)

func main() {
	fmt.Println("starting api server")
	//
	// t, _ := os.ReadFile("html_test2.txt")
	// text := utils.ExtractText(string(t))
	// fmt.Println(text)
	// c, err := utils.GetDbClient()
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	return
	// }
	// fmt.Println("client: ", c)
	// err = c.CreateDatabase(c.Context, "localsearch")
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }
	// dbs, err := c.ListDatabases(c.Context)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	return
	// }
	// for _, db := range dbs {
	// 	fmt.Println("DB name: ", db.Name)
	// }
	// colls, err := c.ListCollections(c.Context)
	// if err != nil {
	// 	fmt.Println("faield to list oclleciotions")
	// 	return
	// }
	// for _, col := range colls {
	// 	fmt.Println("collection name: ", col.Name)
	// 	fmt.Println("collection Schema: ", col.Schema)
	// }
	// err = c.DropCollection(c.Context, "localsearch")
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }
	// err = utils.MaybeCreateCollection()
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }
	// c.ReleaseCollection(c.Context, "localsearch")
	startServer()
}
