package config

import (
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type DailyLogWriter struct {
	mu         sync.Mutex
	currentDay string
	location   *time.Location
	file       *os.File
	writer     io.Writer
}

func NewDailyLogWriter(logsPath string, loc *time.Location) (*DailyLogWriter, error) {
	d := &DailyLogWriter{location: loc}
	if err := d.rotate(logsPath); err != nil {
		return nil, err
	}
	go d.autoRotate(logsPath)
	return d, nil
}

func (d *DailyLogWriter) Write(p []byte) (n int, err error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.writer.Write(p)
}

func (d *DailyLogWriter) rotate(logsPath string) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	now := time.Now().In(d.location)
	day := now.Format("2006-01-02")
	if day == d.currentDay {
		return nil
	}

	if d.file != nil {
		d.file.Close()
		go compressLogFile(filepath.Join(logsPath, fmt.Sprintf("gorgon-%s.log", d.currentDay)))
	}

	if err := os.MkdirAll(logsPath, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create logs directory: %v", err)
	}

	filename := fmt.Sprintf("gorgon-%s.log", day)
	fullPath := filepath.Join(logsPath, filename)
	file, err := os.OpenFile(fullPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
	if err != nil {
		return fmt.Errorf("failed to open log file: %v", err)
	}

	d.file = file
	d.currentDay = day
	d.writer = io.MultiWriter(os.Stdout, file)
	return nil
}

func (d *DailyLogWriter) autoRotate(logsPath string) {
	for {
		time.Sleep(30 * time.Second)
		_ = d.rotate(logsPath)
	}
}

// Compact the file .log to .gz
func compressLogFile(path string) {
	in, err := os.Open(path)
	if err != nil {
		log.Printf("compressLogFile: could not open file: %v", err)
		return
	}
	defer in.Close()

	outPath := path + ".gz"
	out, err := os.Create(outPath)
	if err != nil {
		log.Printf("compressLogFile: could not create gz file: %v", err)
		return
	}
	defer out.Close()

	gw := gzip.NewWriter(out)
	defer gw.Close()

	if _, err := io.Copy(gw, in); err != nil {
		log.Printf("compressLogFile: failed to copy: %v", err)
		return
	}

	// Remove original .log file after compression
	if err := os.Remove(path); err != nil {
		log.Printf("compressLogFile: could not remove original file: %v", err)
	}
}

func NewLogger() *slog.Logger {
	loc, err := time.LoadLocation("America/Sao_Paulo")
	if err != nil {
		log.Fatalf("failed to load location: %v", err)
	}

	writer, err := NewDailyLogWriter(LogsPath, loc)
	if err != nil {
		log.Fatalf("failed to create log writer: %v", err)
	}

	handler := slog.NewJSONHandler(writer, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	return slog.New(handler)
}
