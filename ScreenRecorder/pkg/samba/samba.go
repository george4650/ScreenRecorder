package samba

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"myapp/config"
	"net"

	"github.com/hirochachacha/go-smb2"
)

func New(cfg config.Samba) (*smb2.Session, error) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", cfg.Host, cfg.Port))
	if err != nil {
		return nil, err
	}
	//defer conn.Close()

	d := &smb2.Dialer{
		Initiator: &smb2.NTLMInitiator{
			User:     cfg.User,
			Password: cfg.Password,
			Domain:   cfg.Domain,
		},
	}

	s, err := d.Dial(conn)
	if err != nil {
		return nil, err
	}

	log.Print("Подключение samba успешно")

	return s, nil
}
