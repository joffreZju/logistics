package main

import (
	"common/lib/keycrypt"
	"flag"
	"fmt"
)

var (
	key        = flag.String("k", "", "crypt key")
	plaintext  = flag.String("e", "", "plaintext")
	ciptertext = flag.String("d", "", "ciptertext")

	helper = `
keycipter help

encode text
	keycipter -k "this is key" -e "this is plaintext"

encode text
	keycipter -k "this is key" -d "this is ciptertext"
	`
)

func main() {
	flag.Parse()
	if *key == "" {
		fmt.Println(helper)
		return
	}
	if len(*plaintext) > 0 {
		fmt.Println(keycrypt.Encode(*key, *plaintext))
		return
	}
	if len(*ciptertext) > 0 {
		fmt.Println(keycrypt.Decode(*key, *ciptertext))
		return
	}
	fmt.Println(helper)
}
