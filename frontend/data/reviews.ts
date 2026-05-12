// Файл: data/reviews.ts
// Назначение: статичный массив отзывов пациентов для секции «Отзывы» на главной — текст переведён на ru/tk.
import type { Review } from '~/types'

// Список отзывов пациентов клиники, отображаемых в карусели.
export const REVIEWS: Review[] = [
  {
    author: 'Айлар М.',
    text: {
      ru: 'Отличная клиника! Посещаю косметолога уже год. Кожа стала намного лучше, результат виден после первой процедуры.',
      tk: 'Ajaýyp klinika! Bir ýyl bäri kosmetologa gatnaýaryn. Derim has gowulandy, ilkinji prosedurdan soň netije görünýär.',
    },
    rating: 5,
    date: '2026-03-15',
  },
  {
    author: 'Мырат А.',
    text: {
      ru: 'Обратился к дерматологу с давней проблемой. Врач подобрал лечение — за месяц всё прошло. Спасибо!',
      tk: 'Köpden bäri bar bolan mesele bilen dermatologa ýüz tutdym. Lukman bejergi saýlady — bir aýda hemme zat geçdi. Sag boluň!',
    },
    rating: 5,
    date: '2026-02-28',
  },
  {
    author: 'Огулджемал Х.',
    text: {
      ru: 'Прохожу курс трихологии. Видны результаты уже после 3 сеансов. Вежливый персонал, уютная обстановка.',
      tk: 'Trihologiýa kursuny geçýärin. 3 sapakdan soň netijeler eýýäm görünýär. Mylakatly işgärler, oňaýly atmosfera.',
    },
    rating: 5,
    date: '2026-04-10',
  },
]
