package sample

type addMoneyRequest struct {
	accessTokenHaving

	amount int64

	moneySpent moneySpentRequest
}
