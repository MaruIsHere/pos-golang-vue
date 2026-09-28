<template>
  <div 
    class="group flex flex-col overflow-hidden cursor-pointer transition-all duration-300 relative h-full rounded-[24px] bg-white dark:bg-slate-800 border border-slate-200/50 dark:border-slate-700/50 shadow-[0_2px_12px_-4px_rgba(0,0,0,0.06)] hover:-translate-y-1 hover:shadow-[0_8px_24px_-8px_rgba(79,70,229,0.2)] hover:border-indigo-400/50 dark:hover:border-indigo-500/50"
    :class="{ 'opacity-50 grayscale-[50%] cursor-not-allowed hover:translate-y-0 hover:shadow-none hover:border-slate-200': product.stock <= 0 }"
    @click="addToCart"
  >
    <div class="relative w-full h-36 sm:h-40 p-2 shrink-0">
      <div class="w-full h-full overflow-hidden rounded-[18px] bg-slate-50 dark:bg-slate-900 shadow-inner">
        <img 
          :src="product.image_url || 'https://images.unsplash.com/photo-1546069901-ba9599a7e63c?w=400'" 
          :alt="product.name" 
          class="w-full h-full object-cover transition-transform duration-500 group-hover:scale-110"
          @error="onImageError"
        />
      </div>
      
      <!-- Badges overlay directly on the image area for elegant look -->
      <AppBadge 
        class="absolute top-4 right-4 text-[0.65rem] px-2.5 py-1 font-bold shadow-sm backdrop-blur-md"
        :variant="stockBadgeVariant"
      >
        {{ product.stock > 0 ? `Stok: ${product.stock}` : 'Habis' }}
      </AppBadge>
    </div>
    
    <div class="flex flex-col px-4 pb-4 pt-1 flex-1">
      <div class="flex items-center gap-1.5 mb-1.5">
        <span class="text-[10px] font-bold text-slate-400 dark:text-slate-500 uppercase tracking-wider line-clamp-1">
          {{ product.category?.name || 'Umum' }}
        </span>
      </div>
      <h3 class="font-bold text-slate-800 dark:text-slate-100 text-[0.95rem] leading-tight line-clamp-2">{{ product.name }}</h3>
      <p class="text-[0.95rem] font-extrabold text-indigo-600 dark:text-indigo-400 mt-auto pt-3">Rp {{ formatPrice(product.price) }}</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import type { Product } from '../types';
import { PlusIcon } from '@heroicons/vue/24/outline';
import AppButton from './ui/AppButton.vue';
import AppBadge from './ui/AppBadge.vue';

const props = defineProps({
  product: { type: Object as () => Product, required: true }
});

const emit = defineEmits(['add-to-cart']);

const formatPrice = (val: number): string => {
  return new Intl.NumberFormat('id-ID').format(val || 0);
};

const stockBadgeVariant = computed(() => {
  if (props.product.stock <= 0) return 'danger';
  if (props.product.stock <= 10) return 'warning';
  return 'success';
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
