package lineworks

import "net/http"

type lineWorks struct {
    client http.Client
    lineWorksToken string
    lineWorksRefreshToken string
    lineWorksBotID string
    lineWorksGroupID string
}

func NewLineWorks(client http.Client, lineWorksToken, lineWorksRefreshToken, lineWorksBotID, lineWorksGroupID string) *lineWorks {
    return &lineWorks{
        client: client,
        lineWorksToken: lineWorksToken,
        lineWorksRefreshToken: lineWorksRefreshToken,
        lineWorksBotID: lineWorksBotID,
        lineWorksGroupID: lineWorksGroupID,
    }
}

type LineWorksInfo struct {
    client http.Client
    lineWorksClientID string
    lineWorksClientSecret string
    lineWorksServiceAccount string
    lineWorksPrivateKey string
    lineWorksDomainID string
    lineWorksAdminID string
}

func NewLineWorksInfo(client http.Client, lineWorksClientID, lineWorksClientSecret, lineWorksServiceAccount, lineWorksPrivateKey, lineWorksDomainID, lineWorksAdminID string) *LineWorksInfo {
    return &LineWorksInfo{
        client: client,
        lineWorksClientID: lineWorksClientID,
        lineWorksClientSecret: lineWorksClientSecret,
        lineWorksServiceAccount: lineWorksServiceAccount,
        lineWorksPrivateKey: lineWorksPrivateKey,
        lineWorksDomainID: lineWorksDomainID,
        lineWorksAdminID: lineWorksAdminID,
    }
}
