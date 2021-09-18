package cmd

import (
	"fmt"
	"strings"

	"github.com/wenchy/grpcio/internal/corepb"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
)

// GenMethodMapping generates cmd -> grpc method mappings.
func GenMethodMapping(protoPackageName string) map[corepb.Cmd]string {
	cmdMap := make(map[corepb.Cmd]string)
	protoPackage := protoreflect.FullName(protoPackageName)
	protoregistry.GlobalFiles.RangeFilesByPackage(protoPackage, func(fd protoreflect.FileDescriptor) bool {
		// fmt.Printf("filepath: %s\n", fd.Path())
		services := fd.Services()
		for i := 0; i < services.Len(); i++ {
			sd := services.Get(i)
			// fmt.Printf("service: %s\n", sd.FullName())
			methods := sd.Methods()
			for j := 0; j < methods.Len(); j++ {
				md := methods.Get(j)
				// fmt.Printf("method: %s\n", md.Name())
				opts := md.Options().(*descriptorpb.MethodOptions)
				cmd := proto.GetExtension(opts, corepb.E_Cmd).(corepb.Cmd)
				if cmd == 0 {
					fmt.Printf("cmd not specified of rpc: %v\n", md.Name())
					continue
				}
				// fmt.Printf("cmd: %v, service: %s, method: %s\n", cmd, sd.FullName(), md.Name())
				rpc := fmt.Sprintf("/%s/%s", sd.FullName(), md.Name())
				cmdMap[cmd] = rpc
				fmt.Printf("%v: %v\n", cmd, rpc)
			}
		}
		return true
	})
	return cmdMap
}

// API is the interface for clien and server to communicate.
type API struct {
	CmdName  string `json:"cmdName"`
	CmdID    int    `json:"cmdID"`
	Request  string `json:"request"`
	Response string `json:"response"`
}

// GenAPIList generates api list.
func GenAPIList(protoPackageName string, onlyClient bool) []API {
	apis := []API{}
	protoPackage := protoreflect.FullName(protoPackageName)
	protoregistry.GlobalFiles.RangeFilesByPackage(protoPackage, func(fd protoreflect.FileDescriptor) bool {
		// fmt.Printf("filepath: %s\n", fd.Path())
		services := fd.Services()
		for i := 0; i < services.Len(); i++ {
			sd := services.Get(i)
			// fmt.Printf("service: %s\n", sd.FullName())
			if onlyClient && !strings.HasPrefix(string(sd.FullName()), "core.Client") {
				continue
			}
			methods := sd.Methods()
			for j := 0; j < methods.Len(); j++ {
				md := methods.Get(j)
				// fmt.Printf("method: %s\n", md.Name())
				opts := md.Options().(*descriptorpb.MethodOptions)
				cmd := proto.GetExtension(opts, corepb.E_Cmd).(corepb.Cmd)
				if cmd == 0 {
					fmt.Printf("cmd not specified of rpc: %v\n", md.Name())
					continue
				}
				// fmt.Printf("cmd: %v, service: %s, method: %s\n", cmd, sd.FullName(), md.FullName())
				// rpc := fmt.Sprintf("/%s/%s", sd.FullName(), md.Name())
				oneAPI := API{
					CmdName:  "core.Cmd." + cmd.String(),
					CmdID:    int(cmd),
					Request:  string(md.Input().FullName()),
					Response: string(md.Output().FullName()),
				}
				apis = append(apis, oneAPI)
				// fmt.Printf("%v: %v\n", cmd, oneAPI)
			}
		}
		return true
	})
	return apis
}
