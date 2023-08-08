package app

import (
	"fmt"
	"myapp/config"
	"myapp/internal/grpc/server"
	"myapp/internal/usecase"
	"myapp/internal/usecase/repository"
	"myapp/pkg/ldap"
	"myapp/pkg/oracle"
	"myapp/pkg/postgres"
	"myapp/pkg/samba"
	"net"

	proto "gitlabnew.nextcontact.ru/giorgiy_03/screenrecordercontracts"

	"google.golang.org/grpc"

	"github.com/rs/zerolog/log"
)

func Run(cfg *config.Config) {
	// Ldap
	ldapConn, err := ldap.ConnectToServerLDAP(cfg.Ldap)
	if err != nil {
		log.Error().Err(err).Msg("app - Run - ldap.ConnectToServerLDAP")
	}
	defer ldapConn.Close()

	// Samba
	smbSession, err := samba.New(cfg.Samba)
	if err != nil {
		log.Error().Err(err).Msg("app - Run - samba.New")
	}
	defer smbSession.Logoff()

	// Repository
	ora, err := oracle.New(cfg.Oracle)
	if err != nil {
		log.Error().Err(err).Msg("app - Run - oracle.New")
	}
	defer ora.Close()

	pg, err := postgres.New(cfg.Postgres)
	if err != nil {
		log.Fatal().Err(err).Msg("app - Run - Postgres.New")
	}
	defer pg.Close()

	callDetailsPostgresRepo := repository.NewScreenRecorderPostgres(pg)
	callDetailsOracleRepo := repository.NewScreenRecorderOracle(ora)

	// Use case
	authUseCases := usecase.NewAuthUseCases(ldapConn, cfg.Auth.Key)
	screenRecorderUseCases := usecase.NewScreenRecorderCases(cfg.Links.Audio, smbSession, callDetailsPostgresRepo, callDetailsOracleRepo)
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Grpc.Port))
	if err != nil {
		log.Fatal().Msgf("StartGRPC - Run - New: %w", err)
	}
	defer listener.Close()

	// Start GRPC Server
	grpcServer := grpc.NewServer()

	proto.RegisterScreenRecorderServiceServer(grpcServer, server.NewServer(authUseCases, screenRecorderUseCases))

	log.Printf("GRPC server listening at %v", listener.Addr())

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal().Msgf("GrpcServer - Run - New: %w", err)
	}
}
