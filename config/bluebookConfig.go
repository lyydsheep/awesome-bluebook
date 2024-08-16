package config

type blueBookConfig struct {
	DB
}

type DB struct {
	DSN string
}
