package dsql

import (
	"context"
	"fmt"
	jidouConfig "jidou/internal/config"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dsql/auth"
	"github.com/aws/aws-sdk-go-v2/service/dsql"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	Host                 string
	Port                 string
	User                 string
	Password             string
	Database             string
	Region               string
	TokenRefreshInterval int
}

type Pool struct {
	Pool            *pgxpool.Pool
	config          Config
	ctx             context.Context
	cancelFunc      context.CancelFunc
	dsqlClient      *dsql.Client
	clusterEndpoint string
	mu              sync.Mutex
}

func NewDSQLClient(ctx context.Context, region string) (*dsql.Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS configuration: %v", err)
	}

	// Create a DSQL client using NewFromConfig
	dsqlClient := dsql.NewFromConfig(cfg)

	return dsqlClient, nil
}

func GenerateDbConnectAuthToken(ctx context.Context, clusterEndpoint, region, user string) (string, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return "", err
	}

	if user == "admin" {
		token, err := auth.GenerateDBConnectAdminAuthToken(ctx, clusterEndpoint, region, cfg.Credentials)
		if err != nil {
			return "", err
		}

		return token, nil
	}

	token, err := auth.GenerateDbConnectAuthToken(ctx, clusterEndpoint, region, cfg.Credentials)
	if err != nil {
		return "", err
	}

	return token, nil
}

func CreateConnectionURL(dbConfig Config) string {
	var sb strings.Builder
	sb.WriteString("postgres://")
	sb.WriteString(dbConfig.User)
	sb.WriteString("@")
	sb.WriteString(dbConfig.Host)
	sb.WriteString(":")
	sb.WriteString(dbConfig.Port)
	sb.WriteString("/")
	sb.WriteString(dbConfig.Database)
	sb.WriteString("?sslmode=verify-full")
	sb.WriteString("&sslnegotiation=direct")
	url := sb.String()
	return url
}

func setPoolSettings(poolConfig *pgxpool.Config) {
	// Configure pool settings
	poolConfig.MaxConns = 10
	poolConfig.MinConns = 2
	poolConfig.MaxConnLifetime = 1 * time.Hour
	poolConfig.MaxConnIdleTime = 30 * time.Minute
	poolConfig.HealthCheckPeriod = 1 * time.Minute
}

func (p *Pool) refreshToken() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Generate new token
	token, err := GenerateDbConnectAuthToken(p.ctx, p.clusterEndpoint, p.config.Region, p.config.User)
	if err != nil {
		return fmt.Errorf("failed to refresh auth token: %v", err)
	}

	// Update all connections in the pool with the new token
	conns := p.Pool.Stat().TotalConns()

	// Reset the pool to force new connections with the updated token
	p.Pool.Reset()

	// Create a new connection config with the updated token
	url := CreateConnectionURL(p.config)

	poolConfig, err := pgxpool.ParseConfig(url)
	if err != nil {
		return fmt.Errorf("unable to parse pool config during token refresh: %v", err)
	}

	// Update the password with the new token
	poolConfig.ConnConfig.Password = token

	setPoolSettings(poolConfig)

	// Create a new pool with the updated token
	newPool, err := pgxpool.NewWithConfig(p.ctx, poolConfig)
	if err != nil {
		return fmt.Errorf("unable to create new connection pool during token refresh: %v", err)
	}

	// Replace the old pool with the new one
	oldPool := p.Pool
	p.Pool = newPool

	// Close the old pool
	oldPool.Close()

	fmt.Printf("Successfully refreshed token and updated %d connections\n", conns)
	return nil
}

func (p *Pool) refreshTokenPeriodically() {
	// Calculate refresh interval (75% of token lifetime to refresh before expiration)
	refreshInterval := time.Duration(p.config.TokenRefreshInterval) * time.Second * 3 / 4

	ticker := time.NewTicker(refreshInterval)
	defer ticker.Stop()

	for {
		select {
		case <-p.ctx.Done():
			return
		case <-ticker.C:
			if err := p.refreshToken(); err != nil {
				fmt.Fprintf(os.Stderr, "Error refreshing token: %v\n", err)
			}
		}
	}
}

func NewPool(ctx context.Context, configuration *jidouConfig.Configuration) (*Pool, error) {
	/*jidouCfg, err := jidouConfig.LoadConfiguration()
	if err != nil {
		log.Fatalf("unable to load Jidou config, %v", err)
	}*/
	poolCtx, cancel := context.WithCancel(ctx)
	region := configuration.Region
	client, err := NewDSQLClient(poolCtx, region)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("failed to create DSQL client: %v", err)
	}
	dbConfig := Config{
		Host:                 configuration.ClusterEndpoint,
		Port:                 configuration.Port,
		User:                 configuration.User,
		Password:             "",
		Database:             configuration.DatabaseName,
		Region:               configuration.Region,
		TokenRefreshInterval: configuration.TokenRefreshInterval,
	}
	token, err := GenerateDbConnectAuthToken(poolCtx, dbConfig.Host, dbConfig.Region, dbConfig.User)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("failed to generate auth token: %v", err)
	}
	url := CreateConnectionURL(dbConfig)
	poolConfig, err := pgxpool.ParseConfig(url)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("unable to parse pool config: %v", err)
	}
	poolConfig.ConnConfig.Password = token
	setPoolSettings(poolConfig)
	pgxPool, err := pgxpool.NewWithConfig(poolCtx, poolConfig)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("failed to create connection pool: %v", err)
	}
	pool := &Pool{
		Pool:            pgxPool,
		config:          dbConfig,
		ctx:             poolCtx,
		cancelFunc:      cancel,
		dsqlClient:      client,
		clusterEndpoint: dbConfig.Host,
	}
	go pool.refreshTokenPeriodically()
	return pool, nil
}

// Close closes the connection pool and cancels the refresh goroutine
func (p *Pool) Close() {
	p.cancelFunc()
	p.Pool.Close()
}

// GetConnectionID returns a unique identifier for a connection in the pool
func (p *Pool) GetConnectionID(ctx context.Context) (string, error) {
	// Retrieve the session variable to confirm it was set
	var connID string
	err := p.Pool.QueryRow(ctx, "select sys.current_session_id();").Scan(&connID)
	if err != nil {
		return "", fmt.Errorf("failed to get connection ID: %v", err)
	}
	return connID, nil
}

// DemonstrateConnectionRefresh shows that connections before and after refresh are different
func (p *Pool) DemonstrateConnectionRefresh(ctx context.Context) error {
	// Get connection ID before refresh
	connIDBefore, err := p.GetConnectionID(ctx)
	if err != nil {
		return err
	}
	fmt.Printf("Connection ID before refresh: %s\n", connIDBefore)

	// Refresh token
	err = p.refreshToken()
	if err != nil {
		return fmt.Errorf("failed to refresh token: %v", err)
	}

	// Get connection ID after refresh
	connIDAfter, err := p.GetConnectionID(ctx)
	if err != nil {
		return err
	}
	fmt.Printf("Connection ID after refresh: %s\n", connIDAfter)

	// Verify that the connections are different
	if connIDBefore == connIDAfter {
		return fmt.Errorf("connection IDs before and after refresh are the same: %s", connIDBefore)
	}

	fmt.Println("Successfully verified that connections before and after refresh are different")
	return nil
}

// GetConnectionPool creates a new connection pool with token refresh capability
func getConnectionPool(ctx context.Context, clusterEndpoint string, region string, configuration *jidouConfig.Configuration) (*pgxpool.Pool, error) {
	pool, err := NewPool(ctx, configuration)
	if err != nil {
		return nil, err
	}

	// Return just the pgxpool.Pool to maintain compatibility with existing code
	return pool.Pool, nil
}
