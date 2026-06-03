// Файл: data/serviceDescriptions.ts
// Назначение: расширенные, понятные пациенту описания услуг (ru/tk) для страницы «Подробнее».
// Ключ — каноничное русское название услуги из БД (как в locales/medical.ts).
// Если услуги нет в словаре — страница показывает короткое service.description (fallback).
//
// ВНИМАНИЕ: туркменские тексты — черновой перевод, нуждаются в проверке носителем.

import type { Locale } from '~/composables/useI18n'

export type LocalizedText = Record<Locale, string>

export const SERVICE_DESCRIPTIONS: Record<string, LocalizedText> = {
  // ─────────────── Косметология ───────────────
  'Консультация косметолога': {
    ru: 'Первичный приём, на котором врач-косметолог осматривает кожу, определяет её тип и состояние и выявляет проблемы — сухость, акне, расширенные поры, пигментацию, возрастные изменения. По итогам вы получаете индивидуальную программу процедур и понятные рекомендации по домашнему уходу. Подходит всем, кто хочет начать заботиться о коже осознанно или обновить свой уход.',
    tk: 'Ilkinji kabul ediş: lukman-kosmetolog deriniň görnüşini we ýagdaýyny barlaýar, gurak deri, akne, giňelen deşikler, pigmentasiýa we ýaşa degişli üýtgemeleri ýüze çykarýar. Netijede size şahsy proseduralar meýilnamasy we öýdäki ideg boýunça düşnükli maslahatlar berilýär. Derisine aňly ideg etmek isleýänleriň hemmesine laýyk.',
  },
  'Чистка лица механическая': {
    ru: 'Глубокое очищение кожи: лицо распаривают, затем вручную удаляют чёрные точки, сальные пробки и комедоны, после чего наносят успокаивающую и сужающую поры маску. Процедура освобождает поры, выравнивает рельеф и возвращает коже свежесть. Рекомендована при жирной и проблемной коже; лёгкое покраснение обычно проходит за 1–2 дня.',
    tk: 'Deriniň çuňňur arassalanmagy: ýüz bug bilen ýumşadylýar, soňra gara nokatlar, ýag dykылары we komedonlar elde aýrylýar, ahyrynda köşeşdiriji we deşikleri daraldyjy maska çalynýar. Prosedura deşikleri açýar, deriniň ýüzüni tekizleýär we täzelik berýär. Ýagly we problemaly deri üçin maslahat berilýär; ýeňil gyzarma adatça 1–2 günde geçýär.',
  },
  'Чистка лица ультразвуковая': {
    ru: 'Бережное бесконтактное очищение: ультразвуковые волны мягко удаляют ороговевшие клетки и загрязнения из пор, не травмируя кожу. Подходит для чувствительной и тонкой кожи, не оставляет покраснений — можно делать даже перед важным событием. Кожа становится гладкой и сияющей сразу после процедуры.',
    tk: 'Näzik, galtaşyksyz arassalama: ultrases tolkunlary ölü öýjükleri we hapalary deşiklerden derini şikeslendirmän aýyrýar. Duýgur we inçe deri üçin amatly, gyzarma galdyrmaýar — wajyp wakadan öň hem edip bolýar. Deri proseduradan soň derrew tekiz we şöhleli bolýar.',
  },
  'Химический пилинг': {
    ru: 'Контролируемое обновление кожи кислотами: состав растворяет верхний ороговевший слой и запускает естественную регенерацию. Помогает при пигментации, тусклом цвете лица, мелких морщинах, постакне и расширенных порах. Глубину пилинга врач подбирает индивидуально; после возможны лёгкое шелушение и временная чувствительность к солнцу.',
    tk: 'Deriniň turşulyklar bilen dolandyrylýan täzelenmegi: serişde ýokarky ölü gatlagy eredýär we tebigy dikeldişi işjeňleşdirýär. Pigmentasiýa, solgun reňk, ownuk ýygyrtlar, postakne we giňelen deşiklerde kömek edýär. Pilingiň çuňlugyny lukman şahsy saýlaýar; soňundan ýeňil soýulma we Güne wagtlaýyn duýgurlyk bolup biler.',
  },
  'Инъекции ботулотоксина (зона)': {
    ru: 'Расслабление мимических мышц микродозами ботулотоксина — разглаживает морщины на лбу, между бровями и вокруг глаз («гусиные лапки»). Эффект развивается за 3–7 дней и держится в среднем 4–6 месяцев. Цена указана за одну зону; процедура занимает 15–20 минут и не требует восстановления.',
    tk: 'Mimiki myşsalary botulotoksiniň mikrodozalary bilen gowşatmak — maňlaýdaky, gaşlaryň arasyndaky we göz töweregindäki ýygyrtlary tekizleýär. Netije 3–7 günde ýüze çykýar we ortaça 4–6 aý saklanýar. Baha bir zola üçin; prosedura 15–20 minut alýar we dikeldiş talap etmeýär.',
  },
  'Контурная пластика (1 мл)': {
    ru: 'Введение филлера на основе гиалуроновой кислоты для восстановления объёма и чёткости контуров — губы, скулы, носогубные складки, подбородок. Результат виден сразу и сохраняется от 6 до 12 месяцев. Процедура обратима и проходит под аппликационной анестезией; цена — за 1 мл препарата.',
    tk: 'Göwrümi we konturlaryň aýdyňlygyny dikeltmek üçin gialuron turşusy esasly filleriň goýberilmegi — dodaklar, ýaňaklar, burun-dodak büküm­leri, eňek. Netije derrew görünýär we 6–12 aý saklanýar. Prosedura yzyna öwrülip bolýar we üst anesteziýasy bilen geçýär; baha 1 ml serişde üçin.',
  },
  'Биоревитализация': {
    ru: 'Инъекционное глубокое увлажнение кожи чистой гиалуроновой кислотой. Возвращает тонус, упругость и сияние, разглаживает мелкие морщинки и борется с сухостью и тусклостью. Обычно проводится курсом из нескольких процедур; подходит для лица, шеи и кожи вокруг глаз.',
    tk: 'Derini arassa gialuron turşusy bilen sanjym arkaly çuňňur nemlendirmek. Tonusy, çeýeligi we şöhläni dikeldýär, ownuk ýygyrtlary tekizleýär, guraklyk we solgunlyk bilen göreşýär. Adatça birnäçe proseduradan ybarat kurs bilen geçirilýär; ýüz, boýun we göz töweregi üçin amatly.',
  },
  'Мезотерапия лица': {
    ru: 'Микроинъекции индивидуально подобранного коктейля из витаминов, аминокислот и пептидов прямо в проблемные зоны. Питает кожу изнутри, улучшает цвет лица, повышает упругость и помогает при усталой, обезвоженной коже. Курс из нескольких сеансов даёт накопительный, заметный эффект.',
    tk: 'Witaminlerden, aminokislotalardan we peptidlerden ybarat şahsy saýlanan kokteýliň gönüden-göni problemaly ýerlere mikrosanjymy. Derini içinden iýmitlendirýär, reňkini gowulandyrýar, çeýeligini artdyrýar we ýadaw, suwsuz deride kömek edýär. Birnäçe seansdan ybarat kurs jemleýji, görnükli netije berýär.',
  },
  'PRP-терапия лица': {
    ru: 'Плазмолифтинг — введение собственной обогащённой тромбоцитами плазмы пациента, полученной из его крови. Запускает естественную регенерацию, повышает плотность и упругость кожи, улучшает цвет лица и сужает поры. Полностью натуральная методика без чужеродных веществ; проводится курсом.',
    tk: 'Plazmolifting — hassanyň öz ganyndan alnan, trombositlere baý plazmasynyň goýberilmegi. Tebigy dikeldişi işjeňleşdirýär, deriniň dykyzlygyny we çeýeligini artdyrýar, reňkini gowulandyrýar we deşikleri daraldýar. Keseki maddasyz, doly tebigy usul; kurs bilen geçirilýär.',
  },
  'Микротоковая терапия': {
    ru: 'Аппаратный лифтинг слабыми импульсными токами, которые тонизируют мышцы лица и стимулируют клетки кожи. Подтягивает овал, уменьшает отёчность и улучшает цвет лица — комфортно и без боли. Хороший вариант «процедуры выходного дня» без восстановительного периода.',
    tk: 'Ýüz myşsalaryny tonuslandyrýan we deri öýjüklerini höweslendirýän gowşak impuls toklary bilen apparat liftingi. Ýüzüň owalyny çekýär, çişmegi azaldýar we reňkini gowulandyrýar — rahat we agyrysyz. Dikeldiş döwri bolmadyk «dynç güni prosedurasy» üçin gowy saýlaw.',
  },

  // ─────────────── Дерматология ───────────────
  'Консультация дерматолога': {
    ru: 'Приём врача-дерматолога: осмотр кожи, постановка диагноза и назначение лечения при высыпаниях, зуде, шелушении, изменениях родинок и других проблемах. При необходимости врач направит на анализы и составит план терапии. Обращайтесь при любых тревожных изменениях кожи, волос или ногтей.',
    tk: 'Lukman-dermatologyň kabul edişi: deriniň barlanmagy, diagnoz goýmak we örgün, gijemek, soýulma, mеňleriň üýtgemegi we beýleki meselelerde bejergi bellemek. Zerur bolsa lukman seljermä ugradar we bejergi meýilnamasyny düzer. Deriniň, saçyň ýa-da dyrnaklaryň islendik aladalandyryjy üýtgemesinde ýüz tutuň.',
  },
  'Лечение акне': {
    ru: 'Комплексный протокол борьбы с угревой болезнью: врач определяет причину высыпаний и составляет план — чистки, лечебные процедуры, наружные и при необходимости системные назначения. Цель — не замаскировать, а устранить причину и предотвратить рубцы и постакне. Лечение проводится курсом под контролем дерматолога.',
    tk: 'Akne keseline garşy toplumlaýyn protokol: lukman örgünleriň sebäbini kesgitleýär we meýilnama düzýär — arassalamalar, bejeriş proseduralary, daşky we zerur bolsa içki bellenilmeler. Maksat — örtmek däl-de, sebäbini aýyrmak we tyrnaklary hem postakneni öňünden almak. Bejergi dermatologyň gözegçiliginde kurs bilen geçirilýär.',
  },
  'Лечение розацеа': {
    ru: 'Диагностика и терапия розацеа — хронического покраснения кожи лица с сосудистыми звёздочками и воспалениями. Врач подбирает щадящий уход и лечение, которые снижают красноту, успокаивают кожу и продлевают ремиссию. Включает рекомендации по образу жизни и триггерам, провоцирующим обострения.',
    tk: 'Rozasea — ýüz derisiniň damar ýyldyzjyklary we sözlemler bilen dowamly gyzarmagy — diagnozy we bejergisi. Lukman gyzarmany azaldýan, derini köşeşdirýän we remissiýany uzaldýan ýumşak ideg we bejergi saýlaýar. Durmuş ýörelgesi we ýitileşmäni döredýän sebäpler boýunça maslahatlary öz içine alýar.',
  },
  'Дерматоскопия (1 элемент)': {
    ru: 'Осмотр родинки или новообразования под дерматоскопом — прибором с многократным увеличением и подсветкой. Позволяет оценить структуру образования и вовремя выявить опасные изменения. Безболезненно и быстро; рекомендуется при появлении новых или изменении старых родинок.',
    tk: 'Meňiň ýa-da täze döremäniň dermatoskop — köp esse ulaldýan we yşyklandyrýan abzal bilen barlanmagy. Döremäniň gurluşyny bahalandyrmaga we howply üýtgemeleri wagtynda ýüze çykarmaga mümkinçilik berýär. Agyrysyz we çalt; täze meňler peýda bolanda ýa-da köne meňler üýtgände maslahat berilýär.',
  },
  'Удаление новообразований (до 5 мм)': {
    ru: 'Удаление папиллом, бородавок, мелких родинок и других образований до 5 мм лазером или жидким азотом. Метод подбирается врачом после осмотра; процедура быстрая, под местной анестезией, с минимальным следом. При показаниях материал может быть направлен на исследование.',
    tk: 'Papillomalary, siňňilleri, ownuk meňleri we 5 mm çenli beýleki döremeleri lazer ýa-da suwuk azot bilen aýyrmak. Usuly lukman barlandan soň saýlaýar; prosedura çalt, ýerli anesteziýa bilen, az yz galdyrýar. Görkezme bolsa material barlaga ugradylyp bilner.',
  },

  // ─────────────── Трихология ───────────────
  'Консультация трихолога': {
    ru: 'Приём специалиста по волосам и коже головы: трихоскопия (осмотр под увеличением), оценка причин выпадения, перхоти, зуда и истончения волос. По результатам — индивидуальный план лечения и ухода. Подходит при любой проблеме с волосами и кожей головы.',
    tk: 'Saç we kelle derisi boýunça hünärmeniň kabul edişi: trihoskopiýa (ulaldyp barlamak), saç düşmeginiň, kepeklenmäniň, gijemegiň we saçyň inçelmeginiň sebäplerini bahalandyrmak. Netijä görä — şahsy bejergi we ideg meýilnamasy. Saç we kelle derisi bilen islendik meselede laýyk.',
  },
  'PRP-терапия волос': {
    ru: 'Плазмолифтинг кожи головы: введение собственной плазмы пациента, богатой факторами роста, чтобы пробудить «спящие» волосяные луковицы. Уменьшает выпадение, укрепляет корни и стимулирует рост новых волос. Натуральная методика, проводится курсом с накопительным эффектом.',
    tk: 'Kelle derisiniň plazmoliftingi: ösüş faktorlaryna baý öz plazmasyny goýbermek arkaly «uklap ýatan» saç soganjyklaryny oýarmak. Saç düşmegini azaldýar, kökleri berkidýär we täze saçlaryň ösmegini höweslendirýär. Tebigy usul, jemleýji netije bilen kurs görnüşinde geçirilýär.',
  },
  'Мезотерапия волосистой части головы': {
    ru: 'Микроинъекции витаминно-питательного коктейля в кожу головы для укрепления и стимуляции роста волос. Улучшает питание корней, уменьшает выпадение, делает волосы гуще и сильнее. Курс процедур подбирается трихологом индивидуально.',
    tk: 'Saçy berkitmek we ösüşini höweslendirmek üçin witamin-iýmit kokteýliniň kelle derisine mikrosanjymy. Kökleriň iýmitlenmegini gowulandyrýar, saç düşmegini azaldýar, saçy goýy we güýçli edýär. Proseduralaryň kursy triholog tarapyndan şahsy saýlanýar.',
  },
  'Озонотерапия волос': {
    ru: 'Насыщение кожи головы активным кислородом (озоном), который улучшает кровообращение и питание волосяных луковиц. Помогает при выпадении, перхоти и жирности кожи головы, снимает воспаление и зуд. Хорошо сочетается с другими процедурами для волос в составе курса.',
    tk: 'Kelle derisini gan aýlanyşygyny we saç soganjyklarynyň iýmitlenmegini gowulandyrýan işjeň kislorod (ozon) bilen baýlaşdyrmak. Saç düşmeginde, kepeklenmede we kelle derisiniň ýaglylygynda kömek edýär, sözlemi we gijemegi aýyrýar. Kursuň düzüminde saç üçin beýleki proseduralar bilen gowy utgaşýar.',
  },

  // ─────────────── Эстетическая медицина ───────────────
  'Консультация эстетиста': {
    ru: 'Беседа и осмотр, на которых специалист оценивает черты лица, пропорции и состояние кожи и вместе с вами планирует деликатную коррекцию. Вы получите понятную «дорожную карту»: какие процедуры, в каком порядке и с каким ожидаемым результатом. Без давления — только честные рекомендации под ваши цели и бюджет.',
    tk: 'Söhbetdeşlik we barlag: hünärmen ýüzüň keşbini, gabarasyny we deriniň ýagdaýyny bahalandyrýar hem siziň bilen bilelikde näzik düzedişi meýilleşdirýär. Size düşnükli «ýol kartasy» berilýär: haýsy proseduralar, haýsy tertipde we haýsy garaşylýan netije bilen. Basyşsyz — diňe siziň maksatlaryňyza we býujetiňize görä dogruçyl maslahatlar.',
  },
  'Нитевой лифтинг (нос-губы)': {
    ru: 'Подтяжка нижней трети лица рассасывающимися нитями PDO: их вводят под кожу, создавая поддерживающий каркас, который подтягивает ткани и разглаживает носогубные складки. Эффект виден сразу и усиливается со временем за счёт выработки коллагена, держится до 1,5–2 лет. Процедура малотравматична, под местной анестезией.',
    tk: 'Ýüzüň aşaky üçden bir böleginiň ereýän PDO sapaklary bilen çekilmegi: olar deriniň aşagyna goýberilip, dokumalary çekýän we burun-dodak büküm­lerini tekizleýän goldaw karkasyny döredýär. Netije derrew görünýär we kollagen öndürilmegi sebäpli wagtyň geçmegi bilen güýçlenýär, 1,5–2 ýyla çenli saklanýar. Prosedura az şikesli, ýerli anesteziýa bilen.',
  },
  'Объёмное моделирование лица': {
    ru: 'Комплексная коррекция контуров и объёмов лица филлерами — скулы, подбородок, овал, носогубная зона — по индивидуальному плану. Восстанавливает гармоничные пропорции, естественно и без операции. Результат заметен сразу и сохраняется до года; объём препарата рассчитывается под ваши задачи.',
    tk: 'Ýüzüň konturlaryny we göwrümlerini fillerler bilen toplumlaýyn düzetmek — ýaňaklar, eňek, owal, burun-dodak zolagy — şahsy meýilnama boýunça. Sazlaşykly gabaralary operasiýasyz, tebigy görnüşde dikeldýär. Netije derrew görünýär we bir ýyla çenli saklanýar; serişdäniň möçberi siziň talaplaryňyza görä hasaplanýar.',
  },
}
