// Package config 提供基于 Viper 的配置加载工具。
package config

import (
	"bytes"
	"errors"
	"io/fs"

	"github.com/spf13/viper"
)

var ErrNilViper = errors.New("viper instance must not be nil")

// Load 从文件系统读取 YAML 配置并写入 Viper 全局实例。
func Load(name string, configs fs.FS) error {
	return LoadInto(viper.GetViper(), name, configs)
}

// LoadInto 从文件系统读取 YAML 配置并写入指定的 Viper 实例。
func LoadInto(target *viper.Viper, name string, configs fs.FS) error {
	if target == nil {
		return ErrNilViper
	}

	target.SetConfigType("yml")
	content, err := fs.ReadFile(configs, name)
	if err != nil {
		return err
	}
	return target.ReadConfig(bytes.NewReader(content))
}

// InitConfig 从文件系统加载 YAML 配置到 Viper 全局实例。
// 为兼容原有服务，加载失败时会 panic；新代码优先使用 Load。
func InitConfig(name string, configs fs.FS) {
	if err := Load(name, configs); err != nil {
		panic("read config error:\n" + err.Error())
	}
}
