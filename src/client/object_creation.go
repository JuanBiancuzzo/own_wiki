package main

import v "github.com/JuanBiancuzzo/own_wiki/src/core/views"

type ObjectCreatorClient struct {
	defaultCreator v.DefaultObjectCreator

	RequestViewInformation v.FnViewRequest
}

func NewObjectCreatorClient(request v.FnViewRequest) *ObjectCreatorClient {
	return &ObjectCreatorClient{
		defaultCreator: v.DefaultObjectCreator{},

		RequestViewInformation: request,
	}
}

func (obc *ObjectCreatorClient) NewScene(view v.View, worldConfig v.WorldConfiguration) *v.Scene {
	return v.NewCustomScene(view, worldConfig, obc.RequestViewInformation, obc)
}

func (obc *ObjectCreatorClient) NewCamera() *v.Camera {
	return obc.defaultCreator.NewCamera()
}

func (obc *ObjectCreatorClient) NewLayout() *v.Layout {
	return obc.defaultCreator.NewLayout()
}
