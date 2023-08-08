package oracle

import (
	"fmt"
	"myapp/config"
	"time"

	go_oracle "gitlabnew.nextcontact.ru/r.alfimov/go-oracle"
)

func New(cfg config.Oracle) (*go_oracle.Oracle, error) {

	ora, err := go_oracle.New(&go_oracle.Config{
		Host:     cfg.Host,
		Port:     cfg.Port,
		Service:  cfg.Service,
		User:     cfg.User,
		Password: cfg.Password,
		NLSQueries: []string{
			"ALTER SESSION SET NLS_DATE_FORMAT = 'dd.mm.yyyy'",
			"ALTER SESSION SET NLS_TIMESTAMP_FORMAT = 'dd.mm.yyyy hh24:mi:ss'",
			"ALTER SESSION SET NLS_TIMESTAMP_TZ_FORMAT = 'dd.mm.yyyy hh24:mi:ss tzr'",
			"ALTER SESSION SET NLS_NUMERIC_CHARACTERS = '. '",
			"ALTER SESSION SET NLS_SORT = 'RUSSIAN'",
			"ALTER SESSION SET NLS_LANGUAGE = 'RUSSIAN'",
			"ALTER SESSION SET NLS_COMP = 'BINARY'",
		},
		CheckConnection: &go_oracle.CheckConnectionConfig{
			ReconnectTryCount:       5,
			ReconnectTryInterval:    time.Second,
			CheckConnectionInterval: 250 * time.Millisecond,
		},
		CustomParams: map[string]string{
			//"TRACE FILE": "trace.log",
			//"DBA PRIVILEGE": "SYSDBA",
		},
		Debug: true,
	})
	if err != nil {
		return nil, fmt.Errorf("oracle - NewOracle - Connect: %w", err)
	}
	return ora, nil
}
