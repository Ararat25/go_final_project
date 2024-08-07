package config

// Config - структура для хранения переменных окружения
type Config struct {
	Port      int
	DbFile    string
	Password  string
	TokenSalt []byte
}

// NewConfig создание нового объекта Config
func NewConfig(port int, tokenSalt []byte, dbFile, password string) *Config {
	return &Config{
		Port:      port,
		DbFile:    dbFile,
		Password:  password,
		TokenSalt: tokenSalt,
	}
}
