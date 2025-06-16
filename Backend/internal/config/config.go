package config

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	_ "github.com/joho/godotenv/autoload"
)

type Configuration struct {
	ParameterName        string
	ClusterEndpoint      string
	Port                 string
	User                 string
	DatabaseName         string
	Region               string
	ApiKey               string
	TokenRefreshInterval int
}

func LoadConfiguration() (*Configuration, error) {
	var TokenRI = 900
	tri, err := strconv.Atoi(os.Getenv("TOKEN_REFRESH_INTERVAL"))
	if err == nil {
		TokenRI = tri
	}
	cfg := &Configuration{
		ParameterName:        os.Getenv("SSM_PARAMETER_NAME"),
		ClusterEndpoint:      os.Getenv("CLUSTER_ENDPOINT"),
		Port:                 os.Getenv("DB_PORT"),
		User:                 os.Getenv("CLUSTER_USER"),
		DatabaseName:         os.Getenv("DB_NAME"),
		Region:               os.Getenv("REGION"),
		ApiKey:               os.Getenv("API_KEY"),
		TokenRefreshInterval: TokenRI,
	}
	awsCfg, err := awsConfig.LoadDefaultConfig(context.TODO(), awsConfig.WithRegion(cfg.Region))
	if err != nil {
		return nil, fmt.Errorf("failed to generate AWS config: %v", err)
	}
	ssmSvc := ssm.NewFromConfig(awsCfg)
	param, err := ssmSvc.GetParameter(context.TODO(), &ssm.GetParameterInput{Name: aws.String("/JIDOU-API/EC2_KEY")})
	if err == nil {
		cfg.ApiKey = strings.Clone(*param.Parameter.Value)
	}
	param, err = ssmSvc.GetParameter(context.TODO(), &ssm.GetParameterInput{Name: aws.String("/JIDOU/CLUSTER_ENDPOINT")})
	if err == nil {
		cfg.ClusterEndpoint = strings.Clone(*param.Parameter.Value)
	}
	param, err = ssmSvc.GetParameter(context.TODO(), &ssm.GetParameterInput{Name: aws.String("/JIDOU/DB_PORT")})
	if err == nil {
		cfg.Port = strings.Clone(*param.Parameter.Value)
	}
	param, err = ssmSvc.GetParameter(context.TODO(), &ssm.GetParameterInput{Name: aws.String("/JIDOU/CLUSTER_USER")})
	if err == nil {
		cfg.User = strings.Clone(*param.Parameter.Value)
	}
	param, err = ssmSvc.GetParameter(context.TODO(), &ssm.GetParameterInput{Name: aws.String("/JIDOU/DB_NAME")})
	if err == nil {
		cfg.DatabaseName = strings.Clone(*param.Parameter.Value)
	}
	param, err = ssmSvc.GetParameter(context.TODO(), &ssm.GetParameterInput{Name: aws.String("/JIDOU/TOKEN_REFRESH_INTERVAL")})
	if err == nil {
		tri, err := strconv.Atoi(*param.Parameter.Value)
		if err == nil {
			cfg.TokenRefreshInterval = tri
		}
	}
	return cfg, nil
}
