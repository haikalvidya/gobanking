package app

import (
	"context"
	"gobanking/internal/bank/models"
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
	return a.mysqlConn.AutoMigrate(
		&models.Account{},
		&models.Transaction{},
		&models.BalanceLedger{},
		&models.Currency{},
	)
}

// seeding currencies
func (a *app) seedCurrencies(ctx context.Context) error {
	var currencies = []models.Currency{
		{
			Code: "USD",
			Name: "US Dollar",
		},
		{
			Code: "EUR",
			Name: "Euro",
		},
		{
			Code: "GBP",
			Name: "British Pound",
		},
		{
			Code: "JPY",
			Name: "Japanese Yen",
		},
		{
			Code: "RUB",
			Name: "Russian Ruble",
		},
		{
			Code: "IDR",
			Name: "Indonesian Rupiah",
		},
	}

	for _, currency := range currencies {
		// create if name not exist
		if err := a.mysqlConn.FirstOrCreate(&currency, models.Currency{Name: currency.Name}).Error; err != nil {
			return errors.Wrap(err, "mysql.FirstOrCreate")
		}
	}

	return nil
}
