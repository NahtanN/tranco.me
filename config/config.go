package config

type Config struct {
	DatabaseDir string `mapstructure:"database_dir"`
	SharedDB    string `mapstructure:"shared_db"`
	UserDBDir   string `mapstructure:"user_db_dir"`
}
