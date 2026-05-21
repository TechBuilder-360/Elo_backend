package actions

import (
	"encoding/csv"
	"os"

	"github.com/Toflex/directory_v2/pkg/log"
)

type Service struct {
	Name           string
	Identifier     string
	Min            interface{}
	Max            interface{}
	Fee            interface{}
	ActiveProvider string
}

var services []Service

func loadServiceFile(l log.Entry) {
	file, err := os.Open("static/service/services.csv")
	if err != nil {
		l.WithError(err).Error("failed to open services.csv file")
	}
	defer file.Close()

	rows, err := csv.NewReader(file).ReadAll()
	if err != nil {
		l.WithError(err).Error("failed to read file")
	}

	for i, row := range rows {
		if i == 0 {
			continue // skip header
		}

		services = append(services, Service{
			Name:           row[0],
			Identifier:     row[1],
			ActiveProvider: row[2],
			Min:            row[3],
			Max:            row[4],
			Fee:            row[5],
		})
	}
}

func LoadServices(l log.Entry) {
	// load files
	loadServiceFile(l)

}
