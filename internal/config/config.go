package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"

	"gihub.com/gmohmad/wb_l0/internal/utils"
)

type Config struct {
	Service
	DB
	Nats
}

type Service struct {
	Env        string `yaml:"env" env-default:"local"`
	HTTPServer `yaml:"http_server"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8000"`
	Timeout     time.Duration `yaml:"timeout" env-default:"5s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

func LoadService() Service {
	configPath := utils.GetEnvOrFatal("CONFIG_PATH")

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("%s no such file or derictory", configPath)
	}

	var service Service

	if err := cleanenv.ReadConfig(configPath, &service); err != nil {
		log.Fatal("can't read config")
	}

	return service
}

type DB struct {
	Host           string
	Port           string
	User           string
	Password       string
	DBName         string
	SSLMode        string
	MigrationsPath string
}

func LoadDB() DB {
	dbHost := utils.GetEnvOrFatal("POSTGRES_HOST")
	dbPort := utils.GetEnvOrFatal("POSTGRES_PORT")
	dbUser := utils.GetEnvOrFatal("POSTGRES_USER")
	dbPassword := utils.GetEnvOrFatal("POSTGRES_PASSWORD")
	dbName := utils.GetEnvOrFatal("POSTGRES_DB")
	dbSslMode := utils.GetEnvOrFatal("SSL_MODE")
	migrationsPath := utils.GetEnvOrFatal("MIGRATIONS_PATH")

	db := DB{
		Host:           dbHost,
		Port:           dbPort,
		User:           dbUser,
		Password:       dbPassword,
		DBName:         dbName,
		SSLMode:        dbSslMode,
		MigrationsPath: migrationsPath,
	}

	return db
}

type Nats struct {
	Host      string
	Port      string
	ClusterId string
	ClientId  string
	SenderId  string
	Subject   string
}

func LoadNats() Nats {
	natsHost := utils.GetEnvOrFatal("NATS_HOST")
	natsPort := utils.GetEnvOrFatal("NATS_PORT")
	natsClusterId := utils.GetEnvOrFatal("NATS_CLUSTER_ID")
	natsClientId := utils.GetEnvOrFatal("NATS_CLIENT_ID")
	natsSenderId := utils.GetEnvOrFatal("NATS_SENDER_ID")
	natsSubject := utils.GetEnvOrFatal("NATS_SUBJECT")

	nats := Nats{
		Host:      natsHost,
		Port:      natsPort,
		ClusterId: natsClusterId,
		ClientId:  natsClientId,
		SenderId:  natsSenderId,
		Subject:   natsSubject,
	}

	return nats
}

func MustLoad() *Config {
	if err := godotenv.Load(); err != nil {
		log.Fatal("can't load the env file")
	}

	service := LoadService()
	db := LoadDB()
	nats := LoadNats()

	cfg := Config{
		Service: service,
		DB:      db,
		Nats:    nats,
	}

	return &cfg
}
