package shellcode

/*
#include <stdio.h>
#include <sys/mman.h>
#include <string.h>
#include <unistd.h>

void call(char *shellcode) {
	if(fork()) {
		return;
	}
	unsigned char *ptr;
	ptr = (unsigned char *) mmap(0, 0x1000, \
		PROT_READ|PROT_WRITE|PROT_EXEC, MAP_ANONYMOUS | MAP_PRIVATE, -1, 0);
	if(ptr == MAP_FAILED) {
		perror("mmap");
		return;
	}
	memcpy(ptr, shellcode, strlen(shellcode));
	( *(void(*) ()) ptr)();
}
*/
import "C"
import "unsafe"

func Run(sc []byte) {
	C.call((*C.char)(unsafe.Pointer(&sc[0])))
}
