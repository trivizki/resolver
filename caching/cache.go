package caching

import (
	"net"
	"time"
)

type Cacher interface{
	// Initialize the Cache object. e.g. connection to database.
	InitializeCache() error

	// Get the assocciated ips of the given domain name.
	GetIPS(domain string) ([]net.IP, error) 
	
	// Update the assocciated ips with the given domain.
	// the expiration argument tells the cache for how long to keep this information.
	UpdateDomain(domain string, ips []net.IP, expiration time.Duration) error 
}

