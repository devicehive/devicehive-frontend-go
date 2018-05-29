package pg

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
)

type PgConnection struct {
	conn *sql.DB
}

func (p *PgConnection) Close() {
	p.conn.Close()
}

func (p *PgConnection) GetDeviceById(id string) (Device, error) {
	device := Device{}
	err := p.conn.QueryRow("SELECT network_id, device_type_id FROM device WHERE device_id=$1", id).Scan(&device.NetworkId, &device.DeviceTypeId)
	if err != nil {
		return device, err
	}
	device.DeviceId = id
	return device, nil

}

func New(connStr string) *PgConnection {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	return &PgConnection{conn: db}
}
