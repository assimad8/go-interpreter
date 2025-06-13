package parser

import (
	"fmt"
	"testing"

	"github.com/assimad8/go-interpreter/internal/ast"
	"github.com/assimad8/go-interpreter/internal/lexer"
)

func TestLetStatements(t *testing.T) {

	tests := []struct{
		input string
		expectedIdentifier string
		expectedValue any
	}{
		{"let x = 5;","x",5},
		{"let y = true;","y",true},
		{"let foobar = y;","foobar","y"},
	}
	for _,tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t,p)

		if program == nil {
			t.Fatalf("ParseProgram() returned nil")
		}
		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d",1, len(program.Statements))
		}
		stmt := program.Statements[0]
		if !testLetStatement(t,stmt,tt.expectedIdentifier){
			return
		}

		val := stmt.(*ast.LetStatement).Value

		if !testLiteralExpression(t,val,tt.expectedValue){
			return
		}
	}
}

func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral not 'let'. got=%q", s.TokenLiteral())
		return false
	}
	
	letStmt, ok := s.(*ast.LetStatement)
	
	if !ok {
		t.Errorf("s not *ast.LetStatement. got=%T", s)
		return false
	}
	
	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value not '%s'. got=%s", name, letStmt.Name.Value)
		return false
	}
	
	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("s.Name not '%s'. got=%s", name, letStmt.Name)
		return false
	}
	
	return true
}

func TestReturnStatements(t *testing.T) {
	
	input := `
	return 5;
	return 10;
	return 838383;
	`
	l := lexer.New(input)
	p := New(l)
	
	program := p.ParseProgram()
	checkParserErrors(t,p)
	
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d", len(program.Statements))
	}

	for _, stmt := range program.Statements {
		returnStmt,ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stmt not *ast.ReturnStatement.got=%T",stmt)
			continue
		}
		if returnStmt.TokenLiteral()!="return"{
			t.Errorf("returnStmt.TokenLiteral not 'return'. got=%q",returnStmt.TokenLiteral())
		}
	}
}

func TestIdentifier(t *testing.T) {
	input := "foobar;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t,p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d.",len(program.Statements))
	}
	stmt,ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] not *ast.ExpressionStatement. got=%T.",program.Statements[0])
	}
	ident,ok := stmt.Expression.(*ast.Identifier)
	if !ok{
		t.Fatalf("exp not *ast.Identifier. got=%T.",stmt.Expression)
	}
	if ident.Value != "foobar" {
		t.Errorf("ident.Value not %s. got=%s","foobar",ident.Value)
	}
	if ident.TokenLiteral() != "foobar" {
		t.Errorf("ident.TokenLiteral() not %s. got=%s","foobar",ident.TokenLiteral())
	} 
}

func TestIntegerLiteralExpression(t *testing.T){
	input := "5;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t,p)

	if len(program.Statements)!=1 {
		t.Fatalf("program has not enough Statements.got=%d",len(program.Statements))
	}
	stmt,ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("progtam.Statements[0] is not ast.ExpressionStatement.got=%T",program.Statements[0])
	}
	literal,ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("exp not ast.IntegerLiteral. got=%T",stmt.Expression)
	}
	if literal.Value != 5 {
		t.Errorf("literal.Value not %d. got=%d",5,literal.Value)
	}
	if literal.TokenLiteral() != "5" {
		t.Errorf("literal.Value not %s. got=%s","5",literal.TokenLiteral())
	}
}

func TestParsingPrefixExpression(t *testing.T) {
	prefixTests := []struct{
		input 	     string
		operator	 string
		integerValue int64
	}{
		{"!5","!",5},
		{"-15","-",15},
	}
	for _,tt := range prefixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t,p)

		if len(program.Statements)!= 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d",1,len(program.Statements))
		}
		stmt,ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",program.Statements[0])
		}
		exp,ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("stmt.Expression not ast.PrefixExpression. got=%T",stmt.Expression)
		}
		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s'.got=%s",tt.operator,exp.Operator)
		}
		if !testIntegerLiteral(t,exp.Right,tt.integerValue){
			return
		}
	}
}

func testIntegerLiteral(t *testing.T,il ast.Expression,v int64)bool {
	integ,ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il not *ast.IntegerLiteral. got=%T",il)
		return false
	}
	if integ.Value != v {
		t.Errorf("integ.Value not %d. got=%d",v,integ.Value)
		return false
	}
	if integ.TokenLiteral() != fmt.Sprintf("%d",v) {
		t.Errorf("integ.TokenLiteral not %q. got=%q",v,integ.TokenLiteral())
		return false
	}
	return true
}

func TestParsingInfixExpression(t *testing.T) {
	infixTests := []struct {
		input		  string
		leftValue     int64
		operator	  string
		rightValue	  int64
	}{
		{"5 + 5",5,"+",5},
		{"5 - 5",5,"-",5},
		{"5 * 5",5,"*",5},
		{"5 / 5",5,"/",5},
		{"5 > 5",5,">",5},
		{"5 < 5",5,"<",5},
		{"5 == 5",5,"==",5},
		{"5 != 5",5,"!=",5},
	}

	for _,tt := range infixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t,p)
		
		if len(program.Statements)!=1{
			t.Fatalf("program.Statements does not contain %d. got=%d",1,len(program.Statements))
		}
		stmt,ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatements. got=%T",program.Statements[0])
		}
		exp,ok := stmt.Expression.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("stmt.Expression is not ast.InfixExpression. got=%T",stmt.Expression)
		}

		if !testIntegerLiteral(t,exp.Left,tt.leftValue){
			return
		}
		
		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got='%s'",tt.operator,exp.Operator)
		}
		if !testIntegerLiteral(t,exp.Right,tt.rightValue){
			return
		}
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input 		string
		exprected 	string
	}{
		{"-a * b","((-a) * b)"},
		{"!-a","(!(-a))"},
		{"a + b + c","((a + b) + c)"},
		{"a + b - c","((a + b) - c)"},
		{"a * b * c","((a * b) * c)"},
		{"a * b / c","((a * b) / c)"},
		{"a + b / c","(a + (b / c))"},
		{"a + b * c + d / e - f","(((a + (b * c)) + (d / e)) - f)"},
		{"3 + 4; -5 * 5","(3 + 4)((-5) * 5)"},
		{"5 > 4 == 3 < 4","((5 > 4) == (3 < 4))"},
		{"5 < 4 != 3 > 4","((5 < 4) != (3 > 4))",},
		{"3 + 4 * 5 == 3 * 1 + 4 * 5","((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",},
		{"3 + 4 * 5 == 3 * 1 + 4 * 5","((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",},
		{"true","true"},
		{"false","false"},
		{"3 < 4 == true","((3 < 4) == true)"},
		{"3 > 4 == false","((3 > 4) == false)"},
		{"(5 + 5) * 2","((5 + 5) * 2)",},
		{"-(5 + 5)","(-(5 + 5))",},
		{"a * [1, 2, 3, 4][b * c] * d","((a * ([1, 2, 3, 4][(b * c)])) * d)",},
		{"add(a * b[2], b[1], 2 * [1, 2][1])","add((a * (b[2])), (b[1]), (2 * ([1, 2][1])))",},
	}

	for _,tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t,p)

		actual := program.String()

		if actual != tt.exprected {
			t.Errorf("expected[%q]. got=[%q]",tt.exprected,actual)
		}

	}
}

func TestIfExpression(t *testing.T) {
	input := "if (x < y) {x}"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t,p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statments not %d statements. got=%d",1,len(program.Statements))
	}

	stmt,ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statement[0] is not ast.ExpressionStatement. got=%T",program.Statements[0])
	}
	exp,ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.IfExpression. gt=%T",stmt)
	}

	if !testInfixExpression(t,exp.Condition,"x","<","y") {
		return
	}

	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statements. got=%d\n",len(exp.Consequence.Statements))
	}
	consequence,ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement) 
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T",exp.Consequence.Statements[0])
	}
	if !testIdentifier(t,consequence.Expression,"x"){
		return 
	}

	if exp.Alternative != nil {
		t.Errorf("exp.Alternative.Statements was not nil. got=%+v",exp.Alternative)
	}
}

func TestFunctionLiteralParsing(t *testing.T) {
	input := "fn(x, y) { x + y }"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t,p)

	if len(program.Statements)!= 1{
		t.Fatalf("program.Body does not contain %d statements. got=%d",1,len(program.Statements))
	}

	stmt,ok := program.Statements[0].(*ast.ExpressionStatement) 
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",program.Statements[0])
	}

	function ,ok := stmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.FunctionLiteral. got=%T",stmt.Expression)
	}
	if len(function.Parameters)!= 2 {
		t.Fatalf("function literal parameters wrong.want 2 got=%d",len(function.Parameters))
	}

	testLiteralExpression(t,function.Parameters[0],"x")
	testLiteralExpression(t,function.Parameters[1],"y")

	if len(function.Body.Statements) != 1 {
		t.Fatalf("function.Body.Statements has not 1 statements.got=%d",len(function.Body.Statements))
	}

	bodyStmt,ok := function.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("function body stmt is not ast.ExpressionStatement. got=%T",function.Body.Statements[0])
	}

	testInfixExpression(t,bodyStmt.Expression,"x","+","y")
}

func TestFunctionParametersParsing(t *testing.T) {
	tests := []struct{
		input string
		expectedParams []string
	}{
		{"fn() {}",[]string{}},
		{"fn(x) {}",[]string{"x"}},
		{"fn(x, y, z) {}",[]string{"x","y","z"}},
	}


	for _,tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t,p)

		stmt := program.Statements[0].(*ast.ExpressionStatement)
		function := stmt.Expression.(*ast.FunctionLiteral)

		if len(function.Parameters) != len(tt.expectedParams) {
			t.Errorf("length parameters wrong.expected %d. gpt=%d", len(tt.expectedParams),len(function.Parameters))
		}

		for i,ident := range tt.expectedParams {
			testLiteralExpression(t,function.Parameters[i],ident)
		}
	}
}

func TestCallExpressionParsing(t *testing.T) {
	input := "add( 1 , 2 * 1 , 4 + 5 );"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t,p)
	
	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d",1,len(program.Statements))
	}

	stmt,ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok{
		t.Fatalf("stmt is not ast.ExpressionStatement. got=%T",program.Statements[0])
	}
	exp,ok := stmt.Expression.(*ast.CallExpression)
	if !ok{
		t.Fatalf("stmt.Expression is not ast.CallExpression. got=%T",stmt.Expression)
	}

	if !testIdentifier(t,exp.Function,"add"){
		return
	}
	if len(exp.Arguments)!=3 {
		t.Fatalf("wrong length of arguments expected [3]. got=%d",len(exp.Arguments))
	}

	testLiteralExpression(t,exp.Arguments[0],1)
	testInfixExpression(t,exp.Arguments[1],2,"*",1)
	testInfixExpression(t,exp.Arguments[2],4,"+",5)
}

func TestStringLiteralExpression(t *testing.T) {
	input := `"hello world"`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t,p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	literal,ok := stmt.Expression.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("exp not *ast.StringLiteral. got=%T",stmt.Expression)
	}
	if literal.Value != "hello world" {
		t.Fatalf("literal.Value not %q. gpt=%q","hello world",literal.Value)
	}
}

func TestParsingArrayLiterals(t *testing.T){
	input := "[1,2*2,3+3]"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t,p)

	stmt,ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",program.Statements[0])
	}
	exp,ok := stmt.Expression.(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.ArrayLiteral. got=%T",stmt.Expression)
	} 
	if len(exp.Elements) != 3 {
		t.Fatalf("len(exp.Elements) is not %d. got=%d",3,len(exp.Elements))
	}

	testIntegerLiteral(t,exp.Elements[0],1)
	testInfixExpression(t,exp.Elements[1],2,"*",2)
	testInfixExpression(t,exp.Elements[2],3,"+",3)
}

func TestParsingIndexExpressions(t *testing.T) {
	input := "myArray[1+1]"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t,p)

	stmt,ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",program.Statements[0])
	}

	indexExp,ok := stmt.Expression.(*ast.IndexExpression) 
	if !ok {
		t.Fatalf("stmt.Expression is not ast.IndexExpression. got=%T(%+v)",stmt.Expression,stmt.Expression)
	}
	if !testIdentifier(t,indexExp.Left,"myArray"){
		return
	}
	if !testInfixExpression(t,indexExp.Index,1,"+",1){
		return
	}
}

func TestParsingHashLiterals(t *testing.T) {
	input := `{"one":1,"two":2}`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t,p)

	stmt,ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",program.Statements[0])
	}
	hash,ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.HashLiteral. got=%T",stmt.Expression)
	}
	if len(hash.Pairs) != 2 {
		t.Errorf("hash.Pairs has wrong length. got=%d",len(hash.Pairs))
	}
	expected := map[string]int64{
		"one":1,
		"two":2,
	}
	
	for key,value := range hash.Pairs {
		literal,ok := key.(*ast.StringLiteral)
		if !ok {
			t.Errorf("key is not ast.StringLiteral. got=%T",key)
		}
		expectedValue := expected[literal.String()]
		testIntegerLiteral(t,value,expectedValue)
	}
}

func TestEmptyHashLiteral(t *testing.T) {
	input := "{}"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t,p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	hash,ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("exp is not ast.HashLiteral. got=%T",stmt.Expression)
	}
	if len(hash.Pairs)!=0{
		t.Errorf("hash.Pairs has wrong length. got=%d",len(hash.Pairs))
	}
}
func TestParsingHashLiteralsWithExpressions(t *testing.T) {
	input := `{"one":0+1,"two":10-8,"three":15/5}`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t,p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	hash,ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("exp is not ast.HashLiteral. got=%T",stmt.Expression)
	}
	if len(hash.Pairs)!=3{
		t.Errorf("hash.Pairs has wrong length. got=%d",len(hash.Pairs))
	}

	tests := map[string]func(ast.Expression){
		"one":func(e ast.Expression){
			testInfixExpression(t,e,0,"+",1)
		},
		"two":func(e ast.Expression){
			testInfixExpression(t,e,10,"-",8)
		},
		"three":func(e ast.Expression){
			testInfixExpression(t,e,15,"/",5)
		},
	}

	for key ,value := range hash.Pairs {
		literal,ok := key.(*ast.StringLiteral)
		if !ok {
			t.Errorf("key is not ast.StringLiteral. got=%T",key)
			continue
		}
		testFunc,ok := tests[literal.String()]
		if !ok {
			t.Errorf("No test function for key %q found",literal.String())
		}
		testFunc(value)
	}
}

func testIdentifier(t *testing.T,exp ast.Expression,value string) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("exp not *ast.Identifier. got=%T",exp)
		return false
	}

	if ident.Value != value {
		t.Errorf("ident.Value not %s . got=%s",value,ident.Value)
		return false
	}
	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral not %s. got=%s",value,ident.TokenLiteral())
		return false
	}
	return true
}

func testLiteralExpression(t *testing.T,exp ast.Expression,expected any) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t,exp,int64(v))
	case int64:
		return testIntegerLiteral(t,exp,v)
	case string:
		return testIdentifier(t,exp,v)
	case bool:
		return testBooleanLiteral(t,exp,v)
	}
	t.Errorf("type of exp not handled. got=%T",exp)
	return false
}

func testInfixExpression(t *testing.T,exp ast.Expression,left any,operator string,right any) bool {
	opExp,ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("exp is not ast.OperatorExpression . got=%T",exp)
		return false
	}
	if !testLiteralExpression(t,opExp.Left,left){
		return false
	}
	if opExp.Operator != operator {
		t.Errorf("exp.Operator is not '%s' .go='%s'",operator,opExp.Operator)
	}
	if !testLiteralExpression(t,opExp.Right,right){
		return false
	}
	return true
}

func testBooleanLiteral(t *testing.T,exp ast.Expression,v bool) bool {
	bo ,ok := exp.(*ast.Boolean)
	if !ok {
		t.Errorf("exp is not ast.Boolean . got=%T",exp)
		return false
	}
	if bo.Value != v {
		t.Errorf("bo.Value not %t. got=%t",v,bo.Value)
		return false
	}
	if bo.TokenLiteral() != fmt.Sprintf("%t",v) {
		t.Errorf("bo.TokenLiteral not %t. got=%t",v,bo.Value)
		return false
	}

	return true
}

func checkParserErrors(t *testing.T,p *Parser) {
	errors := p.Errors()
	if len(errors)==0 {
		return
	}
	t.Errorf("parser has %d errors",len(errors))
	for _,msg := range errors {
		t.Errorf("parser error: %q",msg)
	}
	t.FailNow()
}