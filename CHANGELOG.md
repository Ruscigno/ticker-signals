## 36
fix: empty symbol

## 35
fix: exchange symbol translation

## 34
feat: always create beats that are no more than 24h old

## 33
feat: avoid querying the DB too much for beats

## 32
fix: create time
feat: set type time

## 31
fix: price null in the signal beats

## 30
chore: remove debug log

## 29
chore: debug log

## 28
fix: update position profit

## 27
fix: CreateTickerBeatsFromOrders

## 26
fix: volume change for close positions trades

## 25
feat: symbolTranslation entity

## 24
feat: rollback last change

## 23
fix: don't select invalid signal trades

## 22
feat: cleaning up code and logs

## 21
feat: several improvements
## 20
feat: new beats generation system

## 19
feat: beats from trade requests

## 18
feat: add PositionsHistoryService
## 17
feat: adding debug logs
## 16
feat: docker push as latest tag
## 13
feat: add jenkins build script
## v0.18.3
feat: update dockerfile to use go 1.17
## v0.18.2
feat: add some debug logs
## v0.18.1
feat: update to go 1.17
## v0.18.0
- feat: use tradetransactions_seq

## v0.17.1.2
- fix: ticker beats query
## v0.17.1
- feat: improve select to CreateTickerBeats
## v0.17.0
- feat: integration test
- feat: stop & loss
- feat: CreateAccount test
- fix: fatal error: concurrent map write
## v0.16.1
- fix: disable nexttime feature

## v0.16.0
- feat: Trade Rules initial code

## v0.15.0
- feat: OrderBoostType feature

## v0.14.4
- fix: avoid duplicated tickes ENTRY_OUT

## v0.14.3
-fix: CloseDeadPositions Error

## v0.14.2
-fix: ConfirmByExternalID
-feat: first slice of NeedToCloseAllPositions

## v0.14.1
-feat: improve logs

## v0.14.0
-feat: order boost

## v0.13.20
-fix: CloseDeadPositions query

## v0.13.19
-fix: signal status calculation

## v0.13.18
-feat: delete invalid close position beats
-feat: close dead positions
-feat: update status of dead positions

## v0.13.17
- feat: improve logs

## v0.13.16
- fix: Ticker Beats Response
- feat: close positions

## v0.13.15
- fix: skiping new transactions

## v0.13.14
- fix: remove expired beats

## v0.13.13
- fix: skip new records
- feat: add GMT time

## v0.13.12
- feat: RemoveDuplicatedSignals

## v0.13.11
- feat: update commission, swap, profit

## v0.13.10
- fix: deal OUT = position closed

## v0.13.9
- rolled back changes...

## v0.13.8
- feat: null fields
- fix: duplicated key on index
- fix: avoid closing newer opened position

## v0.13.7
- feat: update orderId on deals model
- fix: ConfirmByExternalID method

## v0.13.6
- fix: api request processing

## v0.13.5
- fix: time expiration conversion

## v0.13.4.2
- feat: merge main

## v0.13.4.1
- fix: build

## v0.13.4
- feat: calc UTC time [TimeGMTOffset]
- feat: UTC everywhere

## v0.13.3
- feat: avoid unnecessary db transactions

## v0.13.2
- feat: Position direction IN working

## v0.13.1
- feat: return TradeRequest as result
## v0.13.0.1
- fix: enum values

## v0.13.0
- feat: add server time to transactions

## v0.12.0
- feat: new position model

## v0.11.0
- feat: working ticker beats API

## v0.10.0
- feat: transactions

## v0.9.7
- fix: updating more fields

## v0.9.6
- fix: magic could be null

## v0.9.5
- feat: ticker beats update

## v0.9.4
- fix: update statement

## v0.9.3
- fix: setting the active position(s)

## v0.9.2
- fix: not found error on ConfirmByExternalID

## v0.9.1
- feat: All set, TickerBeats Services

## v0.9.0
- feat: Signal service

## v0.8.1
- fix: stop updating delete time if it was alread set before
- fix: use ticket as position ID as it isn't recicled

## v0.8.0
- feat: updating current price

## v0.7.0
- feat: server time on account info

## v0.6.3
- feat: AccountId as main field
- feat: Keep only current positions
## v0.6.2
- feat: mark positions as active
  
## v0.6.1
- feat: improved checks before send request

## v0.6.0
- feat: TickerBeats

## v0.5.0
- feat: TradeTransaction processing
- 
## v0.4.4
- Fix: please don't panic

## v0.4.3
- Improving logs format

## v0.4.2
- Fix: aborting when there are no rows in result set

## v0.4.1
- Implementing GetByID on Repos to be used instead of using the upsert method
- Fix: space on env variables

## v0.4.0
- synchronizing versions between client and server