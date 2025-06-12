package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

const AccessControl string = "access-control"
const FTv2 string = "42-v2"
const FTAttendance string = "42-attendance"
const FTFreeze string = "42-freeze"

var ConfigData configFile

type configFile struct {
	AccessControl struct {
		Endpoint string `yaml:"endpoint"`
		TestPath string `yaml:"testpath"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"AccessControl"`
	ApiV2 struct {
		TokenUrl           string   `yaml:"tokenUrl"`
		Endpoint           string   `yaml:"endpoint"`
		TestPath           string   `yaml:"testpath"`
		Uid                string   `yaml:"uid"`
		Secret             string   `yaml:"secret"`
		Scope              string   `yaml:"scope"`
		CampusID           string   `yaml:"campusId"`
		ApprenticeProjects []string `yaml:"apprenticeProjects"`
	} `yaml:"42apiV2"`
	Freeze42 struct {
		AutoPost bool   `yaml:"autoPost"`
		TokenUrl string `yaml:"tokenUrl"`
		Endpoint string `yaml:"endpoint"`
		TestPath string `yaml:"testpath"`
		Uid      string `yaml:"uid"`
		Secret   string `yaml:"secret"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"42Freeze"`
	Attendance42 struct {
		AutoPost bool   `yaml:"autoPost"`
		TokenUrl string `yaml:"tokenUrl"`
		Endpoint string `yaml:"endpoint"`
		TestPath string `yaml:"testpath"`
		Uid      string `yaml:"uid"`
		Secret   string `yaml:"secret"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"42Attendance"`
}

func LoadConfig(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&ConfigData)
	if err != nil {
		return err
	}
	return nil
}
