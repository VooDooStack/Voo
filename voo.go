package voo

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

const version = "0.0.1"

type Voo struct {
	AppName  string
	Debug    bool
	Version  string
	ErrorLog *log.Logger
	InfoLog  *log.Logger
	RootPath string
	config   config
}

type config struct {
	port     string
	renderer string
}

func (v *Voo) New(rootPath string) error {
	pathConfig := initPaths{
		rootPath:    rootPath,
		folderNames: []string{"migrations", "handlers", "logs", "data", "tmp", "views", "public", "tmp", "middlewares"},
	}
	err := v.Init(pathConfig)
	if err != nil {

		return err
	}

	err = v.checkDotEnv(rootPath)
	if err != nil {
		return err
	}

	//! read .env file
	err = godotenv.Load(rootPath + "/.env")
	if err != nil {
		return err
	}

	//! create loggers
	infoLog, errorLog := v.startLoggers()
	v.InfoLog = infoLog
	v.ErrorLog = errorLog
	v.Debug, _ = strconv.ParseBool(os.Getenv("DEBUG"))
	v.Version = version
	v.RootPath = rootPath

	v.config = config{
		port:     os.Getenv("PORT"),
		renderer: os.Getenv("RENDERER"),
	}

	return nil
}

func (v *Voo) Init(p initPaths) error {
	root := p.rootPath
	for _, path := range p.folderNames {
		//! if folder does not exits, create it
		err := v.CreateDirIfNotExits(root + "/" + path)
		if err != nil {
			return err
		}
	}

	return nil
}

func (v *Voo) checkDotEnv(path string) error {
	err := v.CreateFileIfNotExists(fmt.Sprintf("%s/.env", path))
	if err != nil {
		return err
	}
	return nil
}

func (v *Voo) startLoggers() (*log.Logger, *log.Logger) {
	var infoLog *log.Logger
	var errorLog *log.Logger

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog = log.New(os.Stdout, "ERROR\t ", log.Ldate|log.Ltime|log.Lshortfile)

	return infoLog, errorLog
}