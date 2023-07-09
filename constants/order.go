package constants

// order status
const OrderTemporary = "temporary"
const OrderWaiting = "menunggu"
const OrderCompleted = "selesai"
const OrderProcess = "proses"
const OrderRejected = "ditolak"
const OrderCanceled = "dibatalkan"

// jenis order
const OrderHomeCalling = "home calling"
const OrderTakeDelivery = "take & delivery"

// temporary order exist
const TemporaryOrderExistMessage = "temporary order exist"

// biaya aplikasi
const AppsChargePercentageLessThan100k = 1.0
const AppsChargePercentageMoreThan100k = 1.5
const CostPerKMLessThan10KM = 1500
const CostPerKMMoreThan10KM = 700

// metode pembayarn
const BankTransferPayment = "bank transfer"
const AutomaticPayment = "bayar otomatis"
const CashPayment = "cash"
