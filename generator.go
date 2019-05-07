package main

import (
	"fmt"
	"strings"

	"github.com/tzmfreedom/land/ast"
)

type Generator struct {
	Indent int
}

func (v *Generator) AddIndent(f func()) {
	v.Indent += 4
	f()
	v.Indent -= 4
}

func (v *Generator) withIndent(src string) string {
	return strings.Repeat(" ", v.Indent) + src
}

func (v *Generator) VisitClassDeclaration(n *ast.ClassDeclaration) (interface{}, error) {
	annotations := make([]string, len(n.Annotations))
	for i, a := range n.Annotations {
		r, err := a.Accept(v)
		if err != nil {
			return nil, err
		}
		annotations[i] = r.(string)
	}
	annotationStr := ""
	if len(annotations) != 0 {
		annotationStr = fmt.Sprintf("%s\n", strings.Join(annotations, "\n"))
	}
	modifiers := make([]string, len(n.Modifiers))
	for i, m := range n.Modifiers {
		r, err := m.Accept(v)
		if err != nil {
			return nil, err
		}
		modifiers[i] = r.(string)
	}
	declarations := make([]string, len(n.Declarations))
	v.AddIndent(func() {
		for i, d := range n.Declarations {
			r, err := d.Accept(v)
			if err != nil {
				panic(err)
			}
			declarations[i] = r.(string)
		}
	})
	super := ""
	if n.SuperClassRef != nil {
		r, err := n.SuperClassRef.Accept(v)
		if err != nil {
			return nil, err
		}
		super = "extends " + r.(string)
	}
	implements := make([]string, len(n.ImplementClassRefs))
	for i, impl := range n.ImplementClassRefs {
		r, err := impl.Accept(v)
		if err != nil {
			return nil, err
		}
		implements[i] = r.(string)
	}
	implString := ""
	if len(implements) != 0 {
		implString = "implements " + strings.Join(implements, ", ")
	}
	body := ""
	if len(declarations) != 0 {
		body = fmt.Sprintf("%s\n", strings.Join(declarations, "\n"))
	}
	return fmt.Sprintf(
		`%s%s class %s %s %s {
%s%s`,
		annotationStr,
		strings.Join(modifiers, " "),
		n.Name,
		super,
		implString,
		body,
		v.withIndent("}"),
	), nil
}

func (v *Generator) VisitModifier(n *ast.Modifier) (interface{}, error) {
	return n.Name, nil
}

func (v *Generator) VisitAnnotation(n *ast.Annotation) (interface{}, error) {
	return n.Name, nil
}

func (v *Generator) VisitInterfaceDeclaration(n *ast.InterfaceDeclaration) (interface{}, error) {
	modifiers := make([]string, len(n.Modifiers))
	for i, m := range n.Modifiers {
		r, err := m.Accept(v)
		if err != nil {
			return nil, err
		}
		modifiers[i] = r.(string)
	}
	methods := make([]string, len(n.Methods))
	v.AddIndent(func() {
		for i, m := range n.Methods {
			r, err := m.Accept(v)
			if err != nil {
				panic(err)
			}
			methods[i] = r.(string)
		}
	})
	body := ""
	if len(methods) != 0 {
		body = fmt.Sprintf("%s\n", strings.Join(methods, "\n"))
	}

	return fmt.Sprintf(
		`%s interface %s {
%s%s`,
		strings.Join(modifiers, " "),
		n.Name,
		body,
		v.withIndent("}"),
	), nil
}

func (v *Generator) VisitIntegerLiteral(n *ast.IntegerLiteral) (interface{}, error) {
	return fmt.Sprintf("%d", n.Value), nil
}

func (v *Generator) VisitParameter(n *ast.Parameter) (interface{}, error) {
	r, err := n.TypeRef.Accept(v)
	if err != nil {
		return nil, err
	}
	return fmt.Sprintf(
		"%s %s",
		r.(string),
		n.Name,
	), nil
}

func (v *Generator) VisitArrayAccess(n *ast.ArrayAccess) (interface{}, error) {
	r, err := n.Receiver.Accept(v)
	if err != nil {
		return nil, err
	}
	k, err := n.Key.Accept(v)
	if err != nil {
		return nil, err
	}
	return fmt.Sprintf(
		"%s[%s]",
		r.(string),
		k.(string),
	), nil
}

func (v *Generator) VisitBooleanLiteral(n *ast.BooleanLiteral) (interface{}, error) {
	val := "false"
	if n.Value {
		val = "true"
	}
	return val, nil
}

func (v *Generator) VisitBreak(n *ast.Break) (interface{}, error) {
	return "break", nil
}

func (v *Generator) VisitContinue(n *ast.Continue) (interface{}, error) {
	return "continue", nil
}

func (v *Generator) VisitDml(n *ast.Dml) (interface{}, error) {
	r, err := n.Expression.Accept(v)
	if err != nil {
		return nil, err
	}
	return fmt.Sprintf("%s %s", n.Type, r.(string)), nil
}

func (v *Generator) VisitDoubleLiteral(n *ast.DoubleLiteral) (interface{}, error) {
	return fmt.Sprintf("%f", n.Value), nil
}

func (v *Generator) VisitFieldDeclaration(n *ast.FieldDeclaration) (interface{}, error) {
	modifiers := make([]string, len(n.Modifiers))
	for i, m := range n.Modifiers {
		r, err := m.Accept(v)
		if err != nil {
			return nil, err
		}
		modifiers[i] = r.(string)
	}
	r, err := n.TypeRef.Accept(v)
	if err != nil {
		return nil, err
	}
	declarators := make([]string, len(n.Declarators))
	for i, decl := range n.Declarators {
		r, err := decl.Accept(v)
		if err != nil {
			return nil, err
		}
		declarators[i] = r.(string)
	}

	return fmt.Sprintf(
		`%s %s %s;`,
		v.withIndent(strings.Join(modifiers, " ")),
		r.(string),
		strings.Join(declarators, ", "),
	), nil

}

func (v *Generator) VisitTry(n *ast.Try) (interface{}, error) {
	stmt := ""
	v.AddIndent(func() {
		r, err := n.Block.Accept(v)
		if err != nil {
			panic(err)
		}
		stmt = r.(string)
	})
	catches := make([]string, len(n.CatchClause))
	for i, c := range n.CatchClause {
		r, err := c.Accept(v)
		if err != nil {
			return nil, err
		}
		catches[i] = r.(string)
	}
	f, err := n.FinallyBlock.Accept(v)
	if err != nil {
		return nil, err
	}
	return fmt.Sprintf(
		`try {
%s%s%s
%s`,
		stmt,
		strings.Join(catches, "\n"),
		f.(string),
		v.withIndent("}"),
	), nil
}

func (v *Generator) VisitCatch(n *ast.Catch) (interface{}, error) {
	t, err := n.TypeRef.Accept(v)
	if err != nil {
		return nil, err
	}
	stmt := ""
	v.AddIndent(func() {
		r, err := n.Block.Accept(v)
		if err != nil {
			panic(err)
		}
		stmt = r.(string)
	})
	if stmt != "" {
		stmt = fmt.Sprintf("%s\n", stmt)
	}
	return fmt.Sprintf(
		` catch (%s %s) {
%s%s`,
		t.(string),
		n.Identifier,
		stmt,
		v.withIndent("}"),
	), nil
}

func (v *Generator) VisitFinally(n *ast.Finally) (interface{}, error) {
	stmt := ""
	v.AddIndent(func() {
		r, err := n.Block.Accept(v)
		if err != nil {
			panic(err)
		}
		stmt = r.(string)
	})
	if stmt != "" {
		stmt = fmt.Sprintf("%s\n", stmt)
	}
	return fmt.Sprintf(
		` finally {
%s%s`,
		stmt,
		v.withIndent("}"),
	), nil
}

func (v *Generator) VisitFor(n *ast.For) (interface{}, error) {
	control, err := n.Control.Accept(v)
	if err != nil {
		return nil, err
	}
	stmt := ""
	v.AddIndent(func() {
		r, err := n.Statements.Accept(v)
		if err != nil {
			panic(err)
		}
		stmt = r.(string)
	})
	if stmt != "" {
		stmt = fmt.Sprintf("%s\n", stmt)
	}
	return fmt.Sprintf(
		`for (%s) {
%s%s`,
		control.(string),
		stmt,
		v.withIndent("}"),
	), nil
}

func (v *Generator) VisitForControl(n *ast.ForControl) (interface{}, error) {
	inits := make([]string, len(n.ForInit))
	for i, forInit := range n.ForInit {
		exp, err := forInit.Accept(v)
		if err != nil {
			return nil, err
		}
		inits[i] = exp.(string)
	}
	exp, err := n.Expression.Accept(v)
	if err != nil {
		return nil, err
	}
	updates := make([]string, len(n.ForUpdate))
	for i, u := range n.ForUpdate {
		r, err := u.Accept(v)
		if err != nil {
			return nil, err
		}
		updates[i] = r.(string)
	}
	return fmt.Sprintf(
		`%s %s; %s`,
		strings.Join(inits, ", "),
		exp.(string),
		strings.Join(updates, ","),
	), nil
}

func (v *Generator) VisitEnhancedForControl(n *ast.EnhancedForControl) (interface{}, error) {
	t, err := n.TypeRef.Accept(v)
	if err != nil {
		return nil, err
	}
	exp, err := n.Expression.Accept(v)
	if err != nil {
		return nil, err
	}
	return fmt.Sprintf(
		`%s %s : %s`,
		t.(string),
		n.VariableDeclaratorId,
		exp.(string),
	), nil
}

func (v *Generator) VisitIf(n *ast.If) (interface{}, error) {
	cond, err := n.Condition.Accept(v)
	if err != nil {
		return nil, err
	}
	ifStmt := ""
	v.AddIndent(func() {
		r, err := n.IfStatement.Accept(v)
		if err != nil {
			panic(err)
		}
		ifStmt = r.(string)
	})
	if ifStmt != "" {
		ifStmt = fmt.Sprintf("%s\n", ifStmt)
	}
	elseStmt := ""
	if n.ElseStatement != nil {
		v.AddIndent(func() {
			r, err := n.IfStatement.Accept(v)
			if err != nil {
				panic(err)
			}
			elseStmt = r.(string)
		})
		if elseStmt != "" {
			elseStmt = fmt.Sprintf("%s\n", elseStmt)
		}
		elseStmt = fmt.Sprintf(
			` else {
%s
%s`,
			elseStmt,
			v.withIndent("}"),
		)
	}
	return fmt.Sprintf(
		`if (%s) {
%s%s%s`,
		cond.(string),
		ifStmt,
		v.withIndent("}"),
		elseStmt,
	), nil
}

func (v *Generator) VisitMethodDeclaration(n *ast.MethodDeclaration) (interface{}, error) {
	annotations := make([]string, len(n.Annotations))
	for i, a := range n.Annotations {
		r, err := a.Accept(v)
		if err != nil {
			return nil, err
		}
		annotations[i] = r.(string)
	}
	annotationStr := ""
	if len(annotations) != 0 {
		annotationStr = fmt.Sprintf("%s\n", strings.Join(annotations, "\n"))
	}
	modifiers := make([]string, len(n.Modifiers))
	for i, m := range n.Modifiers {
		r, err := m.Accept(v)
		if err != nil {
			return nil, err
		}
		modifiers[i] = r.(string)
	}
	returnType := "void"
	if n.ReturnType != nil {
		r, err := n.ReturnType.Accept(v)
		if err != nil {
			return nil, err
		}
		returnType = r.(string)
	}
	parameters := make([]string, len(n.Parameters))
	for i, p := range n.Parameters {
		r, err := p.Accept(v)
		if err != nil {
			return nil, err
		}
		parameters[i] = r.(string)
	}
	block := ""
	v.AddIndent(func() {
		r, err := n.Statements.Accept(v)
		if err != nil {
			panic(err)
		}
		block = r.(string)
	})
	if block != "" {
		block = fmt.Sprintf("%s\n", block)
	}
	return fmt.Sprintf(
		`%s%s %s %s (%s) {
%s%s`,
		annotationStr,
		v.withIndent(strings.Join(modifiers, " ")),
		returnType,
		n.Name,
		strings.Join(parameters, ", "),
		block,
		v.withIndent("}"),
	), nil
}

func (v *Generator) VisitMethodInvocation(n *ast.MethodInvocation) (interface{}, error) {
	exp, err := n.NameOrExpression.Accept(v)
	if err != nil {
		return nil, err
	}
	parameters := make([]string, len(n.Parameters))
	for i, p := range n.Parameters {
		r, err := p.Accept(v)
		if err != nil {
			return nil, err
		}
		parameters[i] = r.(string)
	}
	return fmt.Sprintf(
		"%s(%s)",
		exp.(string),
		strings.Join(parameters, ", "),
	), nil
}

func (v *Generator) VisitNew(n *ast.New) (interface{}, error) {
	t, err := n.TypeRef.Accept(v)
	if err != nil {
		return nil, err
	}
	parameters := make([]string, len(n.Parameters))
	for i, p := range n.Parameters {
		r, err := p.Accept(v)
		if err != nil {
			return nil, err
		}
		parameters[i] = r.(string)
	}
	return fmt.Sprintf(
		"new %s(%s)",
		t.(string),
		strings.Join(parameters, ", "),
	), nil
}

func (v *Generator) VisitNullLiteral(n *ast.NullLiteral) (interface{}, error) {
	return "null", nil
}

func (v *Generator) VisitUnaryOperator(n *ast.UnaryOperator) (interface{}, error) {
	val, err := n.Expression.Accept(v)
	if err != nil {
		return nil, err
	}
	if n.IsPrefix {
		return fmt.Sprintf("%s%s", n.Op, val.(string)), nil
	}
	return fmt.Sprintf("%s%s", val.(string), n.Op), nil
}

func (v *Generator) VisitBinaryOperator(n *ast.BinaryOperator) (interface{}, error) {
	l, err := n.Left.Accept(v)
	if err != nil {
		return nil, err
	}
	r, err := n.Right.Accept(v)
	if err != nil {
		return nil, err
	}
	return fmt.Sprintf("%s %s %s", l.(string), n.Op, r.(string)), nil
}

func (v *Generator) VisitReturn(n *ast.Return) (interface{}, error) {
	if n.Expression != nil {
		exp, err := n.Expression.Accept(v)
		if err != nil {
			return nil, err
		}
		return fmt.Sprintf("return %s", exp.(string)), nil
	}
	return "return", nil
}

func (v *Generator) VisitThrow(n *ast.Throw) (interface{}, error) {
	if n.Expression != nil {
		exp, err := n.Expression.Accept(v)
		if err != nil {
			return nil, err
		}
		return fmt.Sprintf("throw %s", exp.(string)), nil
	}
	return "throw", nil
}

func (v *Generator) VisitSoql(n *ast.Soql) (interface{}, error) {
	where := ""
	fields := make([]string, len(n.SelectFields))
	from := ""
	v.AddIndent(func() {
		v.AddIndent(func() {
			for i, f := range n.SelectFields {
				switch val := f.(type) {
				case *ast.SelectField:
					fields[i] = v.withIndent(strings.Join(val.Value, "."))
				case *ast.SoqlFunction:
					fields[i] = v.withIndent(val.Name + "()")
				}
			}

			from = v.withIndent(n.FromObject)

			if n.Where != nil {
				where = v.withIndent(v.createWhere(n.Where))
			}
		})
	})

	indent := ""
	v.AddIndent(func() {
		indent = v.withIndent("")
	})
	if where != "" {
		where = "\n" + indent + "WHERE\n" + where
	}
	orderBy := ""
	groupBy := ""
	limit := ""
	if n.Limit != nil {
		i, err := n.Limit.Accept(v)
		if err != nil {
			return nil, err
		}
		v.AddIndent(func() {
			v.AddIndent(func() {
				limit = "\n" + indent + "LIMIT\n" + v.withIndent(i.(string))
			})
		})
	}

	return fmt.Sprintf(`[
%sSELECT
%s
%sFROM
%s%s%s%s%s%s`,
		indent,
		strings.Join(fields, ",\n"),
		indent,
		from,
		where,
		orderBy,
		groupBy,
		limit,
		"\n"+v.withIndent("]"),
	), nil
}

func (v *Generator) createWhere(n ast.Node) string {
	switch val := n.(type) {
	case *ast.WhereCondition:
		var field string
		switch f := val.Field.(type) {
		case *ast.SelectField:
			field = strings.Join(f.Value, ".")
		case *ast.SoqlFunction:
			field = f.Name + "()"
		}
		value, err := val.Expression.Accept(v)
		if err != nil {
			panic(err)
		}
		return fmt.Sprintf("%s %s %s", field, val.Op, value.(string))
	case *ast.WhereBinaryOperator:
		where := ""
		if val.Left != nil {
			where += v.createWhere(val.Left)
		}
		if val.Right != nil {
			where += fmt.Sprintf("\n%s %s", v.withIndent(val.Op), v.createWhere(val.Right))
		}
		return where
	}
	return ""
}

func (v *Generator) VisitSosl(n *ast.Sosl) (interface{}, error) {
	return ast.VisitSosl(v, n)
}

func (v *Generator) VisitStringLiteral(n *ast.StringLiteral) (interface{}, error) {
	return "'" + n.Value + "'", nil
}

func (v *Generator) VisitSwitch(n *ast.Switch) (interface{}, error) {
	exp, err := n.Expression.Accept(v)
	if err != nil {
		return nil, err
	}
	whenStmts := make([]string, len(n.WhenStatements))
	v.AddIndent(func() {
		for i, stmt := range n.WhenStatements {
			r, err := stmt.Accept(v)
			if err != nil {
				panic(err)
			}
			whenStmts[i] = r.(string)
		}
	})
	elseStmt := ""
	v.AddIndent(func() {
		r, err := n.ElseStatement.Accept(v)
		if err != nil {
			panic(err)
		}
		elseStmt = r.(string)
	})
	if elseStmt != "" {
		elseStmt = fmt.Sprintf(
			` else {
%s
%s`,
			elseStmt,
			v.withIndent("}"),
		)
	}
	return fmt.Sprintf(
		`switch on %s {
%s
%s
%s`,
		exp.(string),
		strings.Join(whenStmts, "\n"),
		elseStmt,
		v.withIndent("}"),
	), nil
}

func (v *Generator) VisitTrigger(n *ast.Trigger) (interface{}, error) {
	timings := make([]string, len(n.TriggerTimings))
	for i, t := range n.TriggerTimings {
		r, err := t.Accept(v)
		if err != nil {
			return nil, err
		}
		timings[i] = r.(string)
	}
	stmt, err := n.Statements.Accept(v)
	if err != nil {
		return nil, err
	}
	return fmt.Sprintf(
		`trigger %s on %s (%s) {
%s
%s`,
		n.Name,
		n.Object,
		strings.Join(timings, ", "),
		stmt.(string),
		v.withIndent("}"),
	), nil
}

func (v *Generator) VisitTriggerTiming(n *ast.TriggerTiming) (interface{}, error) {
	return fmt.Sprintf("%s %s", n.Timing, n.Dml), nil
}

func (v *Generator) VisitVariableDeclaration(n *ast.VariableDeclaration) (interface{}, error) {
	t, err := n.TypeRef.Accept(v)
	if err != nil {
		return nil, err
	}
	declarators := make([]string, len(n.Declarators))
	for i, decl := range n.Declarators {
		r, err := decl.Accept(v)
		if err != nil {
			return nil, err
		}
		declarators[i] = r.(string)
	}
	return fmt.Sprintf(
		"%s %s",
		t.(string),
		strings.Join(declarators, ", "),
	), nil
}

func (v *Generator) VisitVariableDeclarator(n *ast.VariableDeclarator) (interface{}, error) {
	if n.Expression == nil {
		return fmt.Sprintf("%s", n.Name), nil
	}
	exp, err := n.Expression.Accept(v)
	if err != nil {
		return nil, err
	}
	return fmt.Sprintf("%s = %s", n.Name, exp.(string)), nil
}

func (v *Generator) VisitWhen(n *ast.When) (interface{}, error) {
	conditions := make([]string, len(n.Condition))
	for i, cond := range n.Condition {
		r, err := cond.Accept(v)
		if err != nil {
			return nil, err
		}
		conditions[i] = r.(string)
	}
	stmt := ""
	v.AddIndent(func() {
		r, err := n.Statements.Accept(v)
		if err != nil {
			panic(err)
		}
		stmt = r.(string)
	})
	return fmt.Sprintf(
		`when %s {
%s
%s`,
		strings.Join(conditions, ", "),
		stmt,
		v.withIndent("}"),
	), nil

	return ast.VisitWhen(v, n)
}

func (v *Generator) VisitWhenType(n *ast.WhenType) (interface{}, error) {
	r, err := n.TypeRef.Accept(v)
	if err != nil {
		return nil, err
	}
	return fmt.Sprintf(
		"%s %s",
		r.(string),
		n.Identifier,
	), nil
}

func (v *Generator) VisitWhile(n *ast.While) (interface{}, error) {
	cond, err := n.Condition.Accept(v)
	if err != nil {
		return nil, err
	}
	statements := ""
	v.AddIndent(func() {
		r, err := n.Statements.Accept(v)
		if err != nil {
			panic(err)
		}
		statements = r.(string)
	})
	return fmt.Sprintf(
		`while (%s) {
%s
%s`,
		cond.(string),
		statements,
		v.withIndent("}"),
	), nil
}

func (v *Generator) VisitNothingStatement(n *ast.NothingStatement) (interface{}, error) {
	return "", nil
}

func (v *Generator) VisitCastExpression(n *ast.CastExpression) (interface{}, error) {
	t, err := n.CastTypeRef.Accept(v)
	if err != nil {
		return nil, err
	}
	exp, err := n.Expression.Accept(v)
	if err != nil {
		return nil, err
	}
	return fmt.Sprintf("(%s)%s", t.(string), exp.(string)), nil
}

func (v *Generator) VisitFieldAccess(n *ast.FieldAccess) (interface{}, error) {
	exp, err := n.Expression.Accept(v)
	if err != nil {
		return nil, err
	}
	return fmt.Sprintf("%s.%s", exp.(string), n.FieldName), nil
}

func (v *Generator) VisitType(n *ast.TypeRef) (interface{}, error) {
	paramString := ""
	params := make([]string, len(n.Parameters))
	for i, param := range n.Parameters {
		r, err := param.Accept(v)
		if err != nil {
			return nil, err
		}
		params[i] = r.(string)
	}
	if len(params) != 0 {
		paramString = fmt.Sprintf("<%s>", strings.Join(params, ", "))
	}
	return fmt.Sprintf(
		"%s%s",
		strings.Join(n.Name, "."),
		paramString,
	), nil
}

func (v *Generator) VisitBlock(n *ast.Block) (interface{}, error) {
	statements := make([]string, len(n.Statements))
	for i, s := range n.Statements {
		r, err := s.Accept(v)
		if err != nil {
			return nil, err
		}
		statements[i] = v.withIndent(r.(string)) + ";"
	}
	return strings.Join(statements, "\n"), nil
}

func (v *Generator) VisitGetterSetter(n *ast.GetterSetter) (interface{}, error) {
	return ast.VisitGetterSetter(v, n)
}

func (v *Generator) VisitPropertyDeclaration(n *ast.PropertyDeclaration) (interface{}, error) {
	return ast.VisitPropertyDeclaration(v, n)
}

func (v *Generator) VisitArrayInitializer(n *ast.ArrayInitializer) (interface{}, error) {
	return ast.VisitArrayInitializer(v, n)
}

func (v *Generator) VisitArrayCreator(n *ast.ArrayCreator) (interface{}, error) {
	return ast.VisitArrayCreator(v, n)
}

func (v *Generator) VisitSoqlBindVariable(n *ast.SoqlBindVariable) (interface{}, error) {
	return ast.VisitSoqlBindVariable(v, n)
}

func (v *Generator) VisitTernalyExpression(n *ast.TernalyExpression) (interface{}, error) {
	return ast.VisitTernalyExpression(v, n)
}

func (v *Generator) VisitMapCreator(n *ast.MapCreator) (interface{}, error) {
	return ast.VisitMapCreator(v, n)
}

func (v *Generator) VisitSetCreator(n *ast.SetCreator) (interface{}, error) {
	return ast.VisitSetCreator(v, n)
}

func (v *Generator) VisitName(n *ast.Name) (interface{}, error) {
	return strings.Join(n.Value, "."), nil
}

func (v *Generator) VisitConstructorDeclaration(n *ast.ConstructorDeclaration) (interface{}, error) {
	annotations := make([]string, len(n.Annotations))
	for i, a := range n.Annotations {
		r, err := a.Accept(v)
		if err != nil {
			return nil, err
		}
		annotations[i] = r.(string)
	}
	annotationStr := ""
	if len(annotations) != 0 {
		annotationStr = fmt.Sprintf("%s\n", strings.Join(annotations, "\n"))
	}
	modifiers := make([]string, len(n.Modifiers))
	for i, m := range n.Modifiers {
		r, err := m.Accept(v)
		if err != nil {
			return nil, err
		}
		modifiers[i] = r.(string)
	}
	returnType := "void"
	if n.ReturnType != nil {
		r, err := n.ReturnType.Accept(v)
		if err != nil {
			return nil, err
		}
		returnType = r.(string)
	}
	parameters := make([]string, len(n.Parameters))
	for i, p := range n.Parameters {
		r, err := p.Accept(v)
		if err != nil {
			return nil, err
		}
		parameters[i] = r.(string)
	}
	block := ""
	v.AddIndent(func() {
		r, err := n.Statements.Accept(v)
		if err != nil {
			panic(err)
		}
		block = r.(string)
	})
	if block != "" {
		block = fmt.Sprintf("%s\n", block)
	}
	return fmt.Sprintf(
		`%s%s %s %s (%s) {
%s%s`,
		annotationStr,
		v.withIndent(strings.Join(modifiers, " ")),
		returnType,
		n.Name,
		strings.Join(parameters, ", "),
		block,
		v.withIndent("}"),
	), nil
}

func Generate(n ast.Node) string {
	visitor := &Generator{}
	r, err := n.Accept(visitor)
	if err != nil {
		panic(err)
	}
	return r.(string)
}
