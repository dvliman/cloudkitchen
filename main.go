package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

func main() {

}

func readOrdersFromFile(filepath string) []Order {
	jsonFile, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}

	defer jsonFile.Close()

	bytes, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		panic(err)
	}
	var orders []Order
	if err := json.Unmarshal(bytes, &orders); err != nil {
		panic(err)
	}

	return orders
}
