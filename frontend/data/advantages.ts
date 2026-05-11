export interface LocalizedString {
  ru: string
  tk: string
}

export interface LocalizedParagraphs {
  ru: string[]
  tk: string[]
}

export interface Advantage {
  slug: string
  icon: string
  title: LocalizedString
  text: LocalizedString
  body: LocalizedParagraphs
}

export const ADVANTAGES: Advantage[] = [
  {
    slug: 'experienced-doctors',
    icon: '👩‍⚕️',
    title: { ru: 'Опытные врачи', tk: 'Tejribeli lukmanlar' },
    text: {
      ru: '8 специалистов с опытом от 5 до 15 лет',
      tk: 'Tejribesi 5–15 ýyl bolan 8 hünärmen',
    },
    body: {
      ru: [
        'В нашей клинике работают восемь врачей с медицинским образованием и сертификацией по основным направлениям эстетической медицины: косметологии, дерматологии, трихологии и эстетической хирургии.',
        'Стаж работы наших специалистов — от 5 до 15 лет. Каждый врач ежегодно подтверждает квалификацию и проходит обучение на международных конгрессах и мастер-классах ведущих производителей косметологического оборудования.',
        'Перед записью на процедуру каждый пациент получает консультацию у профильного специалиста. Мы не назначаем услуг «впрок» и подбираем протокол исходя из реального состояния кожи, волос и противопоказаний.',
      ],
      tk: [
        'Klinikamyzda lukmançylyk bilimi we estetiki lukmançylygyň esasy ugurlary boýunça şahadatnamasy bolan sekiz lukman işleýär: kosmetologiýa, dermatologiýa, trihologiýa we estetiki hirurgiýa.',
        'Hünärmenlerimiziň iş tejribesi — 5–15 ýyl. Her lukman ýyl saýyn hünärini tassyklap, halkara kongreslerinde we kosmetologiýa enjamlaryny öndüriji öňdebaryjy kärhanalaryň ussatlyk sapaklarynda okuw geçýär.',
        'Prosedura ýazylmazdan öň her hassa degişli hünärmen bilen maslahatlaşýar. Biz hyzmatlary «artykmaç» bellemeýäris we protokoly derisi, saçy bilen baglanyşykly hakyky ýagdaýdan we garşy görkezmelerden ugur alyp saýlaýarys.',
      ],
    },
  },
  {
    slug: 'modern-equipment',
    icon: '🏥',
    title: { ru: 'Современное оборудование', tk: 'Häzirki zaman enjamlar' },
    text: {
      ru: 'Аппаратура последнего поколения',
      tk: 'Iň soňky nesil enjamlary',
    },
    body: {
      ru: [
        'Клиника оснащена аппаратурой ведущих европейских и южнокорейских производителей: лазерами для удаления новообразований и эпиляции, ультразвуковыми системами для чисток и SMAS-лифтинга, аппаратами микротоковой и RF-терапии.',
        'Все приборы проходят регулярное техническое обслуживание и имеют сертификаты соответствия. Расходные материалы (иглы, насадки, фильтры) — только одноразовые, поставляются от официальных дистрибьюторов.',
        'Стерилизационная зона оборудована автоклавом класса B, что соответствует санитарным требованиям к медицинским учреждениям.',
      ],
      tk: [
        'Klinika öňdebaryjy ýewropa we günorta-koreý öndürijileriniň enjamlary bilen üpjün edilen: täze döremeleri aýyrmak we epilýasiýa üçin lazerler, arassalama we SMAS-lifting üçin ultrases ulgamlary, mikrotok we RF-terapiýa enjamlary.',
        'Ähli enjamlar yzygiderli tehniki hyzmaty geçýär we degişlilik şahadatnamalary bar. Sarp ediljek serişdeler (iňňeler, ujlar, süzgüçler) diňe bir gezeklikdir we resmi paýlaýjylardan iberilýär.',
        'Sterilizasiýa zolagy B klassly awtoklaw bilen üpjün edilen — bu lukmançylyk edaralary üçin sanitariýa talaplaryna laýyk gelýär.',
      ],
    },
  },
  {
    slug: 'online-booking',
    icon: '📱',
    title: { ru: 'Онлайн-запись', tk: 'Onlaýn bellige durmak' },
    text: {
      ru: 'Запишитесь к врачу за 2 минуты, без звонков',
      tk: 'Lukmana 2 minutda jaňsyz bellige duruň',
    },
    body: {
      ru: [
        'Система онлайн-записи позволяет выбрать специалиста, удобную дату и время за пару минут — без звонков и долгих ожиданий на линии.',
        'После выбора услуги вы увидите только свободные слоты выбранного врача. Подтверждение происходит по SMS-коду, а напоминание о визите придёт за день до приёма.',
        'В личном кабинете вы можете перенести или отменить запись (не позднее чем за 2 часа), посмотреть историю приёмов и скачать справку об оплате медицинских услуг.',
      ],
      tk: [
        'Onlaýn bellige durmak ulgamy hünärmeni, amatly senäni we wagty birnäçe minutda saýlamaga mümkinçilik berýär — jaňsyz we uzak garaşmasyz.',
        'Hyzmat saýlandan soň diňe saýlanan lukmanyň boş wagtlaryny görersiňiz. Tassyklama SMS-kod arkaly geçýär, ýatlatma — kabul edilmezden bir gün öň iberilýär.',
        'Şahsy hasabyňyzda bellige durmagy geçirip ýa-da goýbolsun edip bilersiňiz (kabul edilmezden 2 sagatdan az däl), kabul taryhyny görüp we lukmançylyk hyzmatlary üçin töleg şahadatnamasyny ýükläp bilersiňiz.',
      ],
    },
  },
  {
    slug: 'individual-approach',
    icon: '💬',
    title: { ru: 'Индивидуальный подход', tk: 'Şahsy çemeleşme' },
    text: {
      ru: 'Программа лечения под каждого пациента',
      tk: 'Her hassa üçin aýratyn bejergi meýilnamasy',
    },
    body: {
      ru: [
        'Эстетическая медицина — это не «универсальные процедуры по прайсу», а индивидуальная стратегия для конкретного человека. Возраст, тип кожи, образ жизни, противопоказания, бюджет — мы учитываем всё.',
        'На первой консультации врач задаёт вопросы об общем состоянии здоровья, перенесённых заболеваниях, аллергиях и принимаемых препаратах. По результатам осмотра составляется план процедур с этапами и ориентировочной стоимостью.',
        'Программа корректируется по ходу лечения: если результат достигается быстрее ожидаемого — мы убираем лишние процедуры, не настаиваем на «полном курсе» ради выручки.',
      ],
      tk: [
        'Estetiki lukmançylyk — bu «prais boýunça umumy prosedyralar» däl-de, eýsem belli bir adam üçin aýratyn strategiýadyr. Ýaş, deriniň görnüşi, ýaşaýyş durmuşy, garşy görkezmeler, býujet — biz hemmesini hasaba alýarys.',
        'Ilkinji maslahatda lukman umumy saglyk ýagdaýy, geçirilen kesellikler, allergiýalar we kabul edilýän derman serişdeleri barada soraýar. Barlagyň netijesinde tapgyrlary we takmyny bahasy bilen prosedyralaryň meýilnamasy düzülýär.',
        'Meýilnama bejerginiň dowamynda düzedilýär: eger netije garaşylanyňdan çalt gazanylsa — artykmaç prosedyralary aýyrýarys, girdeji üçin «doly kurs» geçmegi tabşyrmaýarys.',
      ],
    },
  },
]
