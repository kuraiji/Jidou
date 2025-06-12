package api_server

import (
	"context"
	"log"
	"net/http"
	"strings"

	jidouConfig "jidou/internal/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/labstack/echo/v4"
)

func ServerLoop() {
	jidouCfg, err := jidouConfig.LoadConfiguration()
	if err != nil {
		log.Fatalf("unable to load Jidou config, %v", err)
	}
	awsCfg, err := awsConfig.LoadDefaultConfig(context.TODO(), awsConfig.WithRegion("us-west-1"))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	dbSvc := dynamodb.NewFromConfig(awsCfg)
	smmSvc := ssm.NewFromConfig(awsCfg)
	param, err := smmSvc.GetParameter(context.TODO(), &ssm.GetParameterInput{
		Name: aws.String(jidouCfg.ParameterName),
	})
	e := echo.New()
	e.GET("/", func(c echo.Context) error { return get(c, dbSvc) })
	e.POST("/", func(c echo.Context) error { return post(c, dbSvc) })
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			key := c.Request().Header.Get(echo.HeaderAuthorization)
			if key == "" {
				return echo.ErrUnauthorized
			}
			if key != *param.Parameter.Value {
				return echo.ErrUnauthorized
			}
			return next(c)
		}
	})
	e.Logger.Fatal(e.Start(":1323"))
}

func get(c echo.Context, db *dynamodb.Client) error {
	resp, err := db.ListTables(context.TODO(), &dynamodb.ListTablesInput{
		Limit: aws.Int32(5),
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	var strBuilder strings.Builder
	strBuilder.WriteString("Tables:\t[")
	for index, tableName := range resp.TableNames {
		strBuilder.WriteString(tableName)
		if index != len(resp.TableNames)-1 {
			strBuilder.WriteString(",")
		} else {
			strBuilder.WriteString("]")
		}
	}
	return c.String(http.StatusOK, strBuilder.String())
}

func post(c echo.Context, db *dynamodb.Client) error {
	return c.String(http.StatusOK, "Hello, World!")
}
