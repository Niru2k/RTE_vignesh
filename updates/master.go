package updates

import (
	//inbuilt package
	"reflect"
)

type Master struct{}

func (Master) Trigger(file string) error {
	values := reflect.ValueOf(Master{}).MethodByName(file).Call(nil)
	if len(values) > 0 && values[0].Type() == reflect.TypeOf((*error)(nil)).Elem() && values[0].Interface() != nil {
		return values[0].Interface().(error)
	}
	return nil
}
