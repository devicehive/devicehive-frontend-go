package notifications

type DeviceNotificationWrapperParameters struct {
	JsonString string `json:"jsonString,omitempty"`
}

type DeviceNotificationWrapper struct {
	Notification string                              `json:"notification"`
	Timestamp    string                              `json:"timestamp,omitempty"`
	Parameters   DeviceNotificationWrapperParameters `json:"parameters,omitempty"`
}
