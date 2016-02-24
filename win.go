// +build windows

package gopass

import "syscall"
import "unsafe"

var getch = func() (byte, error) {
	modkernel32 := syscall.NewLazyDLL("kernel32.dll")
	procReadFile := modkernel32.NewProc("ReadFile")
	procGetConsoleMode := modkernel32.NewProc("GetConsoleMode")
	procSetConsoleMode := modkernel32.NewProc("SetConsoleMode")

	var mode uint32
	pMode := &mode
	procGetConsoleMode.Call(uintptr(syscall.Stdin), uintptr(unsafe.Pointer(pMode)))

	var echoMode, lineMode, procMode uint32
	echoMode = 4
	lineMode = 2
	procMode = 1
	var newMode uint32
	newMode = mode &^ (echoMode | lineMode | procMode)

	procSetConsoleMode.Call(uintptr(syscall.Stdin), uintptr(newMode))
	defer procSetConsoleMode.Call(uintptr(syscall.Stdin), uintptr(mode))

	line := make([]byte, 1)
	pLine := &line[0]
	var result uintptr
	var err error
	var n uint16
	for n < 1 {
		result, _, err = procReadFile.Call(uintptr(syscall.Stdin),
			uintptr(unsafe.Pointer(pLine)),
			uintptr(len(line)),
			uintptr(unsafe.Pointer(&n)),
			uintptr(unsafe.Pointer(nil)))
	}

	if result > 0 {
		return line[0], nil
	} else {
		return 13, err
	}
}
