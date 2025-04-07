package storage

//
//import (
//	"context"
//	"fmt"
//	"github.com/google/uuid"
//	_ "github.com/jackc/pgx/v5/stdlib"
//	"github.com/sheinsviatoslav/shortener/internal/common"
//	"github.com/sheinsviatoslav/shortener/internal/storage"
//	"github.com/stretchr/testify/assert"
//	"github.com/stretchr/testify/require"
//	"regexp"
//	"testing"
//)
//
//var defaultDSN = "user=postgres"
//
//func dropDB(p *PgStorage) {
//	p.DB.Exec("DROP TABLE urls")
//}
//
//func TestPgStorage_Connect(t *testing.T) {
//	type args struct {
//		dsn string
//	}
//	tests := []struct {
//		name    string
//		args    args
//		wantErr assert.ErrorAssertionFunc
//	}{
//		{name: "success connected", args: args{dsn: defaultDSN}, wantErr: assert.NoError},
//		{name: "error", args: args{}, wantErr: assert.Error},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			p := NewPgStorage()
//			tt.wantErr(t, p.Connect(tt.args.dsn))
//			dropDB(p)
//		})
//	}
//}
//
//func TestPgStorage_AddNewURL(t *testing.T) {
//	type args struct {
//		ctx         context.Context
//		originalURL string
//		shortURL    string
//		userID      string
//	}
//	tests := []struct {
//		name    string
//		args    args
//		wantErr assert.ErrorAssertionFunc
//	}{
//		{
//			name: "success",
//			args: args{
//				ctx:         context.Background(),
//				originalURL: "https://yandex.ru/",
//				shortURL:    "7IENelKX",
//				userID:      uuid.New().String(),
//			},
//			wantErr: assert.NoError,
//		},
//		{
//			name: "wrong userID",
//			args: args{
//				ctx:         context.Background(),
//				originalURL: "https://yandex.ru/",
//				shortURL:    "7IENelKX",
//				userID:      "",
//			},
//			wantErr: assert.Error,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			p := NewPgStorage()
//			assert.NoError(t, p.Connect(defaultDSN))
//			tt.wantErr(t, p.AddNewURL(tt.args.ctx, tt.args.originalURL, tt.args.shortURL, tt.args.userID))
//			dropDB(p)
//		})
//	}
//}
//
//func TestPgStorage_GetOriginalURLByShortURL(t *testing.T) {
//	type args struct {
//		ctx      context.Context
//		shortURL string
//	}
//	tests := []struct {
//		name    string
//		args    args
//		want    string
//		want1   bool
//		wantErr assert.ErrorAssertionFunc
//	}{
//		{
//			name:    "success",
//			args:    args{ctx: context.Background(), shortURL: "7IENelKX"},
//			want:    "https://yandex.ru/",
//			want1:   false,
//			wantErr: assert.NoError,
//		},
//		{
//			name:    "no value",
//			args:    args{ctx: context.Background(), shortURL: "7IENelKX"},
//			want:    "",
//			want1:   false,
//			wantErr: assert.Error,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			p := NewPgStorage()
//			assert.NoError(t, p.Connect(defaultDSN))
//			if tt.want != "" {
//				assert.NoError(t, p.AddNewURL(tt.args.ctx, tt.want, tt.args.shortURL, uuid.New().String()))
//			}
//			got, got1, err := p.GetOriginalURLByShortURL(tt.args.ctx, tt.args.shortURL)
//			if !tt.wantErr(t, err) {
//				return
//			}
//			assert.Equal(t, tt.want, got)
//			assert.Equal(t, tt.want1, got1)
//			dropDB(p)
//		})
//	}
//}
//
//func TestPgStorage_GetShortURLByOriginalURL(t *testing.T) {
//	type args struct {
//		ctx         context.Context
//		originalURL string
//	}
//	tests := []struct {
//		name    string
//		args    args
//		want    string
//		want1   bool
//		wantErr assert.ErrorAssertionFunc
//	}{
//		{
//			name:    "success",
//			args:    args{ctx: context.Background(), originalURL: "https://yandex.ru/"},
//			want:    "7IENelKX",
//			want1:   true,
//			wantErr: assert.NoError,
//		},
//		{
//			name:    "no value",
//			args:    args{ctx: context.Background()},
//			want:    "",
//			want1:   false,
//			wantErr: assert.NoError,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			p := NewPgStorage()
//			assert.NoError(t, p.Connect(defaultDSN))
//			if tt.want != "" {
//				assert.NoError(t, p.AddNewURL(tt.args.ctx, tt.args.originalURL, tt.want, uuid.New().String()))
//			}
//			got, got1, err := p.GetShortURLByOriginalURL(tt.args.ctx, tt.args.originalURL)
//			if !tt.wantErr(t, err) {
//				return
//			}
//			assert.Equal(t, tt.want, got)
//			assert.Equal(t, tt.want1, got1)
//			dropDB(p)
//		})
//	}
//}
//
//func TestPgStorage_AddManyUrls(t *testing.T) {
//	type args struct {
//		ctx    context.Context
//		urls   storage.InputManyUrls
//		userID string
//	}
//	tests := []struct {
//		name    string
//		args    args
//		want    storage.OutputManyUrls
//		wantErr assert.ErrorAssertionFunc
//	}{
//		{
//			name: "success",
//			args: args{
//				ctx: context.Background(),
//				urls: storage.InputManyUrls{
//					{CorrelationID: "e4b54da9-edab-4954-a822-1dd4fc4b7de5", OriginalURL: "https://yandex.ru/"},
//					{CorrelationID: "e4b54da9-edab-4954-a822-1dd4fc4b7de6", OriginalURL: "https://practicum.ru/"},
//				},
//				userID: uuid.New().String(),
//			},
//			want: storage.OutputManyUrls{
//				{CorrelationID: "e4b54da9-edab-4954-a822-1dd4fc4b7de5", ShortURL: "http://localhost:8080/SDNPRPub"},
//				{CorrelationID: "e4b54da9-edab-4954-a822-1dd4fc4b7de6", ShortURL: "http://localhost:8080/SDNPRPub"},
//			},
//			wantErr: assert.NoError,
//		},
//		{
//			name: "no original url",
//			args: args{
//				ctx: context.Background(),
//				urls: storage.InputManyUrls{
//					{CorrelationID: "e4b54da9-edab-4954-a822-1dd4fc4b7de5", OriginalURL: ""},
//					{CorrelationID: "e4b54da9-edab-4954-a822-1dd4fc4b7de6", OriginalURL: "https://practicum.ru/"},
//				},
//				userID: uuid.New().String(),
//			},
//			want: storage.OutputManyUrls{
//				{CorrelationID: "e4b54da9-edab-4954-a822-1dd4fc4b7de5", ShortURL: "http://localhost:8080/SDNPRPub"},
//				{CorrelationID: "e4b54da9-edab-4954-a822-1dd4fc4b7de6", ShortURL: "http://localhost:8080/SDNPRPub"},
//			},
//			wantErr: assert.Error,
//		},
//		{
//			name: "invalid original url",
//			args: args{
//				ctx: context.Background(),
//				urls: storage.InputManyUrls{
//					{CorrelationID: "e4b54da9-edab-4954-a822-1dd4fc4b7de5", OriginalURL: "invalid"},
//					{CorrelationID: "e4b54da9-edab-4954-a822-1dd4fc4b7de6", OriginalURL: "https://practicum.ru/"},
//				},
//				userID: uuid.New().String(),
//			},
//			want: storage.OutputManyUrls{
//				{CorrelationID: "e4b54da9-edab-4954-a822-1dd4fc4b7de5", ShortURL: "http://localhost:8080/SDNPRPub"},
//				{CorrelationID: "e4b54da9-edab-4954-a822-1dd4fc4b7de6", ShortURL: "http://localhost:8080/SDNPRPub"},
//			},
//			wantErr: assert.Error,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			p := new(PgStorage)
//			assert.NoError(t, p.Connect(defaultDSN))
//			got, err := p.AddManyUrls(tt.args.ctx, tt.args.urls, tt.args.userID)
//			if !tt.wantErr(t, err) {
//				return
//			}
//			for _, item := range got {
//				isMatch, _ := regexp.MatchString(fmt.Sprintf(
//					"http://localhost:8080/[0-9a-zA-Z]{%d}",
//					common.DefaultHashLength,
//				), item.ShortURL)
//				require.NoError(t, err)
//				assert.Equal(t, true, isMatch)
//			}
//			dropDB(p)
//		})
//	}
//}
//
//func TestPgStorage_GetUserUrls(t *testing.T) {
//	type args struct {
//		ctx    context.Context
//		userID string
//	}
//	tests := []struct {
//		name    string
//		args    args
//		want    storage.UserUrls
//		wantErr assert.ErrorAssertionFunc
//	}{
//		{
//			name: "success",
//			args: args{
//				ctx:    context.Background(),
//				userID: uuid.New().String(),
//			},
//			want: storage.UserUrls{
//				{ShortURL: "http://localhost:8080/SDNPRPub", OriginalURL: "https://practicum.ru/"},
//				{ShortURL: "http://localhost:8080/7IENelKX", OriginalURL: "https://yandex.ru/"},
//			},
//			wantErr: assert.NoError,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			p := NewPgStorage()
//			assert.NoError(t, p.Connect(defaultDSN))
//			assert.NoError(t, p.AddNewURL(tt.args.ctx, "https://practicum.ru/", "SDNPRPub", tt.args.userID))
//			assert.NoError(t, p.AddNewURL(tt.args.ctx, "https://yandex.ru/", "7IENelKX", tt.args.userID))
//			got, err := p.GetUserUrls(tt.args.ctx, tt.args.userID)
//			if !tt.wantErr(t, err) {
//				return
//			}
//			assert.Equal(t, tt.want, got)
//			dropDB(p)
//		})
//	}
//}
//
//func TestPgStorage_DeleteUserUrls(t *testing.T) {
//	type args struct {
//		ctx       context.Context
//		shortUrls []string
//		userID    string
//	}
//	tests := []struct {
//		name    string
//		args    args
//		wantErr assert.ErrorAssertionFunc
//	}{
//		{
//			name: "success",
//			args: args{
//				ctx:       context.Background(),
//				shortUrls: []string{"SDNPRPub", "7IENelKX"},
//				userID:    uuid.New().String(),
//			},
//			wantErr: assert.NoError,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			p := NewPgStorage()
//			assert.NoError(t, p.Connect(defaultDSN))
//			assert.NoError(t, p.AddNewURL(tt.args.ctx, "https://practicum.ru/", "SDNPRPub", tt.args.userID))
//			assert.NoError(t, p.AddNewURL(tt.args.ctx, "https://yandex.ru/", "7IENelKX", tt.args.userID))
//			tt.wantErr(t, p.DeleteUserUrls(tt.args.ctx, tt.args.shortUrls, tt.args.userID))
//			_, isDeleted, _ := p.GetOriginalURLByShortURL(tt.args.ctx, tt.args.shortUrls[0])
//			assert.True(t, isDeleted)
//			_, isDeleted, _ = p.GetOriginalURLByShortURL(tt.args.ctx, tt.args.shortUrls[1])
//			assert.True(t, isDeleted)
//			dropDB(p)
//		})
//	}
//}
