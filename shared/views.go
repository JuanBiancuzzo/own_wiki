package shared

import (
	"bytes"
	"encoding/gob"
	"reflect"

	v "github.com/JuanBiancuzzo/own_wiki/view"
)

type ViewInformation struct {
	name   string
	entity string
}

func GetViewInformation[V v.View, T any]() *ViewInformation {
	var view V
	viewInfo := reflect.TypeOf(view)

	var t T
	typeInfo := reflect.TypeOf(t)

	return &ViewInformation{
		name:   viewInfo.Name(),
		entity: typeInfo.Name(),
	}
}

func (vi *ViewInformation) GobEncode() ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	value := struct {
		Name   string
		Entity string
	}{Name: vi.name, Entity: vi.entity}

	if err := enc.Encode(value); err != nil {
		return []byte{}, err
	}

	return buf.Bytes(), nil
}

func (vi *ViewInformation) GobDecode(buf []byte) error {
	value := struct {
		Name   string
		Entity string
	}{}

	dec := gob.NewDecoder(bytes.NewReader(buf))
	if err := dec.Decode(&value); err != nil {
		return err
	}

	vi.name = value.Name
	vi.entity = value.Entity

	return nil
}
