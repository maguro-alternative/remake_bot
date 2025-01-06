package lineworks

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type lineWorksMemberList struct {
	Members []string `json:"members"`
	ResponseMetaData struct {
		NextCursor string `json:"nextCursor"`
	} `json:"responseMetaData"`
}

func (l *lineWorks) GetMembers(ctx context.Context) (*lineWorksMemberList, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://www.worksapis.com/bots/%s/channels/%s/messages", l.lineWorksBotID, l.lineWorksGroupID), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+l.lineWorksToken)

	res, err := l.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var memberList lineWorksMemberList
	if err := json.NewDecoder(res.Body).Decode(&memberList); err != nil {
		return nil, err
	}

	return &memberList, nil
}
