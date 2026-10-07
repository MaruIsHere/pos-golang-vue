// Shared API types — mirror of backend/models/models.go (JSON shape).

export interface StoreSetting {
  id?: number;
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

export interface Store {
  id: number;
  name: string;
  code: string;
  address?: string;
  phone?: string;
  is_active?: boolean;
  total_cashiers?: number;
  total_products?: number;
  created_at?: string;
  updated_at?: string;
}

export interface UserProfile {
  id: number;
  username: string;
  name: string;
  profile_photo?: string;
  role: string;
  store_id?: number | null;
  store?: Store;
}

export interface Category {
  id: number;
  name: string;
  icon?: string;
  parent_id?: number | null;
  created_at?: string;
}

export interface Artist {
  id: number;
  name: string;
  created_at?: string;
}

export interface ProductType {
  id: number;
  name: string;
  created_at?: string;
}

export interface Product {
  id: number;
  category_id: number;
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
  master_product_id?: number | null;
  store_id?: number | null;
  store?: Store;
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
  product_id: number;
  quantity: number;
  notes?: string;
}

export interface CreateOrderPayload {
  store_id?: number | null;
  customer_name: string;
  cashier_name?: string;
  payment_method: string;
  paid_amount: number;
  discount: number;
  tax: number;
  items: OrderItemInput[];
}

export interface OrderItem {
  id: number;
  order_id: number;
  product_id: number;
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
  id: number;
  invoice_no: string;
  total_amount: number;
  discount: number;
  tax: number;
  grand_total: number;
  paid_amount: number;
  change_amount: number;
  payment_method: string;
  status: string;
  cashier_name?: string;
  customer_name?: string;
  created_at?: string;
  order_items?: OrderItem[];
}

export interface Customer {
  id: number;
  name: string;
  phone?: string;
  email?: string;
  address?: string;
  points?: number;
  created_at?: string;
}

export interface Voucher {
  id?: number;
  code: string;
  type: string;
  value: number;
  description?: string;
  is_active?: boolean;
  created_at?: string;
}

export interface StockMovement {
  id: number;
  product_id: number;
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
  product_id: number;
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
