package app

import (
	"context"
	"gobanking/pkg/mysql"
	"time"

	"github.com/avast/retry-go"
	"github.com/pkg/errors"
)

func (a *app) connectMysql(ctx context.Context) error {

	retryOptions := []retry.Option{
		retry.Attempts(a.cfg.Timeouts.MysqlInitRetryCount),
		retry.Delay(time.Duration(a.cfg.Timeouts.MysqlInitMilliseconds) * time.Millisecond),
		retry.DelayType(retry.BackOffDelay),
		retry.LastErrorOnly(true),
		retry.Context(ctx),
		retry.OnRetry(func(n uint, err error) {
			a.log.Errorf("retry connect mysql err: %v", err)
		}),
	}

	return retry.Do(func() error {
		mysqlConn, err := mysql.NewMysqlConn(a.cfg.Mysql)
		if err != nil {
			return errors.Wrap(err, "mysql.NewMysqlConn")
		}
		a.mysqlConn = mysqlConn
		statsSql, _ := a.mysqlConn.DB()
		a.log.Infof("(mysql connected) poolStat: %s", mysql.GetMysqlStats(statsSql))
		return nil
	}, retryOptions...)
}

func (a *app) migrateMysql(ctx context.Context) error {
	return a.mysqlConn.AutoMigrate()
}
