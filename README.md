# BeautyMed

Веб-приложение клиники красоты и здоровья (Туркменабат): онлайн-запись пациентов
через SMS-код, личный кабинет, рабочее место врача, админ-панель.

**Стек:** Go (Gin) + Nuxt 3 (Vue 3) + PostgreSQL + Redis + nginx, всё в Docker.

Полное человекочитаемое описание архитектуры — в файле `EXPLANATION.txt`
(если он у вас есть локально; в git его нет — спросите коллегу).

---

## 1. Что установить (Windows)

Поставьте по порядку. Все программы бесплатные.

### 1.1. WSL2 (Linux внутри Windows) — обязательно

WSL даёт вам полноценный Ubuntu прямо в винде. Все наши скрипты написаны для
Linux, и WSL — самый простой способ их запускать.

1. Откройте **PowerShell от имени администратора** (правая кнопка по «Пуск» →
   «Терминал (Администратор)»).
2. Выполните одну команду:
   ```powershell
   wsl --install
   ```
3. Перезагрузите компьютер.
4. После перезагрузки автоматически откроется окно Ubuntu — задайте имя
   пользователя и пароль (запишите их).
5. Если окно не появилось, нажмите «Пуск» → «Ubuntu».

После этого у вас в системе есть терминал Ubuntu. Все команды ниже выполняйте
**в нём**, не в PowerShell и не в cmd.

### 1.2. Docker Desktop — обязательно

Это программа, которая позволяет запускать «контейнеры» — изолированные
приложения. У нас в контейнерах живут Postgres, Redis, бэкенд и фронтенд.

1. Скачайте: <https://www.docker.com/products/docker-desktop/>
2. Установите, при установке оставьте галочку «Use WSL 2 instead of Hyper-V».
3. Запустите Docker Desktop. В трее (правый нижний угол) появится синяя
   иконка-кит. Дождитесь, пока кит перестанет «шевелиться» — это значит,
   Docker готов.
4. Откройте Docker Desktop → Settings → Resources → WSL Integration. Включите
   галочку для вашего Ubuntu. Нажмите «Apply & restart».

Проверка: в терминале Ubuntu выполните
```bash
docker --version
docker compose version
```
Обе команды должны напечатать номер версии. Если выдают «command not found» —
интеграция WSL не включилась, повторите шаг 4.

### 1.3. Git — обычно уже есть

Проверьте в Ubuntu:
```bash
git --version
```
Если выдаёт «command not found»:
```bash
sudo apt update && sudo apt install -y git
```

### 1.4. VS Code + расширение «WSL» — рекомендуется

1. Скачайте VS Code: <https://code.visualstudio.com/>
2. Откройте VS Code → слева в боковой панели иконка «Extensions» (Ctrl+Shift+X) →
   найдите и установите расширение **«WSL»** от Microsoft.
3. Теперь в Ubuntu-терминале вы можете писать `code .` и проект откроется в
   VS Code, но «думать» он будет внутри Linux — никаких проблем с путями
   и переводами строк.

### 1.5. Утилиты для деплоя

В Ubuntu один раз выполните:
```bash
sudo apt update && sudo apt install -y rsync sshpass curl
```

---

## 2. Клонирование проекта

В терминале Ubuntu:

```bash
cd ~
git clone https://github.com/Shomen-ai/medical.git BeautyMed
cd BeautyMed
```

Откройте проект в VS Code:
```bash
code .
```

> **Важно:** клонируйте именно внутрь `~` (домашний каталог Ubuntu), **не**
> в `/mnt/c/Users/...`. Внутри `~` (`/home/<имя>/...`) файлы хранятся в Linux-FS,
> где Git настроен на LF и всё работает корректно. На `/mnt/c/...` Windows
> агрессивно меняет права доступа и переводы строк — будут случайные баги.

---

## 3. Локальный запуск

### 3.1. Подготовка `.env`

В корне проекта создайте файл `.env` (он в `.gitignore`, в репозиторий не
попадёт). За образец — `.env.example`. Минимально достаточный набор:

```env
# Database
DATABASE_URL=postgres://beautymed:beautymed@postgres:5432/beautymed?sslmode=disable

# Redis
REDIS_URL=redis://redis:6379

# JWT — обязательно поменяйте на любую длинную случайную строку
JWT_SECRET=local-dev-secret-change-me-to-32-chars-please

# SMS — оставьте пустым для локалки. SMS реально не отправятся,
# но OTP-код для входа можно подсмотреть в логах бэкенда.
SMSC_LOGIN=
SMSC_PASSWORD=

# App
PORT=8080
ENV=development
```

Создать файл одной командой:
```bash
cp .env.example .env
nano .env   # отредактируйте JWT_SECRET, сохраните Ctrl+O, Enter, Ctrl+X
```

### 3.2. Запуск всех контейнеров

```bash
docker compose up -d --build
```

Первый раз будет долго (5-15 минут) — скачаются образы Postgres/Redis/Node/Go,
соберутся бэкенд и фронтенд. Дальнейшие запуски — секунды.

Проверка, что всё поднялось:
```bash
docker compose ps
```
Все 5 контейнеров (postgres, redis, api, frontend, nginx) должны быть в
статусе `running`.

### 3.3. Открыть в браузере

- Главная: <http://localhost/>
- Вход для сотрудников: <http://localhost/staff-login>
  - Админ:  `admin` / `admin123` → `/admin`
  - Врачи:  `doctor1` … `doctor15` / `doctor123` → `/doctor`
- API health: <http://localhost/health> → должно вернуть `{"status":"ok"}`

### 3.4. Как войти как пациент (SMS-код)

В локалке SMS реально не отправляется. Чтобы войти:
1. Нажмите «Записаться», заполните шаги, дойдите до ввода телефона.
2. Введите любой номер `+99365XXXXXXX`, нажмите «Получить код».
3. В терминале Ubuntu посмотрите логи бэкенда:
   ```bash
   docker compose logs api | grep OTP
   ```
   Там будет строка вроде `OTP for +99365123456 is 4821`.
4. Введите этот код в окне.

### 3.5. Просмотр логов и остановка

```bash
docker compose logs -f          # все логи, потоком (Ctrl+C — выйти)
docker compose logs -f api      # только бэкенд
docker compose down             # остановить всё (данные Postgres сохранятся)
docker compose down -v          # остановить + СТЕРЕТЬ данные БД
```

---

## 4. Структура проекта (где что лежит)

```
BeautyMed/
├── backend/                  ← Go-сервис (API)
│   ├── cmd/api/main.go       ← точка входа
│   └── internal/
│       ├── handler/          ← HTTP-обработчики (Gin)
│       ├── service/          ← бизнес-логика
│       ├── repository/       ← SQL-запросы (sqlx)
│       ├── model/            ← структуры данных
│       ├── middleware/       ← JWT-аутентификация
│       ├── router/           ← регистрация роутов
│       └── db/migrations/    ← SQL-миграции (применяются автоматом при старте)
├── frontend/                 ← Nuxt 3 (Vue 3 + TypeScript)
│   ├── pages/                ← страницы (одна страница = один файл)
│   ├── components/           ← переиспользуемые компоненты
│   ├── stores/               ← Pinia (общая память Vue)
│   ├── composables/          ← вспомогательные функции
│   ├── locales/              ← переводы ru/tk
│   ├── public/doctors/       ← фото врачей
│   └── nuxt.config.ts        ← конфигурация Nuxt
├── nginx/nginx.conf          ← маршрутизация HTTP
├── docker-compose.yml        ← описание всех контейнеров
├── deploy.sh                 ← скрипт деплоя на сервер
├── .env                      ← локальные секреты (НЕ в git)
└── .deploy.env               ← пароль к серверу (НЕ в git)
```

В каждом файле наверху есть русский комментарий с описанием назначения.

---

## 5. Внесение изменений и push в GitHub

### 5.1. Один раз настроить git

```bash
git config --global user.name "Ваше Имя"
git config --global user.email "you@example.com"
```

### 5.2. Цикл работы

```bash
git status                              # что изменилось
git add <файлы>                         # подготовить файлы к коммиту
git commit -m "fix: что-то поправил"    # сделать коммит
git pull --rebase origin main           # подтянуть свежие изменения (если кто-то ещё пушил)
git push origin main                    # отправить на GitHub
```

**Правила сообщений** (по соглашению Conventional Commits, используется в проекте):
- `feat: …` — новая фича
- `fix: …` — исправление бага
- `refactor: …` — переписать без смены поведения
- `chore: …` — рутина (зависимости, форматирование, доки)
- `docs: …` — только документация

После `git push` зайдите на <https://github.com/Shomen-ai/medical/commits/main>
и убедитесь, что ваш коммит сверху.

---

## 6. Деплой на продакшен-сервер

### 6.1. Что такое деплой

Это процесс «отправить ваши изменения с локального компьютера на боевой
сервер». На сервере уже крутятся 5 контейнеров (Postgres/Redis/api/frontend/nginx).
Скрипт `deploy.sh`:
1. rsync'ит свежий код на сервер,
2. пересобирает контейнеры backend/frontend (`docker compose up --build`),
3. проверяет, что всё поднялось (5 контейнеров running + API отвечает 200).

### 6.2. Один раз: создать `.deploy.env`

В корне проекта создайте файл `.deploy.env` (он в `.gitignore`):

```env
SERVER_PASS=<пароль_который_вам_передал_коллега_отдельно>
```

⚠️ **Не присылайте этот пароль через GitHub, Discord, Slack-каналы.** Только
через прямые личные сообщения (Telegram, Signal, голос).

### 6.3. Запуск

```bash
source .deploy.env && ./deploy.sh
```

Что должно произойти:
- Скрипт выведет статус: проверка зависимостей → rsync → docker compose build
  → проверка контейнеров → проверка `/health` и `/api/specialties`.
- В конце: `✓ Деплой успешен. Все контейнеры работают, API отвечает.`

Длительность первого деплоя — 1-3 минуты. Последующие (с кешем образов) —
~30-60 секунд.

### 6.4. Полезные флаги

```bash
DRY_RUN=1     ./deploy.sh   # показать что будет передано, без действий
SKIP_BUILD=1  ./deploy.sh   # только синхронизация, без пересборки
SKIP_VERIFY=1 ./deploy.sh   # без healthcheck
```

### 6.5. Что делать, если деплой упал

В конце скрипт сам подскажет команду для просмотра логов:
```bash
sshpass -p "$SERVER_PASS" ssh root@85.198.80.245 \
  'cd /opt/beautymed && docker compose logs --tail=80'
```

Чаще всего проблемы — это:
- забыли создать `.deploy.env` или пароль неверный → ошибка SSH;
- ошибка компиляции Go или сборки фронта → виден в логе `docker compose build`;
- миграция упала с SQL-ошибкой → виден в логе `docker compose logs api`.

---

## 7. Частые ошибки и решения

| Ошибка | Причина | Решение |
|---|---|---|
| `docker compose: command not found` | WSL-интеграция Docker Desktop выключена | Docker Desktop → Settings → Resources → WSL Integration → включить Ubuntu |
| `/usr/bin/env: 'bash\r': No such file` | Windows подсунул CRLF в `.sh` | `dos2unix deploy.sh` или клонируйте заново внутрь `~` (не `/mnt/c/...`) |
| `permission denied` при `./deploy.sh` | Нет executable-бита | `chmod +x deploy.sh` или запустить как `bash deploy.sh` |
| API недоступен после `docker compose up` | Не отработали миграции / опечатка в .env | `docker compose logs api` — посмотреть первые 50 строк |
| Frontend показывает «502 Bad Gateway» | api ещё не успел стартовать | подождите 10 секунд, обновите страницу |
| `port 80 already in use` | На винде или WSL уже что-то слушает 80 | остановите Skype/IIS/другой Apache, или поменяйте в `docker-compose.yml` порт nginx с `80:80` на `8080:80` |

---

## 8. Полезные ссылки

- Репозиторий: <https://github.com/Shomen-ai/medical>
- Боевой сайт: <http://85.198.80.245/>
- Документация Nuxt: <https://nuxt.com/docs>
- Документация Gin: <https://gin-gonic.com/docs/>
- Документация Docker Compose: <https://docs.docker.com/compose/>
- Conventional Commits: <https://www.conventionalcommits.org/ru/>

---

## 9. Контакты

Если что-то непонятно — пишите тому, кто передал вам этот проект. Не стесняйтесь
задавать «глупые» вопросы; настройка окружения с нуля на чужой машине почти
всегда чем-то да удивляет.
