package config

import (
	"bytes"
	"encoding/json"
	"os"
	"text/template"

	goapi "github.com/gdcorp-domains/fulfillment-go-api"
)

type Config struct {
	goapi.Config
	MinWorkers          int32  `json:"minWorkers"`
	MaxWorkers          int32  `json:"maxWorkers"`
	ScaleUpCooldown     int32  `json:"scaleUpCooldown"`
	ScaleDownCooldown   int32  `json:"scaleDowncooldown"`
	MetricsInterval     int32  `json:"metricsInterval"`
	MaxErrorThreshold   int64  `json:"maxErrorThreshold"`
	ErrorWindowDuration int32  `json:"errorWindowDuration"`
	MaxProcessingTime   int32  `json:"maxProcessingTime"`
	VisibilityTimeout   int32  `json:"visibilityTimeout"`
	MaxRetries          int32  `json:"maxRetries"`
	MaxMessages         int32  `json:"maxMessages"`
	WaitTimeSeconds     int32  `json:"waitTimeSeconds"`
	QueueURL            string `json:"queueUrl"`
	DLQURL              string `json:"dlqUrl"`
	QueueName           string `json:"queueName"`
}

func (conf *Config) Load(configPath string) (err error) {
	var fileContents []byte
	fileContents, err = os.ReadFile(configPath)
	if err != nil {
		return err
	}
	envConfig := struct {
		Env    string
		Region string
		EnvDNS string
	}{
		Env:    os.Getenv("ENV"),
		Region: os.Getenv("AWS_REGION"),
		EnvDNS: os.Getenv("ENV"),
	}
	if envConfig.Env == "dev-private" {
		envConfig.EnvDNS = "dp"
	}

	var t *template.Template
	t, err = template.New("").Parse(string(fileContents))
	if err != nil {
		return err
	}

	buf := bytes.Buffer{}
	if err := t.Execute(&buf, envConfig); err != nil {
		return err
	}

	if err := json.Unmarshal(buf.Bytes(), conf); err != nil {
		return err
	}
	return nil
}
