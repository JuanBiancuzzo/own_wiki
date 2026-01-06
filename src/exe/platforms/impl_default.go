//go:build !terminal
// +build !terminal

package platforms

import p "github.com/JuanBiancuzzo/own_wiki/src/core/platform"

func GetPlatformImplementation() p.Platform {
	panic("There is no implementation for platform")
}
