package main

import (
	"fmt"
	"github.com/Alex-zs/oucjwparser/jwmodel"
)

func main() {
	session := jwmodel.JwSession{}

	for true {
		session.Login("17020031002", "chen1234")
		fmt.Println("success:", jwmodel.LoginSuccess, " failure:", jwmodel.LoginFailure)
	}
}
