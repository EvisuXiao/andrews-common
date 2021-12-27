package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/EvisuXiao/andrews-common/config"
	"github.com/EvisuXiao/andrews-common/exception"
	"github.com/EvisuXiao/andrews-common/logging"
	"github.com/EvisuXiao/andrews-common/utils"
)

const (
	DriverStandalone = "standalone"
	DriverCluster    = "cluster"
	DriverSentinel   = "sentinel"
)

type Client struct {
	Name   string
	Ctx    context.Context
	redis  redis.Cmdable
	prefix string
}

var (
	configs config.Redises
	clients []*Client
)

func Init() {
	configs = config.GetRedisConfigs()
	for _, c := range clients {
		c.setup()
	}
}

func RegisterRedis(c *Client) {
	clients = append(clients, c)
}

func (c *Client) setup() {
	cfg, ok := configs[c.Name]
	if !ok {
		logging.Fatal("Init: redis(%s) connection name not found", c.Name)
	}
	switch cfg.Driver {
	case DriverStandalone:
		c.connClient(cfg)
	case DriverCluster:
		c.connClusterClient(cfg)
	case DriverSentinel:
		c.connSentinelClient(cfg)
	default:
		logging.Fatal("Init: unknown redis driver: %s", cfg.Driver)
	}
	pong, err := c.redis.Ping(c.Ctx).Result()
	if utils.HasErr(err) {
		logging.Fatal("Init: redis(%s) connection failed: %+v", c.Name, err)
	}
	if utils.IsEmpty(c.Ctx) {
		c.Ctx = context.Background()
	}
	logging.Info("Redis init(%s) successfully: %s", c.Name, pong)
}

func (c *Client) connClient(cfg *config.Redis) {
	c.redis = redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password:     cfg.Password,
		DB:           cfg.Database,
		ReadTimeout:  cfg.Timeout.Read,
		WriteTimeout: cfg.Timeout.Write,
	})
	c.prefix = cfg.Prefix
}

func (c *Client) connSentinelClient(cfg *config.Redis) {
	c.redis = redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:    cfg.MasterName,
		SentinelAddrs: cfg.Hosts,
		Password:      cfg.Password,
		DB:            cfg.Database,
		ReadTimeout:   cfg.Timeout.Read,
		WriteTimeout:  cfg.Timeout.Write,
	})
	c.prefix = cfg.Prefix
}

func (c *Client) connClusterClient(cfg *config.Redis) {
	c.redis = redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:        cfg.Hosts,
		Password:     cfg.Password,
		ReadTimeout:  cfg.Timeout.Read,
		WriteTimeout: cfg.Timeout.Write,
	})
	c.prefix = cfg.Prefix
}

func (c *Client) GetClient() redis.Cmdable {
	return c.redis
}

func (c *Client) Set(key string, value interface{}, expired time.Duration) error {
	if !utils.IsSimpleValue(value) {
		value, _ = json.Marshal(value)
	}
	return exception.DbErrWrapper(c.GetClient().Set(c.Ctx, c.getCacheKey(key), value, expired).Err())
}

func (c *Client) SetNX(key string, value interface{}, expire time.Duration) (bool, error) {
	if !utils.IsSimpleValue(value) {
		value, _ = json.Marshal(value)
	}
	result, err := c.GetClient().SetNX(c.Ctx, c.getCacheKey(key), value, expire).Result()
	return result, exception.DbErrWrapper(err)
}

func (c *Client) GetString(key string) (string, error) {
	result, err := c.getCmd(key).Result()
	return result, exception.DbErrWrapper(err)
}

func (c *Client) GetInt(key string) (int, error) {
	result, err := c.getCmd(key).Int()
	return result, exception.DbErrWrapper(err)
}

func (c *Client) GetTime(key string) (time.Time, error) {
	result, err := c.getCmd(key).Time()
	return result, exception.DbErrWrapper(err)
}

func (c *Client) GetScan(key string, out interface{}) error {
	b, err := c.getCmd(key).Bytes()
	if utils.HasErr(err) {
		return exception.DbErrWrapper(err)
	}
	return json.Unmarshal(b, out)
}

func (c *Client) Exists(key string) bool {
	return c.GetClient().Exists(c.Ctx, c.getCacheKey(key)).Val() > 0
}

func (c *Client) Delete(key string) error {
	return exception.DbErrWrapper(c.GetClient().Del(c.Ctx, c.getCacheKey(key)).Err())
}

func (c *Client) getCmd(key string) *redis.StringCmd {
	return c.GetClient().Get(c.Ctx, c.getCacheKey(key))
}

func (c *Client) getCacheKey(key string) string {
	return c.prefix + key
}

func IsNilErr(err error) bool {
	return err.Error() == redis.Nil.Error()
}
