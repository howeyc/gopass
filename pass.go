package gopass

import (
	"errors"
	"fmt"
	"io"
	"os"
)

type passwd interface {
	ReadPasswd() error
	ReadPasswdMasked() error
	ReadPasswdPrompt(prompt string, mask bool, r FdReader, w io.Writer) error
	GetPasswd() []byte
	Clean()
}

type FdReader interface {
	io.Reader
	Fd() uintptr
}

type Shadow struct {
	chars []byte
}

var defaultGetCh = func(r io.Reader) (byte, error) {
	buf := make([]byte, 1)
	if n, err := r.Read(buf); n == 0 || err != nil {
		if err != nil {
			return 0, err
		}
		return 0, io.EOF
	}
	return buf[0], nil
}

var (
	maxLength            = 512
	ErrInterrupted       = errors.New("interrupted")
	ErrMaxLengthExceeded = fmt.Errorf("maximum byte limit (%v) exceeded", maxLength)

	// Provide variable so that tests can provide a mock implementation.
	getch = defaultGetCh
)

// getPasswd returns the input read from terminal.
// If prompt is not empty, it will be output as a prompt to the user
// If masked is true, typing will be matched by asterisks on the screen.
// Otherwise, typing will echo nothing.
func getPasswd(prompt string, masked bool, r FdReader, w io.Writer) ([]byte, error) {
	var err error
	var pass, bs, mask []byte
	if masked {
		bs = []byte("\b \b")
		mask = []byte("*")
	}

	if isTerminal(r.Fd()) {
		if oldState, err := makeRaw(r.Fd()); err != nil {
			return pass, err
		} else {
			defer func() {
				restore(r.Fd(), oldState)
				fmt.Fprintln(w)
			}()
		}
	}

	if prompt != "" {
		fmt.Fprint(w, prompt)
	}

	// Track total bytes read, not just bytes in the password.  This ensures any
	// errors that might flood the console with nil or -1 bytes infinitely are
	// capped.
	var counter int
	for counter = 0; counter <= maxLength; counter++ {
		if v, e := getch(r); e != nil {
			err = e
			break
		} else if v == 127 || v == 8 {
			if l := len(pass); l > 0 {
				pass = pass[:l-1]
				fmt.Fprint(w, string(bs))
			}
		} else if v == 13 || v == 10 {
			break
		} else if v == 3 {
			err = ErrInterrupted
			break
		} else if v != 0 {
			pass = append(pass, v)
			fmt.Fprint(w, string(mask))
		}
	}

	if counter > maxLength {
		err = ErrMaxLengthExceeded
	}

	return pass, err
}

// ReadPasswd returns the password read from the terminal without echoing input.
// The returned byte array does not include end-of-line characters.
func (pass *Shadow) ReadPasswd() (err error) {
	pass.chars, err = getPasswd("", false, os.Stdin, os.Stdout)
	return err
}

// ReadPasswdMasked returns the password read from the terminal, echoing asterisks.
// The returned byte array does not include end-of-line characters.
func (pass *Shadow) ReadPasswdMasked() (err error) {
	pass.chars, err = getPasswd("", true, os.Stdin, os.Stdout)
	return err
}

// ReadPasswdPrompt prompts the user and returns the password read from the terminal.
// If mask is true, then asterisks are echoed.
// The returned byte array does not include end-of-line characters.
func (pass *Shadow) ReadPasswdPrompt(prompt string, mask bool, r FdReader, w io.Writer) (err error) {
	pass.chars, err = getPasswd(prompt, mask, r, w)
	return err
}

func (pass *Shadow) GetPasswd() []byte {
	return pass.chars
}

func (pass *Shadow) Clean() {
	for i := range pass.chars {
		pass.chars[i] = 0
	}
}
