package main

import (
	"time"

	"github.com/Pallinder/go-randomdata"
	"github.com/antonfisher/nested-logrus-formatter"
	"github.com/sirupsen/logrus"
)

const (
	size = 10
)

func main() {
	defaultLogger := logrus.New()
	defaultLogger.SetFormatter(&formatter.Formatter{
		TimestampFormat: "2006-01-02 | 15:04:05",
		HideKeys:        true,
	})
	defultEntries := prepare(defaultLogger, size)

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

	// Data uploaded in past 15s: sum_over_time({container_name="log-generator"} |= "finished uploading" | logfmt bytes | unwrap bytes [15s])
	go info1(logfmtEntries)
	// Status from pending to running in past 1m in each host: sum by (host_name, host_ip) (count_over_time({container_name="log-generator"} |= "status updated" | pattern "<_> [<_>] [<job_name>] [<job_uuid>]<_>status updated from <before> to <after>\n" | before="pending" | after="running" [1m]))
	// Status from running to paused in past 1m in each host: sum by (host_name, host_ip) (count_over_time({container_name="log-generator"} |= "status updated" | pattern "<_> [<_>] [<job_name>] [<job_uuid>]<_>status updated from <before> to <after>\n" | before="running" | after="paused" [1m]))
	// Status from paused to running in past 1m in each host: sum by (host_name, host_ip) (count_over_time({container_name="log-generator"} |= "status updated" | pattern "<_> [<_>] [<job_name>] [<job_uuid>]<_>status updated from <before> to <after>\n" | before="paused" | after="running" [1m]))
	go info2(defultEntries)
	// Worker current status in last 15s: last_over_time({container_name="log-generator"} |= "current status" | pattern "<_> [<_>] [<job_name>] [<job_uuid>]<_>worker <worker> current status is <status>\n" | label_format status = `{{ regexReplaceAll "running" .status "1" }}` | label_format status = `{{ regexReplaceAll "paused" .status "0.5" }}` | unwrap status [15s])
	go info3(defultEntries)
	// Error count in past 1m: count_over_time({container_name="log-generator"} |= "error occurred in worker" | logfmt worker, error_type [1m])
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

func info1(entries []*logrus.Entry) {
	for {
		index := randomdata.Number(0, len(entries))
		entries[index].
			WithFields(map[string]interface{}{
				"worker": randomdata.Number(0, 10),
				"bytes":  randomdata.Number(1, 1024*1024),
			}).
			Info("finished uploading")
		time.Sleep(time.Duration(randomdata.Number(1000*1000, 1000*1000*1000)))
	}
}

func info2(entries []*logrus.Entry) {
	for {
		index := randomdata.Number(0, len(entries))
		entries[index].Infof(
			"status updated from %s to %s",
			randomdata.StringSample("pending", "running", "paused"),
			randomdata.StringSample("pending", "running", "paused"),
		)
		time.Sleep(time.Duration(randomdata.Number(1000*1000, 1000*1000*1000)))
	}
}

func info3(entries []*logrus.Entry) {
	for {
		index := randomdata.Number(0, len(entries))
		entries[index].Infof(
			"worker %d current status is %s",
			randomdata.Number(0, 10),
			randomdata.StringSample("running", "paused"),
		)
		time.Sleep(time.Duration(randomdata.Number(1000*1000, 1000*1000*1000)))
	}
}

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
