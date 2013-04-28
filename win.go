// +build windows

package gopass

import "syscall"
import "unsafe"
import "unicode/utf16"

// Returns password byte array read from terminal without input being echoed.
// Array of bytes does not include end-of-line characters.
func getch() byte {
	modkernel32 := syscall.NewLazyDLL("kernel32.dll")
	procReadConsole := modkernel32.NewProc("ReadConsoleW")
	procGetConsoleMode := modkernel32.NewProc("GetConsoleMode")
	procSetConsoleMode := modkernel32.NewProc("SetConsoleMode")

	var mode uint32
	pMode := &mode
	procGetConsoleMode.Call(uintptr(syscall.Stdin), uintptr(unsafe.Pointer(pMode)))

	var echoMode, lineMode uint32
	echoMode = 4
	lineMode = 2
	var newMode uint32
	newMode = mode ^ (echoMode | lineMode)

	procSetConsoleMode.Call(uintptr(syscall.Stdin), uintptr(newMode))

	line := make([]uint16, 1)
	pLine := &line[0]
	var n uint16
	procReadConsole.Call(uintptr(syscall.Stdin), uintptr(unsafe.Pointer(pLine)), uintptr(len(line)), uintptr(unsafe.Pointer(&n)))

	// For some reason n returned seems to big by 2 (Null terminated maybe?)
	if n > 2 {
		n -= 2
	}

	b := []byte(string(utf16.Decode(line[:n])))

	procSetConsoleMode.Call(uintptr(syscall.Stdin), uintptr(mode))

	return b[0]
}
