package config

func LoadAppConfig() App {
	return App{
		Global: loadGlobalConfig(),
		Domain: Domain{},
	}
}
