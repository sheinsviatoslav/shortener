package grpcserv

import (
	"context"
	"errors"
	"github.com/sheinsviatoslav/shortener/internal/common"
	"github.com/sheinsviatoslav/shortener/internal/storage"
	"github.com/sheinsviatoslav/shortener/internal/utils/hash"
	pb "github.com/sheinsviatoslav/shortener/proto"
)

type UrlsServer struct {
	pb.UnimplementedUrlsServer
	storage storage.Storage
}

func (s *UrlsServer) CreateShortURL(ctx context.Context, in *pb.CreateShortURLRequest) (*pb.CreateShortURLResponse, error) {
	var response pb.CreateShortURLResponse

	if in.OriginalUrl == "" {
		return &response, errors.New("url is required")
	}

	shortURL, isExists, err := s.storage.GetShortURLByOriginalURL(ctx, in.OriginalUrl)
	if err != nil {
		return &response, err
	}

	if !isExists {
		shortURL = hash.Generator(common.DefaultHashLength)
		if err = s.storage.AddNewURL(ctx, in.OriginalUrl, shortURL, in.UserId); err != nil {
			return &response, err
		}
	}

	response.ShortUrl = shortURL

	return &response, nil
}

func (s *UrlsServer) GetOriginalURL(ctx context.Context, in *pb.GetOriginalURLRequest) (*pb.GetOriginalURLResponse, error) {
	var response pb.GetOriginalURLResponse

	if in.ShortUrl == "" {
		return &response, errors.New("empty short url")
	}

	originalURL, isDeleted, err := s.storage.GetOriginalURLByShortURL(ctx, in.ShortUrl)
	if err != nil {
		return &response, err
	}

	if isDeleted {
		return &response, errors.New("url is already deleted")
	}

	if originalURL == "" {
		return &response, errors.New("empty original url")
	}

	response.OriginalUrl = originalURL

	return &response, nil
}

func (s *UrlsServer) ShortenBatch(ctx context.Context, in *pb.ShortenBatchRequest) (*pb.ShortenBatchResponse, error) {
	var response pb.ShortenBatchResponse
	var urlsInput storage.InputManyUrls

	if len(in.Urls) == 0 {
		return &response, errors.New("empty urls")
	}

	for _, url := range in.Urls {
		urlsInput = append(urlsInput, storage.InputManyUrlsItem{
			CorrelationID: url.CorrelationId,
			OriginalURL:   url.OriginalUrl,
		})
	}

	urls, err := s.storage.AddManyUrls(ctx, urlsInput, in.UserId)
	if err != nil {
		return &response, err
	}

	for _, url := range urls {
		response.Urls = append(response.Urls, &pb.OutputManyUrlsItem{
			CorrelationId: url.CorrelationID,
			ShortUrl:      url.ShortURL,
		})
	}

	return &response, nil
}

func (s *UrlsServer) GetUserUrls(ctx context.Context, in *pb.GetUserUrlsRequest) (*pb.GetUserUrlsResponse, error) {
	var response pb.GetUserUrlsResponse

	urls, err := s.storage.GetUserUrls(ctx, in.UserId)
	if err != nil {
		return &response, err
	}

	for _, url := range urls {
		response.UserUrls = append(response.UserUrls, &pb.UserUrlsItem{
			OriginalUrl: url.OriginalURL,
			ShortUrl:    url.ShortURL,
		})
	}

	return &response, nil
}

func (s *UrlsServer) DeleteUserUrls(ctx context.Context, in *pb.DeleteUserUrlsRequest) (*pb.DeleteUserUrlsResponse, error) {
	var response pb.DeleteUserUrlsResponse

	go func(urls []string, id string) {
		_ = s.storage.DeleteUserUrls(ctx, in.ShortUrls, id)
	}(in.ShortUrls, in.UserId)

	if err := s.storage.DeleteUserUrls(ctx, in.ShortUrls, in.UserId); err != nil {
		return &response, err
	}

	return &response, nil
}
