package mutasync

import (
	"errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

func parseYaml(filepath string,filetype string) (*Compose,*Sync, error) {
	var (
		bytes   []byte
		err     error
	)

	// Parse it as YML
	dataSync := &Sync{}
	dataCompose := &Compose{}

	bytes, err = ioutil.ReadFile(filepath)
	if err != nil {
		log.Fatal(err)
		return dataCompose, dataSync, err
	}

	switch filetype {
	case "docker-sync":
			err = yaml.Unmarshal(bytes, &dataSync)
			if err != nil {
				log.Fatal(err)
			}
		break
	case "docker-compose":
			err = yaml.Unmarshal(bytes, &dataCompose)
			if err != nil {
				log.Fatal(err)
			}
		break
	default:
		err := errors.New("No supported type")
		return dataCompose, dataSync, err
	}

	return dataCompose,dataSync,err
}

func ParseCompose(filepath string) (*Compose, error) {
	docker, _, err := parseYaml(filepath, "docker-compose")
	return docker, err
}

func ParseSync(filepath string) (*Sync, error) {
	_, sync, err := parseYaml(filepath, "docker-sync")
	return sync, err
}