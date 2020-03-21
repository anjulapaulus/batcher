package main

import "fmt"

func main(){
	s:= make([]int,3)
	fmt.Println(s)
	s[0] = 1
	s[1] = 2
	s[2] = 3
	fmt.Println(s)
}
