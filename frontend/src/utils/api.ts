import { useAuthStore } from '../stores/auth'

export interface ApiOptions extends RequestInit {
  params?: Record<string, string | number | boolean | undefined | null>
}

export interface ApiResponse<T = any> {
  data: T
  status: number
  ok: boolean
}

async function request<T = any>(endpoint: string, options: ApiOptions = {}): Promise<ApiResponse<T>> {
  const authStore = useAuthStore()
  
  let url = endpoint.startsWith('/api') ? endpoint : `/api${endpoint.startsWith('/') ? '' : '/'}${endpoint}`

  if (options.params) {
    const searchParams = new URLSearchParams()
    Object.entries(options.params).forEach(([key, val]) => {
      if (val !== undefined && val !== null) {
        searchParams.append(key, String(val))
      }
    })
    const queryString = searchParams.toString()
    if (queryString) {
      url += (url.includes('?') ? '&' : '?') + queryString
    }
  }

  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    ...((options.headers as Record<string, string>) || {})
  }

  if (authStore.token) {
    headers['Authorization'] = `Bearer ${authStore.token}`
  }

  const res = await fetch(url, {
    ...options,
    headers
  })

  let data: any = null
  const contentType = res.headers.get('content-type')
  if (contentType && contentType.includes('application/json')) {
    try {
      data = await res.json()
    } catch {
      data = null
    }
  } else {
    data = await res.text()
  }

  if (!res.ok) {
    if (res.status === 401 && !endpoint.includes('/auth/')) {
      authStore.logout()
    }
    const errorObj: any = new Error((typeof data === 'object' && data?.error) || `HTTP error ${res.status}`)
    errorObj.response = { status: res.status, data }
    throw errorObj
  }

  return { data, status: res.status, ok: res.ok }
}

export const api = {
  get: <T = any>(url: string, options: ApiOptions = {}) => request<T>(url, { ...options, method: 'GET' }),
  post: <T = any>(url: string, body?: any, options: ApiOptions = {}) => 
    request<T>(url, { ...options, method: 'POST', body: body ? JSON.stringify(body) : undefined }),
  put: <T = any>(url: string, body?: any, options: ApiOptions = {}) => 
    request<T>(url, { ...options, method: 'PUT', body: body ? JSON.stringify(body) : undefined }),
  delete: <T = any>(url: string, options: ApiOptions = {}) => request<T>(url, { ...options, method: 'DELETE' })
}

export default api
