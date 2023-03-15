package main

import (
	"zctl/cmd/z"
)

func main() {
	if err := z.GetCommand().Execute(); err != nil {
		panic(err)
	}
}
