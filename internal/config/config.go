package config

import (
	"crypto/rand"
	"fmt"
	"github.com/adrg/xdg"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/structs"
	"github.com/knadh/koanf/v2"
	"math/big"
	"os"
	"strings"
	"sync"
)

type Config struct {
	WebDAV     WebDAVConfig     `koanf:"webdav"`
	Encryption EncryptionConfig `koanf:"encryption"`
	Logging    LoggingConfig    `koanf:"logging"`
	Style      StyleConfig      `koanf:"style"`
}

type WebDAVConfig struct {
	URL        string `koanf:"url"`
	RemotePath string `koanf:"remote_path"`
	Username   string `koanf:"username"`
	Password   string `koanf:"password"`
}

type EncryptionConfig struct {
	Enabled  bool   `koanf:"enabled"`
	Password string `koanf:"password"`
}

type LoggingConfig struct {
	LogDir   string `koanf:"logdir"`
	LogLevel string `koanf:"loglevel"`
}

type styleBaseConfig struct {
	Background string `koanf:"background"`
	Surface    string `koanf:"surface"`
	Text       string `koanf:"text"`
	Border     string `koanf:"border"`
	Selection  string `koanf:"selection"`
	Link       string `koanf:"link"`
	Dimmed     string `koanf:"dimmed"`
}

type StyleConfig struct {
	Success styleBaseConfig `koanf:"success"`
	Error   styleBaseConfig `koanf:"error"`
	Info    styleBaseConfig `koanf:"info"`
	Normal  styleBaseConfig `koanf:"normal"`
}

const (
	configFileName = "icali.toml"
	xdgDir         = "icali"
	envVarPrefix   = "ICALI_"
)

var (
	once   sync.Once
	k      = koanf.New(".")
	config Config
)

func Get() Config {
	once.Do(func() {
		// Load the default configuration
		conf := defaultConfig()
		err := k.Load(structs.Provider(conf, ""), nil)
		if err != nil {
			panic(err)
		}

		filePath, err := configFilePath()
		if err != nil {
			filePath = saveCurrentConfig()
		}

		// Load the configuration file
		err = k.Load(file.Provider(filePath), toml.Parser())
		if err != nil {
			panic(fmt.Errorf("error loading configuration file: %w", err))
		}

		// Load configuration from environment files
		// Example: ICALI_WEBDAV_USERNAME
		err = k.Load(env.Provider(envVarPrefix, ".", func(s string) string {
			return strings.Replace(strings.ToLower(
				strings.TrimPrefix(s, envVarPrefix)), "_", ".", -1)
		}), nil)
		if err != nil {
			panic(fmt.Errorf("error loading configuration from env vars: %w", err))
		}

		// Parse configuration into the struct
		if err := k.Unmarshal("", &config); err != nil {
			panic(fmt.Errorf("error unmarshaling configuration: %w", err))
		}
	})

	return config
}

func saveCurrentConfig() string {
	path, err := xdg.ConfigFile(xdgDir + "/" + configFileName)
	if err != nil {
		panic(err)
	}

	// Marshal the default configuration into TOML format
	marshal, err := k.Marshal(toml.Parser())
	if err != nil {
		panic(fmt.Errorf("failed to marshall configuration: %w", err))
	}

	// Write the configuration data to the specified file
	if err := os.WriteFile(path, marshal, 0700); err != nil {
		panic(fmt.Errorf("failed to write configuration file: %w", err))
	}

	return path
}

func defaultConfig() Config {
	logDir, err := xdg.StateFile(xdgDir + "/logs")
	if err != nil {
		panic(fmt.Errorf("failed to get log directory: %w", err))
	}

	return Config{
		WebDAV: WebDAVConfig{
			RemotePath: "",
			Username:   "",
			Password:   "",
		},
		Encryption: EncryptionConfig{
			Enabled:  true,
			Password: generatePassword(),
		},
		Logging: LoggingConfig{
			LogDir:   logDir,
			LogLevel: "info",
		},
		Style: styleConfigFromCat(mocha),
	}
}

func styleConfigFromCat(cat Catppuccin) StyleConfig {
	normal := styleBaseConfig{
		Background: cat.Base,
		Surface:    cat.Surface2,
		Text:       cat.Text,
		Border:     cat.Mauve,
		Selection:  cat.Overlay2,
		Link:       cat.Red,
		Dimmed: 	cat.Subtext0,
	}

	success := styleBaseConfig{
		Background: cat.Base,
		Surface:    cat.Surface2,
		Text:       cat.Text,
		Border:     cat.Green,
		Selection:  cat.Overlay2,
		Link:       cat.Red,
		Dimmed: 	cat.Subtext0,
	}

	error := styleBaseConfig{
		Background: cat.Base,
		Surface:    cat.Surface2,
		Text:       cat.Text,
		Border:     cat.Red,
		Selection:  cat.Overlay2,
		Link:       cat.Red,
		Dimmed: 	cat.Subtext0,
	}

	info := styleBaseConfig{
		Background: cat.Base,
		Surface:    cat.Surface2,
		Text:       cat.Text,
		Border:     cat.Sapphire,
		Selection:  cat.Overlay2,
		Link:       cat.Red,
		Dimmed: 	cat.Subtext0,
	}

	return StyleConfig{
		Normal:  normal,
		Error:   error,
		Info:    info,
		Success: success,
	}
}

func configFilePath() (string, error) {
	// Check if the file exists in the local directory
	if _, err := os.Stat(configFileName); err == nil {
		return configFileName, nil
	}

	// Use XDG config directory fallback
	return xdg.SearchConfigFile(xdgDir + "/" + configFileName)
}

func generatePassword() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()-_=+[]{}|;:,.<>?/`~"
	passwordLength := 32
	bigLen := big.NewInt(int64(len(charset)))

	password := make([]byte, passwordLength)
	for i := range password {
		randomInt, err := rand.Int(rand.Reader, bigLen)

		password[i] = charset[randomInt.Int64()]

		if err != nil {
			panic(err)
		}
	}

	return string(password)
}
