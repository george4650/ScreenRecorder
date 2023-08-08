package main

import (
	"context"
	"fmt"
	"myapp/config"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	proto "gitlabnew.nextcontact.ru/giorgiy_03/screenrecordercontracts"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func main() {

	// Настройка логгера
	output := zerolog.ConsoleWriter{
		TimeFormat: "02.01.2006 15:04:05",
		Out:        os.Stdout,
	}
	log.Logger = log.Output(output)

	// Configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal().Msgf("Config error: %s", err)
	}

	// Start GRPC Client
	conn, err := grpc.Dial(fmt.Sprintf("localhost:%d", cfg.Grpc.Port), grpc.WithInsecure())
	if err != nil {
		log.Fatal().Err(err).Msg("app - Run - Start grpc Client")
	}
	defer conn.Close()

	client := proto.NewScreenRecorderServiceClient(conn)

	// Методы

	// Aвторизация
	_, err = client.Authorization(context.Background(), &proto.AuthorizationRequest{Login: "login??", Password: "password??"})
	if err != nil {
		log.Info().Msgf("Не удачная попытка авторизации: %s", err)
	} else {
		log.Info().Msg("Вы авторизованы")
	}

	//Добавить новое видео
	Video := proto.Video{
		Name:        "Какое-то видео",
		Login:       "1234",
		SessionId:   "1234",
		CreatedAt:   timestamppb.Now(),
		Fullpath:    "1234",
		MacAddr:     "1234",
		IpAddr:      "1234",
		GetArgs:     "1234",
		FileName:    "1234",
		ProjectId:   "1234",
		ProjectName: "1234",
	}

	respAddVideo, err := client.AddVideo(context.Background(), &proto.AddVideoRequest{Video: &Video})
	if err != nil {
		log.Fatal().Msgf("client.AddVideo: %s", err)
	}
	log.Info().Err(err).Msg("Добавлено новое видео")
	log.Info().Err(err).Msg(respAddVideo.String())

	// Список всех видео
	respListVideos, err := client.ListVideos(context.Background(), &proto.ListVideosRequest{})
	if err != nil {
		log.Fatal().Msgf("client.ListVideos: %s", err)
	}

	log.Info().Err(err).Msg("Список видео: ")
	for _, video := range respListVideos.Videos {
		log.Info().Err(err).Msgf("%v", video)
	}

	//// Подсчитать количество видео
	//respCountVideos, err := client.CountVideos(context.Background(), &proto.CountVideosRequest{})
	//if err != nil {
	//	log.Fatal().Msgf("client.ListVideos: %s", err)
	//}
	//count := respCountVideos.Count
	//log.Info().Err(err).Msgf("Количество видео: %d", count)

}
