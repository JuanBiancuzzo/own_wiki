package scene

type LayoutDirection uint

const (
	VERTICAL_DIR = iota
	HORIZONTAL_DIR
)

type LayoutConfig struct {
	X, Y, W, H DimensionValue
	Dir        LayoutDirection
	Padding    DimConfig
	Margin     DimConfig
}

type Layout struct {
	Config  LayoutConfig
	Objects []Object
}

func NewLayout(config LayoutConfig) *Layout {
	return &Layout{
		Config: config,
	}
}

// This is to make a layout an object
func (*Layout) isObject() {}

func (l *Layout) Update(config LayoutConfig) {
	l.Config = config
}

func (l *Layout) Init() {}

func (l *Layout) AddObject(object Object) {
	l.Objects = append(l.Objects, object)
}

func (l *Layout) End() {}

func (l *Layout) GenerateCameraDescription() []*CameraDescription {
	// The layout its always an orthogonal camera, not given by the camera that its attach?
	// The z component is always discarted because its 2D
	return []*CameraDescription{}
}
