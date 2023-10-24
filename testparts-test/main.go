package main

import (
	"fmt"

	"github.com/abaskin/testparts"
)

func main() {
	head := testparts.SectionHeadSt{
		Title: "Test Section",
	}

	fmt.Println(head.Title)
}
