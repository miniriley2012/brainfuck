package compile

import (
	"bytes"
	"github.com/miniriley2012/brainfuck/parse"
	bfToken "github.com/miniriley2012/brainfuck/token"
	"go/ast"
	"go/printer"
	"go/token"
	"strconv"
)

type GC struct{}

var Go GC

var idents = map[string]*ast.Ident{
	"f":       ast.NewIdent("f"),
	"c":       ast.NewIdent("c"),
	"p":       ast.NewIdent("p"),
	"b":       ast.NewIdent("b"),
	"discard": ast.NewIdent("_"),
}

var exprs = map[string]ast.Expr{
	"index": &ast.IndexExpr{
		X: &ast.BasicLit{
			Kind:  token.IDENT,
			Value: "c",
		},
		Index: &ast.BasicLit{
			Kind:  token.IDENT,
			Value: "p",
		},
	},
}

func fCall(n int) *ast.CallExpr {
	return &ast.CallExpr{
		Fun: idents["f"],
		Args: []ast.Expr{
			&ast.BasicLit{
				ValuePos: 0,
				Kind:     token.INT,
				Value:    strconv.Itoa(n),
			},
		},
	}
}

func incDecFunc(X ast.Expr, n int, useF bool) ast.Stmt {
	if n == 1 || n == -1 {
		var tok token.Token
		if n == 1 {
			tok = token.INC
		} else {
			tok = token.DEC
		}
		return &ast.IncDecStmt{
			X:   X,
			Tok: tok,
		}
	}

	var value ast.Expr
	if n > 0 || !useF {
		value = &ast.BasicLit{Kind: token.INT, Value: strconv.Itoa(n)}
	} else {
		value = fCall(n)
	}

	return &ast.AssignStmt{
		Lhs: []ast.Expr{X},
		Tok: token.ADD_ASSIGN,
		Rhs: []ast.Expr{value},
	}
}

func (GC) Compile(src []byte) ([]byte, error) {
	toks := parse.Parse(src)
	output := &ast.ExprStmt{X: &ast.CallExpr{
		Fun: ast.NewIdent("print"),
		Args: []ast.Expr{
			&ast.CallExpr{
				Fun: ast.NewIdent("string"),
				Args: []ast.Expr{
					exprs["index"],
				},
			},
		},
	}}
	input := []ast.Stmt{
		&ast.ExprStmt{
			X: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X: &ast.SelectorExpr{
						X:   ast.NewIdent("os"),
						Sel: ast.NewIdent("Stdin"),
					},
					Sel: ast.NewIdent("Read"),
				},
				Args: []ast.Expr{
					ast.NewIdent("b"),
				},
			},
		},
		&ast.AssignStmt{
			Lhs: []ast.Expr{
				exprs["index"],
			},
			Tok: token.ASSIGN,
			Rhs: []ast.Expr{
				&ast.IndexExpr{
					X: ast.NewIdent("b"),
					Index: &ast.BasicLit{
						Kind:  token.INT,
						Value: "0",
					},
				},
			},
		},
	}
	f := &ast.FuncDecl{
		Name: ast.NewIdent("f"),
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{{
					Names: []*ast.Ident{ast.NewIdent("i")},
					Type:  ast.NewIdent("int"),
				}},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{{
					Type: ast.NewIdent("byte"),
				}},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.CallExpr{
							Fun: ast.NewIdent("byte"),
							Args: []ast.Expr{
								ast.NewIdent("i"),
							},
						},
					},
				},
			},
		},
	}
	main := &ast.FuncDecl{
		Name: ast.NewIdent("main"),
		Type: &ast.FuncType{},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						idents["c"],
						idents["p"],
						idents["b"],
					},
					TokPos: 0,
					Tok:    token.DEFINE,
					Rhs: []ast.Expr{
						&ast.ArrayType{
							Len: ast.NewIdent("30000"),
							Elt: &ast.BasicLit{
								Kind:  token.TYPE,
								Value: "byte{}",
							},
						},
						&ast.BasicLit{
							Kind:  token.INT,
							Value: "0",
						},
						&ast.CallExpr{
							Fun: ast.NewIdent("make"),
							Args: []ast.Expr{
								&ast.ArrayType{
									Elt: ast.NewIdent("byte"),
								},
								&ast.BasicLit{
									Kind:  token.INT,
									Value: "1",
								},
							},
						},
					},
				},
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						idents["discard"],
						idents["discard"],
						idents["discard"],
					},
					Tok: token.ASSIGN,
					Rhs: []ast.Expr{
						idents["c"],
						idents["p"],
						idents["b"],
					},
				},
			},
		},
	}
	type Scope struct {
		body   *ast.BlockStmt
		parent *Scope
	}
	scope := &Scope{
		main.Body,
		nil,
	}
	var p int
	check := func(n int) bool {
		p += n
		return !(p < 0 || p > maxCells)
	}
	for i := 0; i < len(toks); i++ {
		switch toks[i] {
		case bfToken.INC_DP:
			j := i
			for toks[i] == bfToken.INC_DP {
				i++
			}
			if !check(i - j) {
				return nil, NewError(i, ErrDPOutOfBounds)
			}
			scope.body.List = append(scope.body.List, incDecFunc(idents["p"], i-j, false))
			i--
		case bfToken.DEC_DP:
			j := i
			for toks[i] == bfToken.DEC_DP {
				i++
			}
			if !check(i - j) {
				return nil, NewError(i, ErrDPOutOfBounds)
			}
			scope.body.List = append(scope.body.List, incDecFunc(idents["p"], -(i-j), false))
			i--
		case bfToken.INC_CELL:
			j := i
			for toks[i] == bfToken.INC_CELL {
				i++
			}
			if !check(i - j) {
				return nil, NewError(i, ErrDPOutOfBounds)
			}
			scope.body.List = append(scope.body.List, incDecFunc(exprs["index"], i-j, true))
			i--
		case bfToken.DEC_CELL:
			j := i
			for toks[i] == bfToken.DEC_CELL {
				i++
			}
			if !check(i - j) {
				return nil, NewError(i, ErrDPOutOfBounds)
			}
			scope.body.List = append(scope.body.List, incDecFunc(exprs["index"], -(i-j), true))
			i--
		case bfToken.OUTPUT:
			scope.body.List = append(scope.body.List, output)
		case bfToken.INPUT:
			scope.body.List = append(scope.body.List, input...)
		case bfToken.JMP_FWRD:
			scope = &Scope{
				body:   &ast.BlockStmt{List: []ast.Stmt{}},
				parent: scope,
			}
		case bfToken.JMP_BACK:
			body := scope.body
			scope = scope.parent
			if scope == nil {
				return nil, NewError(i, ErrRParen)
			}
			scope.body.List = append(scope.body.List, &ast.ForStmt{
				Cond: &ast.BinaryExpr{
					X:  exprs["index"],
					Op: token.NEQ,
					Y: &ast.BasicLit{
						ValuePos: 0,
						Kind:     token.INT,
						Value:    "0",
					},
				},
				Body: body,
			})
		}
	}
	var buf bytes.Buffer
	_ = printer.Fprint(&buf, token.NewFileSet(), &ast.File{
		Name:  ast.NewIdent("main"),
		Decls: []ast.Decl{f, main},
	})
	return buf.Bytes(), nil
}
