# getpasswd in Go [![GoDoc](https://godoc.org/github.com/howeyc/gopass?status.svg)](https://godoc.org/github.com/howeyc/gopass) [![Build Status](https://secure.travis-ci.org/howeyc/gopass.png?branch=master)](http://travis-ci.org/howeyc/gopass)

Retrieve password from user terminal or piped input without echo.

Verified on BSD, Linux, and Windows.

Example:
```go
package main

import "fmt"
import "github.com/howeyc/gopass"

func main() {
  var passwd gopass.Shadow
	fmt.Printf("Password: ")

	// Silent. For printing *'s use gopass.GetPasswdMasked()
	pass, err := passwd.ReadPasswd(), passwd.GetPasswd()
	if err != nil {
		// Handle gopass.ErrInterrupted or getch() read error
	}

	// Do something with pass

  passwd.Clean()
}
```

Caution: Multi-byte characters not supported!
