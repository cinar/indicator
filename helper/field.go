// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

import (
	"errors"
	"reflect"
)

// Field extracts a specific field from a channel of struct pointers and
// delivers it through a new channel.
func Field[T, S any](c <-chan *S, name string) (<-chan T, error) {
	st := reflect.TypeOf((*S)(nil)).Elem()
	if st.Kind() != reflect.Struct {
		return nil, errors.New("type not a struct")
	}

	f, ok := st.FieldByName(name)
	if !ok {
		return nil, errors.New("field is not found")
	}

	result := make(chan T, cap(c))

	go func() {
		defer close(result)

		for n := range c {
			v := reflect.ValueOf(n).Elem()
			result <- v.FieldByIndex(f.Index).Interface().(T)
		}
	}()

	return result, nil
}
