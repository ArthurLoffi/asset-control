package dto

import (
	"fmt"
	"reflect"
	"strings"
)

// Describe renders a request DTO for [execute] logs, e.g.
// MachineRequestDTO(page=0, page_size=20, employee_id=[4]). Nil filters are
// omitted and pointers are dereferenced. Never pass DTOs that hold secrets.
func Describe(request any) string {
	value := reflect.ValueOf(request)
	var fields []string
	describeFields(value, &fields)
	return fmt.Sprintf("%s(%s)", value.Type().Name(), strings.Join(fields, ", "))
}

func describeFields(value reflect.Value, fields *[]string) {
	for index := range value.NumField() {
		field, fieldValue := value.Type().Field(index), value.Field(index)
		if field.Anonymous && fieldValue.Kind() == reflect.Struct {
			describeFields(fieldValue, fields)
			continue
		}
		name := strings.Split(field.Tag.Get("form"), ",")[0]
		if name == "" || (fieldValue.Kind() == reflect.Pointer || fieldValue.Kind() == reflect.Slice) && fieldValue.IsNil() {
			continue
		}
		*fields = append(*fields, fmt.Sprintf("%s=%v", name, reflect.Indirect(fieldValue)))
	}
}
