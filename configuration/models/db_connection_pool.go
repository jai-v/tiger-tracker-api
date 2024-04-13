package models

type DbConnectionPool struct {
	MaxOpenConnections             int `json:"maxOpenConnections"`
	MaxIdleConnections             int `json:"maxIdleConnections"`
	MaxConnectionLifetimeInMinutes int `json:"maxConnectionLifetimeInMinutes"`
}
