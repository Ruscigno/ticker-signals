CREATE SEQUENCE tickerbeats.accountsinfo_seq;
CREATE SEQUENCE tickerbeats.tickerbeats_seq;
CREATE SEQUENCE tickerbeats.tradetransactions_seq;

-- DROP TABLE tickerbeats.accounts;

CREATE TABLE tickerbeats.accounts (
	accountid int8 NOT NULL,
    description varchar (255),
	trademode varchar(50) NOT NULL,
	leverage int8 NOT NULL,
	marginmode varchar(50) NULL,
	stopoutmode varchar(50) NULL,
	tradeallowed bool NULL,
	tradeexpert bool NULL,
	limitorders int8 NULL,
	"name" varchar(255) NOT NULL,
	"server" varchar(100) NOT NULL,
	currency varchar(50) NOT NULL,
	company varchar(100) NULL,
	created timestamp(0) NOT NULL,
	updated timestamp(0) NOT NULL,
	deleted timestamp(0) NULL,
	CONSTRAINT accounts_pk PRIMARY KEY (accountid)
);
commit;
-- DROP TABLE tickerbeats.accountsinfo;

CREATE TABLE tickerbeats.accountsinfo (
	accountid int8 NOT NULL,
	infoid int8 NOT NULL,
	balance numeric(30,5) NULL,
	credit numeric(30,5) NULL,
	profit numeric(30,5) NULL,
	equity numeric(30,5) NULL,
	margin numeric(30,5) NULL,
	freemargin numeric(30,5) NULL,
	marginlevel numeric(30,5) NULL,
	margincall numeric(30,5) NULL,
	marginstopout numeric(30,5) NULL,
	timetradeserver timestamp(0) NULL,
	timecurrent timestamp(0) NULL,
	timelocal timestamp(0) NULL,
	timegmt timestamp(0) NOT NULL,
	localtimegmtoffset int4 NULL,
	servertimegmtoffset int4 NULL,
	CONSTRAINT accountsinfo_pk PRIMARY KEY (accountid, infoid),
	CONSTRAINT accountsinfo_fk FOREIGN KEY (accountid) REFERENCES tickerbeats.accounts(accountid)
);
CREATE UNIQUE INDEX accountsinfo_accountid_idx ON tickerbeats.accountsinfo USING btree (accountid, timegmt);
commit;
-- DROP TABLE tickerbeats.deals;

CREATE TABLE tickerbeats.deals (
	accountid int8 NOT NULL,
	symbol varchar(50) NOT NULL,
	ticket int8 NOT NULL,
	positionid int8 NOT NULL,
	magic int8 NULL,
	orderid int8 NULL,
	dealtime timestamp(0) NOT NULL,
	dealtype varchar(50) NULL,
	entry varchar(50) NULL,
	volume numeric(30,5) NULL,
	price numeric(30,5) NULL,
	commission numeric(30,5) NULL,
	swap numeric(30,5) NULL,
	profit numeric(30,5) NULL,
	"comment" varchar(255) NULL,
	externalid varchar(255) NULL,
	dealid int8 NOT NULL,
	created timestamp(0) NOT NULL,
	updated timestamp(0) NOT NULL,
	deleted timestamp(0) NULL,
	reason int8 NULL,
	dealfee numeric(30,5) NULL,
	CONSTRAINT deals_pk PRIMARY KEY (accountid, ticket),
	CONSTRAINT deals_fk FOREIGN KEY (accountid) REFERENCES tickerbeats.accounts(accountid)
);
commit;
-- DROP TABLE tickerbeats.orders;

CREATE TABLE tickerbeats.orders (
	accountid int8 NOT NULL,
	symbol varchar(50) NOT NULL,
	ticket int8 NOT NULL,
	positionid int8 NULL,
	timesetup timestamp(0) NOT NULL,
	ordertype varchar(50) NULL,
	state varchar(50) NULL,
	timeexpiration timestamp(0) NULL,
	timedone timestamp(0) NOT NULL,
	typefilling varchar(50) NULL,
	typetime varchar(50) NULL,
	magic int8 NULL,
	reason varchar(50) NULL,
	volumeinitial numeric(30,5) NULL,
	volumecurrent numeric(30,5) NULL,
	priceopen numeric(30,5) NULL,
	stoploss numeric(30,5) NULL,
	takeprofit numeric(30,5) NULL,
	pricecurrent numeric(30,5) NULL,
	pricestoplimit numeric(30,5) NULL,
	"comment" varchar(255) NULL,
	externalid varchar(255) NULL,
	orderid int8 NOT NULL,
	created timestamp(0) NOT NULL,
	updated timestamp(0) NOT NULL,
	deleted timestamp(0) NULL,
	positionbyid int8 NULL,
	CONSTRAINT orders_pk PRIMARY KEY (accountid, ticket),
	CONSTRAINT orders_fk FOREIGN KEY (accountid) REFERENCES tickerbeats.accounts(accountid)
);

commit;
-- DROP TABLE tickerbeats.positions;

CREATE TABLE tickerbeats.positions (
	accountid int8 NOT NULL,
	positionid int8 NOT NULL,
	ticket int8 NOT NULL,
	symbol varchar(50) NOT NULL,
	positiontime timestamp(0) NULL,
	positiontype varchar(50) NULL,
	volume numeric(30,5) NULL,
	priceopen numeric(30,5) NULL,
	stoploss numeric(30,5) NULL,
	takeprofit numeric(30,5) NULL,
	pricecurrent numeric(30,5) NULL,
	commission numeric(30,5) NULL,
	swap numeric(30,5) NULL,
	profit numeric(30,5) NULL,
	"comment" varchar(255) NULL,
	created timestamp(0) NOT NULL,
	updated timestamp(0) NOT NULL,
	deleted timestamp(0) NULL,
	reason varchar(50) NOT NULL,
	externalid varchar(255) NULL,
	magic int8 NULL,
	positionupdate timestamp(0) NULL,
	CONSTRAINT positions_pk PRIMARY KEY (accountid, ticket),
	CONSTRAINT positions_fk FOREIGN KEY (accountid) REFERENCES tickerbeats.accounts(accountid) ON UPDATE RESTRICT ON DELETE RESTRICT
);
commit;
-- DROP TABLE tickerbeats.signals;

CREATE TABLE tickerbeats.signals (
	signalid int8 NOT NULL,
	sourceaccountid int8 NOT NULL,
	destinationaccountid int8 NOT NULL,
	active bool NOT NULL,
	maxdepositpercent int4 NOT NULL,
	stopiflessthan int8 NOT NULL,
	maxspread int8 NOT NULL,
	minutestoexpire int4 NOT NULL,
	orderboost numeric(30, 5) NULL,
	orderboosttype int4 NULL,
	CONSTRAINT signals_pk PRIMARY KEY (signalid),
	CONSTRAINT signals_un UNIQUE (sourceaccountid, destinationaccountid),
	CONSTRAINT signals_account_fk1 FOREIGN KEY (sourceaccountid) REFERENCES tickerbeats.accounts(accountid) ON DELETE RESTRICT ON UPDATE RESTRICT,
	CONSTRAINT signals_account_fk2 FOREIGN KEY (destinationaccountid) REFERENCES tickerbeats.accounts(accountid) ON DELETE RESTRICT ON UPDATE RESTRICT
);
commit;
-- DROP TABLE tickerbeats.signalsresult;

CREATE TABLE tickerbeats.signalsresult (
	sourceaccountid int8 NOT NULL,
	destinationaccountid int8 NOT NULL,
	sourcebeatsid int8 NULL,
	destinationbeatsid int8 NULL,
	sourcepositionid int8 NOT NULL,
	destinationpositionid int8 NULL,
	signaltype varchar(30) NOT NULL,
	signalstatus int2 NOT NULL,
	externalid int8 NOT NULL,
	groupid varchar(50) NOT NULL,
	entry varchar(50) NOT NULL,
	senttime timestamp(0) NULL,
	confirmationtime timestamp(0) NULL,
	expireat timestamp(0) NULL,
	created timestamp(0) NOT NULL,
	updated timestamp(0) NOT NULL,
	CONSTRAINT signalsresult_externalid_un UNIQUE (destinationaccountid, externalid),
	CONSTRAINT signalsresult_pk PRIMARY KEY (externalid),
	CONSTRAINT signalsresult_account_destination_fk FOREIGN KEY (destinationaccountid) REFERENCES tickerbeats.accounts(accountid) ON DELETE RESTRICT ON UPDATE RESTRICT,
	CONSTRAINT signalsresult_account_source_fk FOREIGN KEY (sourceaccountid) REFERENCES tickerbeats.accounts(accountid) ON DELETE RESTRICT ON UPDATE RESTRICT,
	CONSTRAINT signalsresult_fk FOREIGN KEY (sourceaccountid,destinationaccountid) REFERENCES tickerbeats.signals(sourceaccountid,destinationaccountid) ON DELETE RESTRICT ON UPDATE RESTRICT,
	CONSTRAINT signalsresult_source_position_fk FOREIGN KEY (sourceaccountid,sourcepositionid) REFERENCES tickerbeats.positions(accountid,ticket) ON DELETE RESTRICT ON UPDATE RESTRICT
);
CREATE INDEX signalsresult_groupid_idx ON tickerbeats.signalsresult USING btree (sourceaccountid, destinationaccountid, groupid);
CREATE UNIQUE INDEX signalsresult_sourceaccountid_idx ON tickerbeats.signalsresult USING btree (sourceaccountid, destinationaccountid, sourcepositionid, destinationpositionid, signaltype, entry);
commit;
-- DROP TABLE tickerbeats.traderequests;

CREATE TABLE tickerbeats.traderequests (
	accountid int8 NOT NULL,
	symbol varchar(50) NULL,
	orderid int8 NOT NULL,
	positionid int8 NULL,
	magic int8 NULL,
	creationorder int8 NOT NULL,
	"action" varchar(50) NULL,
	volume numeric(30,5) NULL,
	price numeric(30,5) NULL,
	stoplimit numeric(30,5) NULL,
	stoploss numeric(30,5) NULL,
	takeprofit numeric(30,5) NULL,
	deviation int8 NULL,
	ordertype varchar(50) NULL,
	typefilling varchar(50) NULL,
	typetime varchar(50) NULL,
	timeexpiration timestamp(0) NULL,
	"comment" varchar(50) NULL,
	positionby int8 NULL,
	created timestamp(0) NOT NULL,
	updated timestamp(0) NOT NULL,
	deleted timestamp(0) NULL,
	entry varchar(50) NULL,
	CONSTRAINT traderequests_pk PRIMARY KEY (accountid, orderid, creationorder),
	CONSTRAINT traderequests_accounts_fk FOREIGN KEY (accountid) REFERENCES tickerbeats.accounts(accountid) ON UPDATE RESTRICT ON DELETE RESTRICT
);
CREATE INDEX traderequests_accountid_idx ON tickerbeats.traderequests USING btree (accountid, deleted, created);
commit;
-- DROP TABLE tickerbeats.traderesults;

CREATE TABLE tickerbeats.traderesults (
	accountid int8 NOT NULL,
	orderid int8 NOT NULL,
	dealid int8 NULL,
	creationorder int8 NOT NULL,
	retcode int8 NULL,
	volume numeric(30,5) NULL,
	price numeric(30,5) NULL,
	bid numeric(30,5) NULL,
	ask numeric(30,5) NULL,
	"comment" varchar(50) NULL,
	requestid int8 NULL,
	retcodeexternal int8 NULL,
	created timestamp(0) NOT NULL,
	updated timestamp(0) NOT NULL,
	deleted timestamp(0) NULL,
	CONSTRAINT traderesults_pk PRIMARY KEY (accountid, orderid, creationorder),
    CONSTRAINT traderesults_accounts_fk FOREIGN KEY (accountid) REFERENCES tickerbeats.accounts(accountid) ON UPDATE RESTRICT ON DELETE RESTRICT
);
commit;
-- DROP TABLE tickerbeats.tradetransactions;

CREATE TABLE tickerbeats.tradetransactions (
	internalid int8 NOT NULL,
	accountid int8 NOT NULL,
	symbol varchar(50) NULL,
	orderid int8 NOT NULL,
	dealid int8 NULL,
	positionid int8 NULL,
	creationorder int8 NOT NULL,
	tradetype varchar(50) NULL,
	ordertype varchar(50) NULL,
	orderstate varchar(50) NULL,
	dealtype varchar(50) NULL,
	timetype varchar(50) NULL,
	timeexpiration timestamp(0) NULL,
	price numeric(30,5) NULL,
	pricetrigger numeric(30,5) NULL,
	pricestoploss numeric(30,5) NULL,
	pricetakeprofit numeric(30,5) NULL,
	volume numeric(30,5) NULL,
	positionby int8 NULL,
	created timestamp(0) NOT NULL,
	updated timestamp(0) NOT NULL,
	deleted timestamp(0) NULL,
	CONSTRAINT tradetransactions_pk PRIMARY KEY (accountid, orderid, creationorder),
    CONSTRAINT tradetransactions_accounts_fk FOREIGN KEY (accountid) REFERENCES tickerbeats.accounts(accountid) ON UPDATE RESTRICT ON DELETE RESTRICT
);
commit;
-- DROP TABLE tickerbeats.traderules;

CREATE TABLE tickerbeats.traderules (
	accountid int8 NOT NULL,
	ruleid int2 NOT NULL,
    rulepriority int2 NOT NULL,
	active bool NOT NULL,
	description varchar(254) NULL,
	symbol varchar(50) NULL,
	ruletype int2 NOT NULL,
    ruleversion int2 NOT NULL,
	rule varchar(5000) NOT NULL,
	CONSTRAINT traderules_pk PRIMARY KEY (accountid, ruleid),
	CONSTRAINT traderules_accounts_fk FOREIGN KEY (accountid) REFERENCES tickerbeats.accounts(accountid) ON DELETE RESTRICT ON UPDATE RESTRICT
);
CREATE INDEX traderules_accountid_idx ON tickerbeats.traderules USING btree (accountid, active, symbol);
commit;