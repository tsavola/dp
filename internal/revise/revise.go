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

			func(node old.Cast) {
				news = append(news, reviseCast(node))
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
					node.At,
					reviseAssignList(node.Targets),
					reviseExprList(node.Values),
					node.EndAt,
				})
			},

			func(node old.Block) {
				news = append(news, new.Block{
					node.At,
					reviseBlock(node.Body),
					node.EndAt,
				})
			},

			func(node old.Break) {
				news = append(news, new.Break{node.At, node.EndAt})
			},

			func(node old.Comment) {
				news = append(news, new.Comment{node.At, node.Source})
			},

			func(node old.Continue) {
				news = append(news, new.Continue{node.At, node.EndAt})
			},

			func(node old.Expression) {
				news = append(news, new.Expression{reviseExpr(node.Expr)})
			},

			func(node old.For) {
				news = append(news, new.For{
					node.At,
					reviseExpr(node.Test),
					node.BodyAt,
					reviseBlock(node.Body),
					node.EndAt,
				})
			},

			func(node old.If) {
				news = append(news, new.If{
					node.At,
					reviseExpr(node.Test),
					node.ThenAt,
					reviseBlock(node.Then),
					node.ThenEndAt,
					reviseBlock(node.Else),
					node.EndAt,
				})
			},

			func(node old.Import) {
				news = append(news, reviseImport(node))
			},

			func(node old.Return) {
				news = append(news, new.Return{
					node.At,
					reviseExprList(node.Values),
					node.EndAt,
				})
			},

			func(node old.VariableDecl) {
				news = append(news, new.VariableDecl{
					node.At,
					node.Names,
					reviseTypeSpecPtr(node.Type),
					node.EndAt,
				})
			},

			func(node old.VariableDef) {
				news = append(news, new.VariableDef{
					node.At,
					node.Names,
					reviseExprList(node.Values),
					node.EndAt,
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
				node.At,
				reviseExpr(node.Expr),
				node.EndAt,
			}
		},

		func(node old.Binary) {
			result = new.Binary{
				reviseExpr(node.Left),
				new.BinaryOp(node.Op),
				reviseExpr(node.Right),
				node.EndAt,
			}
		},

		func(node old.Boolean) {
			result = new.Boolean{node.At, node.Source}
		},

		func(node old.Call) {
			result = reviseCall(node)
		},

		func(node old.Cast) {
			result = reviseCast(node)
		},

		func(node old.Character) {
			result = new.Character{node.At, node.Source}
		},

		func(node old.Clone) {
			result = new.Clone{
				node.At,
				reviseExpr(node.Expr),
				node.EndAt,
			}
		},

		func(node old.Empty) {
			result = new.Empty{node.At, node.EndAt}
		},

		func(node old.Index) {
			result = reviseIndex(node)
		},

		func(node old.Integer) {
			result = new.Integer{node.At, node.Source}
		},

		func(node old.Nil) {
			result = new.Nil{node.At, node.EndAt}
		},

		func(node old.PointerDereference) {
			result = new.PointerDereference{
				node.At,
				reviseExpr(node.Expr),
				node.EndAt,
			}
		},

		func(node old.Selector) {
			result = reviseSelector(node)
		},

		func(node old.String) {
			result = new.String{node.At, node.Source}
		},

		func(node old.Unary) {
			result = new.Unary{
				node.At,
				new.UnaryOp(node.Op),
				reviseExpr(node.Expr),
				node.EndAt,
			}
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
				news = append(news, new.Comment{node.At, node.Source})
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
				news = append(news, new.Comment{node.At, node.Source})
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
				news = append(news, new.Comment{node.At, node.Source})
			},

			func(node old.ConstantDef) {
				news = append(news, new.ConstantDef{
					node.At,
					node.Public,
					node.ConstName,
					reviseExpr(node.Value),
					node.EndAt,
				})
			},

			func(node old.FunctionDef) {
				news = append(news, new.FunctionDef{
					node.At,
					node.Public,
					node.ReceiverName,
					reviseTypeSpecPtr(node.ReceiverType),
					node.FuncName,
					reviseParamList(node.Params),
					node.ParamsEndAt,
					reviseTypeList(node.Results),
					node.BodyAt,
					reviseBlock(node.Body),
					node.EndAt,
				})
			},

			func(node old.Import) {
				news = append(news, reviseImport(node))
			},

			func(node old.Imports) {
				news = append(news, new.Imports{
					node.At,
					reviseImportList(node.Imports),
					node.EndAt,
				})
			},

			func(node old.TypeDef) {
				news = append(news, new.TypeDef{
					node.At,
					node.Public,
					node.TypeName,
					reviseFieldList(node.Fields),
					node.EndAt,
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
				news = append(news, new.Comment{node.At, node.Source})
			},

			func(node old.Identifier) {
				news = append(news, new.Identifier{
					node.At,
					new.QualifiedName(node.Name),
					node.EndAt,
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
				news = append(news, new.Comment{node.At, node.Source})
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
				news = append(news, new.Comment{node.At, node.Source})
			},

			func(node old.Parameter) {
				news = append(news, new.Parameter{
					node.At,
					node.ParamName,
					reviseTypeSpec(node.Type),
					node.EndAt,
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
				news = append(news, new.Comment{node.At, node.Source})
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
		node.At,
		node.FieldName,
		reviseTypeSpec(node.Type),
		newfield.Access(node.Access),
		node.EndAt,
	}
}

func reviseImport(node old.Import) new.Import {
	return new.Import{
		node.At,
		node.Path,
		reviseIdentList(node.Names),
		node.EndAt,
	}
}

func reviseAssignerDereference(node old.AssignerDereference) new.AssignerDereference {
	return new.AssignerDereference{
		node.At,
		node.Name,
		node.EndAt,
	}
}

func reviseCall(node old.Call) new.Call {
	return new.Call{
		reviseSelector(node.Name),
		reviseExprList(node.Args),
		node.EndAt,
	}
}

func reviseCast(node old.Cast) new.Cast {
	return new.Cast{
		node.At,
		node.Name,
		reviseExpr(node.Expr),
		node.EndAt,
	}
}

func reviseIndex(node old.Index) new.Index {
	return new.Index{
		reviseSelector(node.Name),
		reviseExpr(node.Index),
		node.EndAt,
	}
}

func reviseSelector(node old.Selector) new.Selector {
	return new.Selector{
		node.At,
		node.Name,
		node.EndAt,
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
		node.At,
		reviseType(node.Type),
		node.EndAt,
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
