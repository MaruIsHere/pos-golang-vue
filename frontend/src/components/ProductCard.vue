<template>
  <div class="product-card glass-panel" :class="{ 'out-of-stock': product.stock <= 0 }" @click="addToCart">
    <div class="card-image-wrapper">
      <img 
        :src="product.image_url || 'https://images.unsplash.com/photo-1546069901-ba9599a7e63c?w=400'" 
        :alt="product.name" 
        class="card-image"
        @error="onImageError"
      />
      <span class="stock-badge" :class="stockBadgeClass">
        {{ product.stock > 0 ? `Stok: ${product.stock}` : 'Habis' }}
      </span>
      <span v-if="product.category" class="category-badge">
        {{ product.category.name }}
      </span>
    </div>

    <div class="card-body">
      <h3 class="product-title">{{ product.name }}</h3>
      <p class="product-barcode" v-if="product.barcode">SKU: {{ product.barcode }}</p>
      
      <div class="card-footer">
        <div class="price-container">
          <span class="currency">Rp</span>
          <span class="price-value">{{ formatPrice(product.price) }}</span>
        </div>
        
        <button 
          class="btn-add" 
          :disabled="product.stock <= 0"
          @click.stop="addToCart"
        >
          <span>+</span>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import type { Product } from '../types';

const props = defineProps({
  product: { type: Object as () => Product, required: true }
});

const emit = defineEmits(['add-to-cart']);

const formatPrice = (val: number): string => {
  return new Intl.NumberFormat('id-ID').format(val || 0);
};

const stockBadgeClass = computed(() => {
  if (props.product.stock <= 0) return 'badge-danger';
  if (props.product.stock <= 10) return 'badge-warning';
  return 'badge-success';
});

const onImageError = (e: Event): void => {
  (e.target as HTMLImageElement).src = 'https://images.unsplash.com/photo-1546069901-ba9599a7e63c?w=400';
};

const addToCart = () => {
  if (props.product.stock > 0) {
    emit('add-to-cart', props.product);
  }
};
</script>

<style scoped>
.product-card {
  display: flex;
  flex-direction: column;
  overflow: hidden;
  cursor: pointer;
  transition: all 0.25s ease;
  position: relative;
  height: 100%;
}

.product-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 12px 24px rgba(0, 0, 0, 0.4), 0 0 15px rgba(99, 102, 241, 0.2);
  border-color: rgba(99, 102, 241, 0.4);
}

.product-card.out-of-stock {
  opacity: 0.6;
  cursor: not-allowed;
}

.card-image-wrapper {
  position: relative;
  width: 100%;
  height: 120px;
  overflow: hidden;
  background: #1e293b;
}

.card-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.3s ease;
}

.product-card:hover .card-image {
  transform: scale(1.05);
}

.stock-badge {
  position: absolute;
  top: 8px;
  right: 8px;
  font-size: 0.65rem;
  padding: 0.2rem 0.5rem;
  border-radius: 999px;
  font-weight: 700;
  backdrop-filter: blur(4px);
}

.category-badge {
  position: absolute;
  bottom: 8px;
  left: 8px;
  background: rgba(15, 23, 42, 0.75);
  color: var(--text-secondary);
  font-size: 0.65rem;
  padding: 0.15rem 0.45rem;
  border-radius: 6px;
  backdrop-filter: blur(4px);
}

.card-body {
  padding: 0.85rem;
  display: flex;
  flex-direction: column;
  flex: 1;
  justify-content: space-between;
}

.product-title {
  font-size: 0.9rem;
  font-weight: 700;
  color: var(--text-primary);
  line-height: 1.3;
  margin-bottom: 0.25rem;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.product-barcode {
  font-size: 0.7rem;
  color: var(--text-muted);
  margin-bottom: 0.5rem;
}

.card-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-top: 0.5rem;
}

.price-container {
  display: flex;
  align-items: baseline;
  gap: 0.15rem;
}

.currency {
  font-size: 0.7rem;
  font-weight: 700;
  color: var(--accent-secondary);
}

.price-value {
  font-size: 0.95rem;
  font-weight: 800;
  color: var(--accent-secondary);
}

.btn-add {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  border: none;
  background: linear-gradient(135deg, var(--accent-primary), var(--accent-purple));
  color: #ffffff;
  font-size: 1.1rem;
  font-weight: 700;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.2s;
  box-shadow: 0 2px 8px rgba(99, 102, 241, 0.4);
}

.btn-add:hover:not(:disabled) {
  transform: scale(1.1);
  box-shadow: 0 4px 12px rgba(99, 102, 241, 0.6);
}

.btn-add:disabled {
  background: #475569;
  cursor: not-allowed;
  box-shadow: none;
}
</style>
