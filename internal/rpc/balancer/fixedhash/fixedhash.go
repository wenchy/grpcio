// Package fixedhash defines a fixedhash balancer.
package fixedhash

import (
	"fmt"
	"sync"

	"github.com/wenchy/grpcio/internal/uidutil"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/metadata"
)

// Name is the name of fixed_hash balancer.
const Name = "fixed_hash"

var logger = grpclog.Component("fixedhash")

// newBuilder creates a new fixedhash balancer builder.
func newBuilder() balancer.Builder {
	return base.NewBalancerBuilder(Name, &fhPickerBuilder{}, base.Config{HealthCheck: true})
}

func init() {
	balancer.Register(newBuilder())
}

type fhPickerBuilder struct{}

func (*fhPickerBuilder) Build(info base.PickerBuildInfo) balancer.Picker {
	logger.Infof("fixedhash: newPicker called with info: %v", info)
	if len(info.ReadySCs) == 0 {
		return base.NewErrPicker(balancer.ErrNoSubConnAvailable)
	}
	var scs []balancer.SubConn
	nodeMap := map[string]balancer.SubConn{}
	for sc, info := range info.ReadySCs {
		scs = append(scs, sc)
		id := info.Address.Attributes.Value("id")
		if id == nil {
			panic("node is not provided when load balancer build")
		}
		idstr, ok := id.(string)
		if !ok {
			panic("cannot parse node id as string when load balancer build")
		}
		fmt.Printf("register node id: %s\n", idstr)
		nodeMap[idstr] = sc
	}
	return &fhPicker{
		subConns: scs,
		nodeMap:  nodeMap,
	}
}

type fhPicker struct {
	// subConns is the snapshot of the fixedhash balancer when this picker was
	// created. The slice is immutable. Each Get() will do a fixed hash
	// selection from it and return the selected SubConn.
	subConns []balancer.SubConn
	nodeMap  map[string]balancer.SubConn

	mu sync.Mutex
}

func (p *fhPicker) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	pr := balancer.PickResult{}
	md, ok := metadata.FromOutgoingContext(info.Ctx)
	if !ok {
		return pr, fmt.Errorf("error: missing metadata when load balancer pick")
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	nodeIDStrs, ok := md["x-node-id"]
	if ok {
		// route by node id
		if len(nodeIDStrs) < 1 {
			return pr, fmt.Errorf("error: missing metadata key 'x-node-id' when load balancer pick")
		}
		id := nodeIDStrs[0]
		sc, ok := p.nodeMap[id]
		if !ok {
			return pr, fmt.Errorf("error: node id(%s) not found when load balancer pick", id)
		}
		pr.SubConn = sc
	} else {
		// route by uid hash
		uidStrs, ok := md["x-packet-uid"]
		if !ok || len(uidStrs) < 1 {
			return pr, fmt.Errorf("error: missing metadata key 'x-packet-uid' when load balancer pick")
		}
		uid, err := uidutil.Parse(uidStrs[0])
		if err != nil {
			return pr, fmt.Errorf("error: parse x-packet-uid failed when load balancer pick")
		}
		nodenum := uint64(len(p.subConns))
		if nodenum == 0 {
			return pr, fmt.Errorf("error: node num is 0 when load balancer pick")
		}
		index := uid.Value() % nodenum
		pr.SubConn = p.subConns[index]
	}
	return pr, nil
}
