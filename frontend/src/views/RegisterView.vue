<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import api from '@/utils/api'
import AppCard from '../components/ui/AppCard.vue'
import AppInput from '../components/ui/AppInput.vue'
import AppButton from '../components/ui/AppButton.vue'

const router = useRouter()

const username = ref('')
const password = ref('')
const role = ref('kasir')
const error = ref('')
const success = ref('')
const isLoading = ref(false)

const handleRegister = async () => {
  error.value = ''
  success.value = ''
  isLoading.value = true
  try {
    const res = await api.post('/auth/register', {
      username: username.value,
      password: password.value,
      role: role.value
    })

    success.value = 'Registration successful. You can now login.'
    setTimeout(() => {
      router.push('/login')
    }, 2000)
  } catch (err: any) {
    error.value = err.response?.data?.error || err.message || 'Gagal mendaftar'
  } finally {
    isLoading.value = false
  }
}
</script>

<template>
  <div class="min-h-screen flex items-center justify-center bg-slate-50 dark:bg-slate-900">
    <AppCard class="w-full max-w-md p-8">
      <div class="text-center mb-6">
        <h1 class="text-2xl font-bold text-slate-900 dark:text-slate-100">Register</h1>
        <p class="text-slate-500 dark:text-slate-400 mt-2">Create new account</p>
      </div>

      <form @submit.prevent="handleRegister" class="space-y-4">
        <div>
          <label class="block text-sm font-medium text-slate-900 dark:text-slate-100 mb-1">Username</label>
          <AppInput v-model="username" type="text" required placeholder="Enter username" />
        </div>

        <div>
          <label class="block text-sm font-medium text-slate-900 dark:text-slate-100 mb-1">Password</label>
          <AppInput v-model="password" type="password" required placeholder="Enter password (min 6 char)" />
        </div>

        <div>
          <label class="block text-sm font-medium text-slate-900 dark:text-slate-100 mb-1">Role</label>
          <select v-model="role" class="w-full bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 text-slate-900 dark:text-slate-100 rounded-md px-3 py-2 outline-none focus:ring-2 focus:ring-accent-primary">
            <option value="kasir">Kasir (Fitur Kasir & Riwayat)</option>
            <option value="kepala_kasir">Kepala Kasir (Kasir, Pelanggan, Barang & Laporan)</option>
            <option value="owner">Owner (Full Akses)</option>
          </select>
        </div>

        <div v-if="error" class="text-accent-danger text-sm text-center">
          {{ error }}
        </div>
        
        <div v-if="success" class="text-green-500 text-sm text-center">
          {{ success }}
        </div>

        <AppButton type="submit" class="w-full" :isLoading="isLoading">
          Register
        </AppButton>
      </form>
      
      <div class="mt-4 text-center text-sm text-slate-500 dark:text-slate-400">
        Already have an account? <router-link to="/login" class="text-accent-primary hover:underline">Login</router-link>
      </div>
    </AppCard>
  </div>
</template>
