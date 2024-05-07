package provider

import (
	"context"
	"github.com/gogf/gf/v2/i18n/gi18n"
	"sync"
)

type Translate struct {
	manager *gi18n.Manager
}

var translateInstance *Translate
var translateOnce sync.Once

func GetTranslate() *Translate {
	translateOnce.Do(func() {
		translateInstance = &Translate{
			manager: gi18n.New(),
		}
		translateInstance.manager.SetLanguage("zh")
	})
	return translateInstance
}

func (t *Translate) SetLanguage(language string) {
	t.manager.SetLanguage(language)
}

func (t *Translate) Tr(content string) string {
	return t.manager.Translate(context.Background(), content)
}

func Tr(content string) string {
	return GetTranslate().Tr(content)
}
