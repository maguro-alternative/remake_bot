package internal

type LineBotByteEntered struct {
	LineNotifyToken  [][]byte
	LineBotToken     [][]byte
	LineBotSecret    [][]byte
	LineGroupID      [][]byte
	LineClientID     [][]byte
	LineClientSecret [][]byte
	LineDebugMode    bool
}

type LineEntered struct {
	LineNotifyToken  string
	LineBotToken     string
	LineBotSecret    string
	LineGroupID      string
	LineClientID     string
	LineClientSecret string
	LineDebugMode    bool
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
	lineEntered.LineDebugMode = lineBotByte.LineDebugMode
	return lineEntered
}