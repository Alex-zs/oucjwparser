package main

import (
	"fmt"
	"github.com/Alex-zs/oucjwparser/jwmodel"
)

func main() {
	session := jwmodel.NewSession()
	if success, _ := session.Login("17020031002", "chen1234"); success {
		creditList, err := session.GetCreditRequire()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(creditList)
	}
}
