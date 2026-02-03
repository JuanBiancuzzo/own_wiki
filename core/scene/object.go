package scene

// Define this better, because this isnt useful as an interface, and the
// logic to generate the camera description isn't the best
type Object interface {
	GenerateCameraDescription(givenCamera *Camera) []*CameraDescription
}
