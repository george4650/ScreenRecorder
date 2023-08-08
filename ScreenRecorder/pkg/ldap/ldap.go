package ldap

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"myapp/config"

	"github.com/go-ldap/ldap/v3"
)

func ConnectToServerLDAP(cfg config.Ldap) (*ldap.Conn, error) {

	var (
		l   *ldap.Conn
		err error
	)

	for _, server := range cfg.Servers {

		l, err = ldap.DialURL(server)

		if err != nil {
			log.Printf("Не удалось подключиться к серверу: %s, Ошибка: %s\n", server, err)
			continue
		} else {
			break
		}
	}

	if l == nil {
		return nil, fmt.Errorf("ldap - ConnectToServer - Connect: %s", "Не удалось подключиться ни к одному из адресов")
	}

	return l, nil
}
