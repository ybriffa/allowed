package allowed

import (
	"fmt"
	"reflect"
	"strings"
)

const (
	// AllowedTag is the structure tag to set to validate the structure
	AllowedTag = "allowed"
)

type explorer struct {
	value string
	ptrs  map[reflect.Type]map[uintptr]reflect.Value
}

// Check defines whether the given value is allowed in the structure fields.
func Check(value string, i interface{}) error {
	e := explorer{
		ptrs:  map[reflect.Type]map[uintptr]reflect.Value{},
		value: strings.ToLower(value),
	}
	return e.explore(reflect.ValueOf(i))
}

func explorable(t reflect.Type) bool {
	switch t.Kind() {
	case reflect.Struct:
		return true
	case reflect.Ptr, reflect.Array, reflect.Slice:
		return explorable(t.Elem())
	case reflect.Map:
		// XXX: we are not handling map's keys for now
		// Can be accessed with t.Key()
		return explorable(t.Elem())
	}
	return false
}

func (e *explorer) explore(v reflect.Value) error {
	if !explorable(v.Type()) {
		return fmt.Errorf("data of type %s cannot be checked", v.Type())
	}

	return e.validate(v)
}

func (e *explorer) allowed(content string) bool {
	for _, s := range strings.Split(content, ",") {
		if strings.ToLower(s) == e.value {
			return true
		}
	}
	return false
}

func (e *explorer) validate(v reflect.Value) error {
	switch v.Kind() {
	case reflect.Invalid:
		return nil
	case reflect.Array, reflect.Slice:
		if v.IsNil() {
			return nil
		}
		for i := 0; i < v.Len(); i++ {
			if err := e.validate(v.Index(i)); err != nil {
				return fmt.Errorf("item %d: %s", i, err)
			}
		}
		return nil
	case reflect.Map:
		if v.IsNil() {
			return nil
		}

		for _, key := range v.MapKeys() {
			// XXX: not checking map's keys for now
			if err := e.validate(v.MapIndex(key)); err != nil {
				return fmt.Errorf("key %s: %s", key, err)
			}
		}
		return nil
	}

	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		if v.IsNil() {
			return nil
		}
		v = v.Elem()
	}

	return e.validateStruct(v)
}

func (e *explorer) validateStruct(v reflect.Value) error {
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		fieldInfos := t.Field(i)
		fieldValue := v.Field(i)

		allowed := e.allowed(fieldInfos.Tag.Get(AllowedTag))
		if !allowed && isSet(fieldValue) {
			return fmt.Errorf("field %q: not allowed to set", fieldInfos.Name)
		}

		if allowed && explorable(fieldValue.Type()) {
			err := e.validate(fieldValue)
			if err != nil {
				return fmt.Errorf("field %q: %s", fieldInfos.Name, err)
			}
		}
	}
	return nil
}

func isSet(v reflect.Value) bool {
	if v.Kind() == reflect.Invalid {
		return false
	}
	return !v.IsZero()
}
