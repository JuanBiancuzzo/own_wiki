package configuration

var UserConfig UserConfiguration

type UserConfiguration struct {
	TargetFrameRate int
}

func LoadUserConfiguration(userConfigurationPath string) error {

	return nil
}
