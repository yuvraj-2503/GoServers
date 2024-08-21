package sms

import (
	"context"
	"encoding/json"
	"net/http"
	"otp-manager/errors"
	"strconv"
	"strings"
)

type Recipient struct {
	Mobiles string `json:"mobiles"`
	Otp     string `json:"otp"`
}

type Msg91sms struct {
	TemplateId string      `json:"template_id"`
	ShortUrl   string      `json:"short_url"`
	Recipients []Recipient `json:"recipients"`
}

type MSG91Sender struct {
	baseUrl    string
	authKey    string
	templateId string
}

func NewMSG91Sender(baseUrl string, authKey string, templateId string) *MSG91Sender {
	return &MSG91Sender{baseUrl, authKey, templateId}
}

func (m *MSG91Sender) Send(ctx *context.Context, sms *Sms) error {
	jsonPayload := m.getPayload(sms)
	payload := strings.NewReader(jsonPayload)
	req, _ := http.NewRequest("POST", m.baseUrl, payload)
	m.setHeaders(req)
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return &errors.NetworkConnectionError{}
	}
	if response.StatusCode != 200 {
		return &errors.OtpSendingError{Message: response.Status}
	}
	defer response.Body.Close()
	return nil
}

func (m *MSG91Sender) setHeaders(req *http.Request) {
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("authkey", m.authKey)
}

func (m *MSG91Sender) getPayload(sms *Sms) string {
	msg := Msg91sms{
		TemplateId: m.templateId,
		ShortUrl:   "0",
		Recipients: getRecipients(sms),
	}

	jsonBytes, _ := json.Marshal(msg)
	return string(jsonBytes)
}

func getRecipients(sms *Sms) []Recipient {
	return []Recipient{{
		Mobiles: sms.Mobile,
		Otp:     strconv.Itoa(int(sms.Otp)),
	}}
}
