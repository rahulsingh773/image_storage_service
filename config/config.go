package config

import "os"

var Config map[string]string

func init() {
	Config = make(map[string]string)

	Config["bind_port"] = os.Getenv("BIND_PORT")     //bind_port to listen
	Config["kafka_host"] = os.Getenv("KAFKA_HOST")   //kafka host to connect
	Config["kafka_topic"] = os.Getenv("KAFKA_TOPIC") //kafka topic to publish event
	Config["host"] = os.Getenv("HOST")               //app host/ip needed to give static url for images

	if Config["bind_port"] == "" {
		Config["bind_port"] = "3000"
	}

	if Config["kafka_host"] == "" {
		Config["kafka_host"] = "localhost"
	}

	if Config["kafka_topic"] == "" {
		Config["kafka_topic"] = "image_server"
	}

	if Config["host"] == "" {
		Config["host"] = "10.198.107.251:3000"
	}
}

func GetConfigParamString(param string) string {
	return Config[param]
}

func SetConfigParam(param, value string) {
	Config[param] = value
}
