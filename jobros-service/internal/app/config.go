package app

// LoggingConfig holds logging-related configuration
type LoggingConfig struct {
	Level int `yaml:"level" envconfig:"LOG_LEVEL" default:"0"`
}

// JWTConfig holds JWT-related configuration
type JWTConfig struct {
	SecretKey string `yaml:"secretKey" envconfig:"JWT_SECRET_KEY" required:"true"`
}

// MongoConfig holds MongoDB-related configuration
type MongoConfig struct {
	URI      string `yaml:"uri" envconfig:"MONGO_URI" required:"true"`
	Database string `yaml:"database" envconfig:"MONGO_DATABASE" required:"true"`
}

// AppConfig represents the configuration for the application
type AppConfig struct {
	Host    string        `yaml:"host" envconfig:"HOST" default:"localhost"`
	Port    int           `yaml:"port" envconfig:"PORT" default:"8080"`
	Mongo   MongoConfig   `yaml:"mongo"`
	Logging LoggingConfig `yaml:"logging"`
	JWT     JWTConfig     `yaml:"jwt"`
}
