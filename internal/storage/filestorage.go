package storage

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/url"
	"os"
	"slices"

	"github.com/sheinsviatoslav/shortener/internal/common"
	"github.com/sheinsviatoslav/shortener/internal/config"
	"github.com/sheinsviatoslav/shortener/internal/utils/hash"
)

// FileData is a type for data stored in file
type FileData map[string]string

// FileStorage is a storage type
type FileStorage struct {
	Producer Producer
	Consumer Consumer
}

// NewFileStorage constructs new FileStorage struct
func NewFileStorage() *FileStorage {
	return &FileStorage{
		Producer: Producer{},
		Consumer: Consumer{},
	}
}

// GetOriginalURLByShortURL returns original url
func (fs *FileStorage) GetOriginalURLByShortURL(_ context.Context, inputShortURL string) (string, bool, error) {
	urlItems, err := fs.ReadURLItems()
	if err != nil {
		return "", false, err
	}

	for originalURL, shortURL := range *urlItems {
		if shortURL == inputShortURL {
			return originalURL, false, nil
		}
	}

	return "", false, nil
}

// GetShortURLByOriginalURL returns short url
func (fs *FileStorage) GetShortURLByOriginalURL(_ context.Context, originalURL string) (string, bool, error) {
	urlItems, err := fs.ReadURLItems()
	if err != nil {
		return "", false, err
	}

	if shortURL, ok := (*urlItems)[originalURL]; ok {
		return shortURL, true, nil
	}

	return "", false, nil
}

// AddNewURL adds originalURL-shortURL pair into the storage
func (fs *FileStorage) AddNewURL(_ context.Context, originalURL string, shortURL string, _ string) error {
	urlItems, err := fs.ReadURLItems()
	if err != nil {
		return err
	}

	if _, ok := (*urlItems)[originalURL]; !ok {
		(*urlItems)[originalURL] = shortURL
		if fsError := fs.WriteURLItem(*urlItems); fsError != nil {
			return err
		}
	}

	return nil
}

// AddManyUrls adds multiple originalURL-shortURL pairs into the storage
func (fs *FileStorage) AddManyUrls(_ context.Context, urls InputManyUrls, _ string) (OutputManyUrls, error) {
	var output OutputManyUrls
	urlItems, readFileErr := fs.ReadURLItems()
	if readFileErr != nil {
		return nil, readFileErr
	}

	for _, item := range urls {
		if item.OriginalURL == "" {
			return nil, errors.New("url is required")
		}

		if _, err := url.ParseRequestURI(item.OriginalURL); err != nil {
			return nil, err
		}

		shortURL, isExists := (*urlItems)[item.OriginalURL]

		if !isExists {
			shortURL = hash.Generator(common.DefaultHashLength)
			(*urlItems)[item.OriginalURL] = shortURL
		}

		u, _ := url.Parse(*config.BaseURL)
		relative, _ := url.Parse(shortURL)

		output = append(output, OutputManyUrlsItem{CorrelationID: item.CorrelationID, ShortURL: u.ResolveReference(relative).String()})
	}

	if fsError := fs.WriteURLItem(*urlItems); fsError != nil {
		return nil, fsError
	}

	return output, nil
}

// GetUserUrls returns multiple originalURL-shortURL pairs of current user
func (fs *FileStorage) GetUserUrls(_ context.Context, _ string) (UserUrls, error) {
	output := make(UserUrls, 0)
	urlItems, err := fs.ReadURLItems()
	if err != nil {
		return nil, err
	}

	for originalURL, shortURL := range *urlItems {
		u, _ := url.Parse(*config.BaseURL)
		relative, _ := url.Parse(shortURL)

		output = append(output, UserUrlsItem{OriginalURL: originalURL, ShortURL: u.ResolveReference(relative).String()})
	}

	return output, nil
}

// DeleteUserUrls deletes multiple short urls
func (fs *FileStorage) DeleteUserUrls(_ context.Context, shortUrls []string, _ string) error {
	urlItems, err := fs.ReadURLItems()
	if err != nil {
		return nil
	}

	for originalURL, shortURL := range *urlItems {
		if slices.Contains(shortUrls, shortURL) {
			delete(*urlItems, originalURL)
		}
	}

	if fsError := fs.WriteURLItem(*urlItems); fsError != nil {
		return fsError
	}

	return nil
}

// WriteURLItem writes multiple originalURL-shortURL pairs into the file
func (fs *FileStorage) WriteURLItem(urlMap FileData) error {
	var fileWriter, err = NewProducer(*config.FileStoragePath)
	if err != nil {
		return err
	}
	defer fileWriter.Close()

	if err = fileWriter.WriteURLItems(&urlMap); err != nil {
		return err
	}

	return nil
}

// ReadURLItems reads multiple originalURL-shortURL pairs from the file
func (fs *FileStorage) ReadURLItems() (*FileData, error) {
	if _, err := os.Stat(*config.FileStoragePath); err == nil {
		fileReader, fileErr := NewConsumer(*config.FileStoragePath)
		if fileErr != nil {
			return nil, fileErr
		}

		defer fileReader.Close()

		byteValue, readErr := io.ReadAll(fileReader.reader)
		if readErr != nil {
			return nil, readErr
		}

		urlItems := FileData{}
		if err = json.Unmarshal(byteValue, &urlItems); err != nil {
			return nil, err
		}

		return &urlItems, nil
	}

	return &FileData{}, nil
}

// Producer is a producer type
type Producer struct {
	file   *os.File
	writer *bufio.Writer
}

// NewProducer constructs new Producer struct
func NewProducer(filename string) (*Producer, error) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	return &Producer{
		file:   file,
		writer: bufio.NewWriter(file),
	}, nil
}

// WriteURLItems writes multiple originalURL-shortURL pairs into the file
func (p *Producer) WriteURLItems(urlItems *FileData) error {
	data, err := json.MarshalIndent(&urlItems, "", "   ")
	if err != nil {
		return err
	}

	if _, err = p.writer.Write(data); err != nil {
		return err
	}

	if err = p.writer.WriteByte('\n'); err != nil {
		return err
	}

	return p.writer.Flush()
}

// Close closes the producer
func (p *Producer) Close() error {
	return p.file.Close()
}

// Consumer is a consumer type
type Consumer struct {
	file   *os.File
	reader *bufio.Reader
}

// NewConsumer constructs new Consumer struct
func NewConsumer(filename string) (*Consumer, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		file:   file,
		reader: bufio.NewReader(file),
	}, nil
}

// Close closes the consumer
func (c *Consumer) Close() error {
	return c.file.Close()
}
