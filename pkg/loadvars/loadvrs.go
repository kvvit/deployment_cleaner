/*
This module cheks environment variables, stop application if they aren't set, and load them in the stucture
with appropiate type. This variable is used in the main function.
*/
package loadvars

import (
	"log"
	"os"
	"strconv"
)

type EnvVars struct {
	WorkStart      int
	WorkEnd        int
	NameSpace      string
	TimeToDelete   int64
	DeploymentName string
}

func LoadVars() EnvVars {
	var envvars EnvVars

	workStartStr := os.Getenv("WORK_START")
	workStart, err := strconv.Atoi(workStartStr)
	if err != nil {
		log.Fatal("Environment variable WORK_START not set: ", err)
	}
	envvars.WorkStart = workStart

	workEndStr := os.Getenv("WORK_END")
	workEnd, err := strconv.Atoi(workEndStr)
	if err != nil {
		log.Fatal("Environment variable WORK_END not set: ", err)
	}
	envvars.WorkEnd = workEnd

	envvars.NameSpace = os.Getenv("NAMESPACE")
	if envvars.NameSpace == "" {
		log.Fatal("Environment variable NAMESPACE not set")
	}

	envvars.DeploymentName = os.Getenv("DEPLOYMENT_NAME")
	if envvars.DeploymentName == "" {
		log.Fatal("Environment variable DEPLOYMENT_NAME not set")
	}

	timeToDeleteStr := os.Getenv("TIME_TO_DELETE")
	timeToDelete, err := strconv.ParseInt(timeToDeleteStr, 10, 64)
	if err != nil {
		log.Fatal("Environment variable TIME_TO_DELETE not set: ", err)
	}
	envvars.TimeToDelete = timeToDelete
	return envvars
}
