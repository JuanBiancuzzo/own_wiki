package scene

import d "github.com/JuanBiancuzzo/own_wiki/core/scene/draw_commands"

type SceneDescription struct {
	Cameras []*CameraDescription
}

func NewSceneDescription(cameras []*CameraDescription) *SceneDescription {
	return &SceneDescription{
		Cameras: cameras,
	}
}

type CameraDescription struct {
	PerspectiveMatrix Mat4[float64]
	DrawCommands      []d.DrawCommand
}

func NewCameraDescription(perspectiveMatrix Mat4[float64], drawCommands []d.DrawCommand) *CameraDescription {
	return &CameraDescription{
		PerspectiveMatrix: perspectiveMatrix,
		DrawCommands:      drawCommands,
	}
}
