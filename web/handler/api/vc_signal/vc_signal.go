package vcsignal

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/maguro-alternative/remake_bot/repository"

	"github.com/maguro-alternative/remake_bot/web/handler/api/vc_signal/internal"
)

type VcSignalHandler struct {
	repo repository.Repository
}

func NewVcSignalHandler(
	repo repository.Repository,
) *VcSignalHandler {
	return &VcSignalHandler{
		repo: repo,
	}
}

func (h *VcSignalHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		slog.ErrorContext(ctx, "/api/vc-signal Method Not Allowed")
		return
	}
	var vcSignalJson internal.VcSignalJson
	var vcSignalChannelNotGuildId repository.VcSignalChannelNotGuildID

	err := json.NewDecoder(r.Body).Decode(&vcSignalJson)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		slog.ErrorContext(ctx, "jsonの読み取りに失敗しました:", "エラー:", err.Error())
		return
	}
	err = vcSignalJson.Validate()
	if err != nil {
		http.Error(w, "Unprocessable Entity", http.StatusUnprocessableEntity)
		slog.ErrorContext(ctx, "jsonのバリデーションに失敗しました:", "エラー:", err.Error())
		return
	}

	for _, vcSignal := range vcSignalJson.VcSignals {
		vcSignalChannelNotGuildId.VcChannelID = vcSignal.VcChannelID
		vcSignalChannelNotGuildId.SendSignal = vcSignal.SendSignal
		vcSignalChannelNotGuildId.SendChannelID = vcSignal.SendChannelId
		vcSignalChannelNotGuildId.JoinBot = vcSignal.JoinBot
		vcSignalChannelNotGuildId.EveryoneMention = vcSignal.EveryoneMention
		err = h.repo.UpdateVcSignalChannel(ctx, vcSignalChannelNotGuildId)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			slog.ErrorContext(ctx, "VcSignalChannelの更新に失敗しました。", "エラー:", err.Error())
			return
		}
		for _, userId := range vcSignal.VcSignalNgUserIDs {
			err = h.repo.InsertVcSignalNgUser(ctx, vcSignal.VcChannelID, vcSignalJson.GuildID, userId)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				slog.ErrorContext(ctx, "NgUserIDの追加に失敗しました。", "エラー:", err.Error())
				return
			}
		}
		err = h.repo.DeleteVcSignalNgUsersNotInProvidedList(ctx, vcSignal.VcChannelID, vcSignal.VcSignalNgUserIDs)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			slog.ErrorContext(ctx, "NgUserIDの削除に失敗しました。", "エラー:", err.Error())
			return
		}
		for _, roleId := range vcSignal.VcSignalNgRoleIDs {
			err = h.repo.InsertVcSignalNgRole(ctx, vcSignal.VcChannelID, vcSignalJson.GuildID, roleId)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				slog.ErrorContext(ctx, "NgRoleIDの追加に失敗しました。", "エラー:", err.Error())
				return
			}
		}
		err = h.repo.DeleteVcSignalNgRolesNotInProvidedList(ctx, vcSignal.VcChannelID, vcSignal.VcSignalNgRoleIDs)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			slog.ErrorContext(ctx, "NgRoleIDの削除に失敗しました。", "エラー:", err.Error())
			return
		}
		for _, userId := range vcSignal.VcSignalMentionUserIDs {
			err = h.repo.InsertVcSignalMentionUser(ctx, vcSignal.VcChannelID, vcSignalJson.GuildID, userId)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				slog.ErrorContext(ctx, "MentionUserIDの追加に失敗しました。", "エラー:", err.Error())
				return
			}
		}
		err = h.repo.DeleteVcSignalMentionUsersNotInProvidedList(ctx, vcSignal.VcChannelID, vcSignal.VcSignalMentionUserIDs)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			slog.ErrorContext(ctx, "MentionUserIDの削除に失敗しました。", "エラー:", err.Error())
			return
		}
		for _, roleId := range vcSignal.VcSignalMentionRoleIDs {
			err = h.repo.InsertVcSignalMentionRole(ctx, vcSignal.VcChannelID, vcSignalJson.GuildID, roleId)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				slog.ErrorContext(ctx, "MentionRoleIDの追加に失敗しました。", "エラー:", err.Error())
				return
			}
		}
		err = h.repo.DeleteVcSignalMentionRolesNotInProvidedList(ctx, vcSignal.VcChannelID, vcSignal.VcSignalMentionRoleIDs)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			slog.ErrorContext(ctx, "MentionRoleIDの削除に失敗しました。", "エラー:", err.Error())
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}
