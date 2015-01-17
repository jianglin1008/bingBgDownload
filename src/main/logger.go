package main

import log "github.com/cihub/seelog"

func initLog() {
	testConfig := `
<seelog>
	<outputs formatid="critical">
		<rollingfile type="size" filename="./run.log" maxsize="10000" maxrolls="5" />
	</outputs>
	<formats>
		<format id="critical" format="%Time %Date [%LEV] %Msg %n"/>
	</formats>
</seelog>
`
	var data []byte = []byte(testConfig)
	logger, _ := log.LoggerFromConfigAsBytes(data)

	log.ReplaceLogger(logger)
}
