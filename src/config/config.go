package config

import (
	"gopkg.in/ini.v1"
	"log"
  )

type Config struct {
	Es_url                 string `ini:"es_url"`
	Es_index               string `ini:"es_index"`
	Es_type                string `ini:"es_type"`
	Log_level              string `ini:"log_level"`
	Request_interval_mills int32  `ini:"request_interval_mills"`
	Request_parallel_size  int16  `ini:"request_parallel_size"`
}

 var config *Config

func GetConfig() *Config {
	if config != nil {
		return config
	}
	config:=new(Config)
	err:=ini.MapTo(config,"../../app.cfg")
	if err!=nil{
		log.Fatalln(err)
	}
	return config
}
