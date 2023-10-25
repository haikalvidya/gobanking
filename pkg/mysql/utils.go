package mysql

import (
	"database/sql"
	"fmt"
)

func GetMysqlStats(db *sql.DB) string {
	stats := db.Stats()
	return fmt.Sprintf("maxOpenConnections: %d, openConnections: %d, inUse: %d, idle: %d, waitCount: %d, waitDuration: %s, maxIdleClosed: %d, maxLifetimeClosed: %d",
		stats.MaxOpenConnections,
		stats.OpenConnections,
		stats.InUse,
		stats.Idle,
		stats.WaitCount,
		stats.WaitDuration,
		stats.MaxIdleClosed,
		stats.MaxLifetimeClosed,
	)
}
