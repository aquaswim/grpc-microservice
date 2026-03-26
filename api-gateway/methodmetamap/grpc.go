package methodmetamap

import (
	"fmt"
	commonv1 "gaman-microservice/api-gateway/gen/common/v1"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func GetMethodMetaFromFileDesc(fd protoreflect.FileDescriptor) MethodMetaMap {
	out := MethodMetaMap{}
	// method: /user.v1.AuthService/Login
	//   /{service full name}/{method name}
	svcs := fd.Services()
	for i := 0; i < svcs.Len(); i++ {
		svc := svcs.Get(i)

		methods := svc.Methods()
		for j := 0; j < methods.Len(); j++ {
			method := methods.Get(j)

			needAuth, _ := proto.GetExtension(method.Options(), commonv1.E_NeedAuth).(bool)

			out[fmt.Sprintf("/%s/%s", svc.FullName(), method.Name())] = MethodMeta{
				NeedAuth: needAuth,
			}
		}
	}

	return out
}
