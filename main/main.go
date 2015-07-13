package main

import (
	"flag"
	stdlog "log"
	"os"

	"bytes"

	_ "bitbucket.org/shell909090/scheme-go/internal"
	"bitbucket.org/shell909090/scheme-go/parser"
	"bitbucket.org/shell909090/scheme-go/scm"
	"bitbucket.org/shell909090/scheme-go/tsfm"
	logging "github.com/op/go-logging"
)

var log = logging.MustGetLogger("main")

var (
	LogFile  string
	LogLevel string
	Parse    bool
	Trans    bool
	Run      bool
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

func PrepareMacro() {
	buf := bytes.NewBufferString(PreDefineMacro)
	code, err := parser.SourceToAST(buf)
	if err != nil {
		panic("impossible")
	}
	err = tsfm.DefineTransformer.Parse(code)
	if err != nil {
		panic(err.Error())
	}
	return
}

func parse(filename string) (code scm.Obj, err error) {
	file, err := os.Open(filename)
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
	return
}

func run() (err error) {
	os.Stdout.WriteString("-------run parse-------\n")
	code, err := parse(flag.Args()[0])
	if err != nil {
		return
	}
	if Parse {
		os.Stdout.WriteString("-------parsed-------\n")
		os.Stdout.WriteString(scm.Format(code))
		os.Stdout.Write([]byte("\n"))
	}

	os.Stdout.WriteString("-------transform-------\n")
	code, err = tsfm.DefineTransformer.Transform(code)
	if err != nil {
		log.Error("%s", err)
		return
	}
	if Trans {
		os.Stdout.WriteString("-------compiled-------\n")
		os.Stdout.WriteString(scm.Format(code))
		os.Stdout.WriteString("\n")
	}

	if !Run {
		return
	}
	os.Stdout.WriteString("-------runtime-------\n")
	result, _, err := scm.RunCode(code)
	if err != nil {
		log.Error("%s", err)
		return
	}

	os.Stdout.WriteString("-------output-------\n")
	os.Stdout.WriteString(scm.Format(result))
	os.Stdout.WriteString("\n")
	return
}

func main() {
	flag.StringVar(&LogLevel, "loglevel", "INFO", "loglevel")
	flag.StringVar(&LogFile, "logfile", "", "logfile")
	flag.BoolVar(&Parse, "parse", false, "print parse result")
	flag.BoolVar(&Trans, "transform", false, "print transform result")
	flag.BoolVar(&Run, "run", true, "run code")

	flag.Parse()
	if len(flag.Args()) < 1 {
		panic("parameters not enough")
	}

	err := SetLogging()
	if err != nil {
		panic(err)
	}

	PrepareMacro()
	run()
}
