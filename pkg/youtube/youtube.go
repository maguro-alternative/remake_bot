package youtube

import (
	"context"
	"encoding/json"
	"io"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/youtube/v3"
)

type YoutubeAPI struct {
	AccessToken  string `json:"access_token"`
	TokenExpiry  string `json:"token_expiry"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RefreshToken string `json:"refresh_token"`
	ProjectID    string `json:"project_id"`
}

func NewYoutubeAPI(
	accessToken,
	tokenExpiry,
	clientID,
	clientSecret,
	refreshToken,
	projectID string,
) YoutubeAPI {
	return YoutubeAPI{
		AccessToken:  accessToken,
		TokenExpiry:  tokenExpiry,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RefreshToken: refreshToken,
		ProjectID:    projectID,
	}
}

func (y *YoutubeAPI) VideoUpload(
	ctx context.Context,
	video io.Reader,
	title, description, categoryID, privacyStatus string,
	tags []string,
) (string, error) {
	// client_secret.jsonとoauth2.jsonのパスを設定
	b, err := createClientSecret(y.ClientID, y.ProjectID, y.ClientSecret)
	if err != nil {
		return "", err
	}

	// OAuth2クライアント作成
	config, err := google.ConfigFromJSON(b, youtube.YoutubeUploadScope)
	if err != nil {
		return "", err
	}
	token, err := getToken(*y)
	if err != nil {
		return "", err
	}
	client := config.Client(ctx, token)

	// YouTube APIサービス作成
	service, err := youtube.New(client)
	if err != nil {
		return "", err
	}

	// 動画アップロード
	upload := &youtube.Video{
		Snippet: &youtube.VideoSnippet{
			Title:       title,
			Description: description,
			CategoryId:  categoryID,
		},
		Status: &youtube.VideoStatus{PrivacyStatus: privacyStatus},
	}

	// APIは、tagsが空文字列の場合、400 Bad Requestレスポンスを返す。
	upload.Snippet.Tags = tags

	call := service.Videos.Insert([]string{"snippet", "status"}, upload)

	response, err := call.Media(video).Do()
	return response.Id, err
}

func getToken(y YoutubeAPI) (*oauth2.Token, error) {
	f, err := createOAuth2(
		y.AccessToken,
		y.ClientID,
		y.ClientSecret,
		y.RefreshToken,
		y.TokenExpiry,
	)
	if err != nil {
		return nil, err
	}
	t := &oauth2.Token{}
	err = json.Unmarshal(f, t)
	return t, err
}
