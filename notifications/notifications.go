package notifications

import (
	"github.com/araddon/dateparse"
	"github.com/devicehive/devicehive-frontend-go/pg"
	"github.com/devicehive/devicehive-frontend-go/request"
	"math/rand"
	"time"
)

type MessagesToDevice interface {
	GetDeviceId() string
}

type DeviceNotification struct {
	Id           int64  `json:"id"`
	DeviceId     string `json:"deviceId"`
	DeviceTypeId int8   `json:"deviceTypeId"`
	NetworkId    int8   `json:"networkId"`
	Timestamp    int64  `json:"timestamp"`
	Notification string `json:"notification"`

	Parameters DeviceNotificationWrapperParameters `json:"parameters"`
}
type Notification struct {
	DeviceNotification DeviceNotification `json:"deviceNotification"`
	Action             request.Action     `json:"a"`
}

func (n *Notification) GetDeviceId() string {
	return n.DeviceNotification.DeviceId
}
func New(device *pg.Device, wrapper *DeviceNotificationWrapper) *Notification {

	n := DeviceNotification{
		Id:           rand.Int63(),
		DeviceId:     device.DeviceId,
		DeviceTypeId: device.DeviceTypeId,
		NetworkId:    device.NetworkId,
		Notification: wrapper.Notification,
		Parameters:   wrapper.Parameters,
	}

	notif := &Notification{
		DeviceNotification: n,
		Action:             request.NOTIFICATION_INSERT_REQUEST,
	}

	if wrapper.Timestamp != "" {
		t, err := dateparse.ParseAny(wrapper.Timestamp)
		if err != nil {
			notif.DeviceNotification.Timestamp = time.Now().Unix()
		}
		notif.DeviceNotification.Timestamp = int64(t.Unix())
	} else {
		notif.DeviceNotification.Timestamp = time.Now().Unix()
	}
	return notif
}
