package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type Files struct {
	WorkDir    string
	EnvFile    string
	ConfigFile string
}

func NewFiles(relativeRoot string) Files {
	return NewFilesFromCaller(2, relativeRoot)
}

func NewFilesFromCaller(skip int, relativeRoot string) Files {
	_, file, _, ok := runtime.Caller(skip)
	if !ok {
		panic("runtime caller cannot get file information")
	}

	return NewFilesFromWorkDir(filepath.Join(filepath.Dir(file), relativeRoot))
}

func NewFilesFromWorkDir(workDir string) Files {
	env := "test"
	if value := os.Getenv("ENV"); value != "" {
		switch value {
		case "test":
			env = value
		default:
			panic(fmt.Sprintf("unsupported ENV for test setup: %s", value))
		}
	}

	return Files{
		WorkDir:    workDir,
		EnvFile:    fmt.Sprintf(".env.%s", env),
		ConfigFile: fmt.Sprintf("config.%s.yml", env),
	}
}

func LoadEnv(files Files) {
	if err := godotenv.Load(filepath.Join(files.WorkDir, files.EnvFile)); err != nil {
		panic(err)
	}

	gin.SetMode(gin.ReleaseMode)
}
