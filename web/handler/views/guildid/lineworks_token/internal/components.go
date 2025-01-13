package internal

type LineWorksBotByteEntered struct {
	LineWorksBotToken     [][]byte        `db:"line_works_bot_token"`
	LineWorksRefreshToken [][]byte        `db:"line_works_refresh_token"`
	LineWorksGroupID      [][]byte        `db:"line_works_group_id"`
	LineWorksBotID        [][]byte        `db:"line_works_bot_id"`
	LineWorksBotSecret    [][]byte        `db:"line_works_bot_secret"`
	DebugMode             bool          `db:"debug_mode"`
}

type LineWorksBotInfoEntered struct {
	LineWorksClientID	 [][]byte        `db:"line_works_client_id"`
	LineWorksClientSecret [][]byte        `db:"line_works_client_secret"`
	LineWorksServiceAccount [][]byte        `db:"line_works_service_account"`
	LineWorksPrivateKey	 [][]byte        `db:"line_works_private_key"`
	LineWorksDomainID	 [][]byte        `db:"line_works_domain_id"`
	LineWorksAdminID	 [][]byte        `db:"line_works_admin_id"`
}

type LineWorksBotInfo struct {
	LineWorksGroupID      string
	LineWorksBotID        string
	LineWorksBotSecret    string
	LineWorksClientID	 string
	LineWorksClientSecret string
	LineWorksServiceAccount string
	LineWorksPrivateKey	 string
	LineWorksDomainID	 string
	LineWorksAdminID	 string
	DebugMode             bool
}

func EnteredLineWorksBotForm(
	lineWorksBotByte LineWorksBotByteEntered,
	lineWorksBotInfoByte LineWorksBotInfoEntered,
) LineWorksBotInfo {
	var lineWorksBotInfo LineWorksBotInfo
	if lineWorksBotByte.LineWorksGroupID != nil {
		lineWorksBotInfo.LineWorksGroupID = "入力済み"
	}
	if lineWorksBotByte.LineWorksBotID != nil {
		lineWorksBotInfo.LineWorksBotID = "入力済み"
	}
	if lineWorksBotByte.LineWorksBotSecret != nil {
		lineWorksBotInfo.LineWorksBotSecret = "入力済み"
	}
	if lineWorksBotInfoByte.LineWorksClientID != nil {
		lineWorksBotInfo.LineWorksClientID = "入力済み"
	}
	if lineWorksBotInfoByte.LineWorksClientSecret != nil {
		lineWorksBotInfo.LineWorksClientSecret = "入力済み"
	}
	if lineWorksBotInfoByte.LineWorksServiceAccount != nil {
		lineWorksBotInfo.LineWorksServiceAccount = "入力済み"
	}
	if lineWorksBotInfoByte.LineWorksPrivateKey != nil {
		lineWorksBotInfo.LineWorksPrivateKey = "入力済み"
	}
	if lineWorksBotInfoByte.LineWorksDomainID != nil {
		lineWorksBotInfo.LineWorksDomainID = "入力済み"
	}
	if lineWorksBotInfoByte.LineWorksAdminID != nil {
		lineWorksBotInfo.LineWorksAdminID = "入力済み"
	}
	lineWorksBotInfo.DebugMode = lineWorksBotByte.DebugMode
	return lineWorksBotInfo
}
