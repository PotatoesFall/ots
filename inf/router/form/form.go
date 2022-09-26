package form

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"
)

func Read(form url.Values, v any) error {
	rv := reflect.ValueOf(v).Elem()
	nf := rv.NumField()
	for i := 0; i < nf; i++ {
		tag := rv.Type().Field(i).Tag.Get(`form`)
		if tag == `` {
			continue
		}

		formVal, ok := form[tag]
		if !ok || len(formVal) == 0 {
			continue
		}

		err := assign(rv.Field(i).Addr().Interface(), formVal[0])
		if err != nil {
			return err
		}
	}

	return nil
}

func assign(to any, value string) error {
	switch dst := to.(type) {
	case *string:
		*dst = value

	case *int:
		i, err := strconv.Atoi(value)
		if err != nil {
			return err
		}

		*dst = i

	case *bool:
		b, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}

		*dst = b

	default:
		panic(fmt.Sprintf(`forms package: type not supported yet: %T`, to))
	}

	return nil
}
