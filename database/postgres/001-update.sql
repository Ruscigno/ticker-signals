ALTER TABLE tickerbeats.accountsinfo ADD TimeTradeServer     timestamp(0) NULL;
ALTER TABLE tickerbeats.accountsinfo ADD TimeCurrent         timestamp(0) NULL;
ALTER TABLE tickerbeats.accountsinfo ADD TimeLocal           timestamp(0) NULL;
ALTER TABLE tickerbeats.accountsinfo ADD TimeGMT             timestamp(0) NULL;
ALTER TABLE tickerbeats.accountsinfo ADD LocalTimeGMTOffset  int4 NULL;
ALTER TABLE tickerbeats.accountsinfo ADD ServerTimeGMTOffset int4 NULL;

update tickerbeats.accountsinfo set TimeGMT = created;

ALTER TABLE tickerbeats.accountsinfo DROP COLUMN created;
ALTER TABLE tickerbeats.accountsinfo DROP COLUMN updated;
ALTER TABLE tickerbeats.accountsinfo DROP COLUMN deleted;