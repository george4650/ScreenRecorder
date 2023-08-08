package server

import (
	"context"
	"time"

	"myapp/internal/models"
	"myapp/internal/usecase"

	"google.golang.org/protobuf/types/known/timestamppb"

	proto "gitlabnew.nextcontact.ru/giorgiy_03/screenrecordercontracts"
)

type Server struct {
	proto.UnimplementedScreenRecorderServiceServer
	us *usecase.ScreenRecorderUseCases
	au *usecase.AuthUseCases
}

func NewServer(au *usecase.AuthUseCases, us *usecase.ScreenRecorderUseCases) *Server {
	return &Server{
		us: us,
		au: au,
	}
}

func (s *Server) Authorization(ctx context.Context, req *proto.AuthorizationRequest) (*proto.AuthorizationResponse, error) {
	err := s.au.Authentication(req.Login, req.Password)
	if err != nil {
		return nil, err
	}
	return &proto.AuthorizationResponse{}, nil
}

func (s *Server) AddVideo(ctx context.Context, req *proto.AddVideoRequest) (*proto.AddVideoResponse, error) {

	var err error

	video := models.Video{
		Name:        req.Video.Name,
		Login:       req.Video.Name,
		SessionId:   req.Video.SessionId,
		Fullpath:    req.Video.Fullpath,
		MacAddr:     req.Video.MacAddr,
		IpAddr:      req.Video.IpAddr,
		GetArgs:     req.Video.GetArgs,
		FileName:    req.Video.FileName,
		ProjectUUID: req.Video.ProjectId,
	}

	loc, _ := time.LoadLocation("Europe/Moscow")
	timeString := req.Video.CreatedAt.AsTime().In(loc).Format("2006.01.02 15:04:05")
	video.CreatedAt, err = time.Parse("2006.01.02 15:04:05", timeString)
	if err != nil {
		return nil, err
	}

	err = s.us.AddVideo(ctx, video)
	if err != nil {
		return nil, err
	}

	return &proto.AddVideoResponse{}, nil
}

func (s *Server) DownloadVideo(ctx context.Context, req *proto.DownloadVideoRequest) (*proto.DownloadVideoResponse, error) {
	fileName, info, err := s.us.DownloadVideo(ctx, req.VideoId)
	if err != nil {
		return nil, err
	}
	return &proto.DownloadVideoResponse{FileName: fileName, Info: info}, nil
}

func (s *Server) DeleteCashe(ctx context.Context, req *proto.DeleteCasheRequest) (*proto.DeleteCasheResponse, error) {

	err := s.us.DeleteCashe(ctx)
	if err != nil {
		return nil, err
	}
	return &proto.DeleteCasheResponse{}, nil
}

func (s *Server) UpdatePostgresProjects(ctx context.Context, req *proto.UpdatePostgresProjectsRequest) (*proto.UpdatePostgresProjectsResponce, error) {

	err := s.us.UpdatePostgresProjects(ctx)
	if err != nil {
		return nil, err
	}
	return &proto.UpdatePostgresProjectsResponce{}, nil
}

// В последних 2-ух не удалось распарсить offset :(
func (s *Server) ListVideos(ctx context.Context, req *proto.ListVideosRequest) (*proto.ListVideosResponse, error) {

	searchValue := []models.SearchValue{}
	for _, search := range req.SearchBy {
		searchBy := models.SearchValue{
			Field: search.Field,
			Value: search.Value,
		}
		searchValue = append(searchValue, searchBy)
	}

	orderValue := []models.OrderValue{}
	for _, order := range req.OrderBy {
		orderBy := models.OrderValue{
			Field: order.Field,
			Value: order.Value,
		}
		orderValue = append(orderValue, orderBy)
	}

	// invalid memory address or nil pointer dereference
	//offset := &req.Offset.Offset
	//log.Info().Msgf("Offset %d", offset)

	videos, err := s.us.ListVideos(ctx, searchValue, orderValue, 0)
	if err != nil {
		return nil, err
	}

	var Videos []*proto.Video

	for _, v := range videos {
		Video := proto.Video{
			Login:       v.Login,
			SessionId:   v.SessionId,
			CreatedAt:   timestamppb.New(v.CreatedAt),
			IpAddr:      v.IpAddr,
			FileName:    v.FileName,
			ProjectId:   v.ProjectUUID,
			ProjectName: v.ProjectName,
		}
		Videos = append(Videos, &Video)
	}
	return &proto.ListVideosResponse{Videos: Videos}, nil
}

func (s *Server) CountVideos(ctx context.Context, req *proto.CountVideosRequest) (*proto.CountVideosResponce, error) {

	searchValue := []models.SearchValue{}
	for _, search := range req.SearchBy {
		searchBy := models.SearchValue{
			Field: search.Field,
			Value: search.Value,
		}
		searchValue = append(searchValue, searchBy)
	}

	orderValue := []models.OrderValue{}
	for _, order := range req.OrderBy {
		orderBy := models.OrderValue{
			Field: order.Field,
			Value: order.Value,
		}
		orderValue = append(orderValue, orderBy)
	}

	// invalid memory address or nil pointer dereference
	//offset := &req.Offset.Offset
	//log.Info().Msgf("Offset %d", offset)

	count, err := s.us.CountVideos(ctx, searchValue, orderValue, 0)
	if err != nil {
		return nil, err
	}
	return &proto.CountVideosResponce{Count: int64(count)}, nil
}
