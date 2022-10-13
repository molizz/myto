package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/molizz/myto"
)

func main() {
	// if len(os.Args) < 2 {
	// 	log.Panicf("must input DDL sql, %d", len(os.Args))
	// }
	// ddl := os.Args[1]
	reader := bufio.NewReader(os.Stdin)
	input, err := io.ReadAll(reader)
	if err != nil {
		log.Panicf("STD input is required, %+v", err)
	}

	output, err := myto.New(string(input), true).ToDMDB()
	if err != nil {
		panic(err)
	}
	fmt.Println(output)
}
