package main

import (
	"flag"
	stdlog "log"
	"os"

	_ "bitbucket.org/shell909090/scheme-go/internal"
	"bitbucket.org/shell909090/scheme-go/parser"
	"bitbucket.org/shell909090/scheme-go/scmgo"
	"bitbucket.org/shell909090/scheme-go/tsfm"
	logging "github.com/op/go-logging"
)

var log = logging.MustGetLogger("")

var (
	LogFile  string
	LogLevel string
	Parse    bool
	Trans    bool
)

func SetLogging() (err error) {
	var file *os.File
	file = os.Stdout

	if LogFile != "" {
		file, err = os.OpenFile(LogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)
		if err != nil {
			log.Fatal(err)
		}
	}
	logBackend := logging.NewLogBackend(file, "",
		stdlog.LstdFlags|stdlog.Lmicroseconds|stdlog.Lshortfile)
	logging.SetBackend(logBackend)

	logging.SetFormatter(logging.MustStringFormatter("%{level}: %{message}"))

	lv, err := logging.LogLevel(LogLevel)
	if err != nil {
		panic(err.Error())
	}
	logging.SetLevel(lv, "")

	return
}

func parse() (code scmgo.SchemeObject, err error) {
	file, err := os.Open(flag.Args()[0])
	if err != nil {
		log.Error("%s", err)
		return
	}
	defer file.Close()

	code, err = parser.SourceToAST(file)
	if err != nil {
		log.Error("%s", err)
		return
	}

	os.Stdout.WriteString(code.Format())
	os.Stdout.Write([]byte("\n"))
	return
}

func Transform(src scmgo.SchemeObject) (code scmgo.SchemeObject, err error) {
	code, err = tsfm.Transform(src)
	if err != nil {
		log.Error("%s", err)
		return
	}

	os.Stdout.WriteString("-------transform-------\n")
	os.Stdout.WriteString(code.Format())
	os.Stdout.WriteString("\n")
	return
}

func run(code scmgo.SchemeObject) (result scmgo.SchemeObject, err error) {
	os.Stdout.WriteString("-------runtime-------\n")
	result, err = scmgo.RunCode(code)
	if err != nil {
		log.Error("%s", err)
		return
	}
	os.Stdout.WriteString("-------output-------\n")

	os.Stdout.WriteString(result.Format())
	os.Stdout.WriteString("\n")
	return
}

func main() {
	flag.StringVar(&LogLevel, "loglevel", "INFO", "loglevel")
	flag.StringVar(&LogFile, "logfile", "", "logfile")
	flag.BoolVar(&Parse, "parse", false, "just parse, not run")
	flag.BoolVar(&Trans, "transform", false, "just parse and transform, not run")

	flag.Parse()
	if len(flag.Args()) < 1 {
		panic("parameters not enough")
	}

	err := SetLogging()
	if err != nil {
		panic(err)
	}

	code, err := parse()
	if err != nil {
		log.Error("%s", err)
		return
	}
	if Parse {
		return
	}

	code, err = Transform(code)
	if err != nil {
		log.Error("%s", err)
		return
	}
	if Trans {
		return
	}

	run(code)
}
