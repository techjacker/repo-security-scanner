package main

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	elastic "gopkg.in/olivere/elastic.v5"
	elogrus "gopkg.in/sohlich/elogrus.v2"

	"github.com/techjacker/diffence"
)

// Log is the log interface for the app
type Log interface {
	Log(v ...interface{})
}

// Log is the log for the app
type Logger struct {
	log *logrus.Logger
}

func (l Logger) Log(v ...interface{}) {
	i := 1
	for filename, rule := range v[0].(diffence.MatchedRules) {
		i++
		l.log.WithFields(logrus.Fields{
			"organisation": v[1],
			"repo":         v[2],
			"url":          v[3],
			"filename":     filename,
			"reason":       rule[0].Caption,
		}).Error(fmt.Sprintf("Violation found in %s:%s", v[1], v[2]))
	}
}

// NewESLogger is a factory for Elasticsearch loggers
func NewESLogger(esUrl, esIndex string) (Logger, error) {
	log := logrus.New()
	client, err := elastic.NewClient(elastic.SetURL(esUrl))
	if err != nil {
		return Logger{}, err
	}
	hook, err := elogrus.NewElasticHook(client, "localhost", logrus.DebugLevel, esIndex)
	if err != nil {
		return Logger{}, err
	}
	log.Hooks.Add(hook)
	return Logger{log: log}, nil
}
