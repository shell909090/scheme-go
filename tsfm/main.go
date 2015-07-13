package main

import (
	"flag"
	stdlog "log"
	"os"

	"bitbucket.org/shell909090/scheme-go/parser"
	"bitbucket.org/shell909090/scheme-go/scm"

	"github.com/op/go-logging"
)

var (
	LogFile       string
	LogLevel      string
	BaseMacroFile string
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

func PrepareMacro() {
	code, err := parse(BaseMacroFile)
	if err != nil {
		panic("impossible")
	}
	err = DefineTransformer.Parse(code)
	if err != nil {
		panic(err.Error())
	}
	return
}

func main() {
	flag.StringVar(&LogLevel, "loglevel", "INFO", "loglevel")
	flag.StringVar(&LogFile, "logfile", "", "logfile")
	flag.StringVar(&BaseMacroFile, "base", "", "base macro file")

	flag.Parse()
	if len(flag.Args()) < 1 {
		panic("parameters not enough")
	}

	err := SetLogging()
	if err != nil {
		panic(err)
	}

	PrepareMacro()

	code, err := parse(flag.Args()[0])
	if err != nil {
		return
	}

	code, err = DefineTransformer.Transform(code)
	if err != nil {
		log.Error("%s", err)
		return
	}

	os.Stdout.WriteString("package main\nconst PreDefineMacro = `")
	os.Stdout.WriteString(scm.Format(code))
	os.Stdout.WriteString("`\n")
	return
}
