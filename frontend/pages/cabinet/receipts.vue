<script setup lang="ts">
const auth = useAuthStore()
const router = useRouter()
const config = useRuntimeConfig()

onMounted(() => {
  auth.init()
  if (!auth.isLoggedIn || auth.isAdmin || auth.isDoctor) {
    router.replace('/')
  }
})

const currentYear = new Date().getFullYear()
const yearOptions = Array.from({ length: 5 }, (_, i) => currentYear - i)
const selectedYear = ref(currentYear)
const downloading = ref(false)
const error = ref('')

const downloadReceipt = async () => {
  if (!auth.token) return
  downloading.value = true
  error.value = ''
  try {
    const base = import.meta.server ? config.apiBase : ''
    const res = await fetch(
      `${base}/api/cabinet/receipts?year=${selectedYear.value}`,
      { headers: { Authorization: `Bearer ${auth.token}` } }
    )
    if (!res.ok) {
      throw new Error(String(res.status))
    }
    const blob = await res.blob()
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `spravka_${selectedYear.value}.pdf`
    document.body.appendChild(a)
    a.click()
    document.body.removeChild(a)
    URL.revokeObjectURL(url)
  } catch {
    error.value = 'Не удалось сформировать справку. Возможно, в этом году не было оплаченных услуг.'
  } finally {
    downloading.value = false
  }
}

useHead({ title: 'Справка об оплате услуг — BeautyMed' })
</script>

<template>
  <div class="max-w-2xl mx-auto px-4 py-10">
    <NuxtLink to="/cabinet" class="inline-flex items-center text-sm text-muted hover:text-primary mb-4">
      ← К записям
    </NuxtLink>

    <h1 class="text-2xl font-bold text-slate mb-2">Справка об оплате медицинских услуг</h1>
    <p class="text-sm text-muted mb-8">
      Справка о фактически оплаченных медицинских услугах за выбранный календарный год
      для предоставления по месту требования (страховая компания, работодатель, государственные органы).
      Включает все завершённые приёмы за указанный период.
    </p>

    <div class="bg-white rounded-2xl border border-border shadow-sm p-6 space-y-5">
      <div>
        <label class="text-xs font-semibold text-slate block mb-1.5">Налоговый период (год)</label>
        <select
          v-model="selectedYear"
          class="w-full border border-border rounded-lg px-4 py-2.5 text-sm text-slate outline-none focus:border-primary"
        >
          <option v-for="y in yearOptions" :key="y" :value="y">{{ y }}</option>
        </select>
      </div>

      <div v-if="error" class="text-sm text-red-500">{{ error }}</div>

      <button
        type="button"
        class="w-full bg-primary text-white py-3 rounded-xl text-sm font-semibold transition-opacity"
        :class="downloading ? 'opacity-60 cursor-not-allowed' : 'hover:bg-primary/90'"
        :disabled="downloading"
        @click="downloadReceipt"
      >
        {{ downloading ? 'Формируем PDF...' : 'Скачать справку (PDF)' }}
      </button>

      <p class="text-xs text-muted">
        Файл генерируется в реальном времени и не сохраняется на сервере.
        Если в выбранном году у вас не было оплаченных приёмов, справка не будет создана.
      </p>
    </div>
  </div>
</template>
