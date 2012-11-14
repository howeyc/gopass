// +build windows

package main

/*
#include <windows.h>
#include <stdio.h>

unsigned char getch() {
        unsigned char ch;
        DWORD con_mode;
        DWORD dwRead;

        HANDLE hIn=GetStdHandle(STD_INPUT_HANDLE);

        GetConsoleMode( hIn, &con_mode );
        SetConsoleMode( hIn, con_mode & ~(ENABLE_ECHO_INPUT | ENABLE_LINE_INPUT) );

        ReadConsoleA( hIn, &ch, 1, &dwRead, NULL)

        SetConsoleMode( hIn, con_mode);
        return ch;
}
*/
import "C"

func getch() byte {
	return byte(C.getch())
}
