package input

import (
	"errors"
	"fmt"
	"os"
)

type validator interface {
	validate(string) error
}

type validatorFunc func(string) error

func (f validatorFunc) validate(i string) error {
	return f(i)
}

func valueOpts(opts ...string) validator {
	return validatorFunc(func(v string) error {
		for _, opt := range opts {
			if v == opt {
				return nil
			}
		}
		return fmt.Errorf("not in value options %s", opts)
	})
}

func required() validator {
	return validatorFunc(func(v string) error {
		if v == "" {
			return errors.New("required but not specified")
		}
		return nil
	})
}

func pathExists(dir bool) validator {
	return validatorFunc(func(v string) error {
		if v == "" {
			return nil
		}
		file, err := os.Stat(v)
		if os.IsNotExist(err) {
			return errors.New("file does not exist")
		} else if err != nil {
			return fmt.Errorf("can't get file info: %v", err)
		}
		if dir && !file.IsDir() {
			return errors.New("not a directory")
		}
		return nil
	})
}
