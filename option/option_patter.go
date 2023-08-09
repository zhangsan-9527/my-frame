package option

import "errors"

type MyStructOption func(myStruct *MyStruct)
type MyStructOptionErr func(myStruct *MyStruct) error

type MyStruct struct {

	// 第一个部分是: 必须用户输入的字段
	id   uint64
	name string

	// 第二部分是: 可选择输入的字段
	address string
}

//var m = MyStruct{}

// NewMyStruct 参数包含所有必传字段
func NewMyStruct(id uint64, name string, opts ...MyStructOption) *MyStruct {
	res := &MyStruct{
		id:   id,
		name: name,
	}

	for _, opt := range opts {
		opt(res)
	}

	return res
}

func WithAddress(address string) MyStructOption {
	return func(myStruct *MyStruct) {
		myStruct.address = address
	}
}

func WithAddressV1(address string) MyStructOptionErr {
	return func(myStruct *MyStruct) error {
		if address == "" {
			return errors.New("地址为空")
		}
		myStruct.address = address
		return nil
	}
}

func WithAddressV2(address string) MyStructOption {
	return func(myStruct *MyStruct) {
		if address == "" {
			panic("地址为空")
		}
		myStruct.address = address
	}
}

func NewMyStructErr(id uint64, name string, opts ...MyStructOptionErr) (*MyStruct, error) {
	res := &MyStruct{
		id:   id,
		name: name,
	}

	for _, opt := range opts {
		err := opt(res)
		if err != nil {
			return nil, err
		}
	}

	return res, nil
}
