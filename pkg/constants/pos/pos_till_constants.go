package pos_constants

type PosTillTransactionStatus string

const (
	PosTillTransactionStatusOpen  PosTillTransactionStatus = "OPEN"
	PosTillTransactionStatusClose PosTillTransactionStatus = "CLOSE"
)

type PosPaymentTransactionMode string

const (
	PosPaymentTransactionModeCash         PosPaymentTransactionMode = "CASH"
	PosPaymentTransactionModeUpi          PosPaymentTransactionMode = "UPI"
	PosPaymentTransactionModeCards        PosPaymentTransactionMode = "CARDS"
	PosPaymentTransactionModeCard         PosPaymentTransactionMode = "CARD"
	PosPaymentTransactionModeLoyalty      PosPaymentTransactionMode = "LOYALTY"
	PosPaymentTransactionModeVoucher      PosPaymentTransactionMode = "VOUCHER"
	PosPaymentTransactionModeCod          PosPaymentTransactionMode = "COD"
	PosPaymentTransactionModeTotalSales   PosPaymentTransactionMode = "TOTAL_SALES"
	PosPaymentTransactionModeTotalAdvance PosPaymentTransactionMode = "TOTAL_ADVANCE"
)

type PosPaymentTransactionType string

const (
	PosPaymentTransactionTypeSales   PosPaymentTransactionType = "SALES"
	PosPaymentTransactionTypeAdvance PosPaymentTransactionType = "ADVANCE"
	PosPaymentTransactionTypeReturn  PosPaymentTransactionType = "RETURN"
	PosPaymentTransactionTypeExpense PosPaymentTransactionType = "EXPENSE"
	PosPaymentTransactionTypeCod     PosPaymentTransactionType = "COD"
	PosPaymentTransactionTypeDeposit PosPaymentTransactionType = "DEPOSIT"
)
