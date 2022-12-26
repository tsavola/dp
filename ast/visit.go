// Copyright (c) 2023 Timo Savola
// SPDX-License-Identifier: BSD-3-Clause

package ast

func VisitAssignListChild(x AssignListChild,
	visitAssigner func(AssignerDereference),
	visitCall func(Call),
	visitIndex func(Index),
	visitSelector func(Selector),
) {
	switch x := x.(type) {
	case AssignerDereference:
		visitAssigner(x)
	case Call:
		visitCall(x)
	case Index:
		visitIndex(x)
	case Selector:
		visitSelector(x)
	default:
		panic("unknown assignment list child node type")
	}
}

func VisitBlockChild(x BlockChild,
	visitAssign func(Assign),
	visitBlock func(Block),
	visitBreak func(Break),
	visitComment func(Comment),
	visitContinue func(Continue),
	visitExpression func(Expression),
	visitFor func(For),
	visitIf func(If),
	visitImport func(Import),
	visitReturn func(Return),
	visitVariableDecl func(VariableDecl),
	visitVariableDef func(VariableDef),
) {
	switch x := x.(type) {
	case Assign:
		visitAssign(x)
	case Block:
		visitBlock(x)
	case Break:
		visitBreak(x)
	case Comment:
		visitComment(x)
	case Continue:
		visitContinue(x)
	case Expression:
		visitExpression(x)
	case For:
		visitFor(x)
	case If:
		visitIf(x)
	case Import:
		visitImport(x)
	case Return:
		visitReturn(x)
	case VariableDecl:
		visitVariableDecl(x)
	case VariableDef:
		visitVariableDef(x)
	default:
		panic("unknown block child node type")
	}
}

func VisitExpr(x ExprChild,
	visitAddress func(Address),
	visitBinary func(Binary),
	visitBoolean func(Boolean),
	visitCall func(Call),
	visitCharacter func(Character),
	visitClone func(Clone),
	visitIndex func(Index),
	visitInteger func(Integer),
	visitNil func(Nil),
	visitPointer func(PointerDereference),
	visitSelector func(Selector),
	visitString func(String),
	visitUnary func(Unary),
	visitZero func(Zero),
) {
	switch x := x.(type) {
	case Address:
		visitAddress(x)
	case Binary:
		visitBinary(x)
	case Boolean:
		visitBoolean(x)
	case Call:
		visitCall(x)
	case Character:
		visitCharacter(x)
	case Clone:
		visitClone(x)
	case Index:
		visitIndex(x)
	case Integer:
		visitInteger(x)
	case Nil:
		visitNil(x)
	case PointerDereference:
		visitPointer(x)
	case Selector:
		visitSelector(x)
	case String:
		visitString(x)
	case Unary:
		visitUnary(x)
	case Zero:
		visitZero(x)
	default:
		if true {
			panic(x)
		}
		panic("unknown expression child node type")
	}
}

func VisitExprListChild(x ExprListChild,
	visitAssigner func(AssignerDereference),
	visitComment func(Comment),
	visitExpression func(Expression),
) {
	switch x := x.(type) {
	case AssignerDereference:
		visitAssigner(x)
	case Comment:
		visitComment(x)
	case Expression:
		visitExpression(x)
	default:
		panic("unknown expression list child node type")
	}
}

func VisitFieldListChild(x FieldListChild,
	visitComment func(Comment),
	visitField func(Field),
	visitImport func(Import),
) {
	switch x := x.(type) {
	case Comment:
		visitComment(x)
	case Field:
		visitField(x)
	case Import:
		visitImport(x)
	default:
		panic("unknown field list child node type")
	}
}

func VisitFileChild(x FileChild,
	visitComment func(Comment),
	visitConstant func(ConstantDef),
	visitFunction func(FunctionDef),
	visitImport func(Import),
	visitImports func(Imports),
	visitType func(TypeDef),
) {
	switch x := x.(type) {
	case Comment:
		visitComment(x)
	case ConstantDef:
		visitConstant(x)
	case FunctionDef:
		visitFunction(x)
	case Import:
		visitImport(x)
	case Imports:
		visitImports(x)
	case TypeDef:
		visitType(x)
	default:
		panic("unknown file child node type")
	}
}

func VisitIdentListChild(x IdentListChild,
	visitComment func(Comment),
	visitIdentifier func(Identifier),
) {
	switch x := x.(type) {
	case Comment:
		visitComment(x)
	case Identifier:
		visitIdentifier(x)
	default:
		panic("unknown identifier list child node type")
	}
}

func VisitImportListChild(x ImportListChild,
	visitComment func(Comment),
	visitImport func(Import),
) {
	switch x := x.(type) {
	case Comment:
		visitComment(x)
	case Import:
		visitImport(x)
	default:
		panic("unknown import list child node type")
	}
}

func VisitParamListChild(x ParamListChild,
	visitComment func(Comment),
	visitParameter func(Parameter),
) {
	switch x := x.(type) {
	case Comment:
		visitComment(x)
	case Parameter:
		visitParameter(x)
	default:
		panic("unknown parameter list child node type")
	}
}

func VisitTypeListChild(x TypeListChild,
	visitComment func(Comment),
	visitType func(TypeSpec),
) {
	switch x := x.(type) {
	case Comment:
		visitComment(x)
	case TypeSpec:
		visitType(x)
	default:
		panic("unknown type list child node type")
	}
}
