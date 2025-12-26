package shared

import (
	"reflect"

	v "github.com/JuanBiancuzzo/own_wiki/view"
)

type SceneInformation struct {
	ViewName   string
	EntityName string
	// Entity     any
	Scene v.Scene
}

// Cambio de escena, Modificacion o eliminacion de componente
type SceneOperation struct {
	// Cambio de escena
	viewName string
	entityID uint64
}

func ChangeScene[V v.View](entity any) *SceneOperation {
	var view V
	viewInfo := reflect.TypeOf(view)

	return &SceneOperation{
		viewName: viewInfo.Name(),
		entityID: 0, // Hacer que sea el hash de los elementos importantes
	}
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
