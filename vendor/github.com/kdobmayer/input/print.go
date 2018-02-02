package input

import (
	"fmt"
	"reflect"
	"strings"
)

// Secret variables are not shown in the printed output
type Secret string

// String implements fmt.Stringer.String.
// When a Secret is printed, it's masking the underlying string with asterisks.
func (s Secret) String() string {
	return strings.Repeat("*", 5)
}

// Print the name of the struct in blue color followed by a newline,
// then print all fields and their respective values in separate lines.
func Print(config interface{}) {
	v := reflect.ValueOf(config)
	t := reflect.TypeOf(config)

	blue, reset, format := "\x1b[34;1m", "\x1b[0m", "%s%s:%s\n"
	fmt.Printf(format, blue, t.Name(), reset)
	for i := 0; i < t.NumField(); i++ {
		fmt.Printf("- %s: %v\n", t.Field(i).Name, v.Field(i).Interface())
	}
}

type multierror []error

// Error implements built-in error.Error. Every error goes into new line,
// preceded by a hyphen and a space. It gives an organized view for the collected errors.
func (m multierror) Error() string {
	s := ""
	for _, err := range m {
		s += fmt.Sprintf("- %s\n", err.Error())
	}
	return s
}

// FieldError contains the invalid input field, the value of the field and the
// original error.
type FieldError struct {
	Field string
	Value string
	Err   error
}

// Error implements builtin errors.Error
func (e FieldError) Error() string {
	segments := []string{e.Field}
	if e.Value != "" {
		segments = append(segments, e.Value)
	}
	segments = append(segments, e.Err.Error())
	return strings.Join(segments, ": ")
	//return fmt.Sprintf("%s: %s%v", e.Field, value, e.Err.Error())
}
