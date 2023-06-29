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
const AppsChargePercentageLessThan100k = 2
const AppsChargePercentageMoreThan100k = 5
const CostPerKMLessThan10KM = 2500
const CostPerKMMoreThan10KM = 1000
