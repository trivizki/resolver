package resolving

import (
	"net"
	"context"
)

// Resolver is responsible to resovel dns domain's ip.
type Resolver struct{
	server *net.Resolver
}


func NewResolver() *Resolver{
	r := &net.Resolver{
			    PreferGo: true,
				Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
				d := net.Dialer{}
				return d.DialContext(ctx, "udp", net.JoinHostPort("8.8.8.8", "53"))
		},
	}
	return &Resolver{
		server : r,
	}
}

//Resolve find the ips assocciate with the given domain name.
func (r *Resolver) Resolve(domain string) ([]net.IP, error){
	var results []net.IP
	ips, err := r.server.LookupIPAddr(context.Background(), domain)
	for _, ip := range ips{
		results = append(results, ip.IP)
	}
	return results, err
}
