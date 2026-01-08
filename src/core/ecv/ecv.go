package ecv

import (
	e "github.com/JuanBiancuzzo/own_wiki/src/core/events"
	q "github.com/JuanBiancuzzo/own_wiki/src/core/query"
)

type ECV struct{}

func (ecv *ECV) Query(q.QueryRequest) (EntityDescription, error) {
	return EntityDescription{}, nil
}

func (ecv *ECV) SendEvent(e.Event) error {
	return nil
}
