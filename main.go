package main

import (
	"fmt"
)

func main() {
	c, err := LoadConfig()
	if err != nil {
		panic(err)
	}
	for _, cc := range c.Routes {
		fmt.Printf("%+v\n", cc)
	}
}
