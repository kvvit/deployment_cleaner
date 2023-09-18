package loadvars

import (
	"log"
	"os"
	"strconv"
)

type EnvVars struct {
	workStart    int
	workEnd      int
	nameSpace    string
	timeToDelete int64
}

func loadVars() EnvVars {
	var envvars EnvVars

	workStartStr := os.Getenv("WORK_START")
	workStart, err := strconv.Atoi(workStartStr)
	if err != nil {
		log.Fatal("Environment variable WORK_START not set: ", err)
	}
	envvars.workStart = workStart

	workEndStr := os.Getenv("WORK_END")
	workEnd, err := strconv.Atoi(workEndStr)
	if err != nil {
		log.Fatal("Environment variable WORK_END not set: ", err)
	}
	envvars.workEnd = workEnd

	envvars.nameSpace = os.Getenv("NAMESPACE")

	timeToDeleteStr := os.Getenv("TIME_TO_DELETE")
	timeToDelete, err := strconv.ParseInt(timeToDeleteStr, 10, 64)
	if err != nil {
		log.Fatal("Environment variable TIME_TO_DELETE not set: ", err)
	}
	envvars.timeToDelete = timeToDelete
	return envvars
}
