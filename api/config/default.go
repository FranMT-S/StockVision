package config

// DefaultConstants defines the default values for pagination and ordering
type DefaultConstantsConfig struct {
	MaxPageSize     int
	DefaultPageSize int
	Page            int
}

var defaultConstantsInstance *DefaultConstantsConfig

func DefaultConstants() *DefaultConstantsConfig {
	if defaultConstantsInstance == nil {
		defaultConstantsInstance = &DefaultConstantsConfig{
			MaxPageSize:     100,
			DefaultPageSize: 10,
			Page:            1,
		}
	}

	return defaultConstantsInstance
}
