package enviroment

import (
	"auth/utils"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"os"
)

type Environment struct {
	DB  *gorm.DB
	RDB *redis.Client
}

var Env *Environment

func NewEnvironment() *Environment {

	err := godotenv.Load()
	if err != nil {
		println("Error loading env file")
		os.Exit(3)
	}

	utils.InitializeLogger(os.Getenv("LOG_PATH") + "students.log")

	env := &Environment{}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB"), os.Getenv("POSTGRES_PORT"))

	log.Info("connecting to Postgre...")
	env.DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.WithError(err).Fatal("Failed to connect to database")
	}

	log.Info("connected to Postgre")

	log.Info("connecting to Redis...")

	env.RDB = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	log.Info("connected to Redis")

	return env
}
