package token

type Token int

const (
	INC_DP Token = iota
	DEC_DP
	INC_CELL
	DEC_CELL
	OUTPUT
	INPUT
	JMP_FWRD
	JMP_BACK
)

func (tok Token) Byte() byte {
	switch tok {
	case INC_DP:
		return '>'
	case DEC_DP:
		return '<'
	case INC_CELL:
		return '+'
	case DEC_CELL:
		return '-'
	case OUTPUT:
		return '.'
	case INPUT:
		return ','
	case JMP_FWRD:
		return '['
	case JMP_BACK:
		return ']'
	}
	return 0
}
