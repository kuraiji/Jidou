package api_server

import (
	"context"
	"log"
	"net/http"
	"time"

	jidouConfig "jidou/internal/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/labstack/echo/v4"
)

const TableName = "JIDOU"

type Post struct {
	Date    string `dynamodbav:"Date"`
	Name    string `dynamodbav:"Name"`
	Message string `dynamodbav:"Message"`
}

func (post Post) GetKey() map[string]types.AttributeValue {
	date, err := attributevalue.Marshal(post.Date)
	if err != nil {
		panic(err)
	}
	/*name, err := attributevalue.Marshal(post.Name)
	if err != nil {
		panic(err)
	}*/
	return map[string]types.AttributeValue{"Date": date /*, "Name": name*/}
}

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
	var err error
	var response *dynamodb.ScanOutput
	var posts []Post

	scanPaginator := dynamodb.NewScanPaginator(db, &dynamodb.ScanInput{
		TableName: aws.String(TableName),
		Limit:     aws.Int32(10),
	})
	for scanPaginator.HasMorePages() {
		response, err = scanPaginator.NextPage(context.TODO())
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		} else {
			var postPage []Post
			err = attributevalue.UnmarshalListOfMaps(response.Items, &postPage)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, err.Error())
			} else {
				posts = append(posts, postPage...)
			}
		}
	}
	return c.JSON(http.StatusOK, posts)
}

func post(c echo.Context, db *dynamodb.Client) error {
	u := new(
		struct {
			Name    string `json:"name"`
			Message string `json:"message"`
		})
	if err := c.Bind(u); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	item, err := attributevalue.MarshalMap(Post{time.Now().Format(time.RFC3339), u.Name, u.Message})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	_, err = db.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(TableName),
		Item:      item,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusCreated, u)
}
