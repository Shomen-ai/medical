<!--
  Файл: pages/staff-login.vue
  Назначение: страница входа сотрудников (администратора и врача) по логину и паролю; после успешного входа перенаправляет в соответствующий портал.
-->
<script setup lang="ts">
const auth = useAuthStore()
const router = useRouter()
const { post } = useApi()

const username = ref('')
const password = ref('')
const loading = ref(false)
const error = ref('')

onMounted(() => {
  auth.init()
  if (auth.isLoggedIn) {
    if (auth.isAdmin) router.replace('/admin')
    else if (auth.isDoctor) router.replace('/doctor')
  }
})

const handleLogin = async () => {
  if (!username.value || !password.value) return
  loading.value = true
  error.value = ''
  try {
    const resp = await post<{ access_token: string; role: string }>(
      '/api/staff/auth/login',
      { username: username.value, password: password.value }
    )
    auth.token = resp.access_token
    auth.role = resp.role as 'admin' | 'doctor'
    if (import.meta.client) {
      localStorage.setItem('auth_token', resp.access_token)
      localStorage.setItem('auth_role', resp.role)
    }
    if (resp.role === 'admin') router.push('/admin')
    else router.push('/doctor')
  } catch {
    error.value = 'Неверный логин или пароль'
  } finally {
    loading.value = false
  }
}

useHead({ title: 'Вход для персонала — BeautyMed' })
</script>

<template>
  <div class="min-h-[80vh] flex items-center justify-center px-4">
    <div class="bg-white rounded-2xl shadow-xl w-full max-w-sm p-8">
      <!-- Logo -->
      <div class="text-center mb-6">
        <span
          class="text-2xl font-extrabold"
          style="background: linear-gradient(135deg,#005A5F,#00959D); -webkit-background-clip:text; -webkit-text-fill-color:transparent;"
        >
          BeautyMed
        </span>
        <p class="text-sm text-gray-500 mt-1">Портал для персонала клиники</p>
      </div>

      <form @submit.prevent="handleLogin" class="space-y-4">
        <div>
          <label class="text-xs font-semibold text-slate block mb-1.5">Логин</label>
          <input
            v-model="username"
            type="text"
            autocomplete="username"
            placeholder="admin"
            class="w-full border border-border rounded-xl px-4 py-2.5 text-sm focus:outline-none focus:border-primary"
          >
        </div>
        <div>
          <label class="text-xs font-semibold text-slate block mb-1.5">Пароль</label>
          <input
            v-model="password"
            type="password"
            autocomplete="current-password"
            placeholder="••••••••"
            class="w-full border border-border rounded-xl px-4 py-2.5 text-sm focus:outline-none focus:border-primary"
          >
        </div>

        <div v-if="error" class="text-xs text-red-500">{{ error }}</div>

        <button
          type="submit"
          class="w-full text-white py-2.5 rounded-xl text-sm font-bold transition-opacity"
          :class="loading ? 'opacity-60 cursor-not-allowed' : 'hover:opacity-90'"
          style="background: linear-gradient(135deg, #005A5F, #00959D)"
          :disabled="loading"
        >
          {{ loading ? 'Входим...' : 'Войти' }}
        </button>
      </form>

      <div class="mt-5 text-center">
        <NuxtLink to="/" class="text-xs text-muted underline">← На главную</NuxtLink>
      </div>
    </div>
  </div>
</template>
