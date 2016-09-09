package main

import (
	_ "encoding/json"
	"time"

	"gopkg.in/mgo.v2-unstable/bson"
)

type Device struct {
	ID          bson.ObjectId `bson:"_id,omitempty"`
	DeviceId    string        `bson:"device_id"`
	InstallDate time.Time     `bson:"install_date"`
}

type DeviceSession struct {
	ID       bson.ObjectId `bson:"_id,omitempty"`
	Token    string        `bson:"token"`
	DeviceId string        `bson:"device_id"`
}

type Devices struct {
	Devices *[]Device `json:"devices"`
	Total   int       `json:total`
}

type PostDate struct {
	Begin time.Time `json:"begin"`
	End   time.Time `json:"end"`
}

type Sample struct {
	ID          string    `json:"id"`
	DeviceId    string    `json:"device-id"`
	InstallDate time.Time `json:"install-date"`
}

type UsePerDay struct {
	ID    bson.ObjectId `bson:"_id,omitempty"`
	Date  time.Time     `bson:"date"`
	Count int           `bson:"count"`
}

type UsePerDayDto struct {
	Date  []int `json:"date"`
	Count []int `json:"count"`
}

type Password struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Session struct {
	ID     string    `json:"id"`
	Expire time.Time `json:"expire"`
}

type Site struct {
	Url  string    `bson:"url" json:"url"`
	Date time.Time `bson:"updated" json:"date"`
}

type Expire struct {
	Date time.Time `bson:"date"`
}
