<template>
  <div 
    class="group flex flex-col cursor-pointer transition-all duration-300 relative rounded-[24px] bg-white dark:bg-slate-800 border border-slate-200/50 dark:border-slate-700/50 shadow-sm hover:-translate-y-1 hover:shadow-lg hover:border-indigo-400/50 dark:hover:border-indigo-500/50 h-full"
    :class="{ 'opacity-50 grayscale-[50%] cursor-not-allowed hover:translate-y-0 hover:shadow-sm': product.stock <= 0 }"
    @click="addToCart"
  >
    <!-- Kontainer Gambar dengan Rasio Paten 4:3 -->
    <div class="relative w-full h-[140px] sm:h-[150px] p-2 shrink-0">
      <div class="w-full h-full overflow-hidden rounded-[18px] bg-slate-100 dark:bg-slate-900 shadow-inner relative">
        <img 
          :src="product.image_url || 'https://images.unsplash.com/photo-1546069901-ba9599a7e63c?w=400'" 
          :alt="product.name" 
          class="w-full h-full object-cover transition-transform duration-500 group-hover:scale-110"
          @error="onImageError"
        />
      </div>
      
      <span v-if="product.is_master" class="absolute top-4 left-4 text-[0.68rem] font-extrabold px-3 py-1 rounded-full bg-amber-500 text-white shadow-md z-10">
        Master Pusat
      </span>
      <span v-else-if="product.outlet" class="absolute top-4 left-4 text-[0.68rem] font-extrabold px-3 py-1 rounded-full bg-indigo-600 text-white shadow-md z-10">
        {{ product.outlet.name }}
      </span>

      <AppBadge 
        class="absolute top-4 right-4 text-[0.65rem] px-2.5 py-1 font-bold shadow-sm backdrop-blur-md"
        :variant="stockBadgeVariant"
      >
        {{ product.stock > 0 ? `Stok: ${formatQuantity(product.stock)} ${unitLabel(product.unit)}` : 'Habis' }}
      </AppBadge>
    </div>
    
    <!-- Kontainer Teks yang Fleksibel -->
    <div class="flex flex-col px-4 pb-4 pt-2 flex-1">
      <div class="flex items-center gap-1.5 mb-1.5 shrink-0">
        <span class="text-[10px] font-bold text-slate-400 dark:text-slate-500 uppercase tracking-wider line-clamp-1">
          {{ product.category?.name || 'Umum' }}
        </span>
      </div>
      
      <!-- Judul dikunci ukurannya setara 2 baris teks (sekitar 2.5rem) agar kartu simetris -->
      <h3 class="font-bold text-slate-800 dark:text-slate-100 text-[0.95rem] leading-tight line-clamp-2 h-[2.5rem] shrink-0">
        {{ product.name }}
      </h3>
      
      <!-- Harga didorong ke paling bawah agar sejajar semua -->
      <p class="text-[0.95rem] font-extrabold text-indigo-600 dark:text-indigo-400 mt-auto pt-2">
        Rp {{ formatPrice(product.price) }} / {{ unitLabel(product.unit) }}
      </p>
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
const formatQuantity = (val: number): string => new Intl.NumberFormat('id-ID', { maximumFractionDigits: 3 }).format(val || 0);
const unitLabel = (unit: Product['unit']): string => unit === 'gram' ? 'gr' : unit === 'liter' ? 'L' : 'pcs';

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
