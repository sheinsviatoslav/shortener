package storage

import (
	"bufio"
	"encoding/json"
	"github.com/sheinsviatoslav/shortener/internal/config"
	"io"
	"os"
)

type FileStorage struct {
	Producer Producer
	Consumer Consumer
	URLMap   URLMap
}

type URLMap map[string]string

type Producer struct {
	file   *os.File
	writer *bufio.Writer
}

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

func (p *Producer) WriteURLItems(urlItem *URLMap) error {
	data, err := json.MarshalIndent(&urlItem, "", "   ")
	if err != nil {
		return err
	}

	if _, err := p.writer.Write(data); err != nil {
		return err
	}

	if err := p.writer.WriteByte('\n'); err != nil {
		return err
	}

	return p.writer.Flush()
}

func WriteURLItemToFile(urlMap URLMap) error {
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

func (p *Producer) Close() error {
	return p.file.Close()
}

type Consumer struct {
	file   *os.File
	reader *bufio.Reader
}

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

func ReadURLItems() (*URLMap, error) {
	if _, err := os.Stat(*config.FileStoragePath); err == nil {
		fileReader, err := NewConsumer(*config.FileStoragePath)
		if err != nil {
			return nil, err
		}

		defer fileReader.Close()

		byteValue, err := io.ReadAll(fileReader.reader)
		if err != nil {
			return nil, err
		}

		urlItem := URLMap{}
		err = json.Unmarshal(byteValue, &urlItem)
		if err != nil {
			return nil, err
		}

		return &urlItem, nil
	}

	return nil, nil
}

func (c *Consumer) Close() error {
	return c.file.Close()
}
