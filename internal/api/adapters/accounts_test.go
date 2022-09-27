package adapters

// func getDefaultAccount(timeNow time.Time) *acc.Account {
// 	return &acc.Account{
// 		AccountID:   1,
// 		Name:        "name1",
// 		TradeMode:   "ACCOUNT_TRADE_MODE_DEMO",
// 		MarginMode:  "ACCOUNT_MARGIN_MODE_RETAIL_NETTING",
// 		StopoutMode: "ACCOUNT_STOPOUT_MODE_PERCENT",
// 		Created:     timeNow,
// 		Updated:     timeNow,
// 		Deleted:     timeNow,
// 	}
// }

// func getDefaultAccountProto(timeNow time.Time) *v1.Account {
// 	pTimeNow, _ := ptypes.TimestampProto(timeNow)
// 	return &v1.Account{
// 		AccountId:   1,
// 		Name:        "name1",
// 		TradeMode:   v1.AccountTradeMode_ACCOUNT_TRADE_MODE_DEMO,
// 		MarginMode:  v1.AccountMarginMode_ACCOUNT_MARGIN_MODE_RETAIL_NETTING,
// 		StopoutMode: v1.AccountStopoutMode_ACCOUNT_STOPOUT_MODE_PERCENT,
// 		Created:     pTimeNow,
// 		Updated:     pTimeNow,
// 		Deleted:     pTimeNow,
// 	}
// }

// func TestAccountToProto(t *testing.T) {
// 	timeNow, _ := ptypes.Timestamp(ptypes.TimestampNow())
// 	tests := []struct {
// 		name string
// 		a    *acc.Account
// 		want *v1.Account
// 	}{
// 		{
// 			name: "success",
// 			a:    getDefaultAccount(timeNow),
// 			want: getDefaultAccountProto(timeNow),
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := AccountToProto(tt.a); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("AccountToProto() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func TestProtoToAccount(t *testing.T) {
// 	timeNow, _ := ptypes.Timestamp(ptypes.TimestampNow())
// 	tests := []struct {
// 		name string
// 		a    *v1.Account
// 		want *acc.Account
// 	}{
// 		{
// 			name: "success",
// 			a:    getDefaultAccountProto(timeNow),
// 			want: getDefaultAccount(timeNow),
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := ProtoToAccount(tt.a); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("ProtoToAccount() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func TestAccountsToProto(t *testing.T) {
// 	timeNow, _ := ptypes.Timestamp(ptypes.TimestampNow())
// 	tests := []struct {
// 		name     string
// 		accounts []*acc.Account
// 		want     []*v1.Account
// 	}{
// 		{
// 			name:     "success",
// 			accounts: append([]*acc.Account{}, getDefaultAccount(timeNow)),
// 			want:     append([]*v1.Account{}, getDefaultAccountProto(timeNow)),
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := AccountsToProto(tt.accounts); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("AccountsToProto() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
