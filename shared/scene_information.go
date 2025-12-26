package shared

import (
	v "github.com/JuanBiancuzzo/own_wiki/view"
)

type SceneInformation struct {
	ViewName   string
	EntityName string
	// Entity     any
	Scene *v.Scene
}

/*
func (si *SceneInformation) GobEncode() ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	value := struct {
		Name   string
		Entity string
	}{Name: si.name, Entity: si.entity}

	if err := enc.Encode(value); err != nil {
		return []byte{}, err
	}

	return buf.Bytes(), nil
}

func (si *SceneInformation) GobDecode(buf []byte) error {
	value := struct {
		Name   string
		Entity string
	}{}

	dec := gob.NewDecoder(bytes.NewReader(buf))
	if err := dec.Decode(&value); err != nil {
		return err
	}

	si.name = value.Name
	si.entity = value.Entity

	return nil
}
*/
