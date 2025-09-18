package typesense

import (
	"errors"
	"fmt"
	"github.com/typesense/typesense-go/v2/typesense/api"
	"reflect"
	"strings"
)

/*
ToFields takes a struct as input and converts its fields into a slice of Typesense field schema definitions.
This function is useful for automatically generating field schemas for Typesense from your Go structs.
It expects a struct type as input and will return an error if the input is not a struct.
Usage example:

	type MyStruct struct {
		Id    string `json:"id,index"`
		Name  string `json:"name"`
		Age   int    `json:"age,facet"`
		Email string `json:"email,optional"`
		UserId        string `json:"user_id,index,join:user.id"` // creates a reference to the collection user
	}

fields, err := ToFields(MyStruct{})
fields now contains the field schema definitions for MyStruct
Supported tags:
index,name,facet,optional,sort,infix
join:{collectionName}.{id} -> e.g. join:user.id -> This will automatically create a reference to the user schema
*/
func ToFields(Struct any) ([]api.Field, error) {
	val := reflect.ValueOf(Struct)
	if val.Kind() == reflect.Ptr || val.Kind() != reflect.Struct {
		return nil, errors.New("input should be a struct")
	}
	return lexField(val.Type())
}

func lexField(typ reflect.Type) ([]api.Field, error) {
	var collectionFields []api.Field
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if field.PkgPath != "" {
			continue // Skip unexported fields
		}
		fieldType := typeAllowed(field.Type)
		if fieldType == api.OBJECT {
			// Check if Object is embedded.
			if field.Anonymous {
				// Get each Individual Field to be Parsed
				composited, err := lexField(field.Type)
				if err != nil {
					return nil, err
				}
				collectionFields = append(collectionFields, composited...)
				continue
			}
		}
		tags := field.Tag.Get("json") // tags save the field_name and Options like facet, index, join, optional etc.
		apiField, err := parseField(fieldType, tags)
		if err != nil {
			return nil, err
		}
		collectionFields = append(collectionFields, apiField)
	}
	return collectionFields, nil
}

func parseField(T api.Type, tag string) (api.Field, error) {
	params := strings.Split(tag, ",")
	var field api.Field
	var True bool = true

	// We need the json Field.Tag
	if len(params) == 0 {
		return api.Field{}, errors.New("field name has to be provided for matching")
	}

	field.Name = params[0]
	field.Type = string(T)

	for _, key := range params[1:] {
		switch key {
		case "optional": // optional fields, can be null
			field.Optional = &True
		case "facet": // If a field is facet its also automatically indexed, correct?
			field.Facet = &True
			field.Index = &True
		case "index":
			field.Index = &True
		case "sort":
			field.Sort = &True
		case "infix":
			field.Infix = &True
		default:
			if ref, ok := strings.CutPrefix(key, "join:"); ok {
				field.Reference = Pointer(ref)
				continue
			}
		}
	}

	return field, nil
}

func typeAllowed(t reflect.Type) api.Type {
	switch t.Kind() {
	case reflect.String:
		return api.STRING
	case reflect.Int32, reflect.Int:
		return api.INT32
	case reflect.Int64:
		return api.INT64
	case reflect.Float32, reflect.Float64:
		return api.FLOAT
	case reflect.Bool:
		return api.BOOL
	case reflect.Slice:
		elemType := typeAllowed(t.Elem())
		if elemType != "" {
			return elemType + "[]"
		}
	case reflect.Struct:
		return api.OBJECT
	case reflect.Pointer:
		return typeAllowed(t.Elem())
	default:
		panic("type not allowed")
	}
	fmt.Println(t.Kind())
	return ""
}

// Pointer returns the Pointer of a Type v
func Pointer[T any](v T) *T {
	return &v
}
