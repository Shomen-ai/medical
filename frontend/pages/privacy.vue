<!--
  Файл: pages/privacy.vue
  Назначение: страница политики конфиденциальности клиники с текстом обработки персональных данных и контактами.
  Двуязычная: рендерится на русском (locale 'ru') либо на туркменском (locale 'tk').
-->
<script setup lang="ts">
const config = useRuntimeConfig()
const { locale } = useI18n()

const clinicName = config.public.clinicName as string
const clinicAddress = config.public.clinicAddress as string
const clinicPhone = config.public.clinicPhone as string

interface Section {
  h: string
  paras?: string[]
  items?: string[]
}
interface Content {
  title: string
  intro: string
  sections: Section[]
}

const ru: Content = {
  title: 'Политика конфиденциальности',
  intro:
    'Настоящая Политика разработана в соответствии с Законом Туркменистана ' +
    'от 20 марта 2017 года «Об информации о личной жизни и её защите» и определяет ' +
    `порядок обработки персональных данных пользователей сайта ${clinicName}.`,
  sections: [
    {
      h: '1. Владелец информации',
      paras: [
        `Владельцем (оператором), обрабатывающим персональные данные, является клиника «${clinicName}», адрес: ${clinicAddress}, телефон: ${clinicPhone}.`,
      ],
    },
    {
      h: '2. Состав обрабатываемой информации',
      paras: ['Владелец обрабатывает следующие категории сведений о личной жизни:'],
      items: [
        'фамилия, имя, отчество;',
        'дата рождения;',
        'пол;',
        'адрес проживания;',
        'номер мобильного телефона;',
        'адрес электронной почты (при наличии);',
        'номер удостоверения личности (паспорта Туркменистана) и сведения о его выдаче — для ведения медицинской карты;',
        'сведения о записях на приём, оказанных медицинских услугах и оплатах.',
      ],
    },
    {
      h: '3. Цели обработки',
      items: [
        'идентификация пользователя при онлайн-записи и в личном кабинете;',
        'оказание медицинских услуг и связанной с ними коммуникации;',
        'отправка SMS-уведомлений о записи, переносе или отмене приёма;',
        'исполнение требований законодательства Туркменистана.',
      ],
    },
    {
      h: '4. Правовые основания',
      paras: [
        'Обработка информации о личной жизни осуществляется на основании согласия субъекта информации о личной жизни, договора об оказании медицинских услуг и требований действующего законодательства Туркменистана.',
      ],
    },
    {
      h: '5. Срок хранения',
      paras: [
        'Сведения хранятся в течение всего срока действия учётной записи и в течение сроков, установленных законодательством Туркменистана для медицинской документации, после чего удаляются или обезличиваются.',
      ],
    },
    {
      h: '6. Права субъекта',
      paras: ['Пользователь имеет право:'],
      items: [
        'получать сведения об обработке своей информации о личной жизни;',
        'требовать уточнения, блокирования или уничтожения сведений;',
        'отозвать согласие на обработку в любой момент, направив письменный запрос Владельцу.',
      ],
    },
    {
      h: '7. Меры защиты',
      paras: [
        'Владелец применяет организационные и технические меры защиты сведений о личной жизни от неправомерного доступа, изменения, уничтожения или распространения, включая ограничение доступа сотрудников и регулярное резервное копирование.',
      ],
    },
    {
      h: '8. Трансграничная передача',
      paras: [
        'Передача информации о личной жизни за пределы Туркменистана осуществляется только в случаях и порядке, предусмотренных Законом Туркменистана «Об информации о личной жизни и её защите».',
      ],
    },
    {
      h: '9. Контакты',
      paras: [
        `По вопросам обработки информации о личной жизни обращайтесь по телефону ${clinicPhone} или по адресу ${clinicAddress}.`,
      ],
    },
  ],
}

const tk: Content = {
  title: 'Gizlinlik syýasaty',
  intro:
    'Şu Syýasat Türkmenistanyň «Şahsy durmuş barada maglumat we ony goramak hakynda» ' +
    '20.03.2017-nji ýyldaky Kanunyna laýyklykda işlenip taýýarlanyldy hem-de ' +
    `${clinicName} saýtynyň ulanyjylarynyň şahsy durmuş baradaky maglumatlaryny işlemegiň tertibini kesgitleýär.`,
  sections: [
    {
      h: '1. Maglumatyň eýesi',
      paras: [
        `Şahsy durmuş barada maglumaty işleýän eýesi (operator) «${clinicName}» klinikidir, salgysy: ${clinicAddress}, telefon: ${clinicPhone}.`,
      ],
    },
    {
      h: '2. Işlenilýän maglumatyň düzümi',
      paras: ['Eýesi şahsy durmuş barada maglumatyň şu görnüşlerini işleýär:'],
      items: [
        'familiýasy, ady, atasynyň ady;',
        'doglan senesi;',
        'jynsy;',
        'ýaşaýan ýeriniň salgysy;',
        'jübi telefonynyň belgisi;',
        'elektron poçtanyň salgysy (bar bolsa);',
        'şahsyýetnamanyň (Türkmenistanyň pasportynyň) belgisi we onuň berlişi baradaky maglumatlar — saglyk kartasyny ýöretmek üçin;',
        'kabul edişlige ýazylyşlar, edilen lukmançylyk hyzmatlary we tölegler baradaky maglumatlar.',
      ],
    },
    {
      h: '3. Işlemegiň maksatlary',
      items: [
        'onlaýn ýazylyşda we şahsy otagda ulanyjyny anyklamak;',
        'lukmançylyk hyzmatlaryny etmek we olar bilen baglanyşykly habarlaşmak;',
        'kabul edişlige ýazylyş, geçirilmegi ýa-da ýatyrylmagy barada SMS-habarnamalary ibermek;',
        'Türkmenistanyň kanunçylygynyň talaplaryny ýerine ýetirmek.',
      ],
    },
    {
      h: '4. Hukuk esaslary',
      paras: [
        'Şahsy durmuş barada maglumaty işlemek şahsy durmuş barada maglumatyň subýektiniň razylygy, lukmançylyk hyzmatlaryny etmek baradaky şertnama hem-de Türkmenistanyň hereket edýän kanunçylygynyň talaplary esasynda amala aşyrylýar.',
      ],
    },
    {
      h: '5. Saklanyş möhleti',
      paras: [
        'Maglumatlar hasap ýazgysynyň hereket edýän bütin möhletinde hem-de lukmançylyk resminamalary üçin Türkmenistanyň kanunçylygy bilen bellenen möhletlerde saklanýar, soňra bolsa pozulýar ýa-da şahssyzlaşdyrylýar.',
      ],
    },
    {
      h: '6. Subýektiň hukuklary',
      paras: ['Ulanyjynyň şu hukuklary bardyr:'],
      items: [
        'öz şahsy durmuş barada maglumatynyň işlenilişi barada maglumat almak;',
        'maglumatlary takyklamagy, bökdemegi ýa-da ýok etmegi talap etmek;',
        'islendik pursatda işlemäge berlen razylygyny Eýesine ýazmaça arza ýollap yzyna almak.',
      ],
    },
    {
      h: '7. Goraýyş çäreleri',
      paras: [
        'Eýesi şahsy durmuş barada maglumaty bikanun elýeterlilikden, üýtgedilmekden, ýok edilmekden ýa-da ýaýradylmakdan goramak üçin guramaçylyk we tehniki çäreleri, şol sanda işgärleriň elýeterliligini çäklendirmegi we yzygiderli ätiýaçlyk nusgalaryny döretmegi ulanýar.',
      ],
    },
    {
      h: '8. Serhetden aşa geçirmek',
      paras: [
        'Şahsy durmuş barada maglumaty Türkmenistanyň çäklerinden daşaryk geçirmek diňe Türkmenistanyň «Şahsy durmuş barada maglumat we ony goramak hakynda» Kanunynda göz öňünde tutulan halatlarda we tertipde amala aşyrylýar.',
      ],
    },
    {
      h: '9. Habarlaşmak',
      paras: [
        `Şahsy durmuş barada maglumaty işlemek baradaky meseleler boýunça ${clinicPhone} telefony arkaly ýa-da ${clinicAddress} salgysy boýunça ýüz tutuň.`,
      ],
    },
  ],
}

const c = computed<Content>(() => (locale.value === 'tk' ? tk : ru))

useHead({
  title: () => `${c.value.title} — ${clinicName}`,
})
</script>

<template>
  <div class="max-w-3xl mx-auto px-4 py-10 text-slate">
    <h1 class="text-2xl font-bold mb-6">{{ c.title }}</h1>
    <p class="text-sm text-muted mb-6">{{ c.intro }}</p>

    <section
      v-for="section in c.sections"
      :key="section.h"
      class="space-y-3 text-sm mb-6"
    >
      <h2 class="text-lg font-bold">{{ section.h }}</h2>
      <p v-for="(para, i) in section.paras" :key="`p-${i}`">{{ para }}</p>
      <ul v-if="section.items" class="list-disc pl-6 space-y-1">
        <li v-for="(item, i) in section.items" :key="`i-${i}`">{{ item }}</li>
      </ul>
    </section>
  </div>
</template>
