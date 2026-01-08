package api

import (
	"fmt"

	"github.com/JuanBiancuzzo/own_wiki/src/core/ecv"
)

type ErrorLoadPath struct {
	HasError    bool
	ErrorReason string
}

func NoErrorLoadPath() ErrorLoadPath {
	return ErrorLoadPath{
		HasError: false,
	}
}

func NewErrorLoadPath(reason string, args ...any) ErrorLoadPath {
	return ErrorLoadPath{
		HasError:    true,
		ErrorReason: fmt.Sprintf(reason, args...),
	}
}

type ReturnRegisterStructure struct {
	HasError    bool
	ErrorReason string
	Ecv         ecv.ECVBuilder
}

func NewErrorRegisterStructure(reason string, args ...any) ReturnRegisterStructure {
	return ReturnRegisterStructure{
		HasError:    true,
		ErrorReason: fmt.Sprintf(reason, args...),
	}
}

func ReturnStructure(system ecv.ECVBuilder) ReturnRegisterStructure {
	return ReturnRegisterStructure{
		HasError: false,
		Ecv:      system,
	}
}
