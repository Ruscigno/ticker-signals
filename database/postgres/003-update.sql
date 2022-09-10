ALTER TABLE tickerbeats.positions ADD positionupd timestamp(0) NULL;
update tickerbeats.positions set positionupd = to_timestamp(positionupdate/1000::int8);
ALTER TABLE tickerbeats.positions DROP COLUMN positionupdate;
ALTER TABLE tickerbeats.positions ADD positionupdate timestamp(0) NULL;
update tickerbeats.positions set positionupdate = positionupd;
ALTER TABLE tickerbeats.positions DROP COLUMN positionupd;