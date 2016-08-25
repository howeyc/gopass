package gopass

import (
	"syscall"

	"golang.org/x/sys/unix"
)

type terminalState struct {
	state *unix.Termios
}

// isTerminal returns true if there is a terminal attached to the given
// file descriptor.
// Source: http://src.illumos.org/source/xref/illumos-gate/usr/src/lib/libbc/libc/gen/common/isatty.c
func isTerminal(fd uintptr) bool {
	var termio unix.Termio
	err := unix.IoctlSetTermio(int(fd), unix.TCGETA, &termio)
	return err == nil
}

// makeRaw puts the terminal connected to the given file descriptor into raw
// mode and returns the previous state of the terminal so that it can be
// restored.
// Source: http://src.illumos.org/source/xref/illumos-gate/usr/src/lib/libast/common/uwin/getpass.c
func makeRaw(fd uintptr) (*terminalState, error) {
	oldTermiosPtr, err := unix.IoctlGetTermios(int(fd), unix.TCGETS)
	if err != nil {
		return nil, err
	}
	oldTermios := *oldTermiosPtr

	newTermios := oldTermios
	newTermios.Lflag &^= syscall.ECHO | syscall.ECHOE | syscall.ECHOK | syscall.ECHONL
	if err := unix.IoctlSetTermios(int(fd), unix.TCSETS, &newTermios); err != nil {
		return nil, err
	}

	return &terminalState{
		state: oldTermiosPtr,
	}, nil
}

func restore(fd uintptr, oldState *terminalState) error {
	return unix.IoctlSetTermios(int(fd), unix.TCSETS, oldState.state)
}
