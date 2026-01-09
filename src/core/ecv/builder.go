package ecv

import "fmt"

type ComponentDescription struct{}

type EntityDescription struct{}

type ViewDescription struct{}

type ECVBuilder struct {
}

func NewECVBuilder() *ECVBuilder {
	return &ECVBuilder{}
}

func (ecv *ECVBuilder) RegisterComponent(component any) error {
	return nil
}

func (ecv *ECVBuilder) RegisterEntity(entity any) error {
	return nil
}

func (ecv *ECVBuilder) RegisterView(view any, main bool) error {
	/* if _, ok := view.(v.View); !ok {
		return fmt.Errorf("The register 'view' (%v) does not implement the view interface", view)
	} */

	// check if the view has a View method, and then check the type of the first parameter
	// and thats how we know the entity for that view

	return nil
}

func (ecv ECVBuilder) Verify() bool {
	return false
}

func (ecv ECVBuilder) BuildECV() (*ECV, error) {
	if !ecv.Verify() {
		return nil, fmt.Errorf("Failed to build ecv")
	}

	return &ECV{}, nil
}
