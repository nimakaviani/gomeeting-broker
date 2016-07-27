package utils

import (
	"encoding/json"
	"os"
	"regexp"
)

var sqlConnectionRegex = regexp.MustCompile(`^.+:\/\/(.+):(.+)@(.+):\d+\/(.+)$`)

type Credentials struct {
	URI string `json:"uri,omitempty"`
}

type ElephantSql struct {
	Credentials Credentials `json:"credentials,omitempty"`
}

type VCAPServices struct {
	Elephantsql []ElephantSql `json:"elephantsql,omitempty"`
}

type VCAPServicesStruct struct {
	VCAP_SERVICES VCAPServices `json:"VCAP_SERVICES"`
}

type DBCredentials struct {
	User     string
	Password string
	Host     string
	Name     string
}

func GetDBCreds() DBCredentials {
	vcapServicesString := os.Getenv("VCAP_SERVICES")

	if vcapServicesString == "" {
		vcapServicesString = "{}"
	}

	vcapServices := VCAPServicesStruct{}

	err := json.Unmarshal([]byte(vcapServicesString), &vcapServices)
	if err != nil {
		panic(err)
	}

	if len(vcapServices.VCAP_SERVICES.Elephantsql) == 0 {
		return DBCredentials{
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Host:     os.Getenv("DB_HOST"),
			Name:     os.Getenv("DB_NAME"),
		}
	}

	matches := sqlConnectionRegex.FindStringSubmatch(vcapServices.VCAP_SERVICES.Elephantsql[0].Credentials.URI)

	return DBCredentials{
		User:     matches[1],
		Password: matches[2],
		Host:     matches[3],
		Name:     matches[4],
	}
}
