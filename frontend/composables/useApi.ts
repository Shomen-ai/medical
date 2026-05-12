// Файл: composables/useApi.ts
// Назначение: обёртка над $fetch с динамическим базовым URL — на SSR используется внутренний docker-адрес api:8080, на клиенте — относительный путь через nginx.

// SSR: uses runtimeConfig.apiBase (internal Docker URL)
// Client: uses '' (relative URL, nginx proxies /api/ to backend)
// Композабл useApi возвращает методы get/post/patch с автоматической вставкой Bearer-токена.
export const useApi = () => {
  const config = useRuntimeConfig()
  const base = import.meta.server ? config.apiBase : ''

  const get = <T>(path: string, token?: string) =>
    $fetch<T>(`${base}${path}`, {
      headers: token ? { Authorization: `Bearer ${token}` } : undefined,
    })

  const post = <T>(path: string, body: unknown, token?: string) =>
    $fetch<T>(`${base}${path}`, {
      method: 'POST',
      body,
      headers: token ? { Authorization: `Bearer ${token}` } : undefined,
    })

  const patch = <T>(path: string, body: unknown, token: string) =>
    $fetch<T>(`${base}${path}`, {
      method: 'PATCH',
      body,
      headers: { Authorization: `Bearer ${token}` },
    })

  return { get, post, patch }
}
