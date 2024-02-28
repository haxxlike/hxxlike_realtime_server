package lua

type StorageRead struct {
	Collection string `json:"collection"`
	Key        string `json:"key"`
	UserID     string `json:"user_id"`
}

type WalletUpdate struct {
	UserID       string                 `json:"user_id"`
	Changeset    map[string]int64       `json:"changeset"`
	Metadata     map[string]interface{} `json:"metadata"`
	UpdateLedger bool                   `json:"updateLedger,omitempty"`
}

type NotificationSend struct {
	UserID     string                 `json:"user_id"`
	Subject    string                 `json:"subject"`
	Content    map[string]interface{} `json:"content"`
	Code       int                    `json:"code"`
	Sender     string                 `json:"sender,omitempty"`
	Persistent bool                   `json:"persistent,omitempty"`
}

type LeaderboardRecordWrite struct {
	ID               string                 `json:"id"`
	OwnerID          string                 `json:"owner"`
	Username         string                 `json:"username"`
	Score            int64                  `json:"score"`
	Subscore         int64                  `json:"subscore"`
	Metadata         map[string]interface{} `json:"metadata"`
	OverrideOperator int                    `json:"override_operator"`
}

type ProfileInfo struct {
	Username string `json:"username"`
}
