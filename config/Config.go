package config

var Config = struct {
	App struct {
		ENV      string
		HttpAddr string
		HttpPort string
	}

	DB struct {
		Driver   string
		Host     string
		Port     string `default:"3306"`
		Name     string
		User     string `default:"admin_zahra_user"`
		Password string `default:"zahra_password1945"`
	}

	JWT struct {
		Secret string
	}
}{}
