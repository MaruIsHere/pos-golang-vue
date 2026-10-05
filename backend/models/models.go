package models

import (
	"time"
)

type UserRole string

const (
	RoleOwner         UserRole = "OWNER"
	RoleAdministrator UserRole = "ADMINISTRATOR"
	RoleKepalaKasir   UserRole = "KEPALA_KASIR"
	RoleAdmin         UserRole = "ADMIN"
	RoleKasir         UserRole = "KASIR"
)

func (r UserRole) IsValid() bool {
	switch r {
	case RoleOwner, RoleAdministrator, RoleAdmin, RoleKepalaKasir, RoleKasir:
		return true
	}
	return false
}

type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Username     string    `gorm:"size:100;not null;unique" json:"username"`
	Name         string    `gorm:"size:100" json:"name"`
	ProfilePhoto string    `gorm:"type:longtext" json:"profile_photo"`
	Password     string    `gorm:"size:255;not null" json:"-"`
	Role         UserRole  `gorm:"type:varchar(20);default:'KASIR'" json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Category struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:100;not null;unique" json:"name"`
	Icon      string    `gorm:"size:50;default:'utensils'" json:"icon"`
	ParentID  *uint     `gorm:"index" json:"parent_id,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type Artist struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:100;not null;uniqueIndex" json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type ProductType struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:100;not null;uniqueIndex" json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type Product struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	CategoryID  uint      `gorm:"index;not null" json:"category_id"`
	Category    Category  `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Name        string    `gorm:"size:150;not null" json:"name"`
	Artist      string    `gorm:"size:100;index" json:"artist"`
	ProductType string    `gorm:"size:100;index" json:"product_type"`
	Price       float64   `gorm:"type:decimal(12,2);not null" json:"price"`
	CostPrice   float64   `gorm:"type:decimal(12,2);default:0" json:"cost_price"`
	Stock       float64   `gorm:"type:decimal(12,3);default:0" json:"stock"`
	Unit        string    `gorm:"size:10;not null;default:'pcs'" json:"unit"`
	Barcode     string    `gorm:"size:100;index" json:"barcode"`
	ImageURL    string    `gorm:"size:255" json:"image_url"`
	IsActive    bool      `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Order struct {
	ID            uint        `gorm:"primaryKey" json:"id"`
	InvoiceNo     string      `gorm:"size:50;uniqueIndex;not null" json:"invoice_no"`
	TotalAmount   float64     `gorm:"type:decimal(12,2);not null" json:"total_amount"`
	Discount      float64     `gorm:"type:decimal(12,2);default:0" json:"discount"`
	Tax           float64     `gorm:"type:decimal(12,2);default:0" json:"tax"`
	GrandTotal    float64     `gorm:"type:decimal(12,2);not null" json:"grand_total"`
	PaidAmount    float64     `gorm:"type:decimal(12,2);not null" json:"paid_amount"`
	ChangeAmount  float64     `gorm:"type:decimal(12,2);default:0" json:"change_amount"`
	PaymentMethod string      `gorm:"size:50;default:'cash'" json:"payment_method"` // cash, qris, debit, credit
	Status        string      `gorm:"size:30;default:'completed'" json:"status"`    // completed, refunded, cancelled
	CashierName   string      `gorm:"size:100;default:'Kasir Utama'" json:"cashier_name"`
	CustomerName  string      `gorm:"size:100;default:'Umum'" json:"customer_name"`
	CreatedAt     time.Time   `json:"created_at"`
	OrderItems    []OrderItem `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE" json:"order_items"`
}

type OrderItem struct {
	ID           uint    `gorm:"primaryKey" json:"id"`
	OrderID      uint    `gorm:"index;not null" json:"order_id"`
	ProductID    uint    `gorm:"index;not null" json:"product_id"`
	ProductName  string  `gorm:"size:150;not null" json:"product_name"`
	Artist       string  `gorm:"size:100" json:"artist"`
	ProductType  string  `gorm:"size:100" json:"product_type"`
	ProductPrice float64 `gorm:"type:decimal(12,2);not null" json:"product_price"`
	Quantity     float64 `gorm:"type:decimal(12,3);not null" json:"quantity"`
	Unit         string  `gorm:"size:10;not null;default:'pcs'" json:"unit"`
	Subtotal     float64 `gorm:"type:decimal(12,2);not null" json:"subtotal"`
	Notes        string  `gorm:"size:255" json:"notes"`
}

type StoreSetting struct {
	ID                       uint    `gorm:"primaryKey" json:"id"`
	StoreName                string  `gorm:"size:150;default:'KASIR POS PRO'" json:"store_name"`
	Address                  string  `gorm:"size:255;default:'Jl. Merdeka No. 45, Jakarta'" json:"address"`
	Phone                    string  `gorm:"size:50;default:'0812-3456-7890'" json:"phone"`
	ReceiptFooter            string  `gorm:"size:255;default:'Terima kasih telah berbelanja!'" json:"receipt_footer"`
	TaxPercentage            float64 `gorm:"default:10" json:"tax_percentage"`
	MemberDiscountPercentage float64 `gorm:"default:5" json:"member_discount_percentage"`
	DbEngine                 string  `gorm:"size:20;default:'sqlite'" json:"db_engine"` // sqlite or mysql
	MysqlDsn                 string  `gorm:"size:255;default:'root:@tcp(127.0.0.1:3306)/pos_db?charset=utf8mb4&parseTime=True&loc=Local'" json:"mysql_dsn"`
	QrisImageUrl             string  `gorm:"type:longtext" json:"qris_image_url"`
}

type Voucher struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Code        string    `gorm:"size:50;not null;unique" json:"code"`
	Type        string    `gorm:"size:20;default:'percent'" json:"type"` // "percent" or "flat"
	Value       float64   `gorm:"type:decimal(12,2);not null" json:"value"`
	Description string    `gorm:"size:150" json:"description"`
	IsActive    bool      `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
}

type Customer struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:100;not null" json:"name"`
	Phone     string    `gorm:"size:50" json:"phone"`
	Email     string    `gorm:"size:100" json:"email"`
	Address   string    `gorm:"size:255" json:"address"`
	Points    int       `gorm:"default:0" json:"points"`
	CreatedAt time.Time `json:"created_at"`
}

type StockMovement struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ProductID uint      `gorm:"index;not null" json:"product_id"`
	Product   Product   `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	Type      string    `gorm:"size:20;not null" json:"type"` // "in" (Penerimaan) or "out" (Pengeluaran)
	Quantity  float64   `gorm:"type:decimal(12,3);not null" json:"quantity"`
	Unit      string    `gorm:"size:10;not null;default:'pcs'" json:"unit"`
	Reason    string    `gorm:"size:50;not null" json:"reason"` // "pembelian_supplier", "barang_rusak", "barang_hilang", "expired", "promosi", "retur_penjualan"
	Notes     string    `gorm:"size:255" json:"notes"`
	CreatedAt time.Time `json:"created_at"`
}
