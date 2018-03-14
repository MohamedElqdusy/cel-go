package interpreter

import (
	"github.com/google/cel-go/interpreter/functions"
	testExpr "github.com/google/cel-go/interpreter/testing"
	"github.com/google/cel-go/interpreter/types"
	"github.com/google/cel-go/interpreter/types/adapters"
	"github.com/google/cel-spec/proto/v1"
	"testing"
)

func TestInterpreter_CallExpr(t *testing.T) {
	interpreter := StandardTestInterpreter()
	program := NewProgram("google.api.expr",
		testExpr.Equality.Expr,
		testExpr.Equality.Info(t.Name()))
	interpretable := interpreter.NewInterpretable(program)
	result, state := interpretable.Eval(
		NewActivation(map[string]interface{}{"a": int64(41)}))
	if result != false {
		t.Errorf("Expected false, go: %v", result)
	}
	if ident, found := state.Value(1); !found || ident != int64(41) {
		t.Errorf("State of ident 'a' != 41, got: %v", ident)
	}
}

func TestInterpreter_SelectExpr(t *testing.T) {
	interpreter := StandardTestInterpreter()
	program := NewProgram("",
		testExpr.Select.Expr,
		testExpr.Select.Info(t.Name()))

	interpretable := interpreter.NewInterpretable(program)
	result, _ := interpretable.Eval(
		NewActivation(map[string]interface{}{"a.b": adapters.NewMapAdapter(map[string]bool{"c": true})}))
	if result != true {
		t.Errorf("Expected true, got: %v", result)
	}
}

func TestInterpreter_ConditionalExpr(t *testing.T) {
	// a ? b < 1.0 : c == ["hello"]
	interpreter := StandardTestInterpreter()
	program := NewProgram("",
		testExpr.Conditional.Expr,
		testExpr.Conditional.Info(t.Name()))

	interpretable := interpreter.NewInterpretable(program)
	result, _ := interpretable.Eval(
		NewActivation(map[string]interface{}{
			"a": true,
			"b": 0.999}))
	if result != true {
		t.Errorf("Expected true, got: %v", result)
	}
}

func TestInterpreter_ComprehensionExpr(t *testing.T) {
	// [1, 1u, 1.0].exists(x, type(x) == uint)
	interpreter := StandardTestInterpreter()
	program := NewProgram("",
		testExpr.Exists.Expr,
		testExpr.Exists.Info(t.Name()))

	interpretable := interpreter.NewInterpretable(program)
	// TODO: make the type identifiers part of the standard declaration set.
	result, _ := interpretable.Eval(
		NewActivation(map[string]interface{}{"uint": types.UintType}))
	if result != true {
		t.Errorf("Expected true, got: %v", result)
	}
}

func BenchmarkInterpreter_ConditionalExpr(b *testing.B) {
	// a ? b < 1.0 : c == ["hello"]
	interpreter := StandardTestInterpreter()
	program := NewProgram("",
		testExpr.Conditional.Expr,
		testExpr.Conditional.Info(b.Name()))
	interpretable := interpreter.NewInterpretable(program)
	activation := NewActivation(map[string]interface{}{
		"a": false,
		"b": 0.999,
		"c": adapters.NewListAdapter([]string{"hello"})})
	for i := 0; i < b.N; i++ {
		interpretable.Eval(activation)
	}
}

func BenchmarkInterpreter_ComprehensionExpr(b *testing.B) {
	// [1, 1u, 1.0].exists(x, type(x) == uint)
	interpreter := StandardTestInterpreter()
	program := NewProgram("",
		testExpr.Exists.Expr,
		testExpr.Exists.Info(b.Name()))
	interpretable := interpreter.NewInterpretable(program)
	activation := NewActivation(map[string]interface{}{"uint": types.UintType})
	for i := 0; i < b.N; i++ {
		interpretable.Eval(activation)
	}
}

func StandardTestInterpreter() Interpreter {
	dispatcher := NewDispatcher()
	dispatcher.Add(functions.StandardBuiltins()...)
	typeProvider := types.NewTypeProvider(&syntax_proto.ParsedExpr{})
	return NewInterpreter(dispatcher, typeProvider)
}
