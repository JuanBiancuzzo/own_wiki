package views

import (
	e "github.com/JuanBiancuzzo/own_wiki/src/core/events"
	u "github.com/JuanBiancuzzo/own_wiki/src/core/plugin/user"
	v "github.com/JuanBiancuzzo/own_wiki/src/core/views"
)

type PluginView struct {
	Plugin u.UserPlugin
}

func NewPluginView(plugin u.UserPlugin) *PluginView {
	return &PluginView{
		Plugin: plugin,
	}
}

func (pv *PluginView) View(world *v.World, creator v.ObjectCreator, yield v.FnYield) v.View {
	if err := pv.Plugin.InitializeViewManeger(world.GetConfiguration()); err != nil {
		// Show error and suggestion
		return nil
	}

	for events := range yield() {
		// Estaria mejor mejorar esta interfaz, es medio confiar que este todo bien
		if createView, ok, rest := getCreationViewEvent(events); ok {
			err := pv.Plugin.InitializeView(createView.ViewName, createView.EntityData)
			if err != nil {
				// Show error and what to do
				return nil
			}

		} else if err := pv.Plugin.SendEvents(rest); err != nil {
			// Show error and what to do
			return nil
		}
	}
	return nil
}

func getCreationViewEvent(events []e.Event) (e.CreateViewEvent, bool, []e.Event) {
	for i, event := range events {
		if createViewEvent, ok := event.(e.CreateViewEvent); ok {
			var rest []e.Event
			if i+1 < len(events) {
				rest = events[i+1:]
			}

			return createViewEvent, true, append(events[:i], rest...)
		}
	}

	return e.CreateViewEvent{}, false, events
}
