package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// BaseModel menyediakan UUID dan timestamp untuk semua entitas SaaS
type BaseModel struct {
	ID        uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate otomatis generate UUID jika belum diset
func (base *BaseModel) BeforeCreate(tx *gorm.DB) error {
	if base.ID == uuid.Nil {
		base.ID = uuid.New()
	}
	return nil
}

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

// ---------------------------------------------------------
// ENTITAS UTAMA SAAS (MULTI-TENANT)
// ---------------------------------------------------------

// 1. Merchant (Tenant Utama / Pemilik Usaha)
type Merchant struct {
	BaseModel
	Name     string   `gorm:"size:150;not null" json:"name"`
	Email    string   `gorm:"size:150;uniqueIndex" json:"email"`
	Phone    string   `gorm:"size:50" json:"phone"`
	IsActive bool     `gorm:"default:true" json:"is_active"`
	Outlets  []Outlet `gorm:"foreignKey:MerchantID" json:"outlets,omitempty"`
}

// 2. Outlet (Lapak / Gudang Cabang)
type Outlet struct {
	BaseModel
	MerchantID  uuid.UUID `gorm:"type:char(36);index;not null" json:"merchant_id"`
	Name        string    `gorm:"size:150;not null" json:"name"`
	Code        string    `gorm:"size:50;index;not null" json:"code"`
	Address     string    `gorm:"size:255" json:"address"`
	Phone       string    `gorm:"size:50" json:"phone"`
	IsWarehouse bool      `gorm:"default:false" json:"is_warehouse"` // Penanda Gudang Utama
	IsActive    bool      `gorm:"default:true" json:"is_active"`
}

// 3. Pencatatan Pengeluaran (Expense Tracker - Standar Akuntansi)
type Expense struct {
	BaseModel
	MerchantID  uuid.UUID `gorm:"type:char(36);index;not null" json:"merchant_id"`
	OutletID    uuid.UUID `gorm:"type:char(36);index;not null" json:"outlet_id"`
	Amount      float64   `gorm:"type:decimal(12,2);not null" json:"amount"`
	Category    string    `gorm:"size:100;not null" json:"category"` // Listrik, Gaji, Sewa
	Description string    `gorm:"size:255" json:"description"`
	ExpenseDate time.Time `json:"expense_date"`
	ReportedBy  string    `gorm:"size:100" json:"reported_by"`
}

// ---------------------------------------------------------
// ENTITAS TRANSAKSIONAL (Terisolasi per Merchant)
// ---------------------------------------------------------

type User struct {
	BaseModel
	MerchantID   *uuid.UUID `gorm:"type:char(36);index" json:"merchant_id,omitempty"`
	OutletID     *uuid.UUID `gorm:"type:char(36);index" json:"outlet_id,omitempty"`
	Username     string     `gorm:"size:100;not null;unique" json:"username"`
	Name         string     `gorm:"size:100" json:"name"`
	ProfilePhoto string     `gorm:"type:longtext" json:"profile_photo"`
	Password     string     `gorm:"size:255;not null" json:"-"`
	Role         UserRole   `gorm:"type:varchar(20);default:'KASIR'" json:"role"`
}

type Category struct {
	BaseModel
	MerchantID *uuid.UUID `gorm:"type:char(36);index" json:"merchant_id,omitempty"`
	Name       string     `gorm:"size:100;not null" json:"name"`
	Icon       string     `gorm:"size:50;default:'utensils'" json:"icon"`
	ParentID   *uuid.UUID `gorm:"type:char(36);index" json:"parent_id,omitempty"`
}

type Artist struct {
	BaseModel
	MerchantID *uuid.UUID `gorm:"type:char(36);index" json:"merchant_id,omitempty"`
	Name       string     `gorm:"size:100;not null" json:"name"`
}

type ProductType struct {
	BaseModel
	MerchantID *uuid.UUID `gorm:"type:char(36);index" json:"merchant_id,omitempty"`
	Name       string     `gorm:"size:100;not null" json:"name"`
}

type Product struct {
	BaseModel
	MerchantID      uuid.UUID  `gorm:"type:char(36);index;not null" json:"merchant_id"`
	OutletID        *uuid.UUID `gorm:"type:char(36);index" json:"outlet_id,omitempty"`
	CategoryID      *uuid.UUID `gorm:"type:char(36);index" json:"category_id"`
	Category        Category   `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Name            string     `gorm:"size:150;not null" json:"name"`
	Artist          string     `gorm:"size:100;index" json:"artist"`
	ProductType     string     `gorm:"size:100;index" json:"product_type"` // 'retail' atau 'fnb'
	Price           float64    `gorm:"type:decimal(12,2);not null" json:"price"`
	CostPrice       float64    `gorm:"type:decimal(12,2);default:0" json:"cost_price"` // HPP Terakhir
	Stock           float64    `gorm:"type:decimal(12,3);default:0" json:"stock"`
	Unit            string     `gorm:"size:10;not null;default:'pcs'" json:"unit"`
	Barcode         string     `gorm:"size:100;index" json:"barcode"`
	ImageURL        string     `gorm:"size:255" json:"image_url"`
	IsActive        bool       `gorm:"default:true" json:"is_active"`
	IsMaster        bool       `gorm:"default:false;index" json:"is_master"`
	MasterProductID *uuid.UUID `gorm:"type:char(36);index" json:"master_product_id,omitempty"`
	OriginalOutletID *uuid.UUID `gorm:"type:char(36);index" json:"original_outlet_id,omitempty"`
	OriginalOutlet   *Outlet    `gorm:"foreignKey:OriginalOutletID" json:"original_outlet,omitempty"`
}

type Order struct {
	BaseModel
	MerchantID    uuid.UUID   `gorm:"type:char(36);index;not null" json:"merchant_id"`
	OutletID      uuid.UUID   `gorm:"type:char(36);index;not null" json:"outlet_id"`
	InvoiceNo     string      `gorm:"size:50;index;not null" json:"invoice_no"` 
	TotalAmount   float64     `gorm:"type:decimal(12,2);not null" json:"total_amount"`
	Discount      float64     `gorm:"type:decimal(12,2);default:0" json:"discount"`
	Tax           float64     `gorm:"type:decimal(12,2);default:0" json:"tax"`
	GrandTotal    float64     `gorm:"type:decimal(12,2);not null" json:"grand_total"`
	PaidAmount    float64     `gorm:"type:decimal(12,2);not null" json:"paid_amount"`
	ChangeAmount  float64     `gorm:"type:decimal(12,2);default:0" json:"change_amount"`
	PaymentMethod    string      `gorm:"size:50;default:'cash'" json:"payment_method"` // cash, qris, debit, credit
	PaymentProof     string      `gorm:"type:longtext" json:"payment_proof"`
	PaymentReference string      `gorm:"size:100" json:"payment_reference"`
	PlatformFee      float64     `gorm:"type:decimal(12,2);default:0" json:"platform_fee"`
	Status           string      `gorm:"size:30;default:'completed'" json:"status"`    
	CashierName   string      `gorm:"size:100;default:'Kasir Utama'" json:"cashier_name"`
	CustomerName  string      `gorm:"size:100;default:'Umum'" json:"customer_name"`
	TableNumber   string      `gorm:"size:50" json:"table_number"`
	OrderItems    []OrderItem `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE" json:"order_items"`
}

type OrderItem struct {
	BaseModel
	OrderID      uuid.UUID `gorm:"type:char(36);index;not null" json:"order_id"`
	ProductID    uuid.UUID `gorm:"type:char(36);index;not null" json:"product_id"`
	ProductName  string    `gorm:"size:150;not null" json:"product_name"`
	Artist       string    `gorm:"size:100" json:"artist"`
	ProductType  string    `gorm:"size:100" json:"product_type"`
	ProductPrice float64   `gorm:"type:decimal(12,2);not null" json:"product_price"`
	Quantity     float64   `gorm:"type:decimal(12,3);not null" json:"quantity"`
	Unit         string    `gorm:"size:10;not null;default:'pcs'" json:"unit"`
	Subtotal     float64   `gorm:"type:decimal(12,2);not null" json:"subtotal"`
	Notes        string    `gorm:"size:255" json:"notes"`
}

type StoreSetting struct {
	BaseModel
	MerchantID               uuid.UUID  `gorm:"type:char(36);index;unique" json:"merchant_id"`
	OutletID                 *uuid.UUID `gorm:"type:char(36);index" json:"outlet_id,omitempty"`
	StoreName                string     `gorm:"size:150;default:'KASIR POS PRO'" json:"store_name"`
	Address                  string     `gorm:"size:255;default:'Jl. Merdeka No. 45, Jakarta'" json:"address"`
	Phone                    string     `gorm:"size:50;default:'0812-3456-7890'" json:"phone"`
	ReceiptFooter            string     `gorm:"size:255;default:'Terima kasih telah berbelanja!'" json:"receipt_footer"`
	TaxPercentage            float64    `gorm:"default:10" json:"tax_percentage"`
	MemberDiscountPercentage float64    `gorm:"default:5" json:"member_discount_percentage"`
	DbEngine                 string     `gorm:"size:20;default:'sqlite'" json:"db_engine"`
	MysqlDsn                 string     `gorm:"size:255;default:'root:@tcp(127.0.0.1:3306)/pos_db?charset=utf8mb4&parseTime=True&loc=Local'" json:"mysql_dsn"`
	QrisImageUrl             string     `gorm:"type:longtext" json:"qris_image_url"`
	EnableRealtimeQris       bool       `gorm:"default:false" json:"enable_realtime_qris"`
	QrisProvider             string     `gorm:"size:50;default:'manual'" json:"qris_provider"` // "manual", "midtrans", "moota"
}

type Voucher struct {
	BaseModel
	MerchantID  uuid.UUID `gorm:"type:char(36);index;not null" json:"merchant_id"`
	Code        string    `gorm:"size:50;index;not null" json:"code"`
	Type        string    `gorm:"size:20;default:'percent'" json:"type"`
	Value       float64   `gorm:"type:decimal(12,2);not null" json:"value"`
	Description string    `gorm:"size:150" json:"description"`
	IsActive    bool      `gorm:"default:true" json:"is_active"`
}

type Customer struct {
	BaseModel
	MerchantID uuid.UUID `gorm:"type:char(36);index;not null" json:"merchant_id"`
	Name       string    `gorm:"size:100;not null" json:"name"`
	Phone      string    `gorm:"size:50" json:"phone"`
	Email      string    `gorm:"size:100" json:"email"`
	Address    string    `gorm:"size:255" json:"address"`
	Points     int       `gorm:"default:0" json:"points"`
}

type StockMovement struct {
	BaseModel
	MerchantID uuid.UUID `gorm:"type:char(36);index;not null" json:"merchant_id"`
	OutletID   uuid.UUID `gorm:"type:char(36);index;not null" json:"outlet_id"`
	ProductID  uuid.UUID `gorm:"type:char(36);index;not null" json:"product_id"`
	Product    Product   `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	Type       string    `gorm:"size:20;not null" json:"type"` // "in" or "out"
	Quantity   float64   `gorm:"type:decimal(12,3);not null" json:"quantity"`
	Unit       string    `gorm:"size:10;not null;default:'pcs'" json:"unit"`
	Reason     string    `gorm:"size:50;not null" json:"reason"` 
	Notes      string    `gorm:"size:255" json:"notes"`
}
