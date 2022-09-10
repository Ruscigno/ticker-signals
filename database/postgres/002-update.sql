ALTER TABLE tickerbeats.accountsinfo DROP CONSTRAINT accountsinfo_pk;
ALTER TABLE tickerbeats.accountsinfo ALTER COLUMN infoid DROP NOT NULL;
update tickerbeats.accountsinfo set infoid = null;
ALTER TABLE tickerbeats.accountsinfo ALTER COLUMN infoid TYPE int8 USING infoid::int8;

create table tickerbeats.accountsinfo2 as
SELECT accountid, ROW_NUMBER () OVER (
           ORDER BY accountid, timeGMT
        ) as infoid, balance, credit, profit, equity, margin, freemargin, marginlevel, margincall, marginstopout, timetradeserver, timecurrent, timelocal, timeGMT, localtimeGMToffset, servertimeGMToffset
FROM tickerbeats.accountsinfo;

delete from tickerbeats.accountsinfo where 2 > 1;

INSERT INTO tickerbeats.accountsinfo
(accountid, infoid, balance, credit, profit, equity, margin, freemargin, marginlevel, margincall, marginstopout, timetradeserver, timecurrent, timelocal, timeGMT, localtimeGMToffset, servertimeGMToffset)
select accountid, infoid, balance, credit, profit, equity, margin, freemargin, marginlevel, margincall, marginstopout, timetradeserver, timecurrent, timelocal, timeGMT, localtimeGMToffset, servertimeGMToffset
from tickerbeats.accountsinfo2;


ALTER TABLE tickerbeats.accountsinfo ALTER COLUMN infoid SET NOT NULL;
ALTER TABLE tickerbeats.accountsinfo ADD CONSTRAINT accountsinfo_pk PRIMARY KEY (accountid,infoid);
drop table tickerbeats.accountsinfo2;

-- SET THE PROPER STARTING VALUE
CREATE SEQUENCE tickerbeats.accountsinfo_seq
	INCREMENT BY 1
	MINVALUE 1
	MAXVALUE 9223372036854775807
	START 
	CACHE 1
	NO CYCLE;

ALTER TABLE tickerbeats.accountsinfo ALTER COLUMN timeGMT SET NOT NULL;
CREATE UNIQUE INDEX accountsinfo_accountid_idx ON tickerbeats.accountsinfo (accountid,timeGMT);