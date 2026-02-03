package scene

import d "github.com/JuanBiancuzzo/own_wiki/core/scene/draw_commands"

type FrameInformation struct {
	// This are inspired from inputs in a ShaderToy shader
	Resolution Vec3[float64] // X and Y are the amount of pixels and Z is the aspect ratio
	Time       float64       // The time from the first call to the view
	DTime      float64       // The delta time between frames
	FrameRate  float64       // The framerate (it would be a fix framerate for now)
	FrameCount uint64        // The amount of frames from the first call to the view
}

// The MainCamera and MainLayout generate a tree like structure where each alter the elements
// that are containt within
type SceneCtx struct {
	FrameInformation

	// Is will contain elements that are going to be alter by the camara perspective matrix
	MainCamera *Camera

	selectedLayout *Layout
}

func NewSceneContext(mainCamera *Camera, sceneInfo FrameInformation) *SceneCtx {
	return &SceneCtx{
		FrameInformation: sceneInfo,
		MainCamera:       mainCamera,

		selectedLayout: nil,
	}
}

func (sCtx *SceneCtx) GenerateDrawCommands() []d.DrawCommand {
	return []d.DrawCommand{}
}

// In the main layout the x, y, width and height are already set, so the values pass
// for those will be ignore
func (sCtx *SceneCtx) AddMainLayout(config LayoutConfig, creation func()) {
	// Overwriting x, y, width and height
	config.X = ValueFix(0)
	config.Y = ValueFix(0)
	config.W = ValueFix(int(sCtx.FrameInformation.Resolution.X))
	config.H = ValueFix(int(sCtx.FrameInformation.Resolution.Y))
	sCtx.MainLayout.Update(config)

	sCtx.selectedLayout = sCtx.MainLayout
	defer func() { sCtx.selectedLayout = nil }()

	sCtx.selectedLayout.Init()
	defer sCtx.selectedLayout.End()

	creation()
}

func (sCtx *SceneCtx) AddLayout(config LayoutConfig, creation func()) {
	if sCtx.selectedLayout == nil {
		// Now we dont support other layouts than the main, as the first
		sCtx.AddMainLayout(config, creation)
		return
	}

	previousLayout := sCtx.selectedLayout
	sCtx.selectedLayout = NewLayout(config)
	defer func() { sCtx.selectedLayout = previousLayout }()

	sCtx.selectedLayout.Init()
	defer sCtx.selectedLayout.End()

	creation()
}

type TextBehaviour uint

const (
	WRAP_TEXT = iota
)

type TitleConfig struct {
	Level     uint // if 0 si the value asign then the title will not be render
	Text      string
	Behaviour TextBehaviour
}

func (sCtx *SceneCtx) AddTitle(config TitleConfig) {}

type DimensionType uint

const (
	DT_NONE = iota
	SHRINK
	SPAND
	FIX        // Pixels (int)
	PORCENTAGE // % (float)
	RELATIVE   // em (float)
)

type DimensionValue struct {
	// The default value is not set, to know when the user
	// didnt add a value, but we would use the Shrink as default
	DimensionType
	// If the type is Shrink o Span, then this is nil
	// If is Fix then is an int
	// If is porcentage o relative then is a float
	Value any
}

func ValueEm(ems float32) DimensionValue {
	return DimensionValue{
		DimensionType: RELATIVE,
		Value:         ems,
	}
}

func ValuePorce(porcentage float32) DimensionValue {
	return DimensionValue{
		DimensionType: PORCENTAGE,
		Value:         porcentage,
	}
}

func ValueFix(pixels int) DimensionValue {
	return DimensionValue{
		DimensionType: FIX,
		Value:         pixels,
	}
}

func Shrink() DimensionValue {
	return DimensionValue{
		DimensionType: SHRINK,
	}
}

func Spand() DimensionValue {
	return DimensionValue{
		DimensionType: SPAND,
	}
}

type DimConfig struct {
	Above, Below, Left, Right DimensionValue
	LeftRight, AboveBelow     DimensionValue
	All                       DimensionValue
}

type HlineConfig struct {
	Padding DimConfig
	Margin  DimConfig
}

func (sCtx *SceneCtx) AddHline(config HlineConfig) {}

type BoxConfig struct {
	X, Y, W, H DimensionValue
	Dir        LayoutDirection
	Padding    DimConfig
	Margin     DimConfig
}

func (sCtx *SceneCtx) AddBox(config BoxConfig, creation func()) {}

type TextConfig struct {
	Text      string
	Behaviour TextBehaviour
}

func (sCtx *SceneCtx) AddText(config TextConfig) {}

type ButtonConfig struct {
	X, Y, W, H     DimensionValue
	Dir            LayoutDirection
	Padding        DimConfig
	Margin         DimConfig
	RoundedCorners DimConfig
}

func (sCtx *SceneCtx) AddButton(config ButtonConfig, creation func()) bool {
	return false
}
