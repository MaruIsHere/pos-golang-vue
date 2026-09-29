<script setup lang="ts">
import { onBeforeUnmount, ref } from 'vue'
import { useRouter } from 'vue-router'
import api from '@/utils/api'
import { useAuthStore } from '../stores/auth'
import AppCard from '../components/ui/AppCard.vue'
import AppInput from '../components/ui/AppInput.vue'
import AppButton from '../components/ui/AppButton.vue'
import { ExclamationTriangleIcon, XMarkIcon } from '@heroicons/vue/24/solid'

const router = useRouter()
const authStore = useAuthStore()

const username = ref('')
const password = ref('')
const error = ref('')
const isLoading = ref(false)
let errorTimeout: ReturnType<typeof setTimeout> | undefined

const dismissError = () => {
  error.value = ''
  if (errorTimeout) clearTimeout(errorTimeout)
}

onBeforeUnmount(dismissError)

const handleLogin = async () => {
  dismissError()
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
    errorTimeout = setTimeout(dismissError, 5000)
  } finally {
    isLoading.value = false
  }
}
</script>

<template>
  <div class="relative flex-1 min-h-0 w-full flex items-center justify-center bg-slate-50 dark:bg-slate-900 p-4">
    <Transition
      enter-active-class="transition duration-200 ease-out"
      enter-from-class="opacity-0 -translate-y-2"
      enter-to-class="opacity-100 translate-y-0"
      leave-active-class="transition duration-150 ease-in"
      leave-from-class="opacity-100 translate-y-0"
      leave-to-class="opacity-0 -translate-y-2"
    >
      <div
        v-if="error"
        role="alert"
        aria-live="assertive"
        class="fixed top-4 left-1/2 z-[200] flex w-[calc(100%-2rem)] max-w-md -translate-x-1/2 items-start gap-3 rounded-lg border border-red-200 bg-white px-4 py-3 text-red-800 shadow-xl dark:border-red-900 dark:bg-slate-800 dark:text-red-200"
      >
        <ExclamationTriangleIcon class="mt-0.5 h-5 w-5 shrink-0 text-red-500" />
        <p class="min-w-0 flex-1 text-sm font-medium">{{ error }}</p>
        <button
          type="button"
          class="shrink-0 rounded p-1 text-red-500 hover:bg-red-50 focus:outline-none focus:ring-2 focus:ring-red-400 dark:hover:bg-red-950"
          aria-label="Tutup notifikasi"
          @click="dismissError"
        >
          <XMarkIcon class="h-4 w-4" />
        </button>
      </div>
    </Transition>

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

        <AppButton type="submit" class="w-full" :isLoading="isLoading">
          Login
        </AppButton>
      </form>
    </AppCard>
  </div>
</template>
