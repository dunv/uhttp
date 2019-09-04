package uhttp

import "log"

func CheckAndLogError(err error) {
	if err != nil {
		if customLog != nil {
			customLog.Errorf(err.Error())
		} else {
			log.Println(err.Error())
		}
	}
}

func CheckAndLogErrorSecondArg(_ interface{}, err error) {
	if err != nil {
		if customLog != nil {
			customLog.Errorf(err.Error())
		} else {
			log.Println(err.Error())
		}
	}
}

func CheckAndLogInfo(err error) {
	if err != nil {
		if customLog != nil {
			customLog.Infof(err.Error())
		} else {
			log.Println(err.Error())
		}
	}
}

func CheckAndLogInfoSecondArg(_ interface{}, err error) {
	if err != nil {
		if customLog != nil {
			customLog.Infof(err.Error())
		} else {
			log.Println(err.Error())
		}
	}
}
