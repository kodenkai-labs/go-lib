package mongo

import "time"

// Config represents the configuration for a database connection.
type Config struct {
	// URI is the URI of the read-write database instance to connect to.
	URI string `mapstructure:"uri"`

	// Name is the name of the database.
	Name string `mapstructure:"name"`

	// MinPoolSize is the minimum number of connections in the connection pool.
	MinPoolSize uint64 `mapstructure:"min_pool_size"`

	// MaxConnIdleTime is the maximum amount of time that a connection can remain idle in the connection pool.
	MaxConnIdleTime time.Duration `mapstructure:"max_conn_idle_time"`
}