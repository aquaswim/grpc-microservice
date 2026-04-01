package methodmetamap

import (
	"fmt"
	commonv1 "gaman-microservice/api-gateway/gen/common/v1"

	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func GetMethodMetaFromFileDesc(fds ...protoreflect.FileDescriptor) (MethodMetaMap, error) {
	out := MethodMetaMap{}
	// method: /user.v1.AuthService/Login
	//   /{service full name}/{method name}
	for _, fd := range fds {
		svcs := fd.Services()
		for i := 0; i < svcs.Len(); i++ {
			svc := svcs.Get(i)

			methods := svc.Methods()
			for j := 0; j < methods.Len(); j++ {
				method := methods.Get(j)
				key := fmt.Sprintf("/%s/%s", svc.FullName(), method.Name())

				needAuth, _ := proto.GetExtension(method.Options(), commonv1.E_NeedAuth).(bool)
				rateLimit, _ := proto.GetExtension(method.Options(), commonv1.E_RateLimit).(int32)

				out[key] = MethodMeta{
					NeedAuth:  needAuth,
					RateLimit: rateLimit,
				}
				log.Debug().Msgf("methodMap detected: %s - %+v", key, out[key])
			}
		}
	}

	return out, nil
}
