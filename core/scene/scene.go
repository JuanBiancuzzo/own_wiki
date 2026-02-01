package scene

type SceneInformation struct {
	Resolution Vec3[float64]
	Time       float64
	DTime      float64
	FrameRate  float64
	FrameCount uint64

	// This are the inputs in a ShaderToy shader
	// uniform vec3      iResolution;           // viewport resolution (in pixels)
	// uniform float     iTime;                 // shader playback time (in seconds)
	// uniform float     iTimeDelta;            // render time (in seconds)
	// uniform float     iFrameRate;            // shader frame rate
	// uniform int       iFrame;                // shader playback frame
	// uniform float     iChannelTime[4];       // channel playback time (in seconds)
	// uniform vec3      iChannelResolution[4]; // channel resolution (in pixels)
	// uniform vec4      iMouse;                // mouse pixel coords. xy: current (if MLB down), zw: click
	// uniform samplerXX iChannel0..3;          // input channel. XX = 2D/Cube
	// uniform vec4      iDate;                 // (year, month, day, time in seconds)
	// uniform float     iSampleRate;           // sound sample rate (i.e., 44100)
}

// We should contamplate the idea of having a way to represent the scene
// instead of sharing all the scene itself
type SceneCtx struct {
	SceneInformation

	MainCamera *Camera
	MainLayout *Layout
}

func NewSceneContext(mainCamera *Camera, mainLayout *Layout, sceneInfo SceneInformation) *SceneCtx {
	return &SceneCtx{
		SceneInformation: sceneInfo,
		MainCamera:       mainCamera,
		MainLayout:       mainLayout,
	}
}

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

// In the main layout the x, y, width and height are already set, so the values pass
// for those will be ignore
func (sCtx *SceneCtx) AddMainLayout(config LayoutConfig, creation func()) {}

func (sCtx *SceneCtx) AddLayout(config LayoutConfig, creation func()) {}

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
