package orders

import uuid "github.com/fossoreslp/go-uuid-v4"

type Order struct {
	ID uuid.UUID
	OrderItem
}

type OrderItem struct {
	OrderUID          string `json:"order_uid" validate:"required"`
	TrackNumber       string `json:"track_number" validate:"required,uppercase"`
	Entry             string `json:"entry" validate:"required,uppercase"`
	Delivery          `json:"delivery" validate:"required"`
	Payment           `json:"payment" validate:"required"`
	Items             []Items `json:"items" validate:"required"`
	Locale            string  `json:"locale" validate:"required"`
	InternalSignature string  `json:"internal_signature" validate:"required"`
	CustomerID        string  `json:"customer_id" validate:"required"`
	DeliveryService   string  `json:"delivery_service" validate:"required"`
	ShardKey          uint32  `json:"shardkey" validate:"required,numeric"`
	SmID              uint32  `json:"sm_id" validate:"required,numeric"`
	DateCreated       string  `json:"date_created" validate:"required"`
	OofShard          uint32  `json:"oof_shard" validate:"required,numeric"`
}

type Delivery struct {
	Name    string `json:"name" validate:"required"`
	Phone   string `json:"phone" validate:"required,e164"`
	Zip     uint32 `json:"zip" validate:"required,numeric"`
	City    string `json:"city" validate:"required"`
	Address string `json:"address" validate:"required"`
	Region  string `json:"region" validate:"required"`
	Email   string `json:"email" validate:"required,email"`
}

type Payment struct {
	Transaction  string `json:"transaction" validate:"required"`
	RequestID    string `json:"request_id" validate:"required"`
	Currency     string `json:"currency" validate:"required,uppercase"`
	Provider     string `json:"provider" validate:"required"`
	Amount       uint32 `json:"amount" validate:"required,numeric"`
	PaymentDT    uint32 `json:"payment_dt" validate:"required,numeric"`
	Bank         string `json:"bank" validate:"required"`
	DeliveryCost uint32 `json:"delivery_cost" validate:"required,numeric"`
	GoodsTotal   uint32 `json:"goods_total" validate:"required,numeric"`
	CustomFee    uint32 `json:"custom_fee" validate:"required,numeric"`
}

type Items struct {
	ChrtID      uint32 `json:"chrt_id" validate:"required,numeric"`
	TrackNumber string `json:"track_number" validate:"required,uppercase"`
	Price       uint32 `json:"price" validate:"required,numeric"`
	RID         string `json:"rid" validate:"required"`
	Name        string `json:"name" validate:"required"`
	Sale        uint32 `json:"sale" validate:"required,numeric"`
	Size        uint32 `json:"size" validate:"required"`
	TotalPrice  uint32 `json:"total_price" validate:"required,numeric"`
	NmID        uint32 `json:"nm_id" validate:"required,numeric"`
	Brand       string `json:"brand" validate:"required"`
	Status      uint32 `json:"status" validate:"required,numeric"`
}
