package storage

import (
	"bufio"
	"encoding/json"
	"io"
	"os"
)

type URLItem struct {
	ID          int    `json:"id"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type URLItems struct {
	Items []URLItem `json:"items"`
}

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

func (p *Producer) WriteURLItems(urlItem *URLItems) error {
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

func (c *Consumer) ReadURLItems() (*URLItems, error) {
	byteValue, err := io.ReadAll(c.reader)
	if err != nil {
		return nil, err
	}

	urlItem := URLItems{}
	err = json.Unmarshal(byteValue, &urlItem)
	if err != nil {
		return nil, err
	}

	return &urlItem, nil
}

func (c *Consumer) Close() error {
	return c.file.Close()
}
