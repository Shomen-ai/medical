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
}

export const medicalDict: Record<'ru' | 'tk', Record<string, string>> = {
  ru: {},     // identity map
  tk: medicalTk,
}
