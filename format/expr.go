// Copyright (c) 2023 Timo Savola
// SPDX-License-Identifier: BSD-3-Clause

package format

import (
	"github.com/tsavola/dp/ast"
)

func formatExpr(w writer, level int, node ast.ExprChild, parentPrec int, tight bool) {
	ast.VisitExpr(node,
		func(node ast.Address) {
			if w.lastRune() == '&' { // Prevent &&
				w.WriteString(" ")
			}
			w.WriteString("&")
			formatExpr(w, level, node.Expr, ast.UltimatePrecedence, tight)
		},

		func(node ast.Binary) {
			prec := node.Op.Precedence()
			if parentPrec > 0 && prec > parentPrec && prec == ast.MaxBinaryPrecedence {
				tight = true
			}
			if parentPrec > 0 && prec != parentPrec {
				w.WriteString("(")
			}

			formatExpr(w, level, node.Left, prec, tight)

			if !tight {
				w.WriteString(" ")
			}
			w.WriteString(node.Op.String())
			if !tight {
				w.WriteString(" ")
			}

			formatExpr(w, level, node.Right, prec, tight)

			if parentPrec > 0 && prec != parentPrec {
				w.WriteString(")")
			}
		},

		func(node ast.Boolean) {
			if node.Value {
				w.WriteString("true")
			} else {
				w.WriteString("false")
			}
		},

		func(node ast.Call) {
			w.WriteString(node.Name.String())
			formatExprList(w, level, node.Pos().Line, node.Args, true)
		},

		func(node ast.Character) {
			w.WriteString(node.Source)
		},

		func(node ast.Clone) {
			w.WriteString("clone ")
			formatExpr(w, level, node.Expr, ast.UltimatePrecedence, tight)
		},

		func(node ast.Index) {
			w.WriteString(node.Name.String())
			w.WriteString("[")
			formatExpr(w, level, node.Index, 0, false)
			w.WriteString("]")
		},

		func(node ast.Integer) {
			w.WriteString(node.Source)
		},

		func(node ast.Nil) {
			w.WriteString("nil")
		},

		func(node ast.PointerDereference) {
			w.WriteString("*")
			formatExpr(w, level, node.Expr, ast.UltimatePrecedence, tight)
		},

		func(node ast.Selector) {
			w.WriteString(node.String())
		},

		func(node ast.String) {
			w.WriteString(node.Source)
		},

		func(node ast.Unary) {
			if node.Op == ast.OpIdentity {
				formatExpr(w, level, node.Expr, parentPrec, tight)
			} else {
				if node.Op == ast.OpComplement && w.lastRune() == '&' { // Prevent &^
					w.WriteString(" ")
				}
				w.WriteString(node.Op.String())
				formatExpr(w, level, node.Expr, ast.UltimatePrecedence, tight)
			}
		},

		func(node ast.Zero) {
			w.WriteString("{}")
		},
	)
}

func formatAssignerDereference(w writer, node ast.AssignerDereference) {
	w.WriteString("(")
	w.WriteString(node.Name)
	w.WriteString(")")
}
