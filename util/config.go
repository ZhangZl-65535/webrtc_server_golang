package util

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"regexp"
)

type Server struct {
	Port string
	Mode string
}

type Config struct {
	Server *Server
}

func (c *Config) Load(file string) error {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	content = bytes.Replace(content, []byte("\\"), []byte("\\\\"), -1)
	content = bytes.Replace(content, []byte("\\\\\""), []byte("\\\""), -1)

	comRegex := regexp.MustCompile(`\s*/\*.*\*/`)
	content = comRegex.ReplaceAll(content, []byte{}) //删除注释

	nRegex := regexp.MustCompile(`\n|\t|\r`)
	content = nRegex.ReplaceAll(content, []byte{})

	jsonObj := make(map[string]interface{})
	err = json.Unmarshal(content, &jsonObj)
	if err != nil {
		log.Println(string(content))
		return errors.New("配置文件格式有误:" + err.Error())
	}

	server, ok := jsonObj["server"].(map[string]interface{})
	if ok {
		c.Server = &Server{}
		c.Server.Port = server["port"].(string)
		c.Server.Mode = server["mode"].(string)
	} else {
		return errors.New("解析server时出错")
	}

	return nil
}
