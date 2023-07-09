package file

import (
	"errors"
	"io"
	"myclipboard/tui"
	"os"
	"strings"

	bubblelist "github.com/charmbracelet/bubbles/list"
)

type Storage struct {
	file         *os.File
	container    []tui.Item
	maxLength    int
	ClipboardSig chan struct{}
}

func (s *Storage) Init(filePath string) error {
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}

	s.file = file
	return nil
}

func (s *Storage) Load(maxLength int) (err error) {
	container := make([]tui.Item, maxLength)
	content, err := io.ReadAll(s.file)
	if err != nil {
		return
	}

	data := strings.Split(string(content), "\n")
	if len(data) == 0 {
		container = make([]tui.Item, s.maxLength)
	}
	for idx, d := range data {
		item := tui.Item{Display: tui.Displayed(d), Value: d}
		container[idx] = item
	}
	s.container = container
	s.ClipboardSig = make(chan struct{})
	return
}

func (s *Storage) Write(text string) (err error) {
	s.container = s.container[:s.maxLength+1]
	s.container = append([]tui.Item{tui.Item{Value: text, Display: tui.Displayed(text)}}, s.container[0:19]...)
	return
}

func (s *Storage) Select(idx int) (err error) {
	t := s.container[idx]
	front := s.container[0:idx]
	back := s.container[idx+1 : len(s.container)]

	var merged []tui.Item
	merged = append(merged, t)
	merged = append(merged, front...)
	merged = append(merged, back...)
	s.container = merged
	return
}

func (s *Storage) Read(idx int) (data string, err error) {
	if len(s.container) == 0 {
		err = errors.New("no list loaded")
		return
	}
	if idx+1 > s.maxLength {
		err = errors.New("index exceed max length of container")
		return
	}

	data = s.container[idx].Value
	return
}

func (s *Storage) ToBubbleList() (ret []bubblelist.Item) {
	for _, item := range s.container {
		if len([]rune(item.Display)) > 80 {
			item.Display = tui.Displayed([]rune(item.Display)[:80])
		}

		if strings.Contains(string(item.Display), "\n") {
			item.Display = tui.Displayed(strings.Split(string(item.Display), "\n")[0]+"\\n"+strings.Split(string(item.Display), "\n")[1]) + "..."
		}

		ret = append(ret, item)
	}
	return
}

func (s *Storage) containerToStringSlice() (ret []string) {
	for _, item := range s.container {
		ret = append(ret, string(item.Display))
	}
	return
}

func (s *Storage) Close() error {
	content := strings.Join(s.containerToStringSlice(), "\n")
	// initialize file
	if err := s.file.Truncate(0); err != nil {
		return err
	}

	// set cursor to line 0
	if _, err := s.file.Seek(0, 0); err != nil {
		return err
	}

	if _, err := s.file.WriteString(content); err != nil {
		return err
	}
	return s.file.Close()
}
