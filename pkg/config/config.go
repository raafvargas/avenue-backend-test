package config

import (
	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

type Configuration struct {
	S3Endpoint     string `env:"S3_ENDPOINT"`
	BucketName     string `env:"BUCKET_NAME"`
	StorageBackend string `env:"STORAGE_BACKEND"`
	GCPProjectID   string `env:"GCP_PROJECT_ID"`
	Host           string `env:"HOST"`
}

func Load(filePath string) (*Configuration, error) {
	err := godotenv.Load(filePath)

	if err != nil {
		return nil, err
	}

	cfg := new(Configuration)

	return cfg, env.Parse(cfg)
}
