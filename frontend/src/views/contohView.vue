<template>
  <!-- Container Utama: Mengunci tinggi layar agar tidak scroll penuh -->
  <div class="h-screen w-screen flex bg-md-surface text-md-on-surface overflow-hidden font-sans">
    
    <!-- 1. AREA KIRI: Grid Produk & Kategori -->
    <div class="flex-1 flex flex-col h-full border-r border-md-outline-variant">
      
      <!-- Header & Pencarian -->
      <div class="p-4 bg-md-surface-container flex gap-4 items-center shadow-sm">
        <input 
          type="text" 
          placeholder="Cari produk atau pindai barcode..." 
          class="flex-grow px-4 py-3 bg-md-surface rounded-md-full border border-md-outline focus:outline-none focus:border-md-primary text-md-body-large"
        />
        <button class="bg-md-secondary-container text-md-on-secondary-container p-3 rounded-md-full">
          🔍
        </button>
      </div>

      <!-- Grid Produk (Scrollable) -->
      <div class="flex-1 overflow-y-auto p-4 grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4 bg-md-surface-container-lowest">
        <!-- Item Produk (Looping v-for) -->
        <div 
          v-for="i in 12" :key="i"
          class="bg-md-surface-container flex flex-col justify-between p-4 rounded-md-large shadow-md-1 hover:shadow-md-2 active:scale-98 transition-all cursor-pointer border border-transparent hover:border-md-primary"
        >
          <div class="font-semibold text-md-title-medium line-clamp-2">Kopi Susu Gula Aren {{ i }}</div>
          <div class="mt-4 flex justify-between items-center">
            <span class="text-md-label-medium text-md-on-surface-variant">Stok: 24</span>
            <span class="font-bold text-md-title-medium text-md-primary">Rp 18.000</span>
          </div>
        </div>
      </div>
    </div>

    <!-- 2. AREA KANAN: Keranjang Belanja & Ringkasan Pembayaran -->
    <div class="w-96 flex flex-col h-full bg-md-surface-container-low shadow-lg">
      <!-- Header Keranjang -->
      <div class="p-4 border-b border-md-outline-variant flex justify-between items-center bg-md-surface-container-high">
        <h2 class="text-md-title-large font-bold">Keranjang</h2>
        <span class="bg-md-primary text-md-on-primary px-2.5 py-0.5 rounded-full text-xs font-medium">3 Item</span>
      </div>

      <!-- Daftar Item Belanja (Scrollable) -->
      <div class="flex-1 overflow-y-auto p-4 flex flex-col gap-3">
        <div v-for="i in 3" :key="i" class="flex justify-between items-center p-3 bg-md-surface rounded-md-medium shadow-sm">
          <div>
            <div class="font-medium text-md-body-large">Kopi Susu Gula Aren</div>
            <div class="text-xs text-md-on-surface-variant">Rp 18.000 x 2</div>
          </div>
          <div class="flex items-center gap-2">
            <button class="w-8 h-8 rounded-full bg-md-surface-container-high flex items-center justify-center font-bold">-</button>
            <span class="font-semibold text-md-body-large">2</span>
            <button class="w-8 h-8 rounded-full bg-md-surface-container-high flex items-center justify-center font-bold">+</button>
          </div>
        </div>
      </div>

      <!-- Ringkasan Total & Tombol Bayar -->
      <div class="p-4 bg-md-surface-container-high border-t border-md-outline-variant flex flex-col gap-4">
        <div class="flex justify-between text-md-body-large">
          <span class="text-md-on-surface-variant">Subtotal</span>
          <span>Rp 54.000</span>
        </div>
        <div class="flex justify-between text-md-body-large">
          <span class="text-md-on-surface-variant">Pajak (11%)</span>
          <span>Rp 5.940</span>
        </div>
        <div class="flex justify-between text-md-headline-small font-bold border-t border-md-outline-variant pt-2">
          <span>Total</span>
          <span class="text-md-primary">Rp 59.940</span>
        </div>
        
        <!-- Tombol Aksi Utama Kasir -->
        <button class="w-full bg-md-primary text-md-on-primary py-4 rounded-md-full text-md-title-large font-bold shadow-md-2 active:scale-95 transition-all mt-2">
          BAYAR (F12)
        </button>
      </div>
    </div>

  </div>
</template>
