package main

import (
	"fmt"
	"log"
	"regexp"
)

// 驗證用的正則表達
const (
	RegexpEmail    = "^[a-zA-Z0-9.!#$%&’*+/=?^_`{|}~-]+@[a-zA-Z0-9-]+(?:\\.[a-zA-Z0-9-]+)*$"
	RegexpPassword = "^[\\w\\d]{6,14}$"
	RegexpPhone    = "^[0-9]{10}$"
)

func main() {
	invalidEmail := "yam8511gmail.com"
	email := "yam8511@gmail.com"
	invalidPassword := "demo12"
	password := "demo123"
	invalidPhone := "0912-345-678"
	phone := "0912345678"

	reg, err := regexp.Compile(RegexpEmail)
	if err != nil {
		log.Fatal(err)
		return
	}
	ok := reg.MatchString(email)
	fmt.Printf("Email: %s ---> %v\n", email, ok)
	ok = reg.MatchString(invalidEmail)
	fmt.Printf("Invalid Email: %s ---> %v\n", invalidEmail, ok)

	reg, err = regexp.Compile(RegexpPassword)
	if err != nil {
		log.Fatal(err)
		return
	}
	ok = reg.MatchString(password)
	fmt.Printf("Password: %s ---> %v\n", password, ok)
	ok = reg.MatchString(invalidPassword)
	fmt.Printf("Invalid Password: %s ---> %v\n", invalidPassword, ok)

	reg, err = regexp.Compile(RegexpPhone)
	if err != nil {
		log.Fatal(err)
		return
	}
	ok = reg.MatchString(phone)
	fmt.Printf("Phone: %s ---> %v\n", phone, ok)
	ok = reg.MatchString(invalidPhone)
	fmt.Printf("Invalid Phone: %s ---> %v\n", invalidPhone, ok)
}
