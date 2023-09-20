package loadvars

import (
	"log"
	"os"
	"strconv"
)

type EnvVars struct {
	WorkStart    int
	WorkEnd      int
	NameSpace    string
	TimeToDelete int64
}

func (e *EnvVars) LoadVars() *EnvVars {
	WSS := os.Getenv("WORK_START")
	WS, err := strconv.Atoi(WSS)
	if err != nil {
		log.Println("Environment variable WORK_START not set, use default value: 9")
		WS = 9
	}
	e.WorkStart = WS

	WES := os.Getenv("WORK_END")
	WE, err := strconv.Atoi(WES)
	if err != nil {
		log.Println("Environment variable WORK_END not set, use default value: 18")
		WE = 18
	}
	e.WorkEnd = WE

	e.NameSpace = os.Getenv("NAMESPACE")
	if e.NameSpace == "" {
		log.Println("Environment variable NAMESPACE not set, use default value: default")
		e.NameSpace = "default"
	}

	TTDS := os.Getenv("TIME_TO_DELETE")
	TTD, err := strconv.ParseInt(TTDS, 10, 64)
	if err != nil {
		log.Println("Environment variable TIME_TO_DELETE not set, use default value: 86400")
		TTD = 86400
	}
	e.TimeToDelete = TTD

	return &EnvVars{
		WorkStart:    e.WorkStart,
		WorkEnd:      e.WorkEnd,
		NameSpace:    e.NameSpace,
		TimeToDelete: e.TimeToDelete,
	}
}
