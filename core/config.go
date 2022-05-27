package core

import (
	"crypto/aes"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"path"
	"reflect"

	remote "github.com/shima-park/agollo/viper-remote"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

type cfg struct {
	Bootstrap string
	Env       string
	ServiceIp string
	Service   ServiceConfig
}
type ServiceConfig struct {
	LogPath     string      `yaml:"logpath"`
	ServiceName string      `yaml:"servicename"`
	Port        string      `yaml:"port"`
	Databases   DBConfig    `yaml:"databases"`
	Redis       RedisConfig `yaml:"redis"`
	EtcdAddress []string    `yaml:"etcaddress"`
	EtcdKey     string      `yaml:"etckey"`
}

type DBConfig struct {
	Drivername string `yaml:"drivername"`
	Host       string `yaml:"host"`
	Port       string `yaml:"port"`
	Database   string `yaml:"database"`
	Password   string `yaml:"password"`
	Dbuser     string `yaml:"dbuser"`
	Charset    string `yaml:"charset"`
}

//  redis配置RedisConfig
type RedisConfig struct {
	Host     string `yaml:"host"`
	Password string `yaml:"password"`
	Db       int    `yaml:"db"`
}

var (
	AppConfig cfg
)

const (
	seed = "LKP_LIU_KUNPENG"
)

// BootStrapConfig 启动文件配置
type BootStrapConfig struct {
	App struct {
		Env string `yaml:"env"`
		ID  string `yaml:"id"`
	} `yaml:"app"`
	Apollo struct {
		NameSpaces string            `yaml:"namespaces"`
		Meta       map[string]string `yaml:"meta"`
	} `yaml:"apollo"`
}

func (g *Garden) bootConfig(filetype string) {
	AppConfig = cfg{}
	bootStrapConfig := new(BootStrapConfig)

	if g.cfg.Bootstrap == "" {
		configErrorHandler("application configure", errors.New("file name is nil"))
	}
	if path.Ext(g.cfg.Bootstrap) != ".yml" {
		configErrorHandler("application configure file name{%s} ", errors.New("suffix must be .yml"))
	}
	confFileStream, err := ioutil.ReadFile(g.cfg.Bootstrap)
	if err != nil {
		configErrorHandler("ioutil.ReadFile(file:%s) = error:%v", err, g.cfg.Bootstrap)
	}

	configErrorHandler("[Unmarshal]init logger error: %v", yaml.Unmarshal(confFileStream, bootStrapConfig))
	remote.SetAppID(bootStrapConfig.App.ID)
	v := viper.New()
	v.SetConfigType("prop")
	configErrorHandler("AddRemoteProvider", v.AddRemoteProvider("apollo", bootStrapConfig.Apollo.Meta[g.cfg.Env], bootStrapConfig.Apollo.NameSpaces))
	configErrorHandler("ReadRemoteConfig", v.ReadRemoteConfig())
	configErrorHandler("Unmarshal", v.Unmarshal(&AppConfig))
	checkSetting(&AppConfig)
}
func configErrorHandler(fmtStr string, err error, param ...interface{}) {
	if err != nil {
		panic(fmt.Errorf(fmtStr, param, err.Error()))
	}
}
func checkSetting(v interface{}) {
	rt := reflect.TypeOf(v)
	rv := reflect.ValueOf(v)

	rte := rt.Elem()
	rve := rv.Elem()

	for i := 0; i < rte.NumField(); i++ {
		switch rte.Field(i).Type.Kind() {
		case reflect.Struct:
			checkSetting(rve.Field(i).Addr().Interface())
		case reflect.String:
			if dec, err := Decode(rve.Field(i).String()); err == nil {
				rve.Field(i).SetString(dec)
			}
		}
	}
}

// Decode 解码加密配置
func Decode(raw string) (string, error) {
	if len(raw) >= 7 && raw[:7] == "CIPHER(" && raw[len(raw)-1:] == ")" {
		base64Str := raw[7 : len(raw)-1]
		crytedByte, err := base64.StdEncoding.DecodeString(base64Str)
		if err != nil {
			return "", err
		}
		key := aesKeySecureRandom(seed)
		return string(pkcs7UnPadding(decryptAes128Ecb(crytedByte, key))), nil
	}
	return "", errors.New("field not encode")
}

func decryptAes128Ecb(data, key []byte) []byte {
	cipher, _ := aes.NewCipher([]byte(key))
	decrypted := make([]byte, len(data))
	size := 16
	for bs, be := 0, size; bs < len(data); bs, be = bs+size, be+size {
		cipher.Decrypt(decrypted[bs:be], data[bs:be])
	}

	return decrypted
}

// SHA1 签名
func SHA1(data []byte) []byte {
	h := sha1.New()
	h.Write(data)
	return h.Sum(nil)
}

func aesKeySecureRandom(keyword string) (key []byte) {
	data := []byte(keyword)
	hashs := SHA1(SHA1(data))
	key = hashs[0:16]
	return key
}

//pkcs7UnPadding 去补码
func pkcs7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:length-unpadding]
}
