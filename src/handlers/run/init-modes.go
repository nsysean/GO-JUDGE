package run_service

import (
	"encoding/xml"
	"io"
	"os"
	"strconv"

	models "app/src/models"
)

func InitModes(task int) models.Problem {
	configFile, err := os.Open("../../tasks/" + strconv.Itoa(task) + "/problem.xml")
	if err != nil {
		panic(err)
	}

	defer configFile.Close()

	byteValue, _ := io.ReadAll(configFile)
	var taskConf models.Problem

	xml.Unmarshal(byteValue, &taskConf)

	return taskConf;
}