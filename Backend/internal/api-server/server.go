package api_server

import (
	"context"
	"fmt"
	jidouConfig "jidou/internal/config"
	jidouDSQL "jidou/internal/dsql"
	"log"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

const TableName = "jidou"

type PostReq struct {
	Name    string `json:"name"`
	Message string `json:"message"`
}
type PostResp struct {
	Date time.Time `json:"date"`
	PostReq
}

func ServerLoop() {
	ctx := context.Background()
	jidouCfg, err := jidouConfig.LoadConfiguration()
	if err != nil {
		log.Fatalf("unable to load Jidou config, %v", err)
	}
	awsCfg, err := awsConfig.LoadDefaultConfig(context.TODO(), awsConfig.WithRegion(jidouCfg.Region))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	poolWrapper, err := jidouDSQL.NewPool(ctx)
	if err != nil {
		log.Fatalf("unable to make new pool, %v", err)
	}
	defer poolWrapper.Close()
	pool := poolWrapper.Pool
	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("unable to ping database, %v", err)
	}
	smmSvc := ssm.NewFromConfig(awsCfg)
	param, err := smmSvc.GetParameter(context.TODO(), &ssm.GetParameterInput{
		Name: aws.String(jidouCfg.ParameterName),
	})
	e := echo.New()
	e.GET("/", func(c echo.Context) error { return get(c, ctx, pool) })
	e.POST("/", func(c echo.Context) error { return post(c, ctx, pool) })
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

func createTableIfNotExists(ctx context.Context, pool *pgxpool.Pool) error {
	query := fmt.Sprintf(`
				CREATE TABLE IF NOT EXISTS %s (
				    date timestamptz PRIMARY KEY DEFAULT CURRENT_TIMESTAMP,
    				name text NOT NULL,
    				message text NOT NULL
				)`, TableName)
	_, err := pool.Exec(ctx, query)
	if err != nil {
		return err
	}
	return nil
}

func get(c echo.Context, ctx context.Context, pool *pgxpool.Pool) error {
	err := createTableIfNotExists(ctx, pool)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	var posts []PostResp
	query := fmt.Sprintf("SELECT * FROM %s ORDER BY date DESC LIMIT 10", TableName)
	rows, err := pool.Query(ctx, query)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	defer rows.Close()
	posts, err = pgx.CollectRows(rows, pgx.RowToStructByName[PostResp])
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, posts)
}

func post(c echo.Context, ctx context.Context, pool *pgxpool.Pool) error {
	err := createTableIfNotExists(ctx, pool)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	p := new(PostReq)
	if err := c.Bind(p); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	query := fmt.Sprintf("INSERT INTO %s (name, message) VALUES ($1, $2)", TableName)
	if _, err := pool.Exec(ctx, query, p.Name, p.Message); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, p)
}
