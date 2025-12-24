package ecv

import (
	"iter"

	e "github.com/JuanBiancuzzo/own_wiki/core/events"
	c "github.com/JuanBiancuzzo/own_wiki/core/system/configuration"
)

type ECV struct {
	EventQueue chan e.Event
	Scene      *Scene

	currentView func() (View, bool)
	stopView    func()
}

func NewECV(config c.UserConfig) *ECV {
	return &ECV{
		EventQueue: make(chan e.Event),
		Scene:      NewScene(config.TargetFrameRate),

		currentView: nil,
		stopView:    nil,
	}
}

func (ecv *ECV) RegisterComponent(component any) {

}

func (ecv *ECV) AssignCurrentView(view View) {
	nextViewChannel := make(chan View, 1)

	iterator := func(yield func(uint8) bool) {
		nextViewChannel <- view.View(ecv.Scene, func() bool { return yield(0) })
	}

	next, stop := iter.Pull(iterator)

	ecv.currentView = func() (View, bool) {
		if _, valid := next(); valid {
			return nil, false
		}
		return <-nextViewChannel, true
	}
	ecv.stopView = stop
}

func (ecv *ECV) GenerateFrame() (SceneRepresentation, bool) {
	if ecv.currentView == nil {
		return nil, false
	}

	if nextView, stopped := ecv.currentView(); stopped && nextView != nil {
		ecv.AssignCurrentView(nextView)

	} else if stopped {
		return nil, false
	}

	return ecv.Scene.GetRepresentation(), true
}

func (ecv *ECV) Close() {
	if ecv.stopView != nil {
		ecv.stopView()
	}
}
