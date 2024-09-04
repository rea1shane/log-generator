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
		TimestampFormat: time.RFC3339Nano,
	})

	logfmtEntries := prepare(logfmtLogger, size)

	go info1(logfmtEntries)
	go error1(logfmtEntries)
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

// Data uploaded in past 15s: sum_over_time({container_name="log-generator"} |= "finished uploading" | logfmt bytes | unwrap bytes [15s])
func info1(entries []*logrus.Entry) {
	for {
		index := randomdata.Number(0, len(entries))
		entries[index].
			WithFields(map[string]interface{}{
				"bytes": randomdata.Number(1, 1024*1024),
			}).
			Info("finished uploading")
		time.Sleep(time.Duration(randomdata.Number(1000*1000, 1000*1000*1000)))
	}
}

// Error count in past 1m: count_over_time({container_name="log-generator"} |= "error occurred in worker" | logfmt worker, error_type [1m])
func error1(entries []*logrus.Entry) {
	for {
		index := randomdata.Number(0, len(entries))
		entries[index].
			WithFields(map[string]interface{}{
				"worker":     randomdata.Number(0, 10),
				"error_type": randomdata.StringSample("network", "core", "disk", "unknown"),
			}).
			Error("error occurred in worker")
		time.Sleep(time.Duration(randomdata.Number(10*1000*1000, 10*1000*1000*1000)))
	}
}
