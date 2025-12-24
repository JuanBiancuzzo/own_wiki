package shared

import (
	"bytes"
	"encoding/gob"
	"reflect"
)

type componentField struct {
	Name string
	Type string
	Tags string
}

type ComponentInformation struct {
	name   string
	fields []componentField
}

func GetComponentInformation[T any]() *ComponentInformation {
	var t T
	typeInfo := reflect.TypeOf(t)

	fieldAmount := typeInfo.NumField()
	fields := make([]componentField, fieldAmount)

	for i := range fieldAmount {
		field := typeInfo.Field(i)
		fields[i].Name = field.Name
		fields[i].Type = field.Type.Name()
		fields[i].Tags = string(field.Tag)
	}

	return &ComponentInformation{
		name:   typeInfo.Name(),
		fields: fields,
	}
}

func (c *ComponentInformation) GobEncode() ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	value := struct {
		Name   string
		Fields []componentField
	}{Name: c.name, Fields: c.fields}

	if err := enc.Encode(value); err != nil {
		return []byte{}, err
	}

	return buf.Bytes(), nil
}

func (c *ComponentInformation) GobDecode(buf []byte) error {
	value := struct {
		Name   string
		Fields []componentField
	}{}

	dec := gob.NewDecoder(bytes.NewReader(buf))
	if err := dec.Decode(&value); err != nil {
		return err
	}

	c.name = value.Name
	c.fields = value.Fields

	return nil
}
