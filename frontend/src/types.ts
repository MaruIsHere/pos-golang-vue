// Shared API types — mirror of backend/models/models.go (JSON shape).

export interface StoreSetting {
  id?: string;
  merchant_id?: string;
  outlet_id?: string | null;
  store_name: string;
  address: string;
  phone: string;
  receipt_footer: string;
  tax_percentage: number;
  member_discount_percentage?: number;
  db_engine: string;
  mysql_dsn?: string;
  qris_image_url?: string;
}

export interface Merchant {
  id: string;
  name: string;
  email?: string;
  phone?: string;
  is_active?: boolean;
}

export interface Outlet {
  id: string;
  merchant_id: string;
  name: string;
  code: string;
  address?: string;
  phone?: string;
  is_active?: boolean;
  is_warehouse?: boolean;
  total_cashiers?: number;
  total_products?: number;
  created_at?: string;
  updated_at?: string;
}

export interface UserProfile {
  id: string;
  username: string;
  name: string;
  profile_photo?: string;
  role: string;
  merchant_id?: string | null;
  outlet_id?: string | null;
  outlet?: Outlet;
}

export interface Category {
  id: string;
  merchant_id?: string;
  name: string;
  icon?: string;
  parent_id?: string | null;
  created_at?: string;
}

export interface Artist {
  id: string;
  merchant_id?: string;
  name: string;
  created_at?: string;
}

export interface ProductType {
  id: string;
  merchant_id?: string;
  name: string;
  created_at?: string;
}

export interface Product {
  id: string;
  merchant_id: string;
  outlet_id?: string | null;
  category_id: string | null;
  category?: Category;
  name: string;
  artist?: string;
  product_type?: string;
  price: number;
  cost_price?: number;
  stock: number;
  unit: 'pcs' | 'gram' | 'liter';
  barcode?: string;
  image_url?: string;
  is_active?: boolean;
  is_master?: boolean;
  master_product_id?: string | null;
  outlet?: Outlet;
  created_at?: string;
  updated_at?: string;
}

export interface ProductFilters {
  artists: string[];
  product_types: string[];
}

export interface CartItem {
  product: Product;
  quantity: number;
  notes: string;
}

export interface OrderItemInput {
  product_id: string;
  quantity: number;
  notes?: string;
}

export interface CreateOrderPayload {
  merchant_id?: string | null;
  outlet_id?: string | null;
  customer_name: string;
  cashier_name?: string;
  payment_method: string;
  payment_proof?: string;
  paid_amount: number;
  discount: number;
  tax: number;
  items: OrderItemInput[];
}

export interface OrderItem {
  id: string;
  order_id: string;
  product_id: string;
  product_name: string;
  artist?: string;
  product_type?: string;
  product_price: number;
  quantity: number;
  unit: 'pcs' | 'gram' | 'liter';
  subtotal: number;
  notes?: string;
}

export interface Order {
  id: string;
  merchant_id: string;
  outlet_id: string;
  invoice_no: string;
  total_amount: number;
  discount: number;
  tax: number;
  grand_total: number;
  paid_amount: number;
  change_amount: number;
  payment_method: string;
  payment_proof?: string;
  status: string;
  cashier_name?: string;
  customer_name?: string;
  created_at?: string;
  order_items?: OrderItem[];
}

export interface Customer {
  id: string;
  merchant_id: string;
  name: string;
  phone?: string;
  email?: string;
  address?: string;
  points?: number;
  created_at?: string;
}

export interface Voucher {
  id?: string;
  merchant_id: string;
  code: string;
  type: string;
  value: number;
  description?: string;
  is_active?: boolean;
  created_at?: string;
}

export interface StockMovement {
  id: string;
  merchant_id: string;
  outlet_id: string;
  product_id: string;
  product?: Product;
  type: string;
  quantity: number;
  unit: 'pcs' | 'gram' | 'liter';
  reason?: string;
  notes?: string;
  created_at?: string;
}

export interface SubGroupStat {
  name: string;
  total_qty: number;
  total_sales: number;
}

export interface ProductSalesStat {
  product_id: string;
  product_name: string;
  artist?: string;
  product_type?: string;
  total_qty: number;
  total_sales: number;
  stock?: number;
  price?: number;
  unit?: 'pcs' | 'gram' | 'liter';
}

export interface DashboardStats {
  total_orders: number;
  total_revenue: number;
  total_items_sold: number;
  top_products: ProductSalesStat[];
  least_products?: ProductSalesStat[];
  all_sold_products?: ProductSalesStat[];
  sales_by_artist?: SubGroupStat[];
  sales_by_type?: SubGroupStat[];
  recent_orders: Order[];
  store_setting?: StoreSetting;
}
