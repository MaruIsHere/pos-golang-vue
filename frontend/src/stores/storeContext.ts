import { defineStore } from 'pinia';
import { ref, computed } from 'vue';
import api from '@/utils/api';
import type { Outlet } from '@/types';

export const useStoreContextStore = defineStore('storeContext', () => {
  const stores = ref<Outlet[]>([]);
  const activeStoreId = ref<string | null>(
    localStorage.getItem('active_store_id') || null
  );
  const isLoading = ref(false);

  const activeStore = computed(() => {
    if (!activeStoreId.value) return stores.value[0] || null;
    return stores.value.find(s => s.id === activeStoreId.value) || stores.value[0] || null;
  });

  const fetchStores = async () => {
    isLoading.value = true;
    try {
      const res = await api.get('/stores');
      stores.value = res.data;
      if (stores.value.length > 0 && (!activeStoreId.value || !stores.value.some(s => s.id === activeStoreId.value))) {
        setActiveStore(stores.value[0].id);
      }
    } catch (err) {
      console.error('Fetch stores error:', err);
    } finally {
      isLoading.value = false;
    }
  };

  const setActiveStore = (id: string) => {
    activeStoreId.value = id;
    localStorage.setItem('active_store_id', id);
  };

  return {
    stores,
    activeStoreId,
    activeStore,
    isLoading,
    fetchStores,
    setActiveStore
  };
});
