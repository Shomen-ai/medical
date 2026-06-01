<!--
  Файл: pages/cabinet/profile.vue
  Назначение: страница редактирования профиля пациента (ФИО, дата рождения, пол, адрес, контакты, удостоверение личности) в личном кабинете.
-->
<script setup lang="ts">
const auth = useAuthStore()
const router = useRouter()
const { get, patch } = useApi()
const { t } = useI18n()

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
  gender: 'm' | 'f' | null
  address: string | null
  id_doc_number: string | null
  id_doc_issued_by: string | null
}

const form = reactive({
  full_name: '',
  birth_date: '',
  email: '',
  gender: '' as '' | 'm' | 'f',
  address: '',
  id_doc_number: '',
  id_doc_issued_by: '',
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
    form.gender = (p.gender ?? '') as '' | 'm' | 'f'
    form.address = p.address ?? ''
    form.id_doc_number = p.id_doc_number ?? ''
    form.id_doc_issued_by = p.id_doc_issued_by ?? ''
  } catch {
    error.value = t('profLoadError')
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
        gender: form.gender || null,
        address: form.address.trim() || null,
        id_doc_number: form.id_doc_number.trim() || null,
        id_doc_issued_by: form.id_doc_issued_by.trim() || null,
      },
      auth.token
    )
    saved.value = true
    setTimeout(() => { saved.value = false }, 2500)
  } catch {
    error.value = t('profSaveError')
  } finally {
    saving.value = false
  }
}

useHead({ title: t('profPageTitle') })
</script>

<template>
  <div class="max-w-2xl mx-auto px-4 py-10">
    <NuxtLink to="/cabinet" class="inline-flex items-center text-sm text-muted hover:text-primary mb-4">
      {{ t('profBackToAppts') }}
    </NuxtLink>

    <h1 class="text-2xl font-bold text-slate mb-2">{{ t('profTitle') }}</h1>
    <p class="text-sm text-muted mb-8">
      {{ t('profIntro') }}
    </p>

    <div v-if="loading" class="text-center py-16 text-muted">{{ t('profLoading') }}</div>

    <form v-else class="bg-white rounded-2xl border border-border shadow-sm p-6 space-y-4" @submit.prevent="save">
      <div>
        <label class="text-xs font-semibold text-slate block mb-1.5">{{ t('profFullName') }}</label>
        <input
          v-model="form.full_name"
          type="text"
          placeholder="Babaýew Begenç Mämmedowiç"
          class="w-full border border-border rounded-lg px-3 py-2 text-sm text-slate outline-none focus:border-primary"
        >
      </div>

      <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
        <div>
          <label class="text-xs font-semibold text-slate block mb-1.5">{{ t('profBirthDate') }}</label>
          <input
            v-model="form.birth_date"
            type="date"
            class="w-full border border-border rounded-lg px-3 py-2 text-sm text-slate outline-none focus:border-primary"
          >
        </div>
        <div>
          <label class="text-xs font-semibold text-slate block mb-1.5">{{ t('profGender') }}</label>
          <div class="flex items-center gap-4 pt-1.5">
            <label class="inline-flex items-center gap-1.5 text-sm text-slate cursor-pointer">
              <input v-model="form.gender" type="radio" value="m" class="accent-primary">
              {{ t('profGenderMale') }}
            </label>
            <label class="inline-flex items-center gap-1.5 text-sm text-slate cursor-pointer">
              <input v-model="form.gender" type="radio" value="f" class="accent-primary">
              {{ t('profGenderFemale') }}
            </label>
          </div>
        </div>
      </div>

      <div>
        <label class="text-xs font-semibold text-slate block mb-1.5">{{ t('profEmail') }}</label>
        <input
          v-model="form.email"
          type="email"
          placeholder="you@example.com"
          class="w-full border border-border rounded-lg px-3 py-2 text-sm text-slate outline-none focus:border-primary"
        >
      </div>

      <div>
        <label class="text-xs font-semibold text-slate block mb-1.5">{{ t('profAddress') }}</label>
        <input
          v-model="form.address"
          type="text"
          :placeholder="t('profAddressPlaceholder')"
          class="w-full border border-border rounded-lg px-3 py-2 text-sm text-slate outline-none focus:border-primary"
        >
      </div>

      <div class="pt-3 border-t border-border">
        <div class="text-sm font-semibold text-slate mb-3">{{ t('profIdDoc') }}</div>

        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <div>
            <label class="text-xs font-semibold text-slate block mb-1.5">{{ t('profPassportNumber') }}</label>
            <input
              v-model="form.id_doc_number"
              type="text"
              maxlength="20"
              placeholder="I-AG № 1234567"
              class="w-full border border-border rounded-lg px-3 py-2 text-sm text-slate outline-none focus:border-primary"
            >
          </div>
          <div>
            <label class="text-xs font-semibold text-slate block mb-1.5">{{ t('profIssuedBy') }}</label>
            <input
              v-model="form.id_doc_issued_by"
              type="text"
              :placeholder="t('profIssuedByPlaceholder')"
              class="w-full border border-border rounded-lg px-3 py-2 text-sm text-slate outline-none focus:border-primary"
            >
          </div>
        </div>
      </div>

      <div v-if="error" class="text-sm text-red-500">{{ error }}</div>
      <div v-else-if="saved" class="text-sm text-emerald-600">{{ t('profSaved') }}</div>

      <button
        type="submit"
        class="w-full bg-primary text-white py-2.5 rounded-lg text-sm font-semibold transition-opacity"
        :class="saving ? 'opacity-60 cursor-not-allowed' : 'hover:bg-primary/90'"
        :disabled="saving"
      >
        {{ saving ? t('profSaving') : t('profSave') }}
      </button>
    </form>
  </div>
</template>
