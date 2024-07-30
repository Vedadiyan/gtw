/*
	Extracted from github.com/vedadiyan/goal/pkg/structutil
	DO NOT EDIT
*/

package structutil

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

var (
	_convertors map[reflect.Kind]func(value any) (any, error)
)

func init() {
	_convertors = make(map[reflect.Kind]func(value any) (any, error))

	_convertors[reflect.Int8] = ToInt8
	_convertors[reflect.Int16] = ToInt16
	_convertors[reflect.Int32] = ToInt32
	_convertors[reflect.Int] = ToInt
	_convertors[reflect.Int64] = ToInt64

	_convertors[reflect.Uint8] = ToUInt8
	_convertors[reflect.Uint16] = ToUInt16
	_convertors[reflect.Uint32] = ToUInt32
	_convertors[reflect.Uint] = ToUInt
	_convertors[reflect.Uint64] = ToUInt64

	_convertors[reflect.Float64] = ToFloat64
	_convertors[reflect.Float32] = ToFloat32

	_convertors[reflect.Complex64] = ToComplex64
	_convertors[reflect.Complex128] = ToComplex128

	_convertors[reflect.Bool] = ToBool
	_convertors[reflect.String] = ToString
}

func ToInt8(value any) (any, error) {
	if output, ok := value.(int8); ok {
		return output, nil
	}
	output, err := strconv.ParseInt(fmt.Sprintf("%v", value), 10, 8)
	if err != nil {
		return nil, err
	}
	return int8(output), nil
}

func ToInt16(value any) (any, error) {
	if output, ok := value.(int16); ok {
		return output, nil
	}
	output, err := strconv.ParseInt(fmt.Sprintf("%v", value), 10, 16)
	if err != nil {
		return nil, err
	}
	return int16(output), nil
}

func ToInt(value any) (any, error) {
	if output, ok := value.(int); ok {
		return output, nil
	}
	output, err := strconv.ParseInt(fmt.Sprintf("%v", value), 10, 32)
	if err != nil {
		return nil, err
	}
	return int(output), nil
}

func ToInt32(value any) (any, error) {
	if output, ok := value.(int32); ok {
		return output, nil
	}
	output, err := strconv.ParseInt(fmt.Sprintf("%v", value), 10, 32)
	if err != nil {
		return nil, err
	}
	return int32(output), nil
}

func ToInt64(value any) (any, error) {
	if output, ok := value.(int64); ok {
		return output, nil
	}
	output, err := strconv.ParseInt(fmt.Sprintf("%v", value), 10, 64)
	if err != nil {
		return nil, err
	}
	return int64(output), nil
}

func ToUInt8(value any) (any, error) {
	if output, ok := value.(uint8); ok {
		return output, nil
	}
	output, err := strconv.ParseUint(fmt.Sprintf("%v", value), 10, 8)
	if err != nil {
		return nil, err
	}
	return uint8(output), nil
}

func ToUInt16(value any) (any, error) {
	if output, ok := value.(uint16); ok {
		return output, nil
	}
	output, err := strconv.ParseUint(fmt.Sprintf("%v", value), 10, 16)
	if err != nil {
		return nil, err
	}
	return uint16(output), nil
}

func ToUInt(value any) (any, error) {
	if output, ok := value.(uint); ok {
		return output, nil
	}
	output, err := strconv.ParseUint(fmt.Sprintf("%v", value), 10, 32)
	if err != nil {
		return nil, err
	}
	return uint(output), nil
}

func ToUInt32(value any) (any, error) {
	if output, ok := value.(uint32); ok {
		return output, nil
	}
	output, err := strconv.ParseUint(fmt.Sprintf("%v", value), 10, 32)
	if err != nil {
		return nil, err
	}
	return uint32(output), nil
}

func ToUInt64(value any) (any, error) {
	if output, ok := value.(uint64); ok {
		return output, nil
	}
	output, err := strconv.ParseUint(fmt.Sprintf("%v", value), 10, 64)
	if err != nil {
		return nil, err
	}
	return uint64(output), nil
}

func ToFloat64(value any) (any, error) {
	if output, ok := value.(float64); ok {
		return output, nil
	}
	output, err := strconv.ParseFloat(fmt.Sprintf("%v", value), 64)
	if err != nil {
		return nil, err
	}
	return float64(output), nil
}

func ToFloat32(value any) (any, error) {
	if output, ok := value.(float64); ok {
		return output, nil
	}
	output, err := strconv.ParseFloat(fmt.Sprintf("%v", value), 32)
	if err != nil {
		return nil, err
	}
	return float64(output), nil
}

func ToComplex64(value any) (any, error) {
	if output, ok := value.(complex64); ok {
		return output, nil
	}
	output, err := strconv.ParseComplex(fmt.Sprintf("%v", value), 64)
	if err != nil {
		return nil, err
	}
	return complex64(output), nil
}

func ToComplex128(value any) (any, error) {
	if output, ok := value.(complex128); ok {
		return output, nil
	}
	output, err := strconv.ParseComplex(fmt.Sprintf("%v", value), 128)
	if err != nil {
		return nil, err
	}
	return complex128(output), nil
}

func ToBool(value any) (any, error) {
	if output, ok := value.(bool); ok {
		return output, nil
	}
	output := strings.ToLower(fmt.Sprintf("%v", value))
	switch output {
	case "true":
		{
			return true, nil
		}
	case "false":
		{
			return false, nil
		}
	default:
		{
			return false, fmt.Errorf("expected true or false but found %v", value)
		}
	}
}

func ToString(value any) (any, error) {
	if output, ok := value.(string); ok {
		return output, nil
	}
	return fmt.Sprintf("%v", value), nil
}
