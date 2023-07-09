package main

import (
	"fmt"
	clipboard2 "myclipboard/clipboard"
	"myclipboard/storage/file"
	"myclipboard/tui"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"golang.design/x/clipboard"
)

// Model : 앱 기본 모델
// cursor 는 UI 상에 표시되는 화살표의 위치를 저장하고
// selector 는 엔터키 입력시 "copied"라는 메시지를 출력시키기 위해 출력 위치를 기록하는데 사용됨
type Model struct {
	choices  list.Model
	storage  *file.Storage
	cursor   int
	selector int
}

// tickMsg : ticker 가 생성하는 tea.Cmd의 리턴타입
type tickMsg struct{}

// ticker : tea.Cmd 를 250 millisecond 마다 생성함 >
// 250 millisecond 마다 생성되는 tea.Cmd 가 tickMsg 를 생성하고, 이 tickMsg 를 Model 의 Update 메서드가 소비하면서 UI가 갱신됨
func ticker() tea.Cmd {
	return tea.Tick(250*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg{}
	})
}

// initialModel : 어플리케이션에서 사용될 기본 Model 을 생성
func initialModel(s *file.Storage) Model {
	_ = s.Load(20)
	clipboard2.InitClipboard(s)
	return Model{
		choices:  list.New(s.ToBubbleList(), list.NewDefaultDelegate(), 0, 0),
		selector: -1,
		storage:  s,
	}
}

// Init : Model 에 대한 작업 시작시 Bubbletea 에서 자체적으로 실행시킴, 여기서는 ticker 를 실행시키도록 함
func (m Model) Init() tea.Cmd {
	return ticker()
}

// Update : Key Input 과 250 millisecond 마다 입력되는 tickMsg 에 대한 행동 정의
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices.Items())-1 {
				m.cursor++
			}
		case "enter", " ":
			m.storage.ClipboardSig <- struct{}{}                                               // ClipboardSig 에 시그널 입력
			m.selector = m.cursor                                                              // "copied" 라는 텍스트 출력을 위해 selector 위치를 현재 커서 위치로 설정
			m.choices.Select(m.cursor)                                                         // m.choices 에 마운트된 클립보드 리스트에서 커서가 위치한 아이템을 선택
			clipboard.Write(clipboard.FmtText, []byte(m.choices.SelectedItem().FilterValue())) // 클립보드에 작성
			_ = m.storage.Select(m.cursor)                                                     // storge 서비스에 마운트된 container (클립보드에 저장된 데이터들의 리스트) 의 순서를 변경
			m.choices.SetItems(m.storage.ToBubbleList())                                       // 순서를 변경한 storage 서비스 container 를 다시 bubbletea 의 리스트로 변경하여 렌더링함
			m.cursor = 0                                                                       // cursor 를 복사된 값이 위치한 최상단으로 변경
			return m, nil
		}
	case tickMsg:
		m.choices.SetItems(m.storage.ToBubbleList())
		m.selector = -1
		return m, ticker()
	}
	return m, nil
}

// View : terminal 에 렌더링되는 실제 문자열
func (m Model) View() string {
	s := ""
	for i, choice := range m.choices.Items() {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		checked := "      "
		if m.cursor == m.selector && m.selector == i {
			checked = "copied"
		}

		choiceItem := choice.(tui.Item)
		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choiceItem.Display)
	}

	return s
}

func main() {
	s := file.Storage{}
	err := s.Init("./storage/data")
	if err != nil {
		return
	}
	defer s.Close()

	p := tea.NewProgram(initialModel(&s))
	if _, err := p.Run(); err != nil {
		os.Exit(1)
	}
}
