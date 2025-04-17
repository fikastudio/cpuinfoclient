package main

import (
	"fmt"

	"github.com/fikastudio/cpuinfoclient"
)

func main() {
	fmt.Println(cpuinfoclient.ProcessorName())
	fmt.Println(cpuinfoclient.NumCores())
}
