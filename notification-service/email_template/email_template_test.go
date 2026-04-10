package emailTemplate

import (
	"gaman-microservice/notification-service/internal/entity"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func writeTemplate(name string, data any) error {
	txt, err := RenderTemplate(name, data)
	if err != nil {
		return err
	}
	return os.WriteFile(name+".gen.html", []byte(txt), 0644)
}

func TestGenerateAllEmail(t *testing.T) {
	tc := []struct {
		TemplateName string
		Data         any
	}{
		{"forgot_password.gohtml", entity.ForgotPasswordNotificationData{
			Token:     "test-token-lala",
			Username:  "test-username",
			Email:     "test@example.com",
			ExpiredAt: time.Now(),
		}},
		{"reset_password_done.gohtml", entity.ResetPasswordSuccess{
			UserId:   "user-id-lala",
			Username: "test-username",
			Email:    "test@example.com",
		}},
	}

	for _, tdata := range tc {
		t.Run("Template: "+tdata.TemplateName, func(t *testing.T) {
			assert.NoErrorf(t, writeTemplate(tdata.TemplateName, tdata.Data), "failed for %s", tdata.TemplateName)
		})
	}
}
