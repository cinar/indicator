// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

import (
	"context"
	"errors"
	"reflect"
)

// Field wraps FieldWithContext for backwards compatibility.
//
// Deprecated: Use FieldWithContext instead.
func Field[T, S any](c <-chan *S, name string) (<-chan T, error) {
	return FieldWithContext[T, S](context.Background(), c, name)
}

// FieldWithContext extracts a specific field from a channel of struct pointers and
// delivers it through a new channel, supporting context cancellation.
func FieldWithContext[T, S any](ctx context.Context, c <-chan *S, name string) (<-chan T, error) {
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

		for {
			select {
			case <-ctx.Done():
				return
			case n, ok := <-c:
				if !ok {
					return
				}
				v := reflect.ValueOf(n).Elem()
				val := v.FieldByIndex(f.Index).Interface().(T)
				select {
				case <-ctx.Done():
					return
				case result <- val:
				}
			}
		}
	}()

	return result, nil
}
