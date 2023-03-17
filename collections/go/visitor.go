package js_to_go

import (
	"fmt"
	"reflect"

	parser "github.com/NekoMaru76/JSCompiler/grammars/lib"
	"github.com/antlr/antlr4/runtime/Go/antlr/v4"
)

type Visitor struct {
	*antlr.BaseParseTreeVisitor
	parser.BaseJavaScriptParserVisitor
}

type Code struct {
	Top  string
	Init string
	Main string
}

func (c *Code) ToFileString(name string) string {
	return "package " + name + "\n\nimport js_to_go \"github.com/NekoMaru76/JSCompiler/collections/go\"" + c.Init + "\n\n" + c.Top + "\n\nfunc main() {\nscopeCtx := js_to_go.Context{}\nfuncCtx := js_to_go.Context{}\n" + c.Main + "\n}"
}

func (v *Visitor) VisitChild(node antlr.RuleNode, i int) Code {
	return v.Visit(node.GetChild(i).(antlr.RuleNode))
}

func (v *Visitor) VisitChildren(node antlr.RuleNode) Code {
	code := Code{}

	for _, child := range node.GetChildren() {
		lineCode := v.Visit(child.(antlr.RuleNode))
		code.Init += lineCode.Init
		code.Main += lineCode.Main
		code.Top += lineCode.Top
	}

	return code
}

func (v *Visitor) Visit(node antlr.RuleNode) Code {
	fmt.Println(reflect.TypeOf(node).String())

	switch ctx := node.(type) {
	case *parser.ProgramContext:
		return v.VisitProgram(ctx)
	case *parser.SourceElementsContext:
		return v.VisitSourceElements(ctx)
	case *parser.SourceElementContext:
		return v.VisitSourceElement(ctx)
	case *parser.NumericLiteralContext:
		return v.VisitNumericLiteral(ctx)
	case *parser.StatementContext:
		return v.VisitStatement(ctx)
	case *parser.ExpressionStatementContext:
		return v.VisitExpressionStatement(ctx)
	case *parser.ExpressionSequenceContext:
		return v.VisitExpressionSequence(ctx)
	case *parser.LiteralExpressionContext:
		return v.VisitLiteralExpression(ctx)
	case *parser.LiteralContext:
		return v.VisitLiteral(ctx)
	case *parser.VariableStatementContext:
		return v.VisitVariableStatement(ctx)
	case *parser.VariableDeclarationListContext:
		return v.VisitVariableDeclarationList(ctx)
	case *parser.VariableDeclarationContext:
		return v.VisitVariableDeclaration(ctx)
	}

	fmt.Println(node.GetText())
	panic("UNKNOWN")
}

func (v *Visitor) VisitProgram(ctx *parser.ProgramContext) Code {
	val := Code{}
	i := 0
	len := ctx.GetChildCount()

	for i < len {
		child := ctx.GetChild(i)

		i++

		switch t := child.(type) {
		case *antlr.TerminalNodeImpl:
			continue
		case antlr.RuleNode:
			val = v.Visit(t)
		}
	}

	return val
}

func (v *Visitor) VisitSourceElement(ctx *parser.SourceElementContext) Code {
	return v.VisitChild(ctx, 0)
}

func (v *Visitor) VisitStatement(ctx *parser.StatementContext) Code {
	return v.VisitChild(ctx, 0)
}

func (v *Visitor) VisitBlock(ctx *parser.BlockContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitStatementList(ctx *parser.StatementListContext) Code {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitImportStatement(ctx *parser.ImportStatementContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitImportFromBlock(ctx *parser.ImportFromBlockContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitModuleItems(ctx *parser.ModuleItemsContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitImportDefault(ctx *parser.ImportDefaultContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitImportNamespace(ctx *parser.ImportNamespaceContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitImportFrom(ctx *parser.ImportFromContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitAliasName(ctx *parser.AliasNameContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitExportDeclaration(ctx *parser.ExportDeclarationContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitExportDefaultDeclaration(ctx *parser.ExportDefaultDeclarationContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitExportFromBlock(ctx *parser.ExportFromBlockContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitDeclaration(ctx *parser.DeclarationContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitVariableStatement(ctx *parser.VariableStatementContext) Code {
	return v.VisitChild(ctx, 0)
}

func defVar(typ string, name string, val Code) Code {
	switch typ {
	case "let":
		return Code{
			Main: "scopeCtx.DefineVariable(\"" + name + "\", " + val.Main + ", false)",
			Top:  val.Top,
			Init: val.Init,
		}
	case "var":
		return Code{
			Main: "funcCtx.DefineVariable(\"" + name + "\", " + val.Main + ", false)",
			Top:  val.Top,
			Init: val.Init,
		}
	case "const":
		return Code{
			Main: "scopeCtx.DefineVariable(\"" + name + "\", " + val.Main + ", true)",
			Top:  val.Top,
			Init: val.Init,
		}
	}

	panic("WHAT TYPE OF VAR IS THIS SHIT")
}

func (v *Visitor) VisitVariableDeclarationList(ctx *parser.VariableDeclarationListContext) Code {
	i := 0
	declars := ctx.AllVariableDeclaration()
	l := len(declars)
	code := Code{}
	mod := ctx.VarModifier()

	for i < l {
		var val Code
		child := declars[i]

		if child.GetChildCount() > 2 {
			val = v.Visit(child.GetChild(2).(antlr.RuleNode))
		} else {
			val = Code{
				Main: "&js_to_go.Undefined{}",
			}
		}

		assignable := child.GetChild(0).(*parser.AssignableContext)
		typ := mod.GetText()

		switch child := assignable.GetChild(0).(type) {
		case *parser.IdentifierContext:
			varCode := defVar(typ, child.GetText(), val)

			code.Init += "\n" + varCode.Init
			code.Main += "\n" + varCode.Main
			code.Top += "\n" + varCode.Top
		case *parser.ArrayLiteralContext:

			/*
				elList := child.ElementList().(*parser.ElementListContext)
				_list := elList.AllArrayElement()
			*/

			panic("NOT SUPPORTED")
		case *parser.ObjectLiteralContext:
			panic("NOT SUPPORTED")
		default:
			fmt.Println(reflect.TypeOf(child).String())
			panic("WAHT")
		}

		i++
	}

	return code
}

func (v *Visitor) VisitVariableDeclaration(ctx *parser.VariableDeclarationContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitEmptyStatement_(ctx *parser.EmptyStatement_Context) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitExpressionStatement(ctx *parser.ExpressionStatementContext) Code {
	return v.VisitChild(ctx, 0)
}

func (v *Visitor) VisitIfStatement(ctx *parser.IfStatementContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitDoStatement(ctx *parser.DoStatementContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitWhileStatement(ctx *parser.WhileStatementContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitForStatement(ctx *parser.ForStatementContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitForInStatement(ctx *parser.ForInStatementContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitForOfStatement(ctx *parser.ForOfStatementContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitVarModifier(ctx *parser.VarModifierContext) Code {
	panic("VISITVARMODIFIER")
}

func (v *Visitor) VisitContinueStatement(ctx *parser.ContinueStatementContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitBreakStatement(ctx *parser.BreakStatementContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitReturnStatement(ctx *parser.ReturnStatementContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitYieldStatement(ctx *parser.YieldStatementContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitWithStatement(ctx *parser.WithStatementContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitSwitchStatement(ctx *parser.SwitchStatementContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitCaseBlock(ctx *parser.CaseBlockContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitCaseClauses(ctx *parser.CaseClausesContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitCaseClause(ctx *parser.CaseClauseContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitDefaultClause(ctx *parser.DefaultClauseContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitLabelledStatement(ctx *parser.LabelledStatementContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitThrowStatement(ctx *parser.ThrowStatementContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitTryStatement(ctx *parser.TryStatementContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitCatchProduction(ctx *parser.CatchProductionContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitFinallyProduction(ctx *parser.FinallyProductionContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitDebuggerStatement(ctx *parser.DebuggerStatementContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitFunctionDeclaration(ctx *parser.FunctionDeclarationContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitClassDeclaration(ctx *parser.ClassDeclarationContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitClassTail(ctx *parser.ClassTailContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitClassElement(ctx *parser.ClassElementContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitMethodDefinition(ctx *parser.MethodDefinitionContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitFormalParameterList(ctx *parser.FormalParameterListContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitFormalParameterArg(ctx *parser.FormalParameterArgContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitLastFormalParameterArg(ctx *parser.LastFormalParameterArgContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitFunctionBody(ctx *parser.FunctionBodyContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitSourceElements(ctx *parser.SourceElementsContext) Code {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitArrayLiteral(ctx *parser.ArrayLiteralContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitElementList(ctx *parser.ElementListContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitArrayElement(ctx *parser.ArrayElementContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitPropertyExpressionAssignment(ctx *parser.PropertyExpressionAssignmentContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitComputedPropertyExpressionAssignment(ctx *parser.ComputedPropertyExpressionAssignmentContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitFunctionProperty(ctx *parser.FunctionPropertyContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitPropertyGetter(ctx *parser.PropertyGetterContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitPropertySetter(ctx *parser.PropertySetterContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitPropertyShorthand(ctx *parser.PropertyShorthandContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitPropertyName(ctx *parser.PropertyNameContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitArguments(ctx *parser.ArgumentsContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitArgument(ctx *parser.ArgumentContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitExpressionSequence(ctx *parser.ExpressionSequenceContext) Code {
	return v.VisitChild(ctx, 0)
}

func (v *Visitor) VisitTemplateStringExpression(ctx *parser.TemplateStringExpressionContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitTernaryExpression(ctx *parser.TernaryExpressionContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitLogicalAndExpression(ctx *parser.LogicalAndExpressionContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitPowerExpression(ctx *parser.PowerExpressionContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitPreIncrementExpression(ctx *parser.PreIncrementExpressionContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitObjectLiteralExpression(ctx *parser.ObjectLiteralExpressionContext) Code {
	return v.VisitChild(ctx, 0)
}

func (v *Visitor) VisitMetaExpression(ctx *parser.MetaExpressionContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitInExpression(ctx *parser.InExpressionContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitLogicalOrExpression(ctx *parser.LogicalOrExpressionContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitOptionalChainExpression(ctx *parser.OptionalChainExpressionContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitNotExpression(ctx *parser.NotExpressionContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitPreDecreaseExpression(ctx *parser.PreDecreaseExpressionContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitArgumentsExpression(ctx *parser.ArgumentsExpressionContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitAwaitExpression(ctx *parser.AwaitExpressionContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitThisExpression(ctx *parser.ThisExpressionContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitFunctionExpression(ctx *parser.FunctionExpressionContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitUnaryMinusExpression(ctx *parser.UnaryMinusExpressionContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitAssignmentExpression(ctx *parser.AssignmentExpressionContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitPostDecreaseExpression(ctx *parser.PostDecreaseExpressionContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitTypeofExpression(ctx *parser.TypeofExpressionContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitInstanceofExpression(ctx *parser.InstanceofExpressionContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitUnaryPlusExpression(ctx *parser.UnaryPlusExpressionContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitDeleteExpression(ctx *parser.DeleteExpressionContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitImportExpression(ctx *parser.ImportExpressionContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitEqualityExpression(ctx *parser.EqualityExpressionContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitBitXOrExpression(ctx *parser.BitXOrExpressionContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitSuperExpression(ctx *parser.SuperExpressionContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitMultiplicativeExpression(ctx *parser.MultiplicativeExpressionContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitBitShiftExpression(ctx *parser.BitShiftExpressionContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitParenthesizedExpression(ctx *parser.ParenthesizedExpressionContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitAdditiveExpression(ctx *parser.AdditiveExpressionContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitRelationalExpression(ctx *parser.RelationalExpressionContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitPostIncrementExpression(ctx *parser.PostIncrementExpressionContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitYieldExpression(ctx *parser.YieldExpressionContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitBitNotExpression(ctx *parser.BitNotExpressionContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitNewExpression(ctx *parser.NewExpressionContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitLiteralExpression(ctx *parser.LiteralExpressionContext) Code {
	return v.VisitChild(ctx, 0)
}

func (v *Visitor) VisitArrayLiteralExpression(ctx *parser.ArrayLiteralExpressionContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitMemberDotExpression(ctx *parser.MemberDotExpressionContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitClassExpression(ctx *parser.ClassExpressionContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitMemberIndexExpression(ctx *parser.MemberIndexExpressionContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitIdentifierExpression(ctx *parser.IdentifierExpressionContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitBitAndExpression(ctx *parser.BitAndExpressionContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitBitOrExpression(ctx *parser.BitOrExpressionContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitAssignmentOperatorExpression(ctx *parser.AssignmentOperatorExpressionContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitVoidExpression(ctx *parser.VoidExpressionContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitCoalesceExpression(ctx *parser.CoalesceExpressionContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitAssignable(ctx *parser.AssignableContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitObjectLiteral(ctx *parser.ObjectLiteralContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitFunctionDecl(ctx *parser.FunctionDeclContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitAnonymousFunctionDecl(ctx *parser.AnonymousFunctionDeclContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitArrowFunction(ctx *parser.ArrowFunctionContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitArrowFunctionParameters(ctx *parser.ArrowFunctionParametersContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitArrowFunctionBody(ctx *parser.ArrowFunctionBodyContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitAssignmentOperator(ctx *parser.AssignmentOperatorContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitLiteral(ctx *parser.LiteralContext) Code {
	child := ctx.GetChild(0)

	switch lit := child.(type) {
	case antlr.RuleNode:
		return v.Visit(lit)
	case *antlr.TerminalNodeImpl:
		return Code{
			Main: "\n&js_to_go.String{Value: " + lit.GetText() + "}",
		}
	}

	panic("WAIJDWODJW")
}

func (v *Visitor) VisitTemplateStringLiteral(ctx *parser.TemplateStringLiteralContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitTemplateStringAtom(ctx *parser.TemplateStringAtomContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitNumericLiteral(ctx *parser.NumericLiteralContext) Code {
	return Code{
		Main: "&js_to_go.Number{Value:" + ctx.GetText() + "}",
	}
}

func (v *Visitor) VisitBigintLiteral(ctx *parser.BigintLiteralContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitGetter(ctx *parser.GetterContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitSetter(ctx *parser.SetterContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitIdentifierName(ctx *parser.IdentifierNameContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitIdentifier(ctx *parser.IdentifierContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitReservedWord(ctx *parser.ReservedWordContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitKeyword(ctx *parser.KeywordContext) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitLet_(ctx *parser.Let_Context) Code {
	return v.Visit(ctx)
}

func (v *Visitor) VisitEos(ctx *parser.EosContext) Code {
	return v.Visit(ctx)
}
