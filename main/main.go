package main

import (
	"flag"
	stdlog "log"
	"os"

	"bitbucket.org/shell909090/scheme-go/scmgo"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("")

func SetLogging() (err error) {
	var file *os.File
	file = os.Stdout

	// if cfg.Logfile != "" {
	// 	file, err = os.OpenFile(cfg.Logfile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// }
	logBackend := logging.NewLogBackend(file, "",
		stdlog.LstdFlags|stdlog.Lmicroseconds|stdlog.Lshortfile)
	logging.SetBackend(logBackend)

	logging.SetFormatter(logging.MustStringFormatter("%{level}: %{message}"))

	lv, err := logging.LogLevel("DEBUG")
	if err != nil {
		panic(err.Error())
	}
	logging.SetLevel(lv, "")

	return
}

func main() {
	flag.Parse()
	if len(flag.Args()) < 1 {
		panic("parameters not enough")
	}

	err := SetLogging()
	if err != nil {
		panic(err)
	}

	file, err := os.Open(flag.Args()[0])
	if err != nil {
		log.Error("%s", err)
		return
	}
	defer file.Close()

	code, err := scmgo.BuildCode(file)
	if err != nil {
		log.Error("%s", err)
		return
	}

	code.Format(os.Stdout, 0)
	os.Stdout.Write([]byte("\n"))
	os.Stdout.Write([]byte("-------output-------\n"))

	result, err := scmgo.RunCode(code)
	if err != nil {
		log.Error("%s", err)
		return
	}

	result.Format(os.Stdout, 0)
	os.Stdout.Write([]byte("\n"))
}
