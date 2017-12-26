package binary

import (
	"fmt"
	"reflect"
)

const (
	//disable BinarySerializer check by default
	defaultSerializer = true
)

//DefaultSerializer return default BinarySerializer check by default
func DefaultSerializer() bool {
	return defaultSerializer
}

// BinarySizer is an interface to define go data Size method.
type BinarySizer interface {
	Size() int
}

// BinaryEncoder is an interface to define go data Encode method.
// buffer is nil-able.
type BinaryEncoder interface {
	Encode(buffer []byte) ([]byte, error)
}

// BinaryDecoder is an interface to define go data Decode method.
type BinaryDecoder interface {
	Decode(buffer []byte) error
}

// BinarySerializer defines the go data Size/Encode/Decode method
type BinarySerializer interface {
	BinarySizer
	BinaryEncoder
	BinaryDecoder
}

var (
	tSizer      reflect.Type //BinarySizer
	tEncoder    reflect.Type //BinaryEncoder
	tDecoder    reflect.Type //BinaryDecoder
	tSerializer reflect.Type //BinarySerializer
)

func init() {
	var sizer BinarySizer
	tSizer = reflect.TypeOf(&sizer).Elem()
	var encoder BinaryEncoder
	tEncoder = reflect.TypeOf(&encoder).Elem()
	var decoder BinaryDecoder
	tDecoder = reflect.TypeOf(&decoder).Elem()
	var serializer BinarySerializer
	tSerializer = reflect.TypeOf(&serializer).Elem()
}

// SerializerSwitch defines switch of BinarySerializer check
type serializerSwitch byte

const (
	serializerDisable    serializerSwitch = iota // disable Serializer
	serializerCheck                              // enable Serializer but need check
	serializerCheckFalse                         // enable and do not need check,result false
	serializerCheckOk                            // enable and do not need check,result true
)

// String return name of this switch
func (ss serializerSwitch) String() string {
	switch ss {
	case serializerDisable:
		return "SerializerDisable"
	case serializerCheck:
		return "SerializerCheck"
	case serializerCheckOk:
		return "SerializerCheckOk"
	}
	panic(fmt.Errorf("SerializerUnknown"))
}

// Enable returns if BinarySerializer check is enable
func (ss serializerSwitch) enable() bool {
	return ss >= serializerCheck
}

// NeedCheck returns if need check BinarySerializer
func (ss serializerSwitch) needCheck() bool {
	return ss == serializerCheck
}

// CheckFail returns if can use BinarySerializer directly
func (ss serializerSwitch) checkFalse() bool {
	return ss == serializerCheckFalse
}

// NeedCheck returns if can use BinarySerializer directly
func (ss serializerSwitch) checkOk() bool {
	return ss == serializerCheckOk
}

// SubSwitch returns SerializerSwitch for sub-data of struct/map/slice/array
func (ss serializerSwitch) subSwitchCheck(t reflect.Type) serializerSwitch {
	if !ss.enable() {
		return serializerDisable
	}
	return ss.subSwitch(querySerializer(indirectType(t)))
}

func (ss serializerSwitch) subSwitch(isSerializer bool) serializerSwitch {
	if !ss.enable() {
		return serializerDisable
	}
	if isSerializer {
		return serializerCheckOk
	}
	return serializerCheckFalse
}

func toplvSerializer(enable bool) serializerSwitch {
	if enable {
		return serializerCheck
	}
	return serializerDisable
}

//CheckSerializer check if t implements BinarySerializer
func CheckSerializer(x interface{}) bool {
	return querySerializer(reflect.TypeOf(x))
}

//CheckSerializerDeep check if t or *t implements BinarySerializer
func CheckSerializerDeep(x interface{}) bool {
	return querySerializer(indirectType(reflect.TypeOf(x)))
}
