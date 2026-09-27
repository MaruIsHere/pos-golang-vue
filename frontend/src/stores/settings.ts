import { defineStore } from 'pinia';
import { ref } from 'vue';
import api from '@/utils/api';
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
      const res = await api.get('/settings');
      settings.value = res.data as StoreSetting;
      loaded.value = true;
    } catch (err: any) {
      console.error('Fetch store settings error:', err.response?.data?.error || err.message || 'Error occurred');
    }
  }

  return { settings, loaded, fetchSettings };
});
