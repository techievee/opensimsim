package config

type Config struct {
	Server Server `json:"server"`
	MySql  MySql  `json:"mysql"`
}

type Server struct {
	Listen   string `json:"listen"`
	TimeOut  string `json:"timeout"`
	LogLevel string `json:"loglevel"`
	Limit    string `json:"limit"`
}

type MySql struct {
	DbDriver string `json:"dbDriver"`
	DbUser   string `json:"dbUser"`
	DbPass   string `json:"dbPass"`
	DbName   string `json:"dbName"`
	DbHost   string `json:"dbHost"`
	DbPort   string `json:"dbPort"`
}
