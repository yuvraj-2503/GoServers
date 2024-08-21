package senders

import "otp-manager/common"

type OtpSender interface {
	Send(contact *common.Contact, otp uint64) error
}
