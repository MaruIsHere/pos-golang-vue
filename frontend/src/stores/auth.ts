import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import router from '../router'

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string | null>(localStorage.getItem('token'))
  let savedUser = null
  try {
    const raw = localStorage.getItem('user')
    if (raw && raw !== 'undefined') savedUser = JSON.parse(raw)
  } catch (e) {}
  const user = ref<any>(savedUser)

  const isAuthenticated = computed(() => !!token.value)
  const userRole = computed(() => { const r = user.value?.role || user.value?.Role; return r ? String(r).toLowerCase() : null })

  const setAuth = (newToken: string, newUser: any) => {
    token.value = newToken
    user.value = newUser
    localStorage.setItem('token', newToken)
    localStorage.setItem('user', JSON.stringify(newUser))
  }

  const updateProfile = (newToken: string, newUser: any) => {
    token.value = newToken
    user.value = newUser
    localStorage.setItem('token', newToken)
    localStorage.setItem('user', JSON.stringify(newUser))
  }

  const updateUser = (newUser: any) => {
    user.value = newUser
    localStorage.setItem('user', JSON.stringify(newUser))
  }

  const logout = () => {
    token.value = null
    user.value = null
    localStorage.removeItem('token')
    localStorage.removeItem('user')
    if (router.currentRoute.value.path !== '/login') {
      window.location.href = '/login'
    }
  }

  return {
    token,
    user,
    isAuthenticated,
    userRole,
    setAuth,
    updateProfile,
    updateUser,
    logout
  }
})
