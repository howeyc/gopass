package gopass

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"
)

// TestGetPasswd tests the password creation and output based on a byte buffer
// as input to mock the underlying getch() methods.
func TestGetPasswd(t *testing.T) {
	type testData struct {
		input []byte

		// Due to how backspaces are written, it is easier to manually write
		// each expected output for the masked cases.
		masked   string
		password string
		byesLeft int
		reason   string
	}

	ds := []testData{
		testData{[]byte("abc\n"), "***\n", "abc", 0, "Password parsing should stop at \\n"},
		testData{[]byte("abc\r"), "***\n", "abc", 0, "Password parsing should stop at \\r"},
		testData{[]byte("a\nbc\n"), "*\n", "a", 3, "Password parsing should stop at \\n"},
		testData{[]byte("*!]|\n"), "****\n", "*!]|", 0, "Special characters shouldn't affect the password."},

		testData{[]byte("abc\r\n"), "***\n", "abc", 1,
			"Password parsing should stop at \\r; Windows LINE_MODE should be unset so \\r is not converted to \\r\\n."},

		testData{[]byte{'a', 'b', 'c', 8, '\n'}, "***\b \b\n", "ab", 0, "Backspace byte should remove the last read byte."},
		testData{[]byte{'a', 'b', 127, 'c', '\n'}, "**\b \b*\n", "ac", 0, "Delete byte should remove the last read byte."},
		testData{[]byte{'a', 'b', 127, 'c', 8, 127, '\n'}, "**\b \b*\b \b\b \b\n", "", 0, "Successive deletes continue to delete."},
		testData{[]byte{8, 8, 8, '\n'}, "\n", "", 0, "Deletes before characters are noops."},
		testData{[]byte{8, 8, 8, 'a', 'b', 'c', '\n'}, "***\n", "abc", 0, "Deletes before characters are noops."},

		testData{[]byte{'a', 'b', 0, 'c', '\n'}, "***\n", "abc", 0,
			"Nil byte should be ignored due; may get unintended nil bytes from syscalls on Windows."},
	}

	// getch methods normally refer to syscalls; replace with a byte buffer.
	var input *bytes.Buffer
	getch = func() (byte, error) {
		b, err := input.ReadByte()
		if err != nil {
			t.Fatal(err.Error())
		}
		return b, nil
	}

	// Redirecting output for tests as they print to os.Stdout but we want to
	// capture and test the output.
	origStdOut := os.Stdout
	for _, masked := range []bool{true, false} {
		for _, d := range ds {
			input = bytes.NewBuffer(d.input)

			r, w, err := os.Pipe()
			if err != nil {
				t.Fatal(err.Error())
			}
			os.Stdout = w

			result, err := getPasswd(masked)
			os.Stdout = origStdOut
			if err != nil {
				t.Errorf("Error getting password:", err.Error())
			}

			// Test output (masked and unmasked).  Delete/backspace actually
			// deletes, overwrites and deletes again.  As a result, we need to
			// remove those from the pipe afterwards to mimic the console's
			// interpretation of those bytes.
			w.Close()
			output, err := ioutil.ReadAll(r)
			if err != nil {
				t.Fatal(err.Error())
			}
			var expectedOutput []byte
			if masked {
				expectedOutput = []byte(d.masked)
			} else {
				expectedOutput = []byte("\n")
			}
			if bytes.Compare(expectedOutput, output) != 0 {
				t.Errorf("Expected output to equal %v (%q) but got %v (%q) instead when masked=%v. %s", expectedOutput, string(expectedOutput), output, string(output), masked, d.reason)
			}

			if string(result) != d.password {
				t.Errorf("Expected %q but got %q instead when masked=%v. %s", d.password, result, masked, d.reason)
			}

			if input.Len() != d.byesLeft {
				t.Errorf("Expected %v bytes left on buffer but instead got %v when masked=%v. %s", d.byesLeft, input.Len(), masked, d.reason)
			}
		}
	}
}
