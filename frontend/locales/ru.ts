// Файл: locales/ru.ts
// Назначение: русские переводы UI — словарь ключей и строк, используемый функцией t() из useI18n. Является эталонным типом для остальных локалей.

// Экспорт словаря русских переводов (источник истины для типизации tk).
export default {
  // Common
  bookOnline: 'Записаться онлайн',
  bookShort: 'Записаться',
  signIn: 'Войти',
  cabinet: 'Кабинет',
  back: '← Назад',
  next: 'Далее →',
  close: 'Закрыть',
  more: 'Подробнее →',
  loading: 'Загружаем...',

  // Header / contacts shared values
  cityName: 'Туркменабад',
  workingHours: 'Часы работы',
  workingHoursValue: 'Пн–Сб: 9:00–20:00, Вс: выходной',

  // Login modal
  loginTitle: 'Вход в личный кабинет',
  loginPhoneLabel: 'Номер телефона',
  loginGetCode: 'Получить код',
  loginSending: 'Отправляем...',
  loginCodeLabel: 'Код из SMS',
  loginCodeSentTo: 'Код отправлен на {phone}',
  loginSubmit: 'Войти',
  loginChecking: 'Проверяем...',
  loginChangeNumber: 'Изменить номер',
  loginStaffLink: 'Вход для персонала клиники',

  // Hero
  heroBadge: 'Клиника красоты · Туркменабад',
  heroTitleLine1: 'Красота и здоровье',
  heroTitleLine2: 'в надёжных руках',
  heroSubtitle: 'Косметология, дерматология, трихология и эстетическая медицина',
  heroOurServices: 'Наши услуги',

  // Advantages
  advantagesTitle: 'Почему выбирают нас',

  // Services
  servicesTitle: 'Наши услуги',
  servicesNone: 'Услуги не найдены',
  serviceDuration: '{n} мин',

  // Doctors
  doctorsTitle: 'Наши врачи',
  doctorsExperience: 'Стаж {n} лет',

  // Reviews
  reviewsTitle: 'Отзывы пациентов',

  // Contacts
  contactsTitle: 'Контакты',
  contactsAddress: 'Адрес',
  contactsPhone: 'Телефон',
  contactsHours: 'Часы работы',

  // Booking modal
  bookingTitle: 'Онлайн-запись',
  bookingStep: 'Шаг {n}',
  step1: 'Специальность',
  step2: 'Врач',
  step3: 'Дата',
  step4: 'Время',
  step5: 'Подтверждение',

  // Booking — confirm
  confirmYourName: 'Ваше имя',
  confirmNamePlaceholder: 'Айна Мырадова',
  confirmPhone: 'Номер телефона',
  confirmGetCode: 'Код',
  confirmResend: 'Повтор',
  confirmOtpLabel: 'Код из SMS',
  confirmSubmit: 'Подтвердить и записаться',
  confirmSubmitting: 'Записываем...',
  confirmCost: 'Стоимость',
  confirmPromoLabel: 'Промокод (необязательно)',
  confirmPromoApply: 'Применить',
  confirmPromoValid: '✓ Промокод применён: скидка {pct}%',
  confirmPromoInvalid: 'Промокод недействителен',
  confirmConsent:
    'Я даю согласие на обработку информации о личной жизни в соответствии с',
  confirmConsentLink: 'Политикой конфиденциальности',
  confirmConsentLaw: '(Закон Туркменистана от 20.03.2017).',
  confirmSuccess: 'Вы записаны!',
  confirmSuccessFooter: 'Ждём вас в клинике BeautyMed',
  confirmToPay: 'К оплате: {price}',
  confirmPhoneErr:
    'Введите телефон в формате +7XXXXXXXXXX (Россия) или +993XXXXXXXX (Туркменистан)',
  confirmCodeErr: 'Неверный или истёкший код. Попробуйте ещё раз.',
  confirmSlotTaken: 'Это время уже занято. Выберите другой слот.',
  confirmGenericErr: 'Ошибка записи. Позвоните нам: ',

  // Footer
  footerPrivacy: 'Политика конфиденциальности',
}
