package mailpit

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"gaman-microservice/notification-service/internal/client/email"
	"net/http"

	"github.com/rs/zerolog/log"
)

const emailFromName = "Gaman"
const emailFromAddress = "noreply@gaman.com"

type client struct {
	httpClient *http.Client
	baseUrl    string
}

func New(
	httpClient *http.Client,
	baseUrl string,
) email.Client {
	return &client{
		httpClient: httpClient,
		baseUrl:    baseUrl,
	}
}

func (c client) SendEmail(ctx context.Context, req *email.SendEmailReq) (*email.SendEmailRes, error) {
	reqBody := &SendEmailRequest{
		From: Recipient{
			Email: emailFromAddress,
			Name:  emailFromName,
		},
		HTML:    req.BodyHtml,
		Subject: req.Subject,
		Tags:    nil,
		Text:    req.BodyText,
		To: []Recipient{
			{
				Email: req.ToEmail,
				Name:  req.ToName,
			},
		},
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseUrl+"/api/v1/send", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")

	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer func() {
		err := response.Body.Close()
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("failed to close response body")
		}
	}()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	resp := &SendEmailResponse{}
	if err := json.NewDecoder(response.Body).Decode(resp); err != nil {
		return nil, err
	}
	return &email.SendEmailRes{EmailId: resp.ID}, nil
}
