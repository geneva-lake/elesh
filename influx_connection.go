package main

import (
	"log"
	"time"

	"github.com/influxdata/influxdb/client/v2"
)

//InfluxConnection for registering api requests
type InfluxConnection struct {
	InfluxClient client.Client
	Point        client.BatchPoints
	mc           MongoConnection
}

const (
	MyDB     = "device_connected"
	MyTestDB = "test_db"
	username = "backend"
	password = "backend"
)

func (con *InfluxConnection) queryDB(cmd string) (res []client.Result, err error) {
	defer func() {
		if e := recover(); e != nil {
			log.Printf("[%v] caught panic: %v", e)
		}
	}()
	q := client.Query{
		Command:  cmd,
		Database: MyDB,
	}
	if response, err := con.InfluxClient.Query(q); err == nil {
		if response.Error() != nil {
			return res, response.Error()
		}
		res = response.Results
	} else {
		return res, err
	}
	return res, nil
}

func (con *InfluxConnection) Init() error {
	ic, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     "http://localhost:8086",
		Username: username,
		Password: password,
	})
	con.InfluxClient = ic
	if err != nil {
		log.Println("Init Error: ", err)
		return err
	}
	p, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  MyDB,
		Precision: "s",
	})
	if err != nil {
		log.Println("Init Error: ", err)
		return err
	}
	con.Point = p
	return nil
}

func (con *InfluxConnection) InitTest() error {
	ic, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     "http://localhost:8086",
		Username: username,
		Password: password,
	})
	con.InfluxClient = ic
	if err != nil {
		log.Println("Error: ", err)
		return err
	}
	p, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  MyTestDB,
		Precision: "s",
	})
	con.Point = p
	if err != nil {
		log.Println("Error: ", err)
		return err
	}
	return nil
}

func (con *InfluxConnection) Write(deviceId string, token string) error {
	tags := map[string]string{"device-id": deviceId}
	fields := map[string]interface{}{
		"token": token,
	}
	pt, err := client.NewPoint("session_time", tags, fields, time.Now())
	if err != nil {
		log.Println("Error: ", err)
		return err
	}
	con.Point.AddPoint(pt)
	err = con.InfluxClient.Write(con.Point)
	if err != nil {
		log.Println("Error: ", err)
		return err
	}
	return nil
}

func (con *InfluxConnection) ApiCount() {
	defer func() {
		if e := recover(); e != nil {
			log.Printf("[%v] caught panic: %v", e)
		}
	}()
	res, err := con.queryDB("SELECT * FROM session_time WHERE time > now() - 1d")
	if err != nil {
		log.Panicln("inf err ", err)
	}
	var count int = 0
	if len(res[0].Series) > 0 {
		for range res[0].Series[0].Values {
			count++
		}
		err = con.mc.WriteUseCount(count, time.Now())
		if err != nil {
			log.Panicln("Error: ", err)
		}
	}
}
