package components

import (
	"fmt"
	"strings"
)

func CreateSelectChennelOptions(
	categoryIDTmps []string,
	defaultChannelID string,
	channelsInCategory map[string][]DiscordChannelSelect,
	categoryPositions map[string]DiscordChannel,
) string {
	htmlSelectChannelBuilders := strings.Builder{}
	categoryOptions := make([]strings.Builder, len(categoryIDTmps)+1)
	var categoryIndex int
	for categoryID, channels := range channelsInCategory {
		for i, categoryIDTmp := range categoryIDTmps {
			if categoryID == "" {
				categoryIndex = len(categoryIDTmps)
				break
			}
			if categoryIDTmp == categoryID {
				categoryIndex = i
				break
			}
		}
		for _, channelSelect := range channels {
			if channelSelect.ID == "" {
				continue
			}
			if defaultChannelID == channelSelect.ID {
				categoryOptions[categoryIndex].WriteString(fmt.Sprintf(`<option value="%s" selected>%s</option>`, channelSelect.ID, channelSelect.Name))
				continue
			}
			categoryOptions[categoryIndex].WriteString(fmt.Sprintf(`<option value="%s">%s</option>`, channelSelect.ID, channelSelect.Name))
		}
	}
	for _, categoryOption := range categoryOptions {
		htmlSelectChannelBuilders.WriteString(categoryOption.String())
	}
	return htmlSelectChannelBuilders.String()
}

func EnteredLineBotForm(
	lineBotByte LineBotByteEntered,
) LineEntered {
	var lineEntered LineEntered
	if lineBotByte.LineNotifyToken != nil {
		lineEntered.LineNotifyToken = "入力済み"
	}
	if lineBotByte.LineBotToken != nil {
		lineEntered.LineBotToken = "入力済み"
	}
	if lineBotByte.LineBotSecret != nil {
		lineEntered.LineBotSecret = "入力済み"
	}
	if lineBotByte.LineClientID != nil {
		lineEntered.LineClientID = "入力済み"
	}
	if lineBotByte.LineClientSecret != nil {
		lineEntered.LineClientSecret = "入力済み"
	}
	if lineBotByte.LineGroupID != nil {
		lineEntered.LineGroupID = "入力済み"
	}
	return lineEntered
}