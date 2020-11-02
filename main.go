package main

import (
	"fmt"

	"github.com/codingtony/ansible-vault-go/cmd"
	"github.com/codingtony/ansible-vault-go/vault"
)

func main() {
	var s, _ = vault.Encrypt("aaa", "bbb")
	fmt.Printf("%s\n", s)
	cmd.Execute()
}
