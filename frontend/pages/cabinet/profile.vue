<script setup lang="ts">
const auth = useAuthStore()
const router = useRouter()
const { get, patch } = useApi()

onMounted(() => {
  auth.init()
  if (!auth.isLoggedIn || auth.isAdmin || auth.isDoctor) {
    router.replace('/')
  }
})

interface Profile {
  id: string
  phone: string
  full_name: string
  birth_date: string | null
  email: string | null
  inn: string | null
  passport_series: string | null
  passport_number: string | null
  passport_issued_at: string | null
  passport_issued_by: string | null
}

const form = reactive({
  full_name: '',
  birth_date: '',
  email: '',
  inn: '',
  passport_series: '',
  passport_number: '',
  passport_issued_at: '',
  passport_issued_by: '',
})

const loading = ref(true)
const saving = ref(false)
const saved = ref(false)
const error = ref('')

const toInputDate = (iso: string | null) => iso ? iso.slice(0, 10) : ''

const loadProfile = async () => {
  if (!auth.token) return
  loading.value = true
  try {
    const p = await get<Profile>('/api/cabinet/profile', auth.token)
    form.full_name = p.full_name ?? ''
    form.birth_date = toInputDate(p.birth_date)
    form.email = p.email ?? ''
    form.inn = p.inn ?? ''
    form.passport_series = p.passport_series ?? ''
    form.passport_number = p.passport_number ?? ''
    form.passport_issued_at = toInputDate(p.passport_issued_at)
    form.passport_issued_by = p.passport_issued_by ?? ''
  } catch {
    error.value = 'Не удалось загрузить профиль'
  } finally {
    loading.value = false
  }
}

onMounted(loadProfile)

const save = async () => {
  if (!auth.token) return
  saving.value = true
  saved.value = false
  error.value = ''
  try {
    await patch(
      '/api/cabinet/profile',
      {
        full_name: form.full_name.trim(),
        birth_date: form.birth_date || null,
        email: form.email.trim() || null,
        inn: form.inn.trim() || null,
        passport_series: form.passport_series.trim() || null,
        passport_number: form.passport_number.trim() || null,
        passport_issued_at: form.passport_issued_at || null,
        passport_issued_by: form.passport_issued_by.trim() || null,
      },
      auth.token
    )
    saved.value = true
    setTimeout(() => { saved.value = false }, 2500)
  } catch {
    error.value = 'Не удалось сохранить'
  } finally {
    saving.value = false
  }
}

useHead({ title: 'Профиль — BeautyMed' })
</script>

<template>
  <div class="max-w-2xl mx-auto px-4 py-10">
    <NuxtLink to="/cabinet" class="inline-flex items-center text-sm text-muted hover:text-primary mb-4">
      ← К записям
    </NuxtLink>

    <h1 class="text-2xl font-bold text-slate mb-2">Профиль</h1>
    <p class="text-sm text-muted mb-8">
      Личные данные для записи и подготовки справок об оплате медицинских услуг.
      Паспортные данные используются при оформлении именной справки.
    </p>

    <div v-if="loading" class="text-center py-16 text-muted">Загружаем профиль...</div>

    <form v-else class="bg-white rounded-2xl border border-border shadow-sm p-6 space-y-4" @submit.prevent="save">
      <div>
        <label class="text-xs font-semibold text-slate block mb-1.5">ФИО</label>
        <input
          v-model="form.full_name"
          type="text"
          placeholder="Иванов Иван Иванович"
          class="w-full border border-border rounded-lg px-3 py-2 text-sm text-slate outline-none focus:border-primary"
        >
      </div>

      <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
        <div>
          <label class="text-xs font-semibold text-slate block mb-1.5">Дата рождения</label>
          <input
            v-model="form.birth_date"
            type="date"
            class="w-full border border-border rounded-lg px-3 py-2 text-sm text-slate outline-none focus:border-primary"
          >
        </div>
        <div>
          <label class="text-xs font-semibold text-slate block mb-1.5">Email</label>
          <input
            v-model="form.email"
            type="email"
            placeholder="you@example.com"
            class="w-full border border-border rounded-lg px-3 py-2 text-sm text-slate outline-none focus:border-primary"
          >
        </div>
      </div>

      <div class="pt-3 border-t border-border">
        <div class="text-sm font-semibold text-slate mb-3">Документы для именной справки</div>

        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <div>
            <label class="text-xs font-semibold text-slate block mb-1.5">ИНН</label>
            <input
              v-model="form.inn"
              type="text"
              maxlength="12"
              placeholder="123456789012"
              class="w-full border border-border rounded-lg px-3 py-2 text-sm text-slate outline-none focus:border-primary"
            >
          </div>
          <div /> <!-- spacer -->

          <div>
            <label class="text-xs font-semibold text-slate block mb-1.5">Паспорт: серия</label>
            <input
              v-model="form.passport_series"
              type="text"
              maxlength="6"
              placeholder="1234"
              class="w-full border border-border rounded-lg px-3 py-2 text-sm text-slate outline-none focus:border-primary"
            >
          </div>
          <div>
            <label class="text-xs font-semibold text-slate block mb-1.5">Паспорт: номер</label>
            <input
              v-model="form.passport_number"
              type="text"
              maxlength="10"
              placeholder="567890"
              class="w-full border border-border rounded-lg px-3 py-2 text-sm text-slate outline-none focus:border-primary"
            >
          </div>

          <div>
            <label class="text-xs font-semibold text-slate block mb-1.5">Дата выдачи паспорта</label>
            <input
              v-model="form.passport_issued_at"
              type="date"
              class="w-full border border-border rounded-lg px-3 py-2 text-sm text-slate outline-none focus:border-primary"
            >
          </div>
          <div>
            <label class="text-xs font-semibold text-slate block mb-1.5">Кем выдан</label>
            <input
              v-model="form.passport_issued_by"
              type="text"
              placeholder="ОВД..."
              class="w-full border border-border rounded-lg px-3 py-2 text-sm text-slate outline-none focus:border-primary"
            >
          </div>
        </div>
      </div>

      <div v-if="error" class="text-sm text-red-500">{{ error }}</div>
      <div v-else-if="saved" class="text-sm text-emerald-600">✓ Сохранено</div>

      <button
        type="submit"
        class="w-full bg-primary text-white py-2.5 rounded-lg text-sm font-semibold transition-opacity"
        :class="saving ? 'opacity-60 cursor-not-allowed' : 'hover:bg-primary/90'"
        :disabled="saving"
      >
        {{ saving ? 'Сохраняем...' : 'Сохранить' }}
      </button>
    </form>
  </div>
</template>
