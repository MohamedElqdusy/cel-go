package adapters

import (
	"bytes"
	"github.com/google/cel-go/interpreter/types/objects"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/struct"
	expr "github.com/google/cel-spec/proto/v1"
	"reflect"
	"testing"
)

func TestExprToProto(t *testing.T) {
	// Core type conversion tests.
	expectExprToProto(t, true, true)
	expectExprToProto(t, int64(-1), int32(-1))
	expectExprToProto(t, int64(2), int64(2))
	expectExprToProto(t, uint64(3), uint32(3))
	expectExprToProto(t, uint64(4), uint64(4))
	expectExprToProto(t, float64(5.5), float32(5.5))
	expectExprToProto(t, float64(-5.5), float64(-5.5))
	expectExprToProto(t, "hello", "hello")
	expectExprToProto(t, []byte("world"), []byte("world"))
	expectExprToProto(t, []int64{1, 2, 3}, []int32{1, 2, 3})
	expectExprToProto(t, NewListAdapter([]int64{1, 2, 3}), []int32{1, 2, 3})
	expectExprToProto(t, map[int64]int64{1: 1, 2: 1, 3: 1},
		map[int32]int32{1: 1, 2: 1, 3: 1})
	expectExprToProto(t, NewMapAdapter(map[int64]int64{1: 1, 2: 1, 3: 1}),
		map[int32]int32{1: 1, 2: 1, 3: 1})

	// Null conversion tests.
	expectExprToProto(t, structpb.NullValue_NULL_VALUE, structpb.NullValue_NULL_VALUE)

	// Proto conversion tests.
	parsedExpr := &expr.ParsedExpr{}
	expectExprToProto(t, NewMsgAdapter(parsedExpr), parsedExpr)
	expectExprToProto(t, NewMsgAdapter(*parsedExpr), *parsedExpr)
}

func TestProtoToExpr(t *testing.T) {
	// Core type conversions.
	expectProtoToExpr(t, true, true)
	expectProtoToExpr(t, int32(-1), int64(-1))
	expectProtoToExpr(t, int64(2), int64(2))
	expectProtoToExpr(t, uint32(3), uint64(3))
	expectProtoToExpr(t, uint64(4), uint64(4))
	expectProtoToExpr(t, float32(5.5), float64(5.5))
	expectProtoToExpr(t, float64(-5.5), float64(-5.5))
	expectProtoToExpr(t, "hello", "hello")
	expectProtoToExpr(t, []byte("world"), []byte("world"))
	expectProtoToExpr(t, []int32{1, 2, 3},
		NewListAdapter([]int32{1, 2, 3}))
	expectProtoToExpr(t, map[int32]int32{1: 1, 2: 1, 3: 1},
		NewMapAdapter(map[int32]int32{1: 1, 2: 1, 3: 1}))

	// Null conversion test.
	expectProtoToExpr(t, structpb.NullValue_NULL_VALUE, structpb.NullValue_NULL_VALUE)

	// Proto conversion test.
	parsedExpr := &expr.ParsedExpr{}
	expectProtoToExpr(t, parsedExpr, NewMsgAdapter(parsedExpr))
}

func TestUnsupportedConversion(t *testing.T) {
	if val, err := ProtoToExpr(NonConvertiable{}); err == nil {
		t.Error("Expected error when converting non-proto struct to proto", val)
	}
	if val, err := ExprToProto(reflect.TypeOf(NonConvertiable{}), NonConvertiable{}); err == nil {
		t.Error("Expected error when converting non-proto expr to proto", val)
	}
}

func expectExprToProto(t *testing.T, in interface{}, out interface{}) {
	t.Helper()
	if val, err := ExprToProto(reflect.TypeOf(out), in); err != nil {
		t.Error(err)
	} else {
		var equals bool
		switch val.(type) {
		case []byte:
			equals = bytes.Equal(val.([]byte), out.([]byte))
		case proto.Message:
			equals = proto.Equal(val.(proto.Message), out.(proto.Message))
		case bool, int32, int64, uint32, uint64, float32, float64, string:
			equals = val == out
		default:
			equals = reflect.DeepEqual(val, out)
		}
		if !equals {
			t.Errorf("Unexpected conversion from expr to proto.\n"+
				"expected: %T, actual: %T", val, out)
		}
	}
}

func expectProtoToExpr(t *testing.T, in interface{}, out interface{}) {
	t.Helper()
	if val, err := ProtoToExpr(in); err != nil {
		t.Error(err)
	} else {
		var equals bool
		switch val.(type) {
		case []byte:
			equals = bytes.Equal(val.([]byte), out.([]byte))
		case proto.Message:
			equals = proto.Equal(val.(proto.Message), out.(proto.Message))
		case bool, int64, uint64, float64, string:
			equals = val == out
		case objects.Equaler:
			equals = val.(objects.Equaler).Equal(out)
		default:
			equals = reflect.DeepEqual(val, out)
		}
		if !equals {
			t.Errorf("Unexpected conversion from expr to proto.\n"+
				"expected: %T, actual: %T", val, out)
		}
	}
}

type NonConvertiable struct {
	Field string
}
