package main

import (
	"fmt"
)

/*
int sum(int a,int b){
	return a+b;
}
*/
import "C"

func main() {
	fmt.Println(C.sum(1, 2))
}
