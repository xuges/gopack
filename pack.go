package gopack

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
)

var (
	files      []file
	unpackPath string
	workDir    string
	envs       []string
	args       []string
)

// AddDependency add a dependency file to pack.
//
// file data must create by `go:embed`
func AddDependency(outPath string, embedData []byte) {
	files = append(files, file{
		typ:  dep,
		path: outPath,
		data: embedData,
	})
}

// AddExecutable add an executable file to pack.
//
// fileData must create by `go:embed`
func AddExecutable(outPath string, embedData []byte) {
	files = append(files, file{
		typ:  exe,
		path: outPath,
		data: embedData,
	})
}

// AddEnv add a runtime environment variant(key=value).
func AddEnv(env string) {
	envs = append(envs, env)
}

// AddArg add a command-line argument.
func AddArg(arg string) {
	args = append(args, arg)
}

// SetUnpackPath set the files unpacked path.
func SetUnpackPath(path string) {
	unpackPath = path
}

// SetWorkerDir set the execute process work dir.
//
// execute process default work dir is `.`
func SetWorkerDir(path string) {
	workDir = path
}

// Unpack all files to unpack-pack.
func Unpack() (err error) {
	if len(unpackPath) != 0 {
		err = os.MkdirAll(unpackPath, os.ModeDir)
		if err != nil {
			return
		}
	}

	for _, file := range files {
		path := filepath.Join(unpackPath, file.path)
		err = os.MkdirAll(filepath.Dir(path), os.ModeDir)
		if err != nil {
			return
		}

		err = os.WriteFile(path, file.data, os.ModePerm)
		if err != nil {
			return
		}
	}

	return
}

// Run the unpacked executable file.
func Run() (exitCode int, err error) {
	for _, file := range files {
		if file.typ == exe {
			exeFile, err := filepath.Abs(filepath.Join(unpackPath, file.path))
			if err != nil {
				return -1, errors.New("get unpacked executable abs path err: " + err.Error())
			}
			if len(workDir) != 0 {
				return runExe(workDir, exeFile)
			}
			return runExe(filepath.Dir(exeFile), exeFile)
		}
	}
	return -1, errors.New("no executable")
}

type file struct {
	path string
	data []byte
	typ  fileType
}

type fileType int

const (
	dep fileType = iota
	exe
)

func runExe(workDir, exeFile string) (code int, err error) {
	cmd := exec.Command(exeFile, args...)
	cmd.Dir = workDir
	if len(envs) != 0 {
		cmd.Env = envs
	}

	err = cmd.Start()
	if err != nil {
		return -1, errors.New("start executable err: " + err.Error())
	}

	procState, err := cmd.Process.Wait()
	if err != nil {
		return -1, errors.New("wait execute process err: " + err.Error())
	}

	return procState.ExitCode(), nil
}
