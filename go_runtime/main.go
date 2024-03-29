package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/haxxlike/hxxlike_realtime_server/lua"
	"github.com/heroiclabs/nakama-common/api"
	"github.com/heroiclabs/nakama-common/runtime"
)

type MoodleUser struct {
	Id       int    `json:"id"`
	UserName string `json:"username"`
}

// noinspection GoUnusedExportedFunction
func encode(v any, logger runtime.Logger, evtCallback any) string {
	var encoded []byte
	var err error
	ev := map[string]interface{}{
		"result":        v,
		"callback_info": evtCallback,
	}
	if encoded, err = json.Marshal(ev); err != nil {
		logger.Error("AccountGetId Handler: %v", err)
		return err.Error()
	}
	logger.Debug("encode: %s", string(encoded))
	return string(encoded)
}

func eventError(ctx context.Context, eventName string, evtMatchId string, nk runtime.NakamaModule, err error, logger runtime.Logger) {
	logger.Error("ERROR EVENT %s: %v", eventName, err)
	nk.MatchSignal(ctx, evtMatchId, "{}")
}

func verifyAndParseJwt(secretKey string, jwt string) error {
	return nil
}

func callCrawler(lambda_URL string) {
	// func callCrawler(lambda_URL string) ([]map[string]interface{}, error) {
	req, err := http.NewRequest("GET", lambda_URL, nil)
	if err != nil {
		fmt.Printf("error!!!!!!!!!!!!!!!!!!!!!1: %v\n", req)
		panic(err)
	}
	// querys := req.URL.Query()
	// querys.Add("wstoken", MOODLE_WEB_SERVICE_TOKEN)
	// querys.Add("moodlewsrestformat", MOODLE_WEB_SERVICE_REST_FORMAT)
	// for k, v := range params {
	// 	querys.Add(k, v.(string))
	// }
	// req.URL.RawQuery = querys.Encode()
	fmt.Printf("req: %v\n", req.URL.Query())
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	fmt.Printf("response Status: %v\n", resp)
	// b, _ := io.ReadAll(resp.Body)
	// body := []map[string]interface{}{}
	// if err := json.Unmarshal(b, &body); err != nil {
	// 	fmt.Printf("error!!!!!!!!!!!!!!!!!!!!!2: %v\n", body)
	// 	panic(err)
	// }
	// fmt.Printf("response Status: %q\n", body)
	// defer resp.Body.Close()
	// return body, nil
}

func InitModule(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, initializer runtime.Initializer) error {
	// logger.Info("Init local cache %v", utils.Cache)
	// utils.Cache.Test(logger)
	env := ctx.Value(runtime.RUNTIME_CTX_ENV).(map[string]string)
	var (
		LAMBDA_URL = env["LAMBDA_URL"]
	)
	if err := initializer.RegisterEvent(func(ctx context.Context, logger runtime.Logger, evt *api.Event) {
		var response any
		var previous any
		var evterror error
		var evtCallback string
		evtProperties := evt.GetProperties()
		evtName := evt.GetName()
		evtCallback = evtProperties["callback_info"]
		evtMatchId := evtProperties["match_id"]
		logger.Info("EVENT =================================================  %s", evtName)
		switch evtName {
		case "account_get_id":
			if response, evterror = nk.AccountGetId(ctx, evtProperties["user_id"]); evterror != nil {
				eventError(ctx, evtName, evtMatchId, nk, evterror, logger)
				return
			}
		case "accounts_get_id":
			var userIDs []string
			_ = json.Unmarshal([]byte(evtProperties["accounts_get_id"]), &userIDs)
			if response, evterror = nk.AccountsGetId(ctx, userIDs); evterror != nil {
				eventError(ctx, evtName, evtMatchId, nk, evterror, logger)
				return
			}
		case "storage_read":
			var reads []*lua.StorageRead
			var converted []*runtime.StorageRead
			logger.Info("evt_properties: %+v", evtProperties["object_ids"])
			_ = json.Unmarshal([]byte(evtProperties["object_ids"]), &reads)
			for _, r := range reads {
				converted = append(converted, (*runtime.StorageRead)(r))
			}
			logger.Info("storage_read: %+v", converted[0])
			if response, evterror = nk.StorageRead(ctx, converted); evterror != nil {
				eventError(ctx, evtName, evtMatchId, nk, evterror, logger)
				return
			}
		case "wallet_update":
			var wallet *lua.WalletUpdate
			_ = json.Unmarshal([]byte(evtProperties["wallet_update"]), &wallet)
			logger.Info("storage_read: %+v", wallet)
			if response, previous, evterror = nk.WalletUpdate(ctx, wallet.UserID, wallet.Changeset, wallet.Metadata, wallet.UpdateLedger); evterror != nil {
				eventError(ctx, evtName, evtMatchId, nk, evterror, logger)
				return
			}
			response = map[string]interface{}{
				"previous": previous,
				"updated":  response,
			}
		case "leaderboard_record_write":
			var record *lua.LeaderboardRecordWrite
			_ = json.Unmarshal([]byte(evtProperties["leaderboard_record_write"]), &record)
			if response, evterror = nk.LeaderboardRecordWrite(ctx, record.ID, record.OwnerID, record.Username, record.Score, record.Subscore, record.Metadata, &record.OverrideOperator); evterror != nil {
				eventError(ctx, evtName, evtMatchId, nk, evterror, logger)
				return
			}
		case "get_profile_info":
			var profileInfo *lua.ProfileInfo
			_ = json.Unmarshal([]byte(evtProperties["get_profile_info"]), &profileInfo)
			logger.Info("profileInfo: %+v", profileInfo)
			logger.Info("env==============================%+v", env)
			logger.Info("LAMBDA_URL==============================%+v", LAMBDA_URL)
			//////////////////////////////////////////////////////////////////////////////
			//call lambda crawler
			callCrawler(LAMBDA_URL)
			//////////////////////////////////////////////////////////////////////////////
			//functions do not return match signal
		case "select_update_account_id":
			decoded := map[string]interface{}{}
			userId := evtProperties["user_id"]
			userName := evtProperties["username"]
			displayName := evtProperties["display_name"]
			if response, evterror = nk.AccountGetId(ctx, userId); evterror != nil {
				eventError(ctx, evtName, evtMatchId, nk, evterror, logger)
				return
			}
			var metadata = response.(*api.Account).User.Metadata
			_ = json.Unmarshal([]byte(metadata), &decoded)
			decoded["last_login"] = time.Now().Unix() + 32400
			_ = nk.AccountUpdateId(ctx, userId, userName, decoded, displayName, "", "", "", "")
			return
		case "notification_send":
			var noti *lua.NotificationSend
			_ = json.Unmarshal([]byte(evtProperties["notification_send"]), &noti)
			logger.Info("notification_send: %+v", noti)
			nk.NotificationSend(ctx, noti.UserID, noti.Subject, noti.Content, noti.Code, noti.Sender, noti.Persistent)
			return
		default:
			logger.Error("unrecognized evt: %+v", evt)
			eventError(ctx, evtName, evtMatchId, nk, evterror, logger)
			return
		}
		nk.MatchSignal(ctx, evtMatchId, encode(response, logger, evtCallback))
	}); err != nil {
		return err
	}

	logger.Info("Server loaded.")
	return nil
}
