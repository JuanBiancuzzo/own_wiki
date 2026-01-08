package ecv

import (
	e "github.com/JuanBiancuzzo/own_wiki/src/core/events"
	q "github.com/JuanBiancuzzo/own_wiki/src/core/query"
)

type ECV struct{}

func (ecv *ECV) Query(q.QueryRequest) (any, error) {
	return nil, nil
}

func (ecv *ECV) SendEvent(e.Event) error {
	return nil
}
