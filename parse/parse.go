package parse

import "github.com/miniriley2012/brainfuck/token"

func Parse(src []byte) (toks []token.Token) {
	for _, b := range src {
		switch b {
		case '>':
			toks = append(toks, token.INC_DP)
		case '<':
			toks = append(toks, token.DEC_DP)
		case '+':
			toks = append(toks, token.INC_CELL)
		case '-':
			toks = append(toks, token.DEC_CELL)
		case '.':
			toks = append(toks, token.OUTPUT)
		case ',':
			toks = append(toks, token.INPUT)
		case '[':
			toks = append(toks, token.JMP_FWRD)
		case ']':
			toks = append(toks, token.JMP_BACK)
		}
	}
	return
}
