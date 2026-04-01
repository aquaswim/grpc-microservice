package methodmetamap

import (
	userv1 "gaman-microservice/api-gateway/gen/user/v1"
	"testing"
)

func TestGetMethodMetaFromFileDesc(t *testing.T) {
	mm, err := GetMethodMetaFromFileDesc(userv1.File_user_v1_auth_proto)
	if err != nil {
		t.Fatalf("failed to get method meta: %v", err)
	}

	t.Logf("methodmeta: %s", mm)
}
