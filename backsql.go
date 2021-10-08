package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

var (
	rmTime        float64
	sqlBackPath   string
	sqlLogPath    string
	sqlPasswd     string
	sqlUser       string
	WarningLogger *log.Logger
	InfoLogger    *log.Logger
	ErrorLogger   *log.Logger
)

func arguments() {
	var back string
	var restore string
	flag.StringVar(&back, "back", "", "这里是备份")
	flag.StringVar(&restore, "restore", "", "这里是还原，后面跟还原的sql文件")
	flag.Parse()
	if back == "" && restore == "" {
		fmt.Println("请在执行脚本的时候传入参数，如不清楚请加-help查看")
	} else {
		if back != "" {
			backup()
		} else {
			restores(restore)
			fmt.Println(time.Now().Format("2006-01-02"))
		}
	}
}

func panJson() {
	ex, err := os.Executable()
	if err != nil {
		os.Exit(3)
	} else {
		workPath := filepath.Dir(ex)
		_, err := os.Stat(workPath + "/sqlBack.json")
		if err == nil {
			jsonFile, err := os.Open(workPath + "/sqlBack.json")
			if err != nil {
				os.Exit(4)
			}
			defer jsonFile.Close()
			byteValue, _ := ioutil.ReadAll(jsonFile)
			var result map[string]interface{}
			json.Unmarshal([]byte(byteValue), &result)
			rmTime = result["rm_time"].(float64)
			sqlBackPath = result["sql_backPath"].(string)
			sqlLogPath = result["sql_logPath"].(string)
			sqlPasswd = result["sql_passwd"].(string)
			sqlUser = result["sql_user"].(string)
		} else {
			ErrorLogger.Println("找不到json文件")
			os.Exit(3)
		}
	}
}

//func oldCe() {
//	t1 := make(map[string]interface{})
//	t1["sql_user"] = "standard"
//	t1["sql_passwd"] = "@fshJo2QFT^mCJ(*')"
//	t1["rm_time"] = 10
//	t1["sql_backPath"] = "/macrosan_backup/mysql_backup/"
//	t1["sql_logPath"] = "/macrosan_backup/mysql_backup/"
//	b, err := json.MarshalIndent(t1, "", "    ")
//	if err != nil {
//		fmt.Println("json err:", err)
//	}
//	fmt.Println(reflect.TypeOf(b), string(b))
//	f, err := os.Create("test.json")
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	l, err := f.Write(b)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	fmt.Println(l)
//}

func logg() {
	logfile := sqlLogPath + "sqlback.log"
	file, err := os.OpenFile(logfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	InfoLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarningLogger = log.New(file, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}
func zhi(ming string) {
	cmd := exec.Command("/bin/bash", "-c", ming)
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()
	err := cmd.Run()
	bytes, _ := ioutil.ReadAll(stderr)
	if err != nil {
		ErrorLogger.Println(ming,string(bytes))
		os.Exit(3)
		return
	}else {
		bytes, _ := ioutil.ReadAll(stdout)
		InfoLogger.Println(ming,string(bytes))
	}
}
func backup() {
	ming := fmt.Sprintf("mysqldump --single-transaction -u%s -p'%s' -h127.0.0.1 -P11306 --all-databases > %s%sback.sql", sqlUser, sqlPasswd, sqlBackPath, time.Now().Format("2006-01-02"))
	zhi(ming)
}

func restores(can string) {
	ming := fmt.Sprintf("mysql -u%s -p'%s' -h127.0.0.1 -P11306 -e \"source %s\"", sqlUser, sqlPasswd, can)
	zhi(ming)
}

func main() {
	logg()
	panJson()
	ming := fmt.Sprintf("rm -rf `find %s -name backup* -ctime +%v`", sqlLogPath, rmTime)
	zhi(ming)
	arguments()
}
