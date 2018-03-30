// +build !solaris

package gopass

import "golang.org/x/crypto/ssh/terminal"

type terminalState struct {
	state *terminal.State
}

func isTerminal(fd uintptr) bool {
	return terminal.IsTerminal(int(fd))
}

func makeRaw(fd uintptr) (*terminalState, error) {
	state, err := terminal.GetState(int(fd))
	if err != nil {
		return nil, err
	}
	terminal.MakeRaw(int(fd))

	return &terminalState{
		state: state,
	}, nil
}

func restore(fd uintptr, oldState *terminalState) error {
	return terminal.Restore(int(fd), oldState.state)
}
