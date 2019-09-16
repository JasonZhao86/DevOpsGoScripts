package main

import (
	"strconv"
	"errors"
	"reflect"
	"strings"
	"io/ioutil"
	"fmt"
)

type config struct{
	Filename string `conf:"file_name" db:"file_name"`
	Filepath string `conf:"file_path"`
	Maxsize int64 `conf:"max_size"`
	Debug string `conf:"debug"`
	Password string `conf:"password"`
}


func parseconfig(file string, result interface{}) (err error){
	t := reflect.TypeOf(result)
	v := reflect.ValueOf(result)

	if t.Kind() != reflect.Ptr{
		err = errors.New("result参数必须是一个指针！")
		return 
	}

	tElem := t.Elem()
	if tElem.Kind() != reflect.Struct{
		err = errors.New("result参数必须是一个机构体指针！")
		return
	}

	data, err := ioutil.ReadFile(file) 
	if err != nil{
		err = fmt.Errorf("打开文件%s失败，原因是：%s", file, err)
		return
	}

	datastr := string(data)
	lines := strings.Split(datastr, "\r\n")    // 切割后，每一行行末尾的换行符都被切没了。

	for linenum, line := range lines{
		line = strings.TrimSpace(line)       // 去掉每一行首位的空格
		if len(line) == 0{
			continue
		}

		if strings.HasPrefix(line, "#"){
			continue
		}

		splitpoint := strings.Index(line, "=")
		if splitpoint == -1{
			err = fmt.Errorf("配置文件的第%d行配置错误，配置项中必须要有=号", linenum+1)
			return err
		}
		key := strings.TrimSpace(line[:splitpoint])
		if len(key) == 0{
			err = fmt.Errorf("配置文件的第%d行配置错误，配置项中必须要有=号前面必须要有配置项", linenum+1)
			return
		}
		value := strings.TrimSpace(line[splitpoint+1:])
		valueindex := strings.LastIndex(value, "#")    // 从右到左找#注释
		if valueindex != -1{        // value部分依然有注释
			value = value[:valueindex]
			value = strings.TrimSpace(value)
		}

		// value可能是用引号引起来的值。
		if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\""){
			value = strings.TrimPrefix(value, "\"")
			value = strings.TrimSuffix(value, "\"")
		}

		for i:=0;i<tElem.NumField();i++{
			fieldname := tElem.Field(i).Tag.Get("conf")  // 配置文件中的key和结构体中的key不同
			if fieldname == key{
				switch tElem.Field(i).Type.Kind() {
				case reflect.String:
					filedvalue := v.Elem().FieldByName(tElem.Field(i).Name)
					filedvalue.SetString(value)
				case reflect.Int64:
					value64, _ := strconv.ParseInt(value, 10, 64) // 获取到的value是字符串，将字符串转换成数字
					v.Elem().Field(i).SetInt(value64)
				}
			}
		}
	}
	return
}

func main(){
	result := &config{}    // 初始值
	err := parseconfig("./demo.conf", result)
	if err != nil{
		panic(err)
	}
	fmt.Printf("%#v", result)
}