<template>
  <div 
    class="group flex flex-col overflow-hidden cursor-pointer transition-all duration-200 relative h-full rounded-xl border border-slate-200 dark:border-slate-700 backdrop-blur-md bg-white/90 dark:bg-slate-900/90 hover:-translate-y-0.5 hover:shadow-md hover:border-indigo-500"
    :class="{ 'opacity-65 cursor-not-allowed hover:translate-y-0 hover:shadow-none hover:border-slate-200 dark:hover:border-slate-700': product.stock <= 0 }"
    @click="addToCart"
  >
    <div class="relative w-full h-[130px] overflow-hidden bg-slate-50 dark:bg-slate-900">
      <img 
        :src="product.image_url || 'https://images.unsplash.com/photo-1546069901-ba9599a7e63c?w=400'" 
        :alt="product.name" 
        class="w-full h-full object-cover transition-transform duration-300 group-hover:scale-105"
        @error="onImageError"
      />
      <AppBadge 
        class="absolute top-2 right-2 text-[0.65rem] px-2 py-0.5 font-bold"
        :variant="stockBadgeVariant"
      >
        {{ product.stock > 0 ? `Stok: ${product.stock}` : 'Habis' }}
      </AppBadge>
      <AppBadge 
        v-if="product.category" 
        variant="secondary"
        class="absolute bottom-2 left-2 text-[0.65rem] font-semibold px-1.5 py-0.5 border border-slate-200 dark:border-slate-700 bg-white/50 dark:bg-slate-800/50 backdrop-blur-sm text-slate-600 dark:text-slate-400 rounded-md"
      >
        {{ product.category.name }}
      </AppBadge>
    </div>

    <div class="p-3.5 flex flex-col flex-1 justify-between">
      <h3 class="text-sm font-bold text-slate-900 dark:text-slate-100 leading-tight mb-1 line-clamp-2">
        {{ product.name }}
      </h3>
      
      <div class="flex flex-wrap gap-1 mb-1.5" v-if="product.artist || product.product_type">
        <AppBadge 
          v-if="product.artist" 
          variant="outline"
          class="text-[0.65rem] font-semibold px-1.5 py-0.5 rounded bg-amber-500/12 text-amber-600 border-amber-500/20"
        >
          {{ product.artist }}
        </AppBadge>
        <AppBadge 
          v-if="product.product_type" 
          variant="outline"
          class="text-[0.65rem] font-semibold px-1.5 py-0.5 rounded bg-pink-500/12 text-pink-600 border-pink-500/20"
        >
          {{ product.product_type }}
        </AppBadge>
      </div>

      <p class="text-[0.7rem] text-slate-500 dark:text-slate-400 mb-2" v-if="product.barcode">
        SKU: {{ product.barcode }}
      </p>
      
      <div class="flex items-center justify-between mt-2">
        <div class="flex items-baseline gap-0.5">
          <span class="text-[0.7rem] font-bold text-indigo-600 dark:text-indigo-400">Rp</span>
          <span class="text-[0.95rem] font-extrabold text-indigo-600 dark:text-indigo-400">{{ formatPrice(product.price) }}</span>
        </div>
        
        <AppButton 
          variant="primary"
          :disabled="product.stock <= 0"
          @click.stop="addToCart"
          class="w-[30px] h-[30px] p-0 flex items-center justify-center rounded-lg shadow-sm hover:scale-105 transition-transform"
        >
          <PlusIcon class="w-4 h-4 text-white" />
        </AppButton>
      </div>
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
