package views

import (
	s "github.com/JuanBiancuzzo/own_wiki/core/scene"
)

// This interfaces allows the user to create layouts and animations in a
// Immediate Mode way.
type View interface {
	// The sCtx *s.SceneCtx is the way to add elements to the scene, and get useful
	// informations regarding the resolution, the time, the platform, etc.
	// The return value of the view should
	//  * If the view should be render the next frame, then it should return itself
	//  * If we need to render another view, then it should return an instances
	// 		of the other view with the corresponding data
	//  * else there is no next view, then it should return nil
	View(sCtx *s.SceneCtx) View
}
