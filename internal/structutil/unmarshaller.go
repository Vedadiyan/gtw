/*
	Extracted from github.com/vedadiyan/goal/pkg/structutil
	DO NOT EDIT
*/

package structutil

import (
	"fmt"
	"reflect"
	_ "unsafe"
)

type (
	Field         = reflect.StructField
	Data          = map[string]any
	Value         = reflect.Value
	RC            int
	UnMarshallers map[int]func(d Data, f Field, v Value, rc RC) error
)

var (
	_u UnMarshallers
)

func UseSignedNumerics(u UnMarshallers) {
	u[int(reflect.Int8)] = UMPrimitive[int8]
	u[int(reflect.Int16)] = UMPrimitive[int16]
	u[int(reflect.Int32)] = UMPrimitive[int32]
	u[int(reflect.Int)] = UMPrimitive[int]
	u[int(reflect.Int64)] = UMPrimitive[int64]
	u[int(reflect.Int8)*100] = UMArray[int8]
	u[int(reflect.Int16)*100] = UMArray[int16]
	u[int(reflect.Int32)*100] = UMArray[int32]
	u[int(reflect.Int)*100] = UMArray[int]
	u[int(reflect.Int64)*100] = UMArray[int64]
}

func UseUnSignedNumerics(u UnMarshallers) {
	u[int(reflect.Uint8)] = UMPrimitive[uint8]
	u[int(reflect.Uint16)] = UMPrimitive[uint16]
	u[int(reflect.Uint32)] = UMPrimitive[uint32]
	u[int(reflect.Uint)] = UMPrimitive[uint]
	u[int(reflect.Uint64)] = UMPrimitive[uint64]
	u[int(reflect.Uint8)*100] = UMArray[uint8]
	u[int(reflect.Uint16)*100] = UMArray[uint16]
	u[int(reflect.Uint32)*100] = UMArray[uint32]
	u[int(reflect.Uint)*100] = UMArray[uint]
	u[int(reflect.Uint64)*100] = UMArray[uint64]
}

func UseFloatingPoints(u UnMarshallers) {
	u[int(reflect.Float32)] = UMPrimitive[float32]
	u[int(reflect.Float64)] = UMPrimitive[float64]
	u[int(reflect.Float32)*100] = UMArray[float32]
	u[int(reflect.Float64)*100] = UMArray[float64]
}

func UsetOtherPrimitives(u UnMarshallers) {
	u[int(reflect.Bool)] = UMPrimitive[bool]
	u[int(reflect.String)] = UMPrimitive[string]
	u[int(reflect.Bool)*100] = UMArray[bool]
	u[int(reflect.String)*100] = UMArray[string]
}

func UsetStructs(u UnMarshallers) {
	u[int(reflect.Struct)] = UMStruct
	u[int(reflect.Struct)*100] = UMStructArray
}

func UseMaps(u UnMarshallers) {
	u[int(reflect.Map)] = UMMap
	u[int(reflect.Map)*100] = UMMapArray
}

func UsePointers(u UnMarshallers) {
	u[int(reflect.Pointer)] = UMPointer
	u[int(reflect.Pointer)*100] = UMPointerArray
	u[int(reflect.Interface)] = func(d Data, f Field, v Value, rc RC) error {
		return nil
	}
}

func UserSlices(u UnMarshallers) {
	u[int(reflect.Slice)*100] = UMSlice
}

func init() {
	_u = make(map[int]func(d Data, f Field, v Value, rc RC) error)
	UseSignedNumerics(_u)
	UseUnSignedNumerics(_u)
	UseFloatingPoints(_u)
	UsetOtherPrimitives(_u)
	UsetStructs(_u)
	UseMaps(_u)
	UsePointers(_u)
	UserSlices(_u)
}

func (rc *RC) Len() int {
	return int(*rc)
}

func (rc *RC) Incr() RC {
	return RC(*rc + 1)
}

func Protect(err *error) {
	if r := recover(); r != nil {
		if r, ok := r.(error); ok {
			*err = r
		}
		*err = fmt.Errorf("%v", r)
	}
}

func GetKind(f reflect.StructField) int {
	field := f.Type
	if field.Kind() == reflect.Slice {
		return int(field.Elem().Kind()) * 100
	}
	return int(field.Kind())
}

func GetKindRaw(f reflect.Kind) int {
	if f == reflect.Slice {
		return int(f) * 100
	}
	return int(f)
}

func GetFieldName(field reflect.StructField) string {
	name := field.Tag.Get("name")
	if len(name) == 0 {
		name = field.Name
	}
	return name
}

func HasValue(v any, b bool) bool {
	return b && v != nil
}

func MustBe[T any](v any) (T, error) {
	value, ok := v.(T)
	if !ok {
		return value, fmt.Errorf("expected %T by found %T", value, v)
	}
	return value, nil
}

func KindOf[T any]() reflect.Kind {
	return reflect.ValueOf(new(T)).Elem().Kind()
}

func IsPointer(rc RC) bool {
	return rc > 0
}

func SameType(f Field, v any) bool {
	return f.Type == reflect.TypeOf(v)
}

func GetDimensions(f reflect.StructField) (reflect.Type, int) {
	dimensions := 0
	t := f.Type.Elem()
	for t.Kind() == reflect.Slice {
		t = t.Elem()
		dimensions++
	}
	return t, dimensions
}

func GetRferenceCount(f reflect.StructField) (reflect.Type, RC) {
	referenceCount := 0
	t := f.Type.Elem()
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
		referenceCount++
	}
	return t, RC(referenceCount)
}

func CreateSlice(value any, dimensions int, typeOfElem reflect.Type, typeOfArray reflect.Type, rc RC, itr int) (*reflect.Value, error) {
	valueOfData := reflect.ValueOf(value)
	switch valueOfData.Kind() {
	case reflect.Map:
		{
			valueOfElem := reflect.New(typeOfElem)
			switch valueOfElem.Elem().Kind() {
			case reflect.Map, reflect.Interface:
				{
					return &valueOfData, nil
				}
			default:
				{
					err := Unmarshal(value.(map[string]any), valueOfElem.Interface())
					if err != nil {
						return nil, err
					}
					output := valueOfElem.Elem()
					return &output, nil
				}
			}
		}
	case reflect.Slice:
		{
			typeOfSlice := reflect.SliceOf(typeOfArray)
			for i := 0; i < dimensions-itr; i++ {
				typeOfSlice = reflect.SliceOf(typeOfSlice)
			}
			valueOfSlice := reflect.MakeSlice(typeOfSlice, 0, 0)
			for i := 0; i < valueOfData.Len(); i++ {
				ref := valueOfData.Index(i).Interface()
				next, err := CreateSlice(ref, dimensions, typeOfElem, typeOfArray, rc, itr+1)
				if err != nil {
					return nil, err
				}
				v := reflect.New(typeOfArray)
				Set(next.Interface(), v.Elem(), rc)
				valueOfSlice = reflect.Append(valueOfSlice, v.Elem())
			}
			return &valueOfSlice, nil
		}
	default:
		{
			return &valueOfData, nil
		}
	}
}

func CreateMap(value Data, tyepOfKey reflect.Type, typeOfValue reflect.Type) (*reflect.Value, error) {
	typeOfMap := reflect.MapOf(tyepOfKey, typeOfValue)
	valueOfMap := reflect.MakeMap(typeOfMap)
	for key, value := range value {
		kind := GetKindRaw(typeOfValue.Kind())
		switch kind {
		case int(reflect.Struct):
			{
				ref := reflect.New(typeOfValue)
				err := Unmarshal(value.(map[string]any), ref.Interface())
				if err != nil {
					return nil, err
				}
				valueOfMap.SetMapIndex(reflect.ValueOf(key), ref.Elem())
			}
		default:
			{
				valueOfMap.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(value))
			}
		}
	}
	return &valueOfMap, nil
}

func CreateStruct(item any, f reflect.Type) (*reflect.Value, error) {
	value, err := MustBe[map[string]any](item)
	if err != nil {
		return nil, err
	}
	ref := reflect.New(f)
	err = Unmarshal(value, ref.Interface())
	if err != nil {
		return nil, err
	}
	return &ref, nil
}

func Set(value any, v Value, rc RC) {
	if IsPointer(rc) {
		typeOfP := reflect.TypeOf(value)
		valueOfP := reflect.New(typeOfP)
		valueOfP.Elem().Set(reflect.ValueOf(reflect.ValueOf(&value).Elem().Interface()))
		for i := 1; i < rc.Len(); i++ {
			typeOfP = reflect.PointerTo(typeOfP)
			tmp := reflect.New(typeOfP)
			tmp.Elem().Set(valueOfP)
			valueOfP = tmp
		}
		v.Set(valueOfP)
		return
	}
	v.Set(reflect.ValueOf(value))
}

func UMPrimitive[T any](d Data, f Field, v Value, rc RC) (_err error) {
	defer Protect(&_err)
	entry, ok := d[GetFieldName(f)]
	if !HasValue(entry, ok) {
		return nil
	}
	prod, err := _convertors[KindOf[T]()](entry)
	if err != nil {
		return err
	}
	Set(prod.(T), v, rc)
	return nil
}

func UMArray[T any](d Data, f Field, v Value, rc RC) (_err error) {
	defer Protect(&_err)
	entry, ok := d[GetFieldName(f)]
	if !HasValue(entry, ok) {
		return nil
	}
	values, err := MustBe[[]any](entry)
	if err != nil {
		return err
	}
	kind := KindOf[T]()
	slice := make([]T, 0)
	for _, item := range values {
		prod, err := _convertors[kind](item)
		if err != nil {
			return err
		}
		slice = append(slice, prod.(T))
	}
	Set(slice, v, rc)
	return nil
}

func UMStruct(d Data, f Field, v Value, rc RC) (_err error) {
	defer Protect(&_err)
	entry, ok := d[GetFieldName(f)]
	if !HasValue(entry, ok) {
		return nil
	}
	if SameType(f, entry) {
		v.Set(reflect.ValueOf(entry))
		return nil
	}
	prod, err := CreateStruct(entry, f.Type)
	if err != nil {
		return err
	}
	Set(prod.Elem().Interface(), v, rc)
	return nil
}

func UMStructArray(d Data, f Field, v Value, rc RC) (_err error) {
	defer Protect(&_err)
	entry, ok := d[GetFieldName(f)]
	if !HasValue(entry, ok) {
		return nil
	}
	if SameType(f, entry) {
		v.Set(reflect.ValueOf(entry))
		return nil
	}
	values, err := MustBe[[]any](entry)
	if err != nil {
		return err
	}
	slice := reflect.MakeSlice(reflect.SliceOf(f.Type.Elem()), 0, 0)
	for _, item := range values {
		if SameType(f, item) {
			slice = reflect.Append(slice, reflect.ValueOf(item))
			continue
		}
		prod, err := CreateStruct(item, f.Type.Elem())
		if err != nil {
			return err
		}
		slice = reflect.Append(slice, prod.Elem())
	}
	v.Set(slice)
	return nil
}

func UMMap(d Data, f Field, v Value, rc RC) (_err error) {
	defer Protect(&_err)
	entry, ok := d[GetFieldName(f)]
	if !HasValue(entry, ok) {
		return nil
	}
	value, err := MustBe[map[string]any](entry)
	if err != nil {
		return err
	}
	prod, err := CreateMap(value, f.Type.Key(), f.Type.Elem())
	if err != nil {
		return err
	}
	Set(prod.Interface(), v, rc)
	return nil
}

func UMMapArray(d Data, f Field, v Value, rc RC) (_err error) {
	defer Protect(&_err)
	entry, ok := d[GetFieldName(f)]
	if !HasValue(entry, ok) {
		return nil
	}
	values, err := MustBe[[]any](entry)
	if err != nil {
		return err
	}
	slice := reflect.MakeSlice(reflect.SliceOf(f.Type.Elem()), 0, 0)
	for _, item := range values {
		prod, err := CreateMap(item.(map[string]any), f.Type.Elem().Key(), f.Type.Elem().Elem())
		if err != nil {
			return err
		}
		slice = reflect.Append(slice, *prod)
	}

	v.Set(slice)
	return nil
}

func UMSlice(d Data, f Field, v Value, rc RC) (_err error) {
	defer Protect(&_err)
	value, ok := d[GetFieldName(f)]
	if !HasValue(value, ok) {
		return nil
	}
	baseType, dimensions := GetDimensions(f)
	slice, err := CreateSlice(value, dimensions, baseType, baseType, 0, 0)
	if err != nil {
		return err
	}
	Set(slice.Interface(), v, rc)
	return nil
}

func UMPointer(d Data, f Field, v Value, rc RC) (_err error) {
	defer Protect(&_err)
	f.Type = f.Type.Elem()
	return _u[GetKindRaw(f.Type.Kind())](d, f, v, rc.Incr())
}

func UMPointerArray(d Data, f Field, v Value, rc RC) (error error) {
	defer Protect(&error)
	value, ok := d[GetFieldName(f)]
	if !HasValue(value, ok) {
		return nil
	}
	baseType, dimensions := GetDimensions(f)
	pointerType, referenceCount := GetRferenceCount(f)
	slice, err := CreateSlice(value, dimensions, pointerType, baseType, referenceCount, 0)
	if err != nil {
		return err
	}
	Set(slice.Interface(), v, rc)
	return nil
}

func Unmarshal(data map[string]any, message any) error {
	p := reflect.TypeOf(message).Elem()
	v := reflect.ValueOf(message).Elem()
	n := p.NumField()
	for i := 0; i < n; i++ {
		field := p.Field(i)
		kind := GetKind(field)
		err := _u[kind](data, field, v.Field(i), 0)
		if err != nil {
			return err
		}
	}
	return nil
}
