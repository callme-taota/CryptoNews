package main

import "os"

func main() {
	list := Start()
	//fmt.Println(len(list), list)

	preFile, err := os.Create("pre_file.csv")
	if err != nil {
		panic(err)
	}
	defer preFile.Close()
	preList := PreDayOnly(list)
	structToCsv(preList, preFile)

	currFile, err := os.Create("curr_file.csv")
	if err != nil {
		panic(err)
	}
	defer currFile.Close()
	currList := CurrDayOnly(list)
	structToCsv(currList, currFile)
}
