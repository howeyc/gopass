// +build freebsd openbsd netbsd darwin linux

package gopass

/*
#include <termios.h>
#include <unistd.h>
#include <stdio.h>

int getch() {
        int ch;
        struct termios t_old, t_new;

        tcgetattr(STDIN_FILENO, &t_old);
        t_new = t_old;
        t_new.c_lflag &= ~(ICANON | ECHO);
        tcsetattr(STDIN_FILENO, TCSANOW, &t_new);

        ch = getchar();

        tcsetattr(STDIN_FILENO, TCSANOW, &t_old);
        return ch;
}
*/
import "C"
import "fmt"
import "os"

var secret = make([]byte, 0)
var mask = byte('*')

func getch() (password byte) {
	password = byte(C.getch())

	// Includes spaces A-Z a-z and special characters
	// Delete * if backspace pressed
	// TODO: Not allow spaces?
	if password >= 32 && password <= 126 {
		secret = append(secret, mask)
		fmt.Print(string(mask))
	} else if password == 127 || password == 8 {
		if len(secret) > 0 {
			secret = secret[:len(secret)-1]
			os.Stdin.Write([]byte("\b \b"))
		}
	}
	return
}

// Returns password byte array read from terminal without input being echoed.
// Array of bytes does not include end-of-line characters.
func GetPasswd() []byte {
	pass := make([]byte, 0)
	for v := getch(); ; v = getch() {
		if v == 127 || v == 8 {
			if len(pass) > 0 {
				pass = pass[:len(pass)-1]
			}
		} else if v == 13 || v == 10 {
			break
		} else {
			pass = append(pass, v)
		}
	}
	println()
	return pass
}
