package vm

import (
	"fmt"
	"monkey/ast"
	"monkey/compiler"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"testing"
)

func parse(input string) *ast.Program {
	l := lexer.New(input)
	p := parser.New(l)
	return p.ParseProgram()
}

func testIntegerObject(expected int64, actual object.Object) error {
	obj, ok := actual.(*object.Integer)
	if !ok {
		return fmt.Errorf("object is not Integer. got=%T (%+v)", actual, actual)
	}
	fmt.Println(obj)
	if obj.Value != expected {
		return fmt.Errorf("object has wrong value. got=%d, want=%d", obj.Value, expected)
	}
	return nil
}

type vmTestCase struct {
	input    string
	expected interface{}
}

func runVmTests(t *testing.T, tests []vmTestCase) {
	t.Helper()
	for _, tt := range tests {
		program := parse(tt.input)
		comp := compiler.New()
		err := comp.Compile(program)
		if err != nil {
			t.Fatalf("compiler error: %s", err)
		}
		machine := New(comp.Bytecode())
		err = machine.Run()
		if err != nil {
			t.Fatalf("vm error: %s", err)
		}
		stackElem := machine.LastPoppedStackElem()
		testExpectedObject(t, tt.expected, stackElem)
	}
}

func testExpectedObject(t *testing.T, expected interface{}, actual object.Object) {
	t.Helper()

	switch expected := expected.(type) {
	case int:
		err := testIntegerObject(int64(expected), actual)
		if err != nil {
			t.Errorf("testIntegerObject failed. %s", err)
		}
	case bool:
		err := testBooleanObject(expected, actual)
		if err != nil {
			t.Errorf("testBooleanObject failed. %s", err)
		}
	case *object.Null:
		if actual != Null {
			t.Errorf("object is not Null, got=%T (%+v)", actual, actual)
		}
	}
}
func testBooleanObject(expected bool, actual object.Object) error {
	result, ok := actual.(*object.Boolean)
	if !ok {
		return fmt.Errorf("object is not boolean, got=%T (%+v)", actual, actual)
	}
	if result.Value != expected {
		return fmt.Errorf("object has wrong value. got=%t, want=%t", result.Value, expected)
	}
	return nil
}

func TestIntegerArithmetic(t *testing.T) {
	tests := []vmTestCase{
		{"1", 1},
		{"2", 2},
		{"1 + 2", 3},
		{"1 - 2", -1},
		{"1 * 2", 2},
		{"5 * (2+10)", 60},
		{"true", true},
		{"1>2", false},
		{"1==2", false},
		{"1==1", true},
		{"1<2", true},
		{"(1>2)==false", true},
		{"-10", -10},
		{"-50+100+-50", 0},
		{"!true", false},
		{"!!false", false},
		{"!!5", true},
		{"!(if (false) {5})", true},
		{"if ((if (false) {10})) {10} else{20}", 20},
		{"let one=1;one;", 1},
		{"let one=1;let two=2;one+two;", 3},
	}

	runVmTests(t, tests)
}

func TestConditionals(t *testing.T) {
	tests := []vmTestCase{
		{"if (1>2) {10} else {20}", 20},
		{"if (2>1) {10} else {20}", 10},
		{"if (2>1) {10}", 10},
		{"if (1) {10}", 10},
		{"if (1>2) {10}", Null},
		{"if (false) {10}", Null},
	}

	runVmTests(t, tests)
}
