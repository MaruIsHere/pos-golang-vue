<template>
  <div class="app-container">
    <Navbar 
      :current-tab="currentTab" 
      :store-setting="storeSetting" 
      @change-tab="currentTab = $event"
    />

    <main class="main-content">
      <div class="view-container">
        <RegisterView 
          v-if="currentTab === 'register'" 
          :store-setting="storeSetting"
        />

        <ProductsView 
          v-else-if="currentTab === 'products'" 
        />

        <InventoryView 
          v-else-if="currentTab === 'inventory'" 
        />

        <CustomersView 
          v-else-if="currentTab === 'customers'" 
        />

        <OrdersView 
          v-else-if="currentTab === 'orders'" 
          :store-setting="storeSetting"
        />

        <ReportsView 
          v-else-if="currentTab === 'reports'" 
        />

        <SettingsView 
          v-else-if="currentTab === 'settings'" 
          :store-setting="storeSetting"
          @refresh-settings="fetchStoreSettings"
        />
      </div>
    </main>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import Navbar from './components/Navbar.vue';
import RegisterView from './views/RegisterView.vue';
import ProductsView from './views/ProductsView.vue';
import InventoryView from './views/InventoryView.vue';
import CustomersView from './views/CustomersView.vue';
import OrdersView from './views/OrdersView.vue';
import ReportsView from './views/ReportsView.vue';
import SettingsView from './views/SettingsView.vue';

const currentTab = ref('register');
const storeSetting = ref({
  store_name: 'KASIR COFFEE & BISTRO',
  address: 'Jl. Boulevard Utama No. 88, Jakarta',
  phone: '0812-9876-5432',
  receipt_footer: 'Terima kasih telah berbelanja!',
  tax_percentage: 10,
  db_engine: 'sqlite'
});

const fetchStoreSettings = async () => {
  try {
    const res = await fetch('/api/settings');
    if (res.ok) {
      storeSetting.value = await res.json();
    }
  } catch (err) {
    console.error('Fetch store settings error:', err);
  }
};

onMounted(fetchStoreSettings);
</script>

<style scoped>
/* App layout container styled in styles.css */
</style>
