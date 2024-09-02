package main

import (
	"time"

	"github.com/Pallinder/go-randomdata"
	"github.com/sirupsen/logrus"
)

const (
	size = 10
)

func main() {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})

	entries := random(logger, size)

	go mockInfoLogs(entries)
	select {}
}

func random(logger *logrus.Logger, size int) (entries []*logrus.Entry) {
	for i := 0; i < size; i++ {
		profile := randomdata.GenerateProfile(randomdata.Male | randomdata.Female | randomdata.RandomGender)
		entry := logger.WithFields(map[string]interface{}{
			"job_name": profile.Login.Username,
			"job_uuid": profile.Login.Md5,
		})
		entries = append(entries, entry)
	}
	return
}

func mockInfoLogs(entries []*logrus.Entry) {
	for {
		index := randomdata.Number(0, len(entries))
		entries[index].Infof("written bytes: %d", randomdata.Number(0, 1024*1024))
		time.Sleep(time.Duration(randomdata.Number(1000*1000, 1000*1000*1000)))
	}
}
