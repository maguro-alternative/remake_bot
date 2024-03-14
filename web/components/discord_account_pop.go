package components

import (
	"fmt"
	"strings"

	"github.com/maguro-alternative/remake_bot/web/shared/session/model"
)

func CreateDiscordAccountPop(discordUser model.DiscordUser) string {
	discordAccountVer := strings.Builder{}
	discordAccountVer.WriteString(fmt.Sprintf(`
	<p>Discordアカウント: %s</p>
	<img src="https://cdn.discordapp.com/avatars/%s/%s.webp?size=64" alt="Discordアイコン">
	<button type="button" id="popover-btn" class="btn btn-primary">
		<a href="/logout/discord" class="btn btn-primary">ログアウト</a>
	</button>
	`, discordUser.Username, discordUser.ID, discordUser.Avatar))
	if discordUser.ID == "" {
		discordAccountVer.Reset()
		discordAccountVer.WriteString(`
		<p>Discordアカウント: 未ログイン</p>
		<button type="button" id="popover-btn" class="btn btn-primary">
			<a href="/login/discord" class="btn btn-primary">ログイン</a>
		</button>
		`)
	}
	return discordAccountVer.String()
}
