package model

type Config struct {
	Server struct {
		Port int    `mapstructure:"port"`
		Addr string `mapstructure:"addr"`
	} `mapstructure:"server"`

	Web struct {
		Open    bool   `mapstructure:"open"`
		Title   string `mapstructure:"title"`
		Desc    string `mapstructure:"desc"`
		NavList []struct {
			Title string `mapstructure:"title"`
			Path  string `mapstructure:"path"`
		} `mapstructure:"nav_list"`
	} `mapstructure:"web"`

	Database struct {
		Host     string `mapstructure:"host"`
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		Name     string `mapstructure:"name"`
	} `mapstructure:"database"`

	Auth struct {
		Id        int64  `mapstructure:"id"`
		User      string `mapstructure:"user"`
		Password  string `mapstructure:"password"`
		SignKey   string `mapstructure:"sign_key"`
		LocalPath string `mapstructure:"local_path"`
		Issuer    string `mapstructure:"issuer"`
	} `mapstructure:"auth"`

	Log struct {
		Level string `mapstructure:"level"`
	} `mapstructure:"log"`
}
