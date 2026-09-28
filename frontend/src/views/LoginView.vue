<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import api from '@/utils/api'
import { useAuthStore } from '../stores/auth'
import AppCard from '../components/ui/AppCard.vue'
import AppInput from '../components/ui/AppInput.vue'
import AppButton from '../components/ui/AppButton.vue'

const router = useRouter()
const authStore = useAuthStore()

const username = ref('')
const password = ref('')
const error = ref('')
const isLoading = ref(false)

const handleLogin = async () => {
  error.value = ''
  isLoading.value = true
  try {
    const res = await api.post('/auth/login', {
      username: username.value,
      password: password.value
    })

    const data = res.data
    authStore.setAuth(data.token, data.user)
    
    const userRole = data.user?.role
    if (userRole === 'owner' || userRole === 'admin') {
      router.push('/reports')
    } else {
      router.push('/')
    }
  } catch (err: any) {
    error.value = err.response?.data?.error || err.message || 'Gagal login'
  } finally {
    isLoading.value = false
  }
}
</script>

<template>
  <div class="min-h-screen flex items-center justify-center bg-slate-50 dark:bg-slate-900 p-4">
    <AppCard class="w-full max-w-md p-6 sm:p-8 shadow-xl">
      <div class="text-center mb-6">
        <h1 class="text-2xl font-bold text-slate-900 dark:text-slate-100">Login</h1>
        <p class="text-slate-500 dark:text-slate-400 mt-2">Welcome to POS System</p>
      </div>

      <form @submit.prevent="handleLogin" class="space-y-4">
        <div>
          <label class="block text-sm font-medium text-slate-900 dark:text-slate-100 mb-1">Username</label>
          <AppInput v-model="username" type="text" required placeholder="Enter username" />
        </div>

        <div>
          <label class="block text-sm font-medium text-slate-900 dark:text-slate-100 mb-1">Password</label>
          <AppInput v-model="password" type="password" required placeholder="Enter password" />
        </div>

        <div v-if="error" class="text-accent-danger text-sm text-center">
          {{ error }}
        </div>

        <AppButton type="submit" class="w-full" :isLoading="isLoading">
          Login
        </AppButton>
      </form>
    </AppCard>
  </div>
</template>
