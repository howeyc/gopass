package main

import "fmt"

func GetPass() []byte {
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
	return pass
}

func main() {
	fmt.Printf("Password: ")
	pass := GetPass()
	fmt.Println(pass)
}
