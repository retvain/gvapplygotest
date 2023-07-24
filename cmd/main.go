package main

import "cmd/processor"

func main() {
	err := processor.Run()
	if err != nil {
		print(err.Error())
	}
}
