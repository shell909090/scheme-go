package main

import (
	// "fmt"
	"flag"
	"os"

	"bitbucket.org/shell909090/scheme-go/scmgo"
)

func main() {
	flag.Parse()
	if len(flag.Args()) < 1 {
		panic("parameters not enough")
	}

	file, err := os.Open(flag.Args()[0])
	if err != nil {
		panic(err)
	}
	defer file.Close()

	code, err := scmgo.BuildCode(file)
	if err != nil {
		panic(err)
	}

	code.Format(os.Stdout, 0)
	// fmt.Println(code)
	os.Stdout.Write([]byte("\n"))
}
