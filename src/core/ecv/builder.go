package ecv

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

	return nil
}

func (ecv *ECVBuilder) BuildECV() (*ECV, error) {
	return &ECV{}, nil
}
