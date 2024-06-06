package config

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Username    string `yaml:"username"`
	Password    string `yaml:"password"`
	PasswordCmd string `yaml:"password_cmd"`
	Host        string `yaml:"host"`
}

func Parse(f string) (*Config, error) {
	b, err := os.ReadFile(f)
	if err != nil {
		return nil, fmt.Errorf("reading config: %w", err)
	}

	c := Config{}
	if err := yaml.Unmarshal(b, &c); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}

	if len(c.Username) == 0 {
		return nil, fmt.Errorf("[ERROR] 'username' must be provided")
	}

	if len(c.Host) == 0 {
		return nil, fmt.Errorf("[ERROR] 'host' must be provided")
	}

	if len(c.Password) == 0 && len(c.PasswordCmd) == 0 {
		return nil, fmt.Errorf("[ERROR] Either 'password' or 'password_cmd' must be provided")
	}

	if len(c.PasswordCmd) > 0 {
		p, err := c.parsePasswordCmd()
		if err != nil {
			return nil, err
		}

		c.Password = p
	}

	return &c, nil
}

func (c *Config) parsePasswordCmd() (string, error) {
	buf := new(bytes.Buffer)
	passwordCmd := strings.Split(c.PasswordCmd, " ")
	cmd := exec.Command(passwordCmd[0], passwordCmd[1:]...)
	cmd.Stdout = buf

	if err := cmd.Run(); err != nil {
		return "", err
	}

	return strings.TrimRight(buf.String(), "\n"), nil
}
