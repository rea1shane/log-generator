package main

import (
	"time"

	"github.com/Pallinder/go-randomdata"
	"github.com/sirupsen/logrus"
)

const (
	metadataSize = 10
)

var (
	metadata, keys = randomBatchMetadata(metadataSize)
)

func randomBatchMetadata(size int) (map[string]string, []string) {
	data := make(map[string]string)
	for i := 0; i < size; i++ {
		profile := randomdata.GenerateProfile(randomdata.Male | randomdata.Female | randomdata.RandomGender)
		data[profile.Login.Username] = profile.Login.Md5
	}

	ks := make([]string, 0, len(data))
	for key := range data {
		ks = append(ks, key)
	}

	return data, ks
}

func gBatchInfoLogs(logger *logrus.Logger) {
	for {
		number := randomdata.Number(0, len(keys))
		logger.WithFields(map[string]interface{}{
			"job_name": keys[number],
			"job_uuid": metadata[keys[number]],
		}).Infof("written bytes: %d", randomdata.Number(0, 1024*1024))
		time.Sleep(time.Duration(randomdata.Number(1000*1000, 1000*1000*1000)))
	}
}
