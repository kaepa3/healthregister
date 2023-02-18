package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/kaepa3/healthplanet"
	"github.com/kaepa3/healthregister/config"
	"github.com/kaepa3/healthregister/db"
	"github.com/kaepa3/healthregister/health"
)

var conf *config.HealthConfig

func main() {
	var err error
	conf, err = config.LoadConfig()
	if err != nil {
		log.Fatal(err)
		return
	}
	data, err := GetData()
	if err != nil {
		fmt.Println(err)
		return
	}
	err = RegisterDB(data)
	if err != nil {
		fmt.Println(err)
	}
}

func RegisterDB(data *healthplanet.JsonResponce) error {
	ctx := context.Background()
	dbOpt := db.RegisterOption{
		Url: "127.0.0.1",
	}
	return db.RegisterDB(data, &dbOpt, ctx)
}

func GetData() (*healthplanet.JsonResponce, error) {
	from, _ := time.Parse("01-02-15-04-05-2006", "01-01-00-00-00-2022")
	ctx := context.Background()
	opt := healthplanet.HealthPlanetOption{
		Format: healthplanet.Json,
		From:   from,
		To:     time.Time{},
	}
	return health.GetHealthData(conf.ClientID, conf.ClientSecret, &opt, ctx)
}
