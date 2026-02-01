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
	Config LayoutConfig
}

func NewLayout(config LayoutConfig) *Layout {
	return &Layout{
		Config: config,
	}
}

func (l *Layout) Update(config LayoutConfig) {
	l.Config = config
}

func (l *Layout) Init() {}

func (l *Layout) End() {}
