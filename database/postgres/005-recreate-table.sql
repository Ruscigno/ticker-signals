DROP TABLE IF EXISTS tickerbeats.signalsresult;
CREATE TABLE tickerbeats.signalsresult (
  SourceAccountID       int8 NOT NULL,
  DestinationAccountID  int8 NOT NULL,
  SourceBeatsID         int8 NOT NULL,
  DestinationBeatsID    int8,
  SourcePositionID      int8 NOT NULL,
  DestinationPositionID int8,
  SignalType            varchar(30) not null,
  SignalStatus          smallint not null,
  ExternalID            int8 NOT NULL,
  GroupID               varchar(50) not null,
  Entry                 varchar(50) not null,
  SentTime              timestamp(0),
  ConfirmationTime      timestamp(0),
  ExpireAt              timestamp(0),
  Created               timestamp(0) NOT NULL,
  Updated               timestamp(0) NOT NULL,
  CONSTRAINT signalsresult_pk                      PRIMARY KEY (SourceAccountID, DestinationAccountID, SourceBeatsID, SignalType),
  CONSTRAINT signalsresult_fk                      FOREIGN KEY (SourceAccountID, DestinationAccountID) REFERENCES tickerbeats.signals(SourceAccountID, DestinationAccountID) ON UPDATE RESTRICT ON DELETE restrict,
  CONSTRAINT signalsresult_externalid_un           UNIQUE      (destinationaccountid,externalid),
  CONSTRAINT signalsresult_account_source_fk       FOREIGN KEY (SourceAccountID)                       REFERENCES tickerbeats.accounts (accountid)        ON UPDATE RESTRICT ON DELETE restrict,
  CONSTRAINT signalsresult_account_destination_fk  FOREIGN KEY (DestinationAccountID)                  REFERENCES tickerbeats.accounts (accountid)        ON UPDATE RESTRICT ON DELETE restrict,
  CONSTRAINT signalsresult_source_position_fk      FOREIGN KEY (SourceAccountID, SourcePositionID)     REFERENCES tickerbeats.positions(accountid,ticket) ON UPDATE RESTRICT ON DELETE restrict
);
CREATE INDEX signalsresult_groupid_idx ON tickerbeats.signalsresult (sourceaccountid,destinationaccountid,groupid);