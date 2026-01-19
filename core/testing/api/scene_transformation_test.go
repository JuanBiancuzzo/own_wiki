package api_test

import (
	"testing"

	"github.com/go-test/deep"

	pb "github.com/JuanBiancuzzo/own_wiki/core/api/proto"
	s "github.com/JuanBiancuzzo/own_wiki/core/scene"
)

// ---+--- tests ---+---

func Test2DScene(t *testing.T) {
	systemScene := s.New2DScene(
		s.NewLayout(
			s.NewTextObject("Test of text"),
			s.NewImageObject("url:to/image", "Image example"),
			s.NewTextObject("Footer"),
		),
	)

	expectScene := pb.New2DScene(
		pb.NewLayout(
			pb.NewSceneTextObject("Test of text"),
			pb.NewSceneImageObject("url:to/image", "Image example"),
			pb.NewSceneTextObject("Footer"),
		),
	)

	if scene, err := pb.ConvertFromSystemScene(systemScene); err != nil {
		t.Errorf("While converting to SceneDescription, got the error: %v", err)

	} else if diff := deep.Equal(expectScene, scene); diff != nil {
		t.Error(diff)

	} else if systemSceneGen, err := scene.ConvertToSystemScene(); err != nil {
		t.Errorf("While converting to System Scene, got the error: %v", err)

	} else if diff := deep.Equal(systemScene, systemSceneGen); diff != nil {
		t.Error(diff)
	}
}

func FuzzSceneObjects(f *testing.F) {
	systems := []s.SceneObject{
		s.NewTextObject("Test text"),
		s.NewImageObject("url:path/to/image", "Test title"),
	}
	protocols := []*pb.SceneObject{
		pb.NewSceneTextObject("Test text"),
		pb.NewSceneImageObject("url:path/to/image", "Test title"),
	}

	amountTestCases := len(systems)
	if amountTestCases != len(protocols) {
		f.Fatalf("Amount of systems objects (%d) is different from protocol objects (%d)", len(systems), len(protocols))
	}

	for i := range amountTestCases {
		f.Add(i)
	}

	f.Fuzz(func(t *testing.T, i int) {
		system, protocol := systems[i], protocols[i]

		if protocolGen, err := pb.ConvertFromSystemObject(system); err != nil {
			t.Errorf("While converting to protocol, got the error: %v", err)

		} else if systemGen, err := protocolGen.ConvertToSystemObject(); err != nil {
			t.Errorf("While converting to system, got the error: %v", err)

		} else if diff := deep.Equal(system, systemGen); diff != nil {
			t.Error(diff)

		} else if diff := deep.Equal(protocol, protocolGen); diff != nil {
			t.Error(diff)
		}
	})
}
