// Copyright (c) 2024 Timo Savola
// SPDX-License-Identifier: BSD-3-Clause

package revise

import (
	new "github.com/tsavola/dp/ast"
	old "github.com/tsavola/dp/ast"
	newfield "github.com/tsavola/dp/field"
	"github.com/tsavola/dp/internal/pan"
	"github.com/tsavola/dp/internal/sync"
	oldlex "github.com/tsavola/dp/lex"
	oldparse "github.com/tsavola/dp/parse"
	"github.com/tsavola/dp/source"

	. "github.com/tsavola/dp/internal/pan/mustcheck"
)

func File(pos source.Position, input string) (news []new.FileChild, err error) {
	err = pan.Recover(func() {
		news = reviseFile(Must(oldparse.File(Must(oldlex.File(pos, input)))))
	})
	return
}

func reviseAssignList(olds []old.AssignListChild) (news []new.AssignListChild) {
	for _, node := range olds {
		old.VisitAssignListChild(node,
			func(node old.AssignerDereference) {
				news = append(news, reviseAssignerDereference(node))
			},

			func(node old.Call) {
				news = append(news, reviseCall(node))
			},

			func(node old.Index) {
				news = append(news, reviseIndex(node))
			},

			func(node old.Selector) {
				news = append(news, reviseSelector(node))
			},
		)
	}
	return
}

func reviseBlock(olds []old.BlockChild) (news []new.BlockChild) {
	for _, node := range olds {
		old.VisitBlockChild(node,
			func(node old.Assign) {
				news = append(news, new.Assign{
					node.Position,
					reviseAssignList(node.Objects),
					reviseExprList(node.Subjects),
					node.End,
				})
			},

			func(node old.Block) {
				news = append(news, new.Block{
					node.Position,
					reviseBlock(node.Body),
					node.End,
				})
			},

			func(node old.Break) {
				news = append(news, new.Break{node.Position, node.End})
			},

			func(node old.Comment) {
				news = append(news, new.Comment{node.Position, node.Source})
			},

			func(node old.Continue) {
				news = append(news, new.Continue{node.Position, node.End})
			},

			func(node old.Expression) {
				news = append(news, new.Expression{reviseExpr(node.Expr)})
			},

			func(node old.For) {
				news = append(news, new.For{
					node.Position,
					reviseExpr(node.Test),
					node.BodyPos,
					reviseBlock(node.Body),
					node.End,
				})
			},

			func(node old.If) {
				news = append(news, new.If{
					node.Position,
					reviseExpr(node.Test),
					node.ThenPos,
					reviseBlock(node.Then),
					node.ThenEnd,
					reviseBlock(node.Else),
					node.End,
				})
			},

			func(node old.Import) {
				news = append(news, reviseImport(node))
			},

			func(node old.Return) {
				news = append(news, new.Return{
					node.Position,
					reviseExprList(node.Values),
					node.End,
				})
			},

			func(node old.VariableDecl) {
				news = append(news, new.VariableDecl{
					node.Position,
					node.Names,
					reviseTypeSpecPtr(node.Type),
					node.End,
				})
			},

			func(node old.VariableDef) {
				news = append(news, new.VariableDef{
					node.Position,
					node.Names,
					reviseExprList(node.Values),
					node.End,
				})
			},
		)
	}
	return
}

func reviseExpr(node old.ExprChild) (result new.ExprChild) {
	old.VisitExpr(node,
		func(node old.Address) {
			result = new.Address{
				node.Position,
				reviseExpr(node.Expr),
				node.End,
			}
		},

		func(node old.Binary) {
			result = new.Binary{
				reviseExpr(node.Left),
				new.BinaryOp(node.Op),
				reviseExpr(node.Right),
				node.End,
			}
		},

		func(node old.Boolean) {
			result = new.Boolean{node.Position, node.Value, node.End}
		},

		func(node old.Call) {
			result = reviseCall(node)
		},

		func(node old.Character) {
			result = new.Character{node.Position, node.Source, node.End}
		},

		func(node old.Clone) {
			result = new.Clone{
				node.Position,
				reviseExpr(node.Expr),
				node.End,
			}
		},

		func(node old.Index) {
			result = reviseIndex(node)
		},

		func(node old.Integer) {
			result = new.Integer{node.Position, node.Source, node.End}
		},

		func(node old.Nil) {
			result = new.Nil{node.Position, node.End}
		},

		func(node old.PointerDereference) {
			result = new.PointerDereference{
				node.Position,
				reviseExpr(node.Expr),
				node.End,
			}
		},

		func(node old.Selector) {
			result = reviseSelector(node)
		},

		func(node old.String) {
			result = new.String{node.Position, node.Source, node.End}
		},

		func(node old.Unary) {
			result = new.Unary{
				node.Position,
				new.UnaryOp(node.Op),
				reviseExpr(node.Expr),
				node.End,
			}
		},

		func(node old.Zero) {
			result = new.Zero{node.Position, node.End}
		},
	)
	return
}

func reviseExprList(olds []old.ExprListChild) (news []new.ExprListChild) {
	for _, node := range olds {
		old.VisitExprListChild(node,
			func(node old.AssignerDereference) {
				news = append(news, reviseAssignerDereference(node))
			},

			func(node old.Comment) {
				news = append(news, new.Comment{node.Position, node.Source})
			},

			func(node old.Expression) {
				news = append(news, new.Expression{reviseExpr(node.Expr)})
			},
		)
	}
	return
}

func reviseFieldList(olds []old.FieldListChild) (news []new.FieldListChild) {
	for _, node := range olds {
		old.VisitFieldListChild(node,
			func(node old.Comment) {
				news = append(news, new.Comment{node.Position, node.Source})
			},

			func(node old.Field) {
				news = append(news, reviseField(node))
			},

			func(node old.Import) {
				news = append(news, reviseImport(node))
			},
		)
	}
	return
}

func reviseFile(olds []old.FileChild) (news []new.FileChild) {
	for _, node := range olds {
		old.VisitFileChild(node,
			func(node old.Comment) {
				news = append(news, new.Comment{node.Position, node.Source})
			},

			func(node old.ConstantDef) {
				news = append(news, new.ConstantDef{
					node.Position,
					node.Public,
					node.ConstName,
					reviseExpr(node.Value),
					node.End,
				})
			},

			func(node old.FunctionDef) {
				news = append(news, new.FunctionDef{
					node.Position,
					node.Public,
					node.ReceiverName,
					reviseTypeSpecPtr(node.ReceiverType),
					node.FuncName,
					reviseParamList(node.Params),
					node.ParamsEnd,
					reviseTypeList(node.Results),
					node.BodyPos,
					reviseBlock(node.Body),
					node.End,
				})
			},

			func(node old.Import) {
				news = append(news, reviseImport(node))
			},

			func(node old.Imports) {
				news = append(news, new.Imports{
					node.Position,
					reviseImportList(node.Imports),
					node.End,
				})
			},

			func(node old.TypeDef) {
				news = append(news, new.TypeDef{
					node.Position,
					node.Public,
					node.TypeName,
					reviseFieldList(node.Fields),
					node.End,
				})
			},
		)
	}
	return
}

func reviseIdentList(olds []old.IdentListChild) (news []new.IdentListChild) {
	for _, node := range olds {
		old.VisitIdentListChild(node,
			func(node old.Comment) {
				news = append(news, new.Comment{node.Position, node.Source})
			},

			func(node old.Identifier) {
				news = append(news, new.Identifier{
					node.Position,
					new.QualifiedName(node.Name),
					node.End,
				})
			},
		)
	}
	return
}

func reviseImportList(olds []old.ImportListChild) (news []new.ImportListChild) {
	for _, node := range olds {
		old.VisitImportListChild(node,
			func(node old.Comment) {
				news = append(news, new.Comment{node.Position, node.Source})
			},

			func(node old.Import) {
				news = append(news, reviseImport(node))
			},
		)
	}
	return
}

func reviseParamList(olds []old.ParamListChild) (news []new.ParamListChild) {
	for _, node := range olds {
		old.VisitParamListChild(node,
			func(node old.Comment) {
				news = append(news, new.Comment{node.Position, node.Source})
			},

			func(node old.Parameter) {
				news = append(news, new.Parameter{
					node.Position,
					node.ParamName,
					reviseTypeSpec(node.Type),
					node.End,
				})
			},
		)
	}
	return
}

func reviseTypeList(olds []old.TypeListChild) (news []new.TypeListChild) {
	for _, node := range olds {
		old.VisitTypeListChild(node,
			func(node old.Comment) {
				news = append(news, new.Comment{node.Position, node.Source})
			},

			func(node old.TypeSpec) {
				news = append(news, reviseTypeSpec(node))
			},
		)
	}
	return
}

func reviseField(node old.Field) new.Field {
	return new.Field{
		node.Position,
		node.FieldName,
		reviseTypeSpec(node.Type),
		newfield.Access(node.Access),
		node.End,
	}
}

func reviseImport(node old.Import) new.Import {
	return new.Import{
		node.Position,
		node.Path,
		reviseIdentList(node.Names),
		node.End,
	}
}

func reviseAssignerDereference(node old.AssignerDereference) new.AssignerDereference {
	return new.AssignerDereference{
		node.Position,
		node.Name,
		node.End,
	}
}

func reviseCall(node old.Call) new.Call {
	return new.Call{
		reviseSelector(node.Name),
		reviseExprList(node.Args),
		node.End,
	}
}

func reviseIndex(node old.Index) new.Index {
	return new.Index{
		reviseSelector(node.Name),
		reviseExpr(node.Index),
		node.End,
	}
}

func reviseSelector(node old.Selector) new.Selector {
	return new.Selector{
		node.Position,
		node.Name,
		node.End,
	}
}

func reviseType(t old.Type) new.Type {
	return new.Type{
		t.Assigner,
		t.Pointer,
		t.Reference,
		t.Shared,
		reviseTypePtr(t.Item),
		new.QualifiedName(t.Name),
	}
}

var typePtrs sync.Map[*old.Type, *new.Type]

func reviseTypePtr(p *old.Type) *new.Type {
	if p == nil {
		return nil
	}
	v := reviseType(*p)
	return typePtrs.LoadOrStore(p, &v)
}

func reviseTypeSpec(node old.TypeSpec) new.TypeSpec {
	return new.TypeSpec{
		node.Position,
		reviseType(node.Type),
		node.End,
	}
}

var typeSpecPtrs sync.Map[*old.TypeSpec, *new.TypeSpec]

func reviseTypeSpecPtr(p *old.TypeSpec) *new.TypeSpec {
	if p == nil {
		return nil
	}
	v := reviseTypeSpec(*p)
	return typeSpecPtrs.LoadOrStore(p, &v)
}
