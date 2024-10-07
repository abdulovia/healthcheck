package server

import "flag"

const (
	defaultFirstEndpoint  = "127.0.0.1:8788"
	defaultSecondEndpoint = "127.0.0.1:8768"
)

type Config struct {
	FirstEndpoint  string
	SecondEndpoint string
}

func NewConfig() func() *Config {
	var (
		firEndpoint = flag.String("fir-endpoint", defaultFirstEndpoint, "Endpoint first")
		secEndpoint = flag.String("sec-endpoint", defaultSecondEndpoint, "Endpoint second")
	)

	return func() *Config {
		return &Config{
			FirstEndpoint:  *firEndpoint,
			SecondEndpoint: *secEndpoint,
		}
	}
}
