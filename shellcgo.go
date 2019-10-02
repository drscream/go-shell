//+build cgo,darwin cgo,linux

package shell

// #include <errno.h>
// #include <pwd.h>
// #include <stdlib.h>
// #include <sys/types.h>
// #include <unistd.h>
import "C"

import (
	"unsafe"
)

func cgoGetUserShell(name string) (string, bool) {
	buflen := C.sysconf(C._SC_GETPW_R_SIZE_MAX)
	if buflen == -1 {
		buflen = 1024
	}
	for {
		var (
			cName  = C.CString(name)
			pwd    C.struct_passwd
			buf    = make([]byte, buflen)
			result *C.struct_passwd
		)
		//nolint:gocritic,staticcheck
		rc := C.getpwnam_r(cName, &pwd, (*C.char)(unsafe.Pointer(&buf[0])), C.ulong(buflen), &result)
		C.free(unsafe.Pointer(cName))
		switch rc {
		case 0:
			return C.GoString(result.pw_shell), true
		case C.ERANGE:
			buflen *= 2
		default:
			return "", false
		}
	}
}
