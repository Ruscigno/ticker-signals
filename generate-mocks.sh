#generate mocks
mockgen -source=internal/transaction/accounts/service/interface.go -destination=internal/transaction/accounts/service/interface_mock.go
mockgen -source=internal/transaction/accounts/repo/interface.go -destination=internal/transaction/accounts/repo/interface_mock.go

mockgen -source=internal/transaction/accountsinfo/service/interface.go -destination=internal/transaction/accountsinfo/service/interface_mock.go
mockgen -source=internal/transaction/accountsinfo/repo/interface.go -destination=internal/transaction/accountsinfo/repo/interface_mock.go

mockgen -source=internal/transaction/deals/service/interface.go -destination=internal/transaction/deals/service/interface_mock.go
mockgen -source=internal/transaction/deals/repo/interface.go -destination=internal/transaction/deals/repo/interface_mock.go

mockgen -source=internal/transaction/orders/service/interface.go -destination=internal/transaction/orders/service/interface_mock.go
mockgen -source=internal/transaction/orders/repo/interface.go -destination=internal/transaction/orders/repo/interface_mock.go

mockgen -source=internal/transaction/positions/service/interface.go -destination=internal/transaction/positions/service/interface_mock.go
mockgen -source=internal/transaction/positions/repo/interface.go -destination=internal/transaction/positions/repo/interface_mock.go

mockgen -source=internal/transaction/traderules/service/interface.go -destination=internal/transaction/traderules/service/interface_mock.go
mockgen -source=internal/transaction/traderules/repo/interface.go -destination=internal/transaction/traderules/repo/interface_mock.go

mockgen -source=internal/transaction/tradetransaction/service/interface.go -destination=internal/transaction/tradetransaction/service/interface_mock.go
mockgen -source=internal/transaction/tradetransaction/repo/interface.go -destination=internal/transaction/tradetransaction/repo/interface_mock.go