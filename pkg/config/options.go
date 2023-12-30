package config

type (
	Config struct {
		EnvFile    optionsFromEnv
		Mongo      mongoOption    `envPrefix:"MONGO_"`
		Postgres   postgresOption `envPrefix:"POSTRGESQL_"`
		Redis      redisOption    `envPrefix:"REDIS_"`
		MQTTClient mqttOption     `envPrefix:"MQTT_"`
		Temporal   temporalOption `envPrefix:"TEMPORAL_"`
		Telegram   telegramOption `envPrefix:"TELEGRAM_"`
		S3         s3Option       `envPrefix:"S3_"`
	}

	Option interface {
		apply(*Config)
	}

	optionsFromEnv struct {
		Path string
	}

	mongoOption struct {
		URI          string `env:"URI"`
		DatabaseName string `env:"DATABASE_NAME"`
	}

	postgresOption struct {
		URI          string `env:"URI"`
		DatabaseName string `env:"DATABASE_NAME"`
	}

	redisOption struct {
		URI      string `env:"URI"`
		Password string `env:"PASSWORD"`
		Database int    `env:"DATABASE"`
	}

	mqttOption struct {
		URI      string `env:"URI"`
		Password string `env:"PASSWORD"`
		Username string `env:"USERNAME"`
		ClientId string `env:"CLIENT_ID"`
	}

	temporalOption struct {
		URI        string   `env:"URI"`
		Namespaces []string `env:"NAMESPACES" envSeparator:":"`
	}

	telegramOption struct {
		APIToken string `env:"API_TOKEN"`
	}

	s3Option struct {
		Address   string `env:"ADDRESS"`
		AccessKey string `env:"ACCESS_KEY"`
		SecretKey string `env:"SECRET_KEY"`
		Bucket    string `env:"BUCKET"`
		Region    string `env:"REGION"`
		Secure    bool   `env:"SECURE"`
	}
)

func (o optionsFromEnv) apply(opts *Config) {
	opts.EnvFile = o
}

func (o mongoOption) apply(opts *Config) {
	opts.Mongo = o
}

func (o postgresOption) apply(opts *Config) {
	opts.Postgres = o
}

func (o redisOption) apply(opts *Config) {
	opts.Redis = o
}

func (o mqttOption) apply(opts *Config) {
	opts.MQTTClient = o
}

func (o temporalOption) apply(opts *Config) {
	opts.Temporal = o
}

func (o telegramOption) apply(opts *Config) {
	opts.Telegram = o
}

func (o s3Option) apply(opts *Config) {
	opts.S3 = o
}

func WithEnvFile(envfile string) Option {
	return optionsFromEnv{
		Path: envfile,
	}
}

func WithMongo(uriData string, databaseName string) Option {
	return mongoOption{
		URI:          uriData,
		DatabaseName: databaseName,
	}
}

func WithPostgres(uriData string, databaseName string) Option {
	return postgresOption{
		URI:          uriData,
		DatabaseName: databaseName,
	}
}

func WithS3(Address, AccessKey, SecretKey, Bucket, Region string, Secure bool) Option {
	return s3Option{
		Address:   Address,
		AccessKey: AccessKey,
		SecretKey: SecretKey,
		Bucket:    Bucket,
		Region:    Region,
		Secure:    Secure,
	}
}

func WithRedis(uriData string, password string, database int) Option {
	return redisOption{
		URI:      uriData,
		Password: password,
		Database: database,
	}
}

func WithMQTT(uriData string, password string, username string, clientID string) Option {
	return mqttOption{
		URI:      uriData,
		Password: password,
		Username: username,
		ClientId: clientID,
	}
}

func WithTemporal(uriData string, namespaces ...string) Option {
	return temporalOption{
		URI:        uriData,
		Namespaces: namespaces,
	}
}

func WithTelegram(apikey string) Option {
	return telegramOption{
		APIToken: apikey,
	}
}
