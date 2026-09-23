import { defineStore } from 'pinia';
import { ref } from 'vue';
import type { StoreSetting } from '../types';

// Pengganti prop-drilling :store-setting ke 6 file.
// Diisi sekali di App.vue onMounted, dibaca reaktif di mana saja.
export const useSettingsStore = defineStore('settings', () => {
  const settings = ref<StoreSetting>({
    store_name: 'KASIR COFFEE & BISTRO',
    address: 'Jl. Boulevard Utama No. 88, Jakarta',
    phone: '0812-9876-5432',
    receipt_footer: 'Terima kasih telah berbelanja!',
    tax_percentage: 10,
    db_engine: 'sqlite',
  });
  const loaded = ref(false);

  async function fetchSettings(): Promise<void> {
    try {
      const res = await fetch('/api/settings');
      if (res.ok) {
        settings.value = (await res.json()) as StoreSetting;
        loaded.value = true;
      }
    } catch (err) {
      console.error('Fetch store settings error:', err);
    }
  }

  return { settings, loaded, fetchSettings };
});
