CREATE INDEX traderequests_accountid_idx ON tickerbeats.traderequests USING btree (accountid, deleted, created);