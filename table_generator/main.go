package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"os"
	"strings"

	"project_b/log"
)

const (
	RuntimeRootDir = "../"
	GenerateSrcDir = "src/csv_readers" // 生成代码目录
	GenerateTabDir = "game_csv"        // 生成csv表目录
	ExcelDir       = "game_excel"      // 拷贝excel文件目录
)

type TableGeneratorConfig struct {
	ExcelPath      string // excel文件源目录
	HeaderIndex    int32  // 表头行索引
	ValueTypeIndex int32  // 值类型行索引
	UserTypeIndex  int32  // 用户类型行索引
	DataStartIndex int32  // 数据开始行索引
}

func create_dirs(dest_path string) (err error) {
	if err = os.MkdirAll(dest_path, os.ModePerm); err != nil {
		return
	}
	if err = os.Chmod(dest_path, os.ModePerm); err != nil {
		return
	}
	return
}

/*
func file_copy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}

	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}
*/

var gslog *log.Logger

func main() {
	gslog = log.NewWithConfig(&log.LogConfig{
		Filename:      "./log/table_generator.log",
		MaxSize:       2,
		MaxBackups:    100,
		MaxAge:        30,
		Compress:      false,
		ConsoleOutput: true,
	}, log.DebugLevel)
	defer gslog.Sync()

	var config_path string
	if len(os.Args) > 1 {
		tmp_path := flag.String("f", "", "config path")
		if nil != tmp_path {
			flag.Parse()
			config_path = *tmp_path
		}
	} else {
		gslog.Warn("arguments not enough")
		return
	}

	data, err := ioutil.ReadFile(config_path)
	if err != nil {
		gslog.Error("read config file [%v] failed: %v", config_path, err.Error())
		return
	}

	var config TableGeneratorConfig
	err = json.Unmarshal(data, &config)
	if err != nil {
		gslog.Error("parse config file [%v] failed: %v", config_path, err.Error())
		return
	}

	var files []os.FileInfo
	excel_path := config.ExcelPath
	files, err = ioutil.ReadDir(excel_path)
	if err != nil {
		gslog.Warn("read excel source dir %v failed %v\n ready to read next excel dir", excel_path, err.Error())

		excel_path = RuntimeRootDir + ExcelDir
		files, err = ioutil.ReadDir(excel_path)
		if err != nil {
			gslog.Error("read local excel dir %v failed: %v", excel_path, err.Error())
			return
		}
	}

	for _, f := range files {
		if f.IsDir() {
			continue
		}
		if !strings.HasSuffix(f.Name(), ".xlsx") {
			continue
		}

		gslog.Info("reading excel file %v >>>", f.Name())

		generate_src_path := RuntimeRootDir + GenerateSrcDir
		generate_tab_path := RuntimeRootDir + GenerateTabDir
		if err = GenSourceAndCsv(excel_path, f.Name(), generate_src_path, generate_tab_path, config.HeaderIndex, config.ValueTypeIndex, config.UserTypeIndex, config.DataStartIndex); err != nil {
			gslog.Error("read excel failed: %v", err)
			return
		}

		gslog.Info("<<< read excel file %v end", f.Name())
	}
}
