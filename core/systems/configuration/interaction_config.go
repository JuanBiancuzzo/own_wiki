package configuration

type UserInteractionConfig struct {
	Protocol string `json:"protocol"`
	Ip       string `json:"ip"`
	Port     uint8  `json:"port"`
}

type SystemInteractionConfig struct {
	Protocol string `json:"protocol"`
	Ip       string `json:"ip"`
	Port     uint8  `json:"port"`
}
