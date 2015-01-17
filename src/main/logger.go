package main

import log "github.com/cihub/seelog"

func initLog() {

	logger, _ := log.LoggerFromParamConfigAsFile("log.xml", nil)

	log.ReplaceLogger(logger)
}
