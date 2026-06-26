package constants

type TillTransactionStatus string

const (
	TillTransactionStatusOpen  TillTransactionStatus = "OPEN"
	TillTransactionStatusClose TillTransactionStatus = "CLOSE"
)

type PaymentTransactionMode string

const (
	PaymentTransactionModeCash         PaymentTransactionMode = "CASH"
	PaymentTransactionModeUpi          PaymentTransactionMode = "UPI"
	PaymentTransactionModeCards        PaymentTransactionMode = "CARDS"
	PaymentTransactionModeCard         PaymentTransactionMode = "CARD"
	PaymentTransactionModeLoyalty      PaymentTransactionMode = "LOYALTY"
	PaymentTransactionModeVoucher      PaymentTransactionMode = "VOUCHER"
	PaymentTransactionModeCod          PaymentTransactionMode = "COD"
	PaymentTransactionModeTotalSales   PaymentTransactionMode = "TOTAL_SALES"
	PaymentTransactionModeTotalAdvance PaymentTransactionMode = "TOTAL_ADVANCE"
)

type PaymentTransactionType string

const (
	PaymentTransactionTypeSales   PaymentTransactionType = "SALES"
	PaymentTransactionTypeAdvance PaymentTransactionType = "ADVANCE"
	PaymentTransactionTypeReturn  PaymentTransactionType = "RETURN"
	PaymentTransactionTypeExpense PaymentTransactionType = "EXPENSE"
	PaymentTransactionTypeCod     PaymentTransactionType = "COD"
	PaymentTransactionTypeDeposit PaymentTransactionType = "DEPOSIT"
)
