package entity

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

type Order struct {
	Id                  string    `gorm:"primaryKey;column:id;"`
	IdUser              string    `gorm:"column:id_user;"`
	IdDesa              string    `gorm:"column:id_desa;"`
	NumberOrder         string    `gorm:"column:number_order;"`
	OrderName           string    `gorm:"column:order_name;"`
	TrxId               int       `gorm:"column:trx_id;"`
	NamaLengkap         string    `gorm:"column:nama_lengkap;"`
	Email               string    `gorm:"column:email;"`
	Phone               string    `gorm:"column:phone;"`
	AlamatPengiriman    string    `gorm:"column:alamat_pengiriman;"`
	Catatan             string    `gorm:"column:catatan;"`
	ShippingCost        float64   `gorm:"column:shipping_cost;"`
	PaymentCash         float64   `gorm:"column:payment_cash;"`
	PaymentPoint        float64   `gorm:"column:payment_point;"`
	PaymentFee          float64   `gorm:"column:payment_fee;"`
	SubTotal            float64   `gorm:"column:sub_total;"`
	TotalBill           float64   `gorm:"column:total_bill;"`
	PaymentMethod       string    `gorm:"column:payment_method;"`
	PaymentChannel      string    `gorm:"column:payment_channel;"`
	PaymentNo           string    `gorm:"column:payment_no;"`
	PaymentStatus       int       `gorm:"column:payment_status;"`
	PaymentName         string    `gorm:"column:payment_name;"`
	PaymentDueDate      null.Time `gorm:"column:payment_due_date;"`
	PaymentSuccessDate  null.Time `gorm:"column:payment_success_date;"`
	PaymentExperiedDate null.Time `gorm:"column:payment_expired_date;"`
	OrderStatus         int       `gorm:"column:order_status;"`
	OrderedDate         time.Time `gorm:"column:order_date;"`
	OrderCompletedDate  null.Time `gorm:"column:order_completed_date;"`
	OrderCanceledDate   null.Time `gorm:"column:order_cancel_date;"`
	FotoBarangSampai    string    `gorm:"column:foto_barang_sampai;"`
	FotoBuktiBayar      string    `gorm:"column:foto_bukti_bayar;"`
}

func (Order) TableName() string {
	return "orders_transaction"
}
