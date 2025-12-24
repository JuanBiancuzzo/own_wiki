package shared

import (
	"bytes"
	"encoding/gob"
	"reflect"
)

type EntityInformation struct {
	name       string
	components []string
}

func GetEntityInformation[T any]() *EntityInformation {
	var t T
	typeInfo := reflect.TypeOf(t)

	fieldAmount := typeInfo.NumField()
	components := make([]string, fieldAmount)

	for i := range fieldAmount {
		components[i] = typeInfo.Field(i).Name
	}

	return &EntityInformation{
		name:       typeInfo.Name(),
		components: components,
	}
}

func (e *EntityInformation) GobEncode() ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	value := struct {
		Name       string
		Components []string
	}{Name: e.name, Components: e.components}

	if err := enc.Encode(value); err != nil {
		return []byte{}, err
	}

	return buf.Bytes(), nil
}

func (e *EntityInformation) GobDecode(buf []byte) error {
	value := struct {
		Name       string
		Components []string
	}{}

	dec := gob.NewDecoder(bytes.NewReader(buf))
	if err := dec.Decode(&value); err != nil {
		return err
	}

	e.name = value.Name
	e.components = value.Components

	return nil
}
