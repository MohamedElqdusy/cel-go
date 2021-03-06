// Copyright 2018 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package checker

import (
	"fmt"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/golang/protobuf/ptypes/struct"
	"github.com/google/cel-spec/proto/checked/v1/checked"
)

const (
	kindUnknown = iota + 1
	kindError
	kindFunction
	kindDyn
	kindPrimitive
	kindWellKnown
	kindWrapper
	kindNull
	kindAbstract // TODO: Update the checked.proto to include abstract
	kindType
	kindList
	kindMap
	kindObject
	kindTypeParam
)

var (
	// Commonly used types.
	Error = &checked.Type{
		TypeKind: &checked.Type_Error{
			Error: &empty.Empty{}}}
	Dyn = &checked.Type{
		TypeKind: &checked.Type_Dyn{
			Dyn: &empty.Empty{}}}
	Null = &checked.Type{
		TypeKind: &checked.Type_Null{
			Null: structpb.NullValue_NULL_VALUE}}
	Bool   = newPrimitive(checked.Type_BOOL)
	Bytes  = newPrimitive(checked.Type_BYTES)
	String = newPrimitive(checked.Type_STRING)
	Double = newPrimitive(checked.Type_DOUBLE)
	Int    = newPrimitive(checked.Type_INT64)
	Uint   = newPrimitive(checked.Type_UINT64)

	// Well-known types.
	// TODO: Replace with an abstract type registry.
	Any       = newWellKnown(checked.Type_ANY)
	Duration  = newWellKnown(checked.Type_DURATION)
	Timestamp = newWellKnown(checked.Type_TIMESTAMP)
)

func newFunction(resultType *checked.Type, argTypes ...*checked.Type) *checked.Type {
	return &checked.Type{
		TypeKind: &checked.Type_Function{
			Function: &checked.Type_FunctionType{
				ResultType: resultType,
				ArgTypes:   argTypes}}}
}

func newTypeParam(name string) *checked.Type {
	return &checked.Type{
		TypeKind: &checked.Type_TypeParam{
			TypeParam: name}}
}

func newType(nested *checked.Type) *checked.Type {
	return &checked.Type{
		TypeKind: &checked.Type_Type{
			Type: nested}}
}

func newPrimitive(primitive checked.Type_PrimitiveType) *checked.Type {
	return &checked.Type{
		TypeKind: &checked.Type_Primitive{
			Primitive: primitive}}
}

func newWellKnown(wellKnown checked.Type_WellKnownType) *checked.Type {
	return &checked.Type{
		TypeKind: &checked.Type_WellKnown{
			WellKnown: wellKnown}}
}

func newWrapper(wrapped *checked.Type) *checked.Type {
	primitive := wrapped.GetPrimitive()
	if primitive == checked.Type_PRIMITIVE_TYPE_UNSPECIFIED {
		// TODO: return an error
		panic("Wrapped type must be a primitive")
	}
	return &checked.Type{
		TypeKind: &checked.Type_Wrapper{
			Wrapper: primitive}}
}

func newList(elem *checked.Type) *checked.Type {
	return &checked.Type{
		TypeKind: &checked.Type_ListType_{
			ListType: &checked.Type_ListType{
				ElemType: elem}}}
}

func newMap(key *checked.Type, value *checked.Type) *checked.Type {
	return &checked.Type{
		TypeKind: &checked.Type_MapType_{
			MapType: &checked.Type_MapType{
				KeyType:   key,
				ValueType: value}}}
}

func newObject(typeName string) *checked.Type {
	return &checked.Type{
		TypeKind: &checked.Type_MessageType{
			MessageType: typeName}}
}

func kindOf(t *checked.Type) int {
	switch t.TypeKind.(type) {
	case *checked.Type_Error:
		return kindError
	case *checked.Type_Function:
		return kindFunction
	case *checked.Type_Dyn:
		return kindDyn
	case *checked.Type_Primitive:
		return kindPrimitive
	case *checked.Type_WellKnown:
		return kindWellKnown
	case *checked.Type_Wrapper:
		return kindWrapper
	case *checked.Type_Null:
		return kindNull
	case *checked.Type_Type:
		return kindType
	case *checked.Type_ListType_:
		return kindList
	case *checked.Type_MapType_:
		return kindMap
	case *checked.Type_MessageType:
		return kindObject
	case *checked.Type_TypeParam:
		return kindTypeParam
	}
	return kindUnknown
}

/** Returns the more general of two types which are known to unify. */
func mostGeneral(t1 *checked.Type, t2 *checked.Type) *checked.Type {
	if isEqualOrLessSpecific(t1, t2) {
		return t1
	}
	return t2
}

func isAssignable(m *mapping, t1 *checked.Type, t2 *checked.Type) *mapping {
	mCopy := m.copy()
	if internalIsAssignable(mCopy, t1, t2) {
		return mCopy
	}
	return nil
}

func isAssignableList(m *mapping, l1 []*checked.Type, l2 []*checked.Type) *mapping {
	mCopy := m.copy()
	if internalIsAssignableList(mCopy, l1, l2) {
		return mCopy
	}
	return nil
}

/**
 * Apply substitution to given type, replacing all direct and indirect occurrences of bound type
 * parameters. Unbound type parameters are replaced by DYN if typeParamToDyn is true.
 */
func substitute(m *mapping, t *checked.Type, typeParamToDyn bool) *checked.Type {
	if tSub, found := m.find(t); found {
		return substitute(m, tSub, typeParamToDyn)
	}
	kind := kindOf(t)
	if typeParamToDyn && kind == kindTypeParam {
		return Dyn
	}
	switch kind {
	case kindType:
		return newType(substitute(m, t.GetType(), typeParamToDyn))
	case kindList:
		return newList(substitute(m, t.GetListType().ElemType, typeParamToDyn))
	case kindMap:
		mt := t.GetMapType()
		return newMap(substitute(m, mt.KeyType, typeParamToDyn),
			substitute(m, mt.ValueType, typeParamToDyn))
	case kindFunction:
		fn := t.GetFunction()
		rt := substitute(m, fn.ResultType, typeParamToDyn)
		args := make([]*checked.Type, len(fn.ArgTypes))
		for i, a := range fn.ArgTypes {
			args[i] = substitute(m, a, typeParamToDyn)
		}
		return newFunction(rt, args...)
	default:
		return t
	}
}

func internalIsAssignableList(m *mapping, l1 []*checked.Type, l2 []*checked.Type) bool {
	if len(l1) != len(l2) {
		return false
	}
	for i, t1 := range l1 {
		if !internalIsAssignable(m, t1, l2[i]) {
			return false
		}
	}
	return true
}

func internalIsAssignable(m *mapping, t1 *checked.Type, t2 *checked.Type) bool {
	// Process type parameters.
	kind1, kind2 := kindOf(t1), kindOf(t2)
	if kind2 == kindTypeParam {
		if t2Sub, found := m.find(t2); found {
			// Adjust the existing substitution to a more common type if possible. This is sound
			// because any previous substitution will be compatible with the common type. This
			// deals with the case the we have e.g. A -> int assigned, but now encounter a test
			// against DYN, and want to widen A to DYN.
			if isEqualOrLessSpecific(t1, t2Sub) && notReferencedIn(t2, t1) {
				m.add(t2, t1)
				return true
			} else {
				// Continue regular process with the assignment for type2.
				return internalIsAssignable(m, t1, t2Sub)
			}
		} else if notReferencedIn(t2, t1) {
			m.add(t2, t1)
			return true
		}
	}

	if kind1 == kindTypeParam {
		// For the lower type bound, we currently do not perform adjustment. The restricted
		// way we use type parameters in lower type bounds, it is not necessary, but may
		// become if we generalize type unification.
		if t1Sub, found := m.find(t1); found {
			return internalIsAssignable(m, t1Sub, t2)
		} else if notReferencedIn(t1, t2) {
			m.add(t1, t2)
			return true
		}
	}

	if kind1 == kindDyn || kind1 == kindError {
		return true
	}
	if kind2 == kindDyn || kind2 == kindError {
		return true
	}
	if kind1 == kindNull && isNullable(kind2) {
		return true
	}
	// Unwrap box types
	if kind1 == kindWrapper {
		return internalIsAssignable(m, newPrimitive(t1.GetWrapper()), t2)
	}
	// Finally check equality and type args recursively.
	if kind1 != kind2 {
		return false
	}

	switch kind1 {
	case kindPrimitive, kindWellKnown, kindObject:
		return proto.Equal(t1, t2)
	case kindType:
		return internalIsAssignable(m, t1.GetType(), t2.GetType())
	case kindList:
		return internalIsAssignable(m, t1.GetListType().ElemType, t2.GetListType().ElemType)
	case kindMap:
		m1 := t1.GetMapType()
		m2 := t2.GetMapType()
		return internalIsAssignableList(m,
			[]*checked.Type{m1.KeyType, m1.ValueType},
			[]*checked.Type{m2.KeyType, m2.ValueType})
	case kindFunction:
		fn1 := t1.GetFunction()
		fn2 := t2.GetFunction()
		return internalIsAssignableList(m,
			append(fn1.ArgTypes, fn1.ResultType),
			append(fn2.ArgTypes, fn2.ResultType))
	default:
		return false
	}
}

/**
 * Check whether one type is equal or less specific than the other one. A type is less specific if
 * it matches the other type using the DYN type.
 */
func isEqualOrLessSpecific(t1 *checked.Type, t2 *checked.Type) bool {
	kind1, kind2 := kindOf(t1), kindOf(t2)
	if kind1 == kindDyn || kind1 == kindTypeParam {
		return true
	}
	if kind2 == kindDyn || kind2 == kindTypeParam {
		return false
	}
	if kind1 != kind2 {
		return false
	}

	switch kind1 {
	case kindObject, kindPrimitive, kindWellKnown, kindWrapper:
		return proto.Equal(t1, t2)
	case kindType:
		return isEqualOrLessSpecific(t1.GetType(), t2.GetType())
	case kindList:
		return isEqualOrLessSpecific(t1.GetListType().ElemType, t2.GetListType().ElemType)
	case kindMap:
		m1 := t1.GetMapType()
		m2 := t2.GetMapType()
		return isEqualOrLessSpecific(m1.KeyType, m2.KeyType) &&
			isEqualOrLessSpecific(m1.KeyType, m2.KeyType)
	case kindFunction:
		fn1 := t1.GetFunction()
		fn2 := t2.GetFunction()
		if len(fn1.ArgTypes) != len(fn2.ArgTypes) {
			return false
		}
		if !isEqualOrLessSpecific(fn1.ResultType, fn2.ResultType) {
			return false
		}
		for i, a1 := range fn1.ArgTypes {
			if !isEqualOrLessSpecific(a1, fn2.ArgTypes[i]) {
				return false
			}
		}
		return true
	default:
		return true
	}
}

func isNullable(kind int) bool {
	switch kind {
	case kindObject, kindWrapper, kindWellKnown:
		return true
	default:
		return false
	}
}

func notReferencedIn(t *checked.Type, withinType *checked.Type) bool {
	if proto.Equal(t, withinType) {
		return false
	}
	withinKind := kindOf(withinType)
	switch withinKind {
	case kindWrapper:
		return notReferencedIn(t, newPrimitive(withinType.GetWrapper()))
	case kindType:
		return notReferencedIn(t, withinType.GetType())
	case kindList:
		return notReferencedIn(t, withinType.GetListType().ElemType)
	case kindMap:
		m := withinType.GetMapType()
		return notReferencedIn(t, m.KeyType) && notReferencedIn(t, m.ValueType)
	case kindFunction:
		fn := withinType.GetFunction()
		types := append(fn.ArgTypes, fn.ResultType)
		for _, a := range types {
			if !notReferencedIn(t, a) {
				return false
			}
		}
		return true
	default:
		return true
	}
}

func typeKey(t *checked.Type) string {
	return fmt.Sprintf("%v:%v", kindOf(t), t.String())
}

func FormatCheckedType(t *checked.Type) string {
	switch kindOf(t) {
	case kindPrimitive:
		switch t.GetPrimitive() {
		case checked.Type_UINT64:
			return "uint"
		case checked.Type_INT64:
			return "int"
		}
		return strings.Trim(strings.ToLower(t.GetPrimitive().String()), " ")
	case kindFunction:
		return formatFunction(t.GetFunction().GetResultType(), t.GetFunction().GetArgTypes(), false)
	case kindWrapper:
		return fmt.Sprintf("wrapper(%s)", FormatCheckedType(newPrimitive(t.GetWrapper())))
	case kindObject:
		return t.GetMessageType()
	case kindList:
		return fmt.Sprintf("list(%s)", FormatCheckedType(t.GetListType().ElemType))
	case kindMap:
		return fmt.Sprintf("map(%s, %s)",
			FormatCheckedType(t.GetMapType().KeyType),
			FormatCheckedType(t.GetMapType().ValueType))
	case kindNull:
		return "null"
	case kindDyn:
		return "dyn"
	case kindType:
		return fmt.Sprintf("type(%s)", FormatCheckedType(t.GetType()))
	case kindError:
		return "!error!"
	}
	return t.String()
}
