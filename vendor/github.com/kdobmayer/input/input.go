package input

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

// New parses a struct using `env` struct tags and loads the appropriate values from
// environment variables. When parsing the struct fields the given validate and options
// constraints are enforced.
func New(i interface{}) error {
	ptr := reflect.ValueOf(i)
	if ptr.Kind() != reflect.Ptr {
		return errors.New("expected a pointer")
	}
	v := ptr.Elem()
	if v.Kind() != reflect.Struct {
		return errors.New("expected a struct")
	}
	t := v.Type()
	var errs multierror
	for i := 0; i < t.NumField(); i++ {
		var err error
		value := os.Getenv(t.Field(i).Tag.Get("env"))

		err = validate(t.Field(i), value)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		if value != "" {
			err = setValue(v.Field(i), value)
			if err != nil {
				errs = append(errs, FieldError{t.Field(i).Name, value, err})
			}
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("invalid input(s):\n%s", errs.Error())
	}
	return nil
}

func validate(field reflect.StructField, value string) error {
	var checks []validator

	if tag, ok := field.Tag.Lookup("validate"); ok {
		validates := strings.Split(tag, ",")
		for _, v := range validates {
			switch v {
			case "required":
				checks = append(checks, required())
			case "path":
				checks = append(checks, pathExists(false))
			case "dir":
				checks = append(checks, pathExists(true))
			default:
				return FieldError{field.Name, value, fmt.Errorf("invalid validate tag: %q", v)}
			}
		}
	}

	if tag, ok := field.Tag.Lookup("opts"); ok {
		opts := strings.Split(tag, ",")
		checks = append(checks, valueOpts(opts...))
	}

	for _, c := range checks {
		if err := c.validate(value); err != nil {
			return FieldError{field.Name, value, err}
		}
	}

	return nil
}

func setValue(field reflect.Value, value string) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Bool:
		field.SetBool(value == "yes" || value == "Yes")
		//		b, err := strconv.ParseBool(value)
		//		if err != nil {
		//			return errors.New("can't convert to int")
		//		}
		//		field.SetBool(b)
	case reflect.Int:
		n, err := strconv.ParseInt(value, 10, 32)
		if err != nil {
			return errors.New("can't convert to bool")
		}
		field.SetInt(n)
	case reflect.Slice:
		field.Set(reflect.ValueOf(strings.Split(value, "|")))
	default:
		return fmt.Errorf("type %q is not supported", field.Kind())
	}
	return nil
}
