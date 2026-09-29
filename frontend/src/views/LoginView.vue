<script setup lang="ts">
import { onBeforeUnmount, ref } from 'vue'
import { useRouter } from 'vue-router'
import api from '@/utils/api'
import { useAuthStore } from '../stores/auth'
import AppInput from '../components/ui/AppInput.vue'
import AppButton from '../components/ui/AppButton.vue'
<<<<<<< HEAD
import { ExclamationTriangleIcon, XMarkIcon } from '@heroicons/vue/24/solid'
=======
import { LockClosedIcon, UserIcon } from '@heroicons/vue/24/outline'
>>>>>>> 9014a46523b945cdafe66acb0a0a95ecf0650ff4

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
    
    const userRole = String(data.user?.role || '').toLowerCase()
    if (userRole === 'owner' || userRole === 'admin') {
      router.push('/reports')
    } else {
      router.push('/')
    }
  } catch (err: any) {
<<<<<<< HEAD
    error.value = err.response?.data?.error || err.message || 'Gagal login'
    errorTimeout = setTimeout(dismissError, 5000)
=======
    error.value = err.response?.data?.error || err.message || 'Username atau password salah!'
>>>>>>> 9014a46523b945cdafe66acb0a0a95ecf0650ff4
  } finally {
    isLoading.value = false
  }
}
</script>

<template>
<<<<<<< HEAD
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
=======
  <div class="min-h-screen w-full flex bg-slate-50 dark:bg-slate-950 font-sans">
    
    <!-- LEFT SIDE: BRANDING & IMAGE (Hidden on mobile) -->
    <div class="hidden lg:flex lg:w-1/2 relative bg-indigo-900 overflow-hidden items-center justify-center">
      <div class="absolute inset-0 bg-gradient-to-br from-indigo-600 to-indigo-900 opacity-90 z-10"></div>
      <img src="https://images.unsplash.com/photo-1556742049-0cfed4f6a45d?q=80&w=1000&auto=format&fit=crop" class="absolute inset-0 w-full h-full object-cover mix-blend-overlay" />
      
      <div class="relative z-20 flex flex-col p-12 text-white max-w-lg">
        <div class="w-16 h-16 bg-white/20 backdrop-blur-md rounded-2xl flex items-center justify-center mb-8 border border-white/30 shadow-xl">
          <svg class="w-8 h-8 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
          </svg>
        </div>
        <h1 class="text-4xl md:text-5xl font-black tracking-tight mb-4 leading-tight">Sistem POS<br/>Kasir Modern</h1>
        <p class="text-indigo-100 text-lg md:text-xl font-medium">Kelola penjualan, pantau stok, dan pantau laporan bisnis Anda dengan cepat dan elegan.</p>
>>>>>>> 9014a46523b945cdafe66acb0a0a95ecf0650ff4
      </div>
    </div>

    <!-- RIGHT SIDE: LOGIN FORM -->
    <div class="w-full lg:w-1/2 flex items-center justify-center p-6 sm:p-12">
      <div class="w-full max-w-md flex flex-col bg-white dark:bg-slate-900 p-8 sm:p-10 rounded-[32px] shadow-2xl border border-slate-100 dark:border-slate-800">
        
        <div class="flex flex-col items-center text-center mb-8">
          <div class="w-12 h-12 bg-indigo-50 dark:bg-indigo-900/30 rounded-xl flex items-center justify-center mb-4 lg:hidden">
            <svg class="w-6 h-6 text-indigo-600 dark:text-indigo-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
            </svg>
          </div>
          <h2 class="text-3xl font-black text-slate-900 dark:text-white tracking-tight">Selamat Datang</h2>
          <p class="text-slate-500 dark:text-slate-400 mt-2 font-medium">Silakan masuk menggunakan kredensial Anda.</p>
        </div>

        <form @submit.prevent="handleLogin" class="flex flex-col gap-5">
          <div class="flex flex-col gap-2">
            <label class="text-sm font-bold text-slate-700 dark:text-slate-300">Username</label>
            <div class="relative">
              <AppInput v-model="username" type="text" required placeholder="admin" class="pl-11 py-3 w-full rounded-2xl bg-slate-50 dark:bg-slate-800 border-slate-200 dark:border-slate-700" />
              <UserIcon class="w-5 h-5 absolute left-4 top-1/2 -translate-y-1/2 text-slate-400" />
            </div>
          </div>

<<<<<<< HEAD
        <AppButton type="submit" class="w-full" :isLoading="isLoading">
          Login
        </AppButton>
      </form>
    </AppCard>
=======
          <div class="flex flex-col gap-2">
            <label class="text-sm font-bold text-slate-700 dark:text-slate-300">Password</label>
            <div class="relative">
              <AppInput v-model="password" type="password" required placeholder="••••••••" class="pl-11 py-3 w-full rounded-2xl bg-slate-50 dark:bg-slate-800 border-slate-200 dark:border-slate-700" />
              <LockClosedIcon class="w-5 h-5 absolute left-4 top-1/2 -translate-y-1/2 text-slate-400" />
            </div>
          </div>

          <div v-if="error" class="p-4 bg-rose-50 dark:bg-rose-900/30 border border-rose-100 dark:border-rose-800 rounded-xl flex items-center gap-3 mt-2">
            <svg class="w-5 h-5 text-rose-600 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>
            <span class="text-sm font-bold text-rose-600 dark:text-rose-400">{{ error }}</span>
          </div>

          <AppButton type="submit" variant="primary" class="w-full py-4 rounded-2xl mt-4 font-black shadow-lg shadow-indigo-600/20 active:scale-[0.98] transition-transform" :isLoading="isLoading">
            Masuk ke Sistem
          </AppButton>
        </form>
        
      </div>
    </div>
>>>>>>> 9014a46523b945cdafe66acb0a0a95ecf0650ff4
  </div>
</template>
