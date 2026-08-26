package config

import (
	"errors"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/spf13/viper"
)

func TestLoadInto(t *testing.T) {
	target := viper.New()
	configs := fstest.MapFS{
		"config.yml": &fstest.MapFile{Data: []byte("service:\n  name: hall\nmysql:\n  port: 3306\n")},
	}

	if err := LoadInto(target, "config.yml", configs); err != nil {
		t.Fatalf("LoadInto() error = %v", err)
	}
	if got := target.GetString("service.name"); got != "hall" {
		t.Fatalf("service.name = %q, want hall", got)
	}
	if got := target.GetInt("mysql.port"); got != 3306 {
		t.Fatalf("mysql.port = %d, want 3306", got)
	}
}

func TestLoadIntoNilViper(t *testing.T) {
	err := LoadInto(nil, "config.yml", fstest.MapFS{})
	if !errors.Is(err, ErrNilViper) {
		t.Fatalf("LoadInto() error = %v, want ErrNilViper", err)
	}
}

func TestInitConfigMissingFilePanics(t *testing.T) {
	viper.Reset()
	t.Cleanup(viper.Reset)

	defer func() {
		value := recover()
		if value == nil {
			t.Fatal("InitConfig() did not panic")
		}
		message, ok := value.(string)
		if !ok || !strings.HasPrefix(message, "read config error:\n") {
			t.Fatalf("panic = %v", value)
		}
	}()

	InitConfig("missing.yml", fstest.MapFS{})
}
