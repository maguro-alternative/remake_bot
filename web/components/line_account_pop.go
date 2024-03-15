package components

import (
	"fmt"
	"strings"

	"github.com/maguro-alternative/remake_bot/web/shared/session/model"
)

func CreateLineAccountVer(lineUser model.LineIdTokenUser) string {
	lineAccountVer := strings.Builder{}
	lineAccountVer.WriteString(fmt.Sprintf(`
	<p>LINEアカウント: %s</p>
	<img src="%s" style="right: 70px;" alt="LINEアイコン">
	<button type="button" id="popover-btn" class="btn btn-primary">
		<a href="/logout/line" class="btn btn-primary">ログアウト</a>
	</button>
	`, lineUser.Name, lineUser.Picture))
	if lineUser.Name == "" {
		lineAccountVer.Reset()
		lineAccountVer.WriteString(`
			<p>LINEアカウント: 未ログイン</p>
			<button type="button" id="popover-btn" class="btn btn-primary">
				<a href="/login/line" class="btn btn-primary">ログイン</a>
			</button>
		`)
	}
	return lineAccountVer.String()
}