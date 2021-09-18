package resolver

import (
	"fmt"
	"log"

	"github.com/wenchy/grpcio/internal/envconf"
	"google.golang.org/grpc/attributes"
	"google.golang.org/grpc/resolver"
)

// Following is an celestial name resolver. It includes a
// ResolverBuilder(https://godoc.org/google.golang.org/grpc/resolver#Builder)
// and a Resolver(https://godoc.org/google.golang.org/grpc/resolver#Resolver).
//
// A ResolverBuilder is registered for a scheme (in this celestial, "celestial" is
// the scheme). When a ClientConn is created for this scheme, the
// ResolverBuilder will be picked to build a Resolver. Note that a new Resolver
// is built for each ClientConn. The Resolver will watch the updates for the
// target, and send updates to the ClientConn.

// celestialResolverBuilder is a
// ResolverBuilder(https://godoc.org/google.golang.org/grpc/resolver#Builder).

const (
	celestialScheme = "pi"
)

type celestialResolverBuilder struct{}

func (*celestialResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	if envconf.Conf == nil {
		// please confirm: envconf has already been inited
		log.Fatalln("envconf not inited")
	}
	addrsStore := map[string][]string{}
	for name, nodes := range envconf.Conf.Nodes {
		addrsStore[name] = make([]string, len(nodes))
		for i, node := range nodes {
			addrsStore[name][i] = node.GRPCAddress
		}
	}
	fmt.Printf("addrsStore: %+v\n", addrsStore)
	r := &celestialResolver{
		target:     target,
		cc:         cc,
		addrsStore: addrsStore,
	}
	mynodes, ok := envconf.Conf.Nodes[r.target.Endpoint]
	if !ok {
		log.Fatalln("server not found: ", r.target.Endpoint)
		return nil, fmt.Errorf("server not found: %s", r.target.Endpoint)
	}
	addrs := make([]resolver.Address, len(mynodes))
	for i, node := range mynodes {
		addrs[i] = resolver.Address{
			Addr:       node.GRPCAddress,
			Attributes: attributes.New("id", node.ID),
		}
	}
	r.cc.UpdateState(resolver.State{Addresses: addrs})

	return r, nil
}
func (*celestialResolverBuilder) Scheme() string { return celestialScheme }

// celestialResolver is a
// Resolver(https://godoc.org/google.golang.org/grpc/resolver#Resolver).
type celestialResolver struct {
	target     resolver.Target
	cc         resolver.ClientConn
	addrsStore map[string][]string
}

func (*celestialResolver) ResolveNow(o resolver.ResolveNowOptions) {}
func (*celestialResolver) Close()                                  {}

func Register() {
	// Register the celestial ResolverBuilder. This is usually done in a package's
	// init() function.
	resolver.Register(&celestialResolverBuilder{})
}

func Scheme() string {
	return celestialScheme
}
