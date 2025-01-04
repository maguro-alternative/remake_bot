package lineworks

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type lineWorksUserProfile struct {
	UserID   string `json:"userId"`
	Email	string  `json:"email"`
	UserName struct {
		LastName string `json:"lastName"`
		FirstName string `json:"firstName"`
		PhoneticLastName string `json:"phoneticLastName"`
		PhoneticFirstName string `json:"phoneticFirstName"`
	} `json:"userName"`
	I18nNames []struct {
		Language string `json:"language"`
		FirstName string `json:"firstName"`
		LastName string `json:"lastName"`
	}
	Organizations []struct {
		DomainId int32 `json:"domainId"`
		Primary bool `json:"primary"`
		UserExternalKey string `json:"userExternalKey"`
		Email string `json:"email"`
		LevelId string `json:"levelId"`
		LevelExternalKey string `json:"levelExternalKey"`
		LevelName string `json:"levelName"`
		Executive bool `json:"executive"`
		OrganizationName string `json:"organizationName"`
		OrgUnits []struct {
			OrgUnitId string `json:"orgUnitId"`
			OrgUnitExternalKey string `json:"orgUnitExternalKey"`
			OrgUnitName string `json:"orgUnitName"`
			OrgUnitEmail string `json:"orgUnitEmail"`
			Primary bool `json:"primary"`
			PositionId string `json:"positionId"`
			PositionExternalKey string `json:"positionExternalKey"`
			PositionName string `json:"positionName"`
			IsManager bool `json:"isManager"`
			Visible bool `json:"visible"`
			UseTeamFeature bool `json:"useTeamFeature"`
		}
	} `json:"organizations"`
	Telephone string `json:"telephone"`
	CellPhone string `json:"cellPhone"`
	Location string `json:"location"`
}

func (l *lineWorks) GetUserProfile(ctx context.Context, userId string) (*lineWorksUserProfile, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://www.worksapis.com/v1.0/users/%s", userId), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+l.lineWorksToken)

	res, err := l.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var userProfile lineWorksUserProfile
	if err := json.NewDecoder(res.Body).Decode(&userProfile); err != nil {
		return nil, err
	}

	return &userProfile, nil
}
