package tui

type Displayed string

// Item : 저장소에서 불러온 clipboard 데이터가 저장되는 모델
// Display > UI 상 표시되는 문자열
// Value > 실제 데이터
type Item struct {
	Display Displayed
	Value   string
}

// FilterValue : Implement of bubble Item interface.
func (i Item) FilterValue() string {
	return i.Value
}
