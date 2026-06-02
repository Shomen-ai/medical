// Файл: locales/medical.ts
// Назначение: медицинский глоссарий — отображает русские названия специальностей и услуг из БД в туркменские эквиваленты для UI; используется через tMed().

// Medical glossary: specialty + service names that come from the DB in Russian
// and need a Turkmen counterpart for the UI. Used via tMed() in components.
//
// The key is always the canonical Russian name from migrations 002+007.
// If a name isn't in the map, tMed() returns it unchanged (graceful fallback).

const medicalTk: Record<string, string> = {
  // Specialties
  'Косметология': 'Kosmetologiýa',
  'Дерматология': 'Dermatologiýa',
  'Трихология': 'Trihologiýa',
  'Эстетическая медицина': 'Estetiki lukmançylyk',

  // Services — Косметология
  'Консультация косметолога': 'Kosmetolog maslahaty',
  'Чистка лица механическая': 'Ýüzüň mehaniki arassalanmagy',
  'Чистка лица ультразвуковая': 'Ýüzüň ultrases bilen arassalanmagy',
  'Химический пилинг': 'Himiki piling',
  'Инъекции ботулотоксина (зона)': 'Botulotoksin sanjymy (bir ýer)',
  'Контурная пластика (1 мл)': 'Konturly plastika (1 ml)',
  'Биоревитализация': 'Biorewitalizasiýa',
  'Мезотерапия лица': 'Ýüz mezoterapiýasy',
  'PRP-терапия лица': 'Ýüzüň PRP-terapiýasy',
  'Микротоковая терапия': 'Mikrotok terapiýasy',

  // Services — Дерматология
  'Консультация дерматолога': 'Dermatolog maslahaty',
  'Лечение акне': 'Akneniň bejergisi',
  'Дерматоскопия (1 элемент)': 'Dermatoskopiýa (1 element)',
  'Удаление новообразований (до 5 мм)': 'Täze döremeleri aýyrmak (5 mm çenli)',
  'Лечение розацеа': 'Rozaseanyň bejergisi',

  // Services — Трихология
  'Консультация трихолога': 'Triholog maslahaty',
  'Мезотерапия волосистой части головы': 'Saç böleginiň mezoterapiýasy',
  'PRP-терапия волос': 'Saçlaryň PRP-terapiýasy',
  'Озонотерапия волос': 'Saçlaryň ozonoterapiýasy',

  // Services — Эстетическая медицина
  'Консультация эстетиста': 'Estetik maslahaty',
  'Нитевой лифтинг (нос-губы)': 'Sapakly lifting (burun-dodak)',
  'Объёмное моделирование лица': 'Ýüzüň göwrümli modelirlemesi',

  // Описания услуг (миграции 002+007) — показываются на странице услуги; ключи строго совпадают с RU из БД
  'Плазмолифтинг кожи головы. Активирует спящие фолликулы.':
    'Kelle derisiniň plazmoliftingi. Uklap ýatan follikullary işjeňleşdirýär.',
  'Плазмолифтинг — введение собственной плазмы. Регенерация и омоложение.':
    'Plazmolifting — öz plazmaňy goýbermek. Dikeldiş we ýaşartma.',
  'Инъекционное увлажнение гиалуроновой кислотой. Улучшает тургор и сияние кожи.':
    'Gialuron turşusy bilen iňýeksion nemlendiriş. Deriniň dartgynlygyny we ýalkymyny gowulandyrýar.',
  'Осмотр родинки или новообразования под дерматоскопом.':
    'Meňiň ýa-da täze döremäniň dermatoskop arkaly barlagy.',
  'Коррекция мимических морщин ботулотоксином. Цена за одну зону.':
    'Botulotoksin bilen mimiki ýygyrtlary düzetmek. Baha bir zolak üçin.',
  'Осмотр, постановка диагноза, назначение лечения.':
    'Barlag, diagnoz goýmak, bejergi bellemek.',
  'Первичный осмотр, анализ состояния кожи, подбор программы ухода.':
    'Ilkinji barlag, deriniň ýagdaýyny seljermek, ideg maksatnamasyny saýlamak.',
  'Трихоскопия, анализ состояния волос и кожи головы.':
    'Trihoskopiýa, saçyň we kelle derisiniň ýagdaýyny seljermek.',
  'Оценка внешности, планирование коррекции.':
    'Daş keşbi bahalandyrmak, düzediş meýilnamasy.',
  'Введение филлера для восстановления объёмов и контуров лица.':
    'Ýüzüň göwrümini we konturlaryny dikeltmek üçin filler goýbermek.',
  'Комплексный протокол лечения угревой болезни. Включает чистку и назначения.':
    'Akne keseliniň toplumlaýyn bejergi teswiri. Arassalaýyşy we bellemeleri öz içine alýar.',
  'Консультация и назначение терапии при розацеа.':
    'Rozasea boýunça maslahat we bejergi bellemek.',
  'Инъекции витаминного коктейля для стимуляции роста волос.':
    'Saçyň ösüşini höweslendirmek üçin witamin kokteýliniň iňýeksiýasy.',
  'Микроинъекции коктейля витаминов и пептидов. Питание и восстановление кожи.':
    'Witaminleriň we peptidleriň kokteýliniň mikroiňýeksiýasy. Deriniň iýmitlenişi we dikeldişi.',
  'Аппаратный лифтинг низкочастотными токами. Улучшает овал лица.':
    'Pes ýygylykly toklar bilen apparatly lifting. Ýüzüň owalyny gowulandyrýar.',
  'Подтяжка нижней трети лица нитями PDO.':
    'Ýüzüň aşaky üçden bir böleginiň PDO sapaklary bilen dartylmagy.',
  'Комплексная коррекция контуров с применением филлеров.':
    'Filler ulanyp konturlary toplumlaýyn düzetmek.',
  'Насыщение кожи головы кислородом. Улучшает питание корней.':
    'Kelle derisini kislorod bilen baýlaşdyrmak. Saç köklüriniň iýmitlenişini gowulandyrýar.',
  'Жидким азотом или лазером. До 5 мм в диаметре.':
    'Suwuk azot ýa-da lazer bilen. Diametri 5 mm çenli.',
  'Поверхностное или срединное отшелушивание. Устраняет пигментацию, мелкие морщины.':
    'Üstki ýa-da orta gatlakly pilling. Pigmentasiýany, ownuk ýygyrtlary aýyrýar.',
  'Глубокое очищение пор, удаление комедонов и акне. Включает распаривание и маску.':
    'Öýjükleriň çuňňur arassalanmagy, komedonlary we akneni aýyrmak. Buglatmagy we maskany öz içine alýar.',
  'Бесконтактное очищение кожи ультразвуком. Подходит для чувствительной кожи.':
    'Ultrases bilen deriniň kontaktsyz arassalanmagy. Duýgur deri üçin amatly.',

  // Образование врачей (миграция 020) — ключи строго совпадают с засеянными RU-строками
  'Государственный медицинский университет Туркменистана им. М. Гаррыева, 2009; ординатура по дерматовенерологии':
    'Myrat Garryýew adyndaky Türkmenistanyň döwlet lukmançylyk uniwersiteti, 2009; dermatowenerologiýa boýunça ordinatura',
  'Первый МГМУ им. И. М. Сеченова (Москва), 2012; курс эстетической медицины':
    'I. M. Seçenow adyndaky Birinji Moskwa döwlet lukmançylyk uniwersiteti, 2012; estetiki lukmançylyk kursy',
  'Государственный медицинский университет Туркменистана им. М. Гаррыева, 2014':
    'Myrat Garryýew adyndaky Türkmenistanyň döwlet lukmançylyk uniwersiteti, 2014',
  'Государственный медицинский университет Туркменистана им. М. Гаррыева, 2007; ординатура по дерматовенерологии':
    'Myrat Garryýew adyndaky Türkmenistanyň döwlet lukmançylyk uniwersiteti, 2007; dermatowenerologiýa boýunça ordinatura',
  'Российский университет дружбы народов (Москва), 2011':
    'Halklaryň dostlugy rus uniwersiteti (Moskwa), 2011',
  'Государственный медицинский университет Туркменистана им. М. Гаррыева, 2013':
    'Myrat Garryýew adyndaky Türkmenistanyň döwlet lukmançylyk uniwersiteti, 2013',
  'Государственный медицинский университет Туркменистана им. М. Гаррыева, 2008':
    'Myrat Garryýew adyndaky Türkmenistanyň döwlet lukmançylyk uniwersiteti, 2008',
  'Hacettepe Üniversitesi (Анкара, Турция), 2010; пластическая и эстетическая медицина':
    'Hacettepe uniwersiteti (Ankara, Türkiýe), 2010; plastiki we estetiki lukmançylyk',
  'Государственный медицинский университет Туркменистана им. М. Гаррыева, 2016':
    'Myrat Garryýew adyndaky Türkmenistanyň döwlet lukmançylyk uniwersiteti, 2016',
  'Государственный медицинский университет Туркменистана им. М. Гаррыева, 2019':
    'Myrat Garryýew adyndaky Türkmenistanyň döwlet lukmançylyk uniwersiteti, 2019',
  'Белорусский государственный медицинский университет (Минск), 2014':
    'Belarus döwlet lukmançylyk uniwersiteti (Minsk), 2014',
  'Государственный медицинский университет Туркменистана им. М. Гаррыева, 2012':
    'Myrat Garryýew adyndaky Türkmenistanyň döwlet lukmançylyk uniwersiteti, 2012',
  'Государственный медицинский университет Туркменистана им. М. Гаррыева, 2017':
    'Myrat Garryýew adyndaky Türkmenistanyň döwlet lukmançylyk uniwersiteti, 2017',
  'Государственный медицинский университет Туркменистана им. М. Гаррыева, 2015; эстетическая медицина':
    'Myrat Garryýew adyndaky Türkmenistanyň döwlet lukmançylyk uniwersiteti, 2015; estetiki lukmançylyk',
  'İstanbul Üniversitesi (Стамбул, Турция), 2010; нитевой лифтинг и эстетическая хирургия':
    'Stambul uniwersiteti (Stambul, Türkiýe), 2010; sapakly lifting we estetiki hirurgiýa',
}

// Словарь медицинских терминов по локалям: ru — тождество (исходные названия), tk — переводы.
export const medicalDict: Record<'ru' | 'tk', Record<string, string>> = {
  ru: {},     // identity map
  tk: medicalTk,
}
