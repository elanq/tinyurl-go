package service_test

import (
	"context"
	"testing"

	"github.com/elanq/pastebin-go/model"
	mock_repository "github.com/elanq/pastebin-go/repository/mock"
	"github.com/elanq/pastebin-go/service"
	mock_service "github.com/elanq/pastebin-go/service/mock"
	"github.com/go-redis/redis/v8"
	"go.uber.org/mock/gomock"
)

func TestCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	urlRepositoryMock := mock_repository.NewMockURL(ctrl)
	cacheServiceMock := mock_service.NewMockCache(ctrl)
	// assert service.create is called with expected url object
	urlRepositoryMock.EXPECT().Create(gomock.Any(), gomock.Any()).Return(&model.URL{}, nil).Times(1)
	cacheServiceMock.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)

	s := service.NewURL(urlRepositoryMock, cacheServiceMock)
	s.Create(context.Background(), model.URL{
		UserId:  "12345",
		LongURL: "http://www.google.com",
	})
}

func TestGetByShortURL(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	urlRepositoryMock := mock_repository.NewMockURL(ctrl)
	cacheServiceMock := mock_service.NewMockCache(ctrl)

	// when cache is missing, it should perform query into url table and create new cache
	s := service.NewURL(urlRepositoryMock, cacheServiceMock)
	cacheServiceMock.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, redis.Nil).Times(1)
	urlRepositoryMock.EXPECT().FindByShortUrl(gomock.Any(), gomock.Any()).Return(&model.URL{}, nil).Times(1)
	cacheServiceMock.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)

	url, err := s.GetByShortURL(context.Background(), "url")
	if err != nil {
		t.Error(err)
	}
	if url == nil {
		t.Error("url should not be null")
	}
}
