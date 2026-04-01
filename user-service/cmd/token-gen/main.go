package main

import (
	"fmt"

	"aidanwoods.dev/go-paseto"
)

func main() {
	fmt.Println("Simple paseto token generator")

	priv := paseto.NewV4AsymmetricSecretKey()
	pub := priv.Public()

	fmt.Println("private:", priv.ExportHex())
	fmt.Println("public: ", pub.ExportHex())
}
