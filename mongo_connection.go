package main

import (
	"fmt"
	"log"
	"time"

	"gopkg.in/mgo.v2-unstable"
	"gopkg.in/mgo.v2-unstable/bson"
)

type MongoConnection struct {
	MongoSession  *mgo.Session
	MongoDataBase *mgo.Database
	User          string
	Password      string
	filters       map[string]string
}

func (mc *MongoConnection) Init() error {
	mc.User = "backend"
	mc.Password = "backend"
	connectS := fmt.Sprintf("%s:%s@localhost:27017/backend", mc.User, mc.Password)
	s, err := mgo.Dial(connectS)
	if err != nil {
		log.Printf(err.Error())
		return err
	}
	s.SetMode(mgo.Monotonic, true)
	mc.MongoSession = s
	mc.MongoDataBase = mc.MongoSession.DB("backend")

	mc.filters = make(map[string]string)
	mc.filters["device-id"] = "device_id"
	mc.filters["install-date"] = "install_date"
	return nil
}

func (mc *MongoConnection) SessionClose() {
	mc.MongoSession.Close()
}

func (mc *MongoConnection) InitTest() error {
	mc.User = "test"
	mc.Password = "test"
	connectS := fmt.Sprintf("%s:%s@localhost:27017/test", mc.User, mc.Password)
	s, err := mgo.Dial(connectS)
	s.SetMode(mgo.Monotonic, true)

	mc.MongoSession = s
	if err != nil {
		log.Printf(err.Error())
		return err
	}
	mc.MongoDataBase = mc.MongoSession.DB("test")
	mc.filters = make(map[string]string)
	mc.filters["device-id"] = "device_id"
	mc.filters["install-date"] = "install_date"
	return nil
}

func (mc *MongoConnection) getSession(deviceId string) (*DeviceSession, error) {
	defer func() {
		if e := recover(); e != nil {
			log.Printf("[%v] caught panic: %v", e)
		}
	}()
	sessionStore := mc.MongoDataBase.C("session")
	deviceSession := DeviceSession{}
	err := sessionStore.Find(bson.M{"device_id": deviceId}).One(&deviceSession)
	if err != nil {
		log.Printf("getSes error ", err.Error())
		return nil, err
	}
	return &deviceSession, nil
}

func (mc *MongoConnection) setSession(deviceId string, token string) (*DeviceSession, error) {
	sessionStore := mc.MongoDataBase.C("session")
	deviceSession := DeviceSession{Token: token, DeviceId: deviceId}
	err := sessionStore.Insert(deviceSession)
	if err != nil {
		log.Printf(err.Error())
		return nil, err
	}
	return &deviceSession, nil
}

func (mc *MongoConnection) getDeviceById(deviceId string) (*Device, error) {
	defer func() {
		if e := recover(); e != nil {
			log.Printf("[%v] caught panic: %v", e)
		}
	}()
	sessionStore := mc.MongoDataBase.C("device")
	device := Device{}
	err := sessionStore.Find(bson.M{"device_id": deviceId}).One(&device)
	if err != nil {
		log.Printf("mongo err", err.Error())
		return nil, err
	}
	return &device, nil
}

func (mc *MongoConnection) getSite() (*Site, error) {
	defer func() {
		if e := recover(); e != nil {
			log.Printf("[%v] caught panic: %v", e)
		}
	}()
	sessionStore := mc.MongoDataBase.C("site")
	site := Site{}
	err := sessionStore.Find(bson.M{}).One(&site)
	if err != nil {
		log.Printf("mongo err site", err.Error())
		return nil, err
	}
	return &site, nil
}

func (mc *MongoConnection) setSite(site Site) error {
	defer func() {
		if e := recover(); e != nil {
			log.Printf("[%v] caught panic: %v", e)
		}
	}()
	sessionStore := mc.MongoDataBase.C("site")
	err := sessionStore.DropCollection()
	if err != nil {
		log.Printf("drop site err ", err.Error())
	}
	err = sessionStore.Insert(site)
	if err != nil {
		log.Printf("insert site err ", err.Error())
		return err
	}
	return nil
}

func (mc *MongoConnection) getPassword() (*Password, error) {
	defer func() {
		if e := recover(); e != nil {
			log.Printf("[%v] caught panic: %v", e)
		}
	}()
	sessionStore := mc.MongoDataBase.C("password")
	psw := Password{}
	err := sessionStore.Find(bson.M{}).One(&psw)
	if err != nil {
		log.Printf("mongo err", err.Error())
		return nil, err
	}
	return &psw, nil
}

func (mc *MongoConnection) getSessionById(deviceId string) (*DeviceSession, error) {
	defer func() {
		if e := recover(); e != nil {
			log.Printf("[%v] caught panic: %v", e)
		}
	}()
	sessionStore := mc.MongoDataBase.C("session")
	deviceSession := DeviceSession{}
	err := sessionStore.Find(bson.M{"device_id": deviceId}).One(&deviceSession)
	if err != nil {
		log.Printf(err.Error())
		return nil, err
	}
	return &deviceSession, nil
}

func (mc *MongoConnection) getAllDevices(skip int, limit int, filter string, order int) (*[]Device, error) {
	defer func() {
		if e := recover(); e != nil {
			log.Printf("[%v] caught panic: %v", e)
		}
	}()
	filter = mc.filters[filter]
	if order == -1 {
		filter = "-" + filter
	}
	deviceCollection := mc.MongoDataBase.C("device")
	var results []Device
	err := deviceCollection.Find(bson.M{}).Skip(skip).Limit(limit).Sort(filter).All(&results)
	if err != nil {
		log.Printf("mongo err ", err.Error())
		return nil, err
	}
	return &results, nil
}

func (mc *MongoConnection) getDevicesCount() (int, error) {
	defer func() {
		if e := recover(); e != nil {
			log.Printf("[%v] caught panic: %v", e)
		}
	}()
	deviceCollection := mc.MongoDataBase.C("device")
	n, err := deviceCollection.Find(bson.M{}).Count()
	if err != nil {
		log.Printf("mongo err count", err.Error())
		return -1, err
	}
	return n, nil
}

func (mc *MongoConnection) CreateDevice(deviceId string, installDate time.Time) (*Device, error) {
	defer func() {
		if e := recover(); e != nil {
			log.Printf("[%v] caught panic: %v", e)
		}
	}()
	sessionStore := mc.MongoDataBase.C("device")
	device := Device{DeviceId: deviceId, InstallDate: installDate}
	err := sessionStore.Insert(device)
	if err != nil {
		log.Printf(err.Error())
		return nil, err
	}
	return &device, nil
}

func (mc *MongoConnection) getUseCount(bd time.Time, ed time.Time) (*[]UsePerDay, error) {

	defer func() {
		if e := recover(); e != nil {
			log.Printf("[%v] caught panic: %v", e)
		}
	}()
	sessionStore := mc.MongoDataBase.C("use_count")
	var results []UsePerDay
	err := sessionStore.Find(bson.M{
		"date": bson.M{"$gt": bd,
			"$lt": ed}}).
		Skip(0).Limit(10).Sort("-use_count").All(&results)
	if err != nil {
		log.Printf(err.Error())
		return nil, err
	}
	return &results, nil
}

func (mc *MongoConnection) WriteUseCount(Count int, installDate time.Time) error {
	defer func() {
		if e := recover(); e != nil {
			log.Printf("[%v] caught panic: %v", e)
		}
	}()
	sessionStore := mc.MongoDataBase.C("use_count")
	usePerDay := UsePerDay{Count: Count, Date: installDate}
	err := sessionStore.Insert(usePerDay)
	if err != nil {
		log.Printf(err.Error())
		return err
	}
	return nil
}

func (mc *MongoConnection) DeleteSession(deviceId string) {
	defer func() {
		if e := recover(); e != nil {
			log.Printf("[%v] caught panic: %v", e)
		}
	}()
	sessionStore := mc.MongoDataBase.C("session")
	_, err := sessionStore.RemoveAll(bson.M{"device_id": deviceId})
	if err != nil {
		log.Printf(err.Error())
	}
}

func (mc *MongoConnection) DeleteSessions() {
	defer func() {
		if e := recover(); e != nil {
			log.Printf("[%v] caught panic: %v", e)
		}
	}()
	sessionStore := mc.MongoDataBase.C("session")
	err := sessionStore.DropCollection()
	if err != nil {
		log.Printf(err.Error())
	}
}

func (mc *MongoConnection) DeleteDevices() {
	defer func() {
		if e := recover(); e != nil {
			log.Printf("[%v] caught panic: %v", e)
		}
	}()
	sessionStore := mc.MongoDataBase.C("device")
	err := sessionStore.DropCollection()
	if err != nil {
		log.Printf(err.Error())
	}
}

func (mc *MongoConnection) DeleteCounts() {
	defer func() {
		if e := recover(); e != nil {
			log.Printf("[%v] caught panic: %v", e)
		}
	}()
	sessionStore := mc.MongoDataBase.C("use_count")
	err := sessionStore.DropCollection()
	if err != nil {
		log.Printf(err.Error())
	}
}

func (mc *MongoConnection) getCookieExp() (*time.Time, error) {
	defer func() {
		if e := recover(); e != nil {
			log.Printf("[%v] caught panic: %v", e)
		}
	}()
	sessionStore := mc.MongoDataBase.C("expire")
	expire := Expire{}
	err := sessionStore.Find(bson.M{}).One(&expire)
	if err != nil {
		log.Printf("mongo err cookie", err.Error())
		return nil, err
	}
	return &expire.Date, nil
}

func (mc *MongoConnection) setCookieExp(expireTime time.Time) error {
	defer func() {
		if e := recover(); e != nil {
			log.Printf("[%v] caught panic: %v", e)
		}
	}()
	sessionStore := mc.MongoDataBase.C("expire")
	err := sessionStore.DropCollection()
	if err != nil {
		log.Printf(err.Error())
	}
	var expire Expire
	expire.Date = expireTime
	err = sessionStore.Insert(expire)
	if err != nil {
		log.Printf(err.Error())
		return err
	}
	return nil
}
