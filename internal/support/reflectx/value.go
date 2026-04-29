package reflectx

import "reflect"

func Interface(inputValue reflect.Value) any {
	if !inputValue.IsValid() {
		return nil
	}

	return inputValue.Interface()
}
