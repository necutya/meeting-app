package config

type Config struct {
	LogLevel   string `mapstructure:"LOG_LEVEL" default:"DEBUG"`
	HTTPConfig HTTP
	RedisCfg   Redis `mapstructure:"REDIS"`
}

type HTTP struct {
	Port            int      `mapstructure:"PORT"  default:"8080"`
	URLPrefix       string   `mapstructure:"URL_PREFIX"  default:"/api"`
	CORSAllowedHost []string `mapstructure:"CORS_ALLOWED_HOST"  default:"*"`
}

// Redis defines configs for Redis.
type Redis struct {
	Address  string `mapstructure:"ADDRESS"  default:"localhost:6379"`
	PoolSize int    `mapstructure:"POOL_SIZE"  default:"10"`
	Password string `mapstructure:"PASSWORD"  default:""`
}
