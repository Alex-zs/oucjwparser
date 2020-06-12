package main

import "github.com/Alex-zs/oucjwparser/jwmodel"

func main() {
	session := jwmodel.JwSession{}
	if session.Login("17020031002", "chen1234") {
		session.GetStuCourse(2019, 2)
	}
}
