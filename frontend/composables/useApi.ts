// Wrapper de $fetch autenticado, com timeout e retry simples apos refresh.
export function useApi() {
  const config = useRuntimeConfig()
  const auth = useAuth()
  const defaultTimeoutMs = 15000

  async function request<T>(path: string, options: Parameters<typeof $fetch>[1] = {}): Promise<T> {
    const url = `${config.public.apiBaseUrl}${path}`

    const makeRequest = () =>
      $fetch.raw<{ data: T }>(url, {
        timeout: defaultTimeoutMs,
        ...options,
        headers: {
          ...auth.authHeaders(),
          ...(options.headers ?? {}),
        },
      })

    try {
      const res = await makeRequest()
      if (res.status === 204) {
        return undefined as T
      }
      return res._data?.data as T
    } catch (error: any) {
      if (error?.status === 401) {
        const restored = await auth.refresh()
        if (restored) {
          const res = await makeRequest()
          if (res.status === 204) {
            return undefined as T
          }
          return res._data?.data as T
        }
      }

      throw error
    }
  }

  async function get<T>(path: string): Promise<T> {
    return request<T>(path)
  }

  async function list<T>(path: string): Promise<{ data: T[]; meta: { page: number; per_page: number; total: number } }> {
    const url = `${config.public.apiBaseUrl}${path}`
    const makeRequest = () =>
      $fetch.raw<{ data: T[]; meta: { page: number; per_page: number; total: number } }>(url, {
        timeout: defaultTimeoutMs,
        headers: auth.authHeaders(),
      })
    try {
      const res = await makeRequest()
      return res._data!
    } catch (error: any) {
      if (error?.status === 401) {
        const restored = await auth.refresh()
        if (restored) {
          const res = await makeRequest()
          return res._data!
        }
      }
      throw error
    }
  }

  async function post<T>(path: string, body: unknown): Promise<T> {
    return request<T>(path, {
      method: 'POST',
      body,
    })
  }

  async function put<T>(path: string, body: unknown): Promise<T> {
    return request<T>(path, {
      method: 'PUT',
      body,
    })
  }

  async function patch<T>(path: string, body: unknown): Promise<T> {
    return request<T>(path, {
      method: 'PATCH',
      body,
    })
  }

  async function del(path: string): Promise<void> {
    const url = `${config.public.apiBaseUrl}${path}`

    const makeRequest = () =>
      $fetch(url, {
        method: 'DELETE',
        timeout: defaultTimeoutMs,
        headers: auth.authHeaders(),
      })

    try {
      await makeRequest()
    } catch (error: any) {
      if (error?.status === 401) {
        const restored = await auth.refresh()
        if (restored) {
          await makeRequest()
          return
        }
      }

      throw error
    }
  }

  return { get, list, post, put, patch, del }
}
