package pgConf

type Connection struct {
	Username,
	Password,
	Host,
	Port,
	Database string
	MaxAttempts int
	Timeout     int
}

func Get() Connection {
	var config Connection
	config.Username = "postgres"
	config.Password = "postgres"
	config.Host = "localhost"
	config.Port = "5432"
	config.Database = "medo_gas_develop_operator"
	config.MaxAttempts = 5
	config.Timeout = 5
	return config
}
