package methodmetamap

import (
	userv1 "gaman-microservice/api-gateway/gen/user/v1"
	"testing"
)

func TestGetMethodMetaFromFileDesc(t *testing.T) {
	mm := GetMethodMetaFromFileDesc(userv1.File_user_v1_auth_proto)

	t.Logf("methodmeta: %s", mm)
}
