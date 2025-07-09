package main

import (
	"backend/pkg/db/sqlite"
	"fmt"
)

func main() {
	_, err := sqlite.ConnectAndMigrate()
if err != nil{
	fmt.Println(err.Error())
}
}
