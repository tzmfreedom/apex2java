package main

import (
	"github.com/tzmfreedom/land/ast"
	"strings"
)

var ImportClasses = map[string]string{
	"system": "com.freedom_man.system.System",
	"database": "com.freedom_man.system.Database",
	"list":     "com.freedom_man.system.List",
	"account":  "com.freedom_man.system.Account",
	"sobject":  "com.freedom_man.system.SObject",
}

type ImportTypeResolver struct {
	importClasses map[string]struct{}
}

func NewImportTypeResolver() *ImportTypeResolver {
	return &ImportTypeResolver{
		importClasses: map[string]struct{}{},
	}
}

func (v *ImportTypeResolver) Resolve(n ast.Node) (interface{}, error) {
	return n.Accept(v)
}

func (v *ImportTypeResolver) VisitClassDeclaration(n *ast.ClassDeclaration) (interface{}, error) {
	for _, d := range n.Declarations {
		d.Accept(v)
	}
	for _, c := range n.InnerClasses {
		c.Accept(v)
	}
	return nil, nil
}

func (v *ImportTypeResolver) VisitModifier(n *ast.Modifier) (interface{}, error) {
	return ast.VisitModifier(v, n)
}

func (v *ImportTypeResolver) VisitAnnotation(n *ast.Annotation) (interface{}, error) {
	return ast.VisitAnnotation(v, n)
}

func (v *ImportTypeResolver) VisitInterfaceDeclaration(n *ast.InterfaceDeclaration) (interface{}, error) {
	return ast.VisitInterfaceDeclaration(v, n)
}

func (v *ImportTypeResolver) VisitIntegerLiteral(n *ast.IntegerLiteral) (interface{}, error) {
	return ast.VisitIntegerLiteral(v, n)
}

func (v *ImportTypeResolver) VisitParameter(n *ast.Parameter) (interface{}, error) {
	return ast.VisitParameter(v, n)
}

func (v *ImportTypeResolver) VisitArrayAccess(n *ast.ArrayAccess) (interface{}, error) {
	return ast.VisitArrayAccess(v, n)
}

func (v *ImportTypeResolver) VisitBooleanLiteral(n *ast.BooleanLiteral) (interface{}, error) {
	return ast.VisitBooleanLiteral(v, n)
}

func (v *ImportTypeResolver) VisitBreak(n *ast.Break) (interface{}, error) {
	return ast.VisitBreak(v, n)
}

func (v *ImportTypeResolver) VisitContinue(n *ast.Continue) (interface{}, error) {
	return ast.VisitContinue(v, n)
}

func (v *ImportTypeResolver) VisitDml(n *ast.Dml) (interface{}, error) {
	return ast.VisitDml(v, n)
}

func (v *ImportTypeResolver) VisitDoubleLiteral(n *ast.DoubleLiteral) (interface{}, error) {
	return ast.VisitDoubleLiteral(v, n)
}

func (v *ImportTypeResolver) VisitFieldDeclaration(n *ast.FieldDeclaration) (interface{}, error) {
	return n.TypeRef.Accept(v)
}

func (v *ImportTypeResolver) VisitTry(n *ast.Try) (interface{}, error) {
	return ast.VisitTry(v, n)
}

func (v *ImportTypeResolver) VisitCatch(n *ast.Catch) (interface{}, error) {
	return ast.VisitCatch(v, n)
}

func (v *ImportTypeResolver) VisitFinally(n *ast.Finally) (interface{}, error) {
	return ast.VisitFinally(v, n)
}

func (v *ImportTypeResolver) VisitFor(n *ast.For) (interface{}, error) {
	return ast.VisitFor(v, n)
}

func (v *ImportTypeResolver) VisitForControl(n *ast.ForControl) (interface{}, error) {
	return ast.VisitForControl(v, n)
}

func (v *ImportTypeResolver) VisitEnhancedForControl(n *ast.EnhancedForControl) (interface{}, error) {
	return ast.VisitEnhancedForControl(v, n)
}

func (v *ImportTypeResolver) VisitIf(n *ast.If) (interface{}, error) {
	return ast.VisitIf(v, n)
}

func (v *ImportTypeResolver) VisitMethodDeclaration(n *ast.MethodDeclaration) (interface{}, error) {
	return n.Statements.Accept(v)
}

func (v *ImportTypeResolver) VisitMethodInvocation(n *ast.MethodInvocation) (interface{}, error) {
	n.NameOrExpression.Accept(v)
	for _, p := range n.Parameters {
		p.Accept(v)
	}
	return nil, nil
}

func (v *ImportTypeResolver) VisitNew(n *ast.New) (interface{}, error) {
	return ast.VisitNew(v, n)
}

func (v *ImportTypeResolver) VisitNullLiteral(n *ast.NullLiteral) (interface{}, error) {
	return ast.VisitNullLiteral(v, n)
}

func (v *ImportTypeResolver) VisitUnaryOperator(n *ast.UnaryOperator) (interface{}, error) {
	return ast.VisitUnaryOperator(v, n)
}

func (v *ImportTypeResolver) VisitBinaryOperator(n *ast.BinaryOperator) (interface{}, error) {
	if _, err := n.Left.Accept(v); err != nil {
		return nil, err
	}
	if _, err := n.Right.Accept(v); err != nil {
		return nil, err
	}
	return nil, nil
}

func (v *ImportTypeResolver) VisitInstanceofOperator(n *ast.InstanceofOperator) (interface{}, error) {
	return ast.VisitInstanceofOperator(v, n)
}

func (v *ImportTypeResolver) VisitReturn(n *ast.Return) (interface{}, error) {
	return ast.VisitReturn(v, n)
}

func (v *ImportTypeResolver) VisitThrow(n *ast.Throw) (interface{}, error) {
	return ast.VisitThrow(v, n)
}

func (v *ImportTypeResolver) VisitSoql(n *ast.Soql) (interface{}, error) {
	if packageName, ok := ImportClasses[strings.ToLower(n.FromObject)]; ok {
		v.importClasses[packageName] = struct{}{}
	}
	v.importClasses["database"] = struct{}{}
	return nil, nil
}

func (v *ImportTypeResolver) VisitSosl(n *ast.Sosl) (interface{}, error) {
	return ast.VisitSosl(v, n)
}

func (v *ImportTypeResolver) VisitStringLiteral(n *ast.StringLiteral) (interface{}, error) {
	return ast.VisitStringLiteral(v, n)
}

func (v *ImportTypeResolver) VisitSwitch(n *ast.Switch) (interface{}, error) {
	return ast.VisitSwitch(v, n)
}

func (v *ImportTypeResolver) VisitTrigger(n *ast.Trigger) (interface{}, error) {
	return ast.VisitTrigger(v, n)
}

func (v *ImportTypeResolver) VisitTriggerTiming(n *ast.TriggerTiming) (interface{}, error) {
	return ast.VisitTriggerTiming(v, n)
}

func (v *ImportTypeResolver) VisitVariableDeclaration(n *ast.VariableDeclaration) (interface{}, error) {
	n.TypeRef.Accept(v)
	for _, d := range n.Declarators {
		d.Accept(v)
	}
	return nil, nil
}

func (v *ImportTypeResolver) VisitVariableDeclarator(n *ast.VariableDeclarator) (interface{}, error) {
	return n.Expression.Accept(v)
}

func (v *ImportTypeResolver) VisitWhen(n *ast.When) (interface{}, error) {
	return ast.VisitWhen(v, n)
}

func (v *ImportTypeResolver) VisitWhenType(n *ast.WhenType) (interface{}, error) {
	return ast.VisitWhenType(v, n)
}

func (v *ImportTypeResolver) VisitWhile(n *ast.While) (interface{}, error) {
	return ast.VisitWhile(v, n)
}

func (v *ImportTypeResolver) VisitNothingStatement(n *ast.NothingStatement) (interface{}, error) {
	return ast.VisitNothingStatement(v, n)
}

func (v *ImportTypeResolver) VisitCastExpression(n *ast.CastExpression) (interface{}, error) {
	return ast.VisitCastExpression(v, n)
}

func (v *ImportTypeResolver) VisitFieldAccess(n *ast.FieldAccess) (interface{}, error) {
	return n.Expression.Accept(v)
}

func (v *ImportTypeResolver) VisitType(n *ast.TypeRef) (interface{}, error) {
	// TODO: impl
	if packageName, ok := ImportClasses[strings.ToLower(n.Name[0])]; ok {
		v.importClasses[packageName] = struct{}{}
	}
	for _, p := range n.Parameters {
		p.Accept(v)
	}
	return nil, nil
}

func (v *ImportTypeResolver) VisitBlock(n *ast.Block) (interface{}, error) {
	for _, s := range n.Statements {
		s.Accept(v)
	}
	return nil, nil
}

func (v *ImportTypeResolver) VisitGetterSetter(n *ast.GetterSetter) (interface{}, error) {
	return ast.VisitGetterSetter(v, n)
}

func (v *ImportTypeResolver) VisitPropertyDeclaration(n *ast.PropertyDeclaration) (interface{}, error) {
	return ast.VisitPropertyDeclaration(v, n)
}

func (v *ImportTypeResolver) VisitArrayInitializer(n *ast.ArrayInitializer) (interface{}, error) {
	return ast.VisitArrayInitializer(v, n)
}

func (v *ImportTypeResolver) VisitArrayCreator(n *ast.ArrayCreator) (interface{}, error) {
	return ast.VisitArrayCreator(v, n)
}

func (v *ImportTypeResolver) VisitSoqlBindVariable(n *ast.SoqlBindVariable) (interface{}, error) {
	return ast.VisitSoqlBindVariable(v, n)
}

func (v *ImportTypeResolver) VisitTernalyExpression(n *ast.TernalyExpression) (interface{}, error) {
	return ast.VisitTernalyExpression(v, n)
}

func (v *ImportTypeResolver) VisitMapCreator(n *ast.MapCreator) (interface{}, error) {
	return ast.VisitMapCreator(v, n)
}

func (v *ImportTypeResolver) VisitSetCreator(n *ast.SetCreator) (interface{}, error) {
	return ast.VisitSetCreator(v, n)
}

func (v *ImportTypeResolver) VisitName(n *ast.Name) (interface{}, error) {
	if packageName, ok := ImportClasses[strings.ToLower(n.Value[0])]; ok {
		v.importClasses[packageName] = struct{}{}
	}
	return ast.VisitName(v, n)
}

func (v *ImportTypeResolver) VisitConstructorDeclaration(n *ast.ConstructorDeclaration) (interface{}, error) {
	return ast.VisitConstructorDeclaration(v, n)
}
