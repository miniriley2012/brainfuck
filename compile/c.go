package compile

import (
	"github.com/miniriley2012/brainfuck/parse"
	"github.com/miniriley2012/brainfuck/token"
	"strconv"
	"strings"
)

type CC struct{}

var C CC

func (c CC) Compile(src []byte) ([]byte, error) {
	toks := parse.Parse(src)
	output := `#include <stdio.h>

int main() {
	char arr[` + strconv.Itoa(maxCells) + `] = {0};
	char *ptr = arr;`
	indents := 1
	var placeColon bool
	for i := range toks {
		if placeColon {
			output += ";"
		} else {
			placeColon = true
		}
		output += "\n" + strings.Repeat("\t", indents)
		switch toks[i] {
		case token.INC_DP:
			output += "++ptr"
		case token.DEC_DP:
			output += "--ptr"
		case token.INC_CELL:
			output += "++*ptr"
		case token.DEC_CELL:
			output += "--*ptr"
		case token.OUTPUT:
			output += "putchar(*ptr)"
		case token.INPUT:
			output += "*ptr = getchar();"
		case token.JMP_FWRD:
			output += "while (*ptr) {"
			placeColon = false
			indents++
		case token.JMP_BACK:
			output = output[:len(output)-1]
			output += "}"
			placeColon = false
			indents--
		}
	}
	output += ";\n}"
	return []byte(output), nil
}
