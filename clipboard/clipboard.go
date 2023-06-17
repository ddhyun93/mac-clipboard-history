package clipboard

import (
	"context"
	"golang.design/x/clipboard"
	"myclipboard/storage/file"
	"time"
)

func InitClipboard(container *file.Storage) {
	_ = clipboard.Init()

	go func() {
	ENTRYPOINT:
		for {
			ctx, cancel := context.WithCancel(context.Background())
			go watchClipboard(container, ctx)

			for {
				if val := <-container.ClipboardSig; val == struct{}{} {
					cancel()
					// enter 메시지가 수신되면 중복 입력을 방지하기 위해 클립보드를 watch 하면서 변경된 값이 있으면, container 에 작성하는 기능을 1초간 정지시킨다.
					time.Sleep(time.Second)
					continue ENTRYPOINT
				}
			}
		}
	}()
}

func watchClipboard(container *file.Storage, ctx context.Context) {
	w := clipboard.Watch(ctx, clipboard.FmtText)
	for data := range w {
		_ = container.Write(string(data))
	}
}
