// +build windows

package main

import "syscall"
import "unsafe"
import "unicode/utf16"

func GetPass() []byte {
	modkernel32 := syscall.NewLazyDLL("kernel32.dll")
	procReadConsole := modkernel32.NewProc("ReadConsoleW")
	procGetConsoleMode := modkernel32.NewProc("GetConsoleMode")
	procSetConsoleMode := modkernel32.NewProc("SetConsoleMode")

	var mode uint32
	pMode := &mode
	procGetConsoleMode.Call(uintptr(syscall.Stdin), uintptr(unsafe.Pointer(pMode)))

	var echoMode uint32
	echoMode = 4
	var newMode uint32
	newMode = mode ^ echoMode

	procSetConsoleMode.Call(uintptr(syscall.Stdin), uintptr(newMode))

	line := make([]uint16, 300)
	pLine := &line[0]
	var n uint16
	procReadConsole.Call(uintptr(syscall.Stdin), uintptr(unsafe.Pointer(pLine)), uintptr(len(line)), uintptr(unsafe.Pointer(&n)))
	if n > 0 {
		n -= 2
	}

	b := []byte(string(utf16.Decode(line[:n])))

	procSetConsoleMode.Call(uintptr(syscall.Stdin), uintptr(mode))

	return b
}
