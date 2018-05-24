package notifications

import (
	"github.com/araddon/dateparse"
	"github.com/devicehive/devicehive-frontend-go/pg"
	"math/rand"
	"time"
)

type MessagesToDevice interface {
	GetDeviceId() string
}

type Notification struct {
	ID           int64
	DeviceId     string
	DeviceTypeId int8
	NetworkId    int8
	Timestamp    int64
	Notification string
	Parameters   DeviceNotificationWrapperParameters
}

func (n *Notification) GetDeviceId() string {
	return n.DeviceId
}
func New(device *pg.Device, wrapper *DeviceNotificationWrapper) *Notification {
	notif := &Notification{
		ID:           rand.Int63(),
		DeviceId:     device.DeviceId,
		DeviceTypeId: device.DeviceTypeId,
		NetworkId:    device.NetworkId,
		Notification: wrapper.Notification,
		Parameters:   wrapper.Parameters,
	}
	if wrapper.Timestamp != "" {
		t, err := dateparse.ParseAny(wrapper.Timestamp)
		if err != nil {
			notif.Timestamp = time.Now().Unix()
		}
		notif.Timestamp = int64(t.Unix())
	} else {
		notif.Timestamp = time.Now().Unix()
	}
	return notif
}
