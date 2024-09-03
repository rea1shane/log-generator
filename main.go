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
	logfmtLogger := logrus.New()
	logfmtLogger.SetFormatter(&logrus.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime: "ts",
		},
	})

	logfmtEntries := prepare(logfmtLogger, size)

	go info1(logfmtEntries)
	select {}
}

func prepare(logger *logrus.Logger, size int) (entries []*logrus.Entry) {
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

func info1(entries []*logrus.Entry) {
	for {
		index := randomdata.Number(0, len(entries))
		entries[index].
			WithField("bytes", randomdata.Number(1, 1024*1024)).
			Info("finished uploading")
		time.Sleep(time.Duration(randomdata.Number(1000*1000, 1000*1000*1000)))
	}
}
