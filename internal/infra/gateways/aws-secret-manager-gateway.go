package gateways

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type AwsSecretManagerGateway struct {
	SecretManager *secretsmanager.Client
}

func (s *AwsSecretManagerGateway) Get(key string) (string, error) {
	if _, ok := os.LookupEnv("AWS_SECRET_MANAGER_NAME"); !ok {
		return "", errors.New("environment variable 'AWS_SECRET_MANAGER_NAME' not set")
	}

	result, err := s.SecretManager.GetSecretValue(context.TODO(), &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(os.Getenv("AWS_SECRET_MANAGER_NAME")),
	})

	if err != nil {
		return "", err
	}

	var secrets map[string]string
	err = json.Unmarshal([]byte(*result.SecretString), &secrets)
	if err != nil {
		return "", err
	}

	value, exists := secrets[key]

	if !exists {
		return "", fmt.Errorf("key %s not found in secrets", key)
	}

	return value, nil
}
