package main

import "fmt"

func main() {
	s := CreateNewServer()
	//s.MountHandlers()
	//http.ListenAndServe(":3000", s.Router)
	fmt.Printf("Hello, world.\n %+v\n", s)
}
