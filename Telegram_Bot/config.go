package main

// .ENV file setting
type Config struct {
	Token       string `env:"TOKEN"`
	Myinfo      string `env:"MYINFO" default:"no personal info"`
	MuGitHub    string `env:"MYGITHUB" default:"https://github.com/"`
	Myemail     string `env:"MYMAIL" default:"nomail@server.com"`
	Botapi      string `env:"BOTAPI" default:"https://api.telegram.org/bot"`
	LogrusLevel string `env:"LOGRUSLEVEL" default:"info"`
}
