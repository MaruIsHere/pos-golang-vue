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

export interface Category {
  id: number;
  name: string;
  icon?: string;
  created_at?: string;
}

export interface Product {
  id: number;
  category_id: number;
  category?: Category;
  name: string;
  price: number;
  cost_price?: number;
  stock: number;
  barcode?: string;
  image_url?: string;
  is_active?: boolean;
  created_at?: string;
  updated_at?: string;
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
  customer_name: string;
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
  product_price: number;
  quantity: number;
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
  reason?: string;
  notes?: string;
  created_at?: string;
}

export interface DashboardStats {
  total_orders: number;
  total_revenue: number;
  total_items_sold: number;
  top_products: Array<{ product_name: string; total_qty: number; total_sales: number }>;
  recent_orders: Order[];
}
