package main

func main() {
	if err := GetCommand().Execute(); err != nil {
		panic(err)
	}
}
