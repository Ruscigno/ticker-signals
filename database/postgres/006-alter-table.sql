ALTER TABLE tickerbeats.positions ALTER COLUMN reason TYPE varchar(50) USING reason::varchar;
ALTER TABLE tickerbeats.positions ALTER COLUMN reason SET NOT NULL;

update tickerbeats.positions set reason = 'POSITION_REASON_CLIENT' where reason = '0';
update tickerbeats.positions set reason = 'POSITION_REASON_MOBILE' where reason = '1';
update tickerbeats.positions set reason = 'POSITION_REASON_WEB' where reason = '2';
update tickerbeats.positions set reason = 'POSITION_REASON_EXPERT' where reason = '3';