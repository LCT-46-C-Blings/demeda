# Demo Medical Database API

REST API для управления медицинской клиникой с функциями работы с пациентами, врачами, приемами и медицинским анамнезом.

## 🚀 Особенности

- **Полный CRUD** для пациентов, врачей, приемов и медицинского анамнеза
- **Автоматическая документация** Swagger/OpenAPI
- **SQLite база данных** с автоматическими миграциями
- **CORS поддержка** для веб-приложений
- **Тестовые данные** для быстрого старта
- **Валидация данных** и обработка ошибок

## 🛠 Технологии

- **Go** 1.19+
- **Gin** - веб-фреймворк
- **GORM** - ORM для работы с базой данных
- **SQLite** - база данных
- **Swagger** - автоматическая документация API
- **Make** - управление сборкой

## 📦 Установка и запуск

### Предварительные требования
- Go 1.19 или выше
- Git

### Клонирование репозитория
```bash
git clone <your-repo-url>
cd demeda
```

### Установка зависимостей
```bash
go mod download
```

### Запуск с помощью Make
```bash
# Сборка и запуск
make

# Или по отдельности
make build
make run
```

### Прямой запуск
```bash
go run main.go
```

Сервер будет доступен по адресу: `http://localhost:8080`

## 📚 Документация API

После запуска сервера документация Swagger доступна по адресу:
- **Swagger UI**: http://localhost:8080/swagger/index.html

### Генерация документации Swagger

Если нужно обновить документацию:

```bash
# Установка swag если не установлен
go install github.com/swaggo/swag/cmd/swag@latest

# Генерация документации
swag init
```

## 🗄 Структура API

### Основные сущности

- **Пациенты** - информация о пациентах клиники
- **Врачи** - данные медицинских специалистов  
- **Приемы** - записи о медицинских приемах
- **Медицинские тесты** - результаты анализов и обследований
- **Медицинский анамнез** - история болезней и состояний

### Эндпоинты

#### Пациенты
- `GET /patients` - список всех пациентов
- `GET /patients/:id` - информация о пациенте
- `POST /patients` - создание пациента
- `PUT /patients/:id` - обновление пациента
- `DELETE /patients/:id` - удаление пациента
- `GET /patients/:id/appointments` - приемы пациента
- `GET /patients/:id/medical-history` - анамнез пациента

#### Врачи
- `GET /doctors` - список врачей
- `GET /doctors/:id` - информация о враче
- `GET /doctors/:id/appointments` - приемы врача

#### Приемы
- `GET /appointments` - список приемов (с фильтрацией)
- `GET /appointments/:id` - информация о приеме
- `POST /appointments` - создание приема
- `PUT /appointments/:id` - обновление приема
- `DELETE /appointments/:id` - удаление приема
- `GET /appointments/:id/tests` - тесты приема

#### Медицинский анамнез
- `GET /medical_history` - список записей анамнеза
- `POST /medical_history` - создание записи
- `DELETE /medical_history/:id` - удаление записи

## 🗃 Модели данных

### Patient (Пациент)
```go
type Patient struct {
    ID             uint
    FullName       string
    BirthDate      time.Time
    Gender         string  // "male" или "female"
    Phone          string
    Email          string
    Appointments   []Appointment
    MedicalHistory []MedicalHistory
}
```

### Doctor (Врач)
```go
type Doctor struct {
    ID             uint
    FullName       string
    Specialization string
    Phone          string
    Email          string
    Appointments   []Appointment
}
```

### Appointment (Прием)
```go
type Appointment struct {
    ID           uint
    PatientID    uint
    DoctorID     uint
    Date         time.Time
    Diagnosis    string
    Treatment    string
    Notes        string
    Patient      Patient
    Doctor       Doctor
    MedicalTests []MedicalTest
}
```

## 🔧 Make команды

```bash
make          # Сборка и запуск
make build    # Сборка проекта
make run      # Запуск собранного приложения
make clean    # Очистка сборки и базы данных
```

## ⚙️ Конфигурация

- **Порт**: 8080
- **База данных**: SQLite (clinic.db)
- **CORS**: разрешены все домены

## 🗂 Структура проекта

```
.
├── main.go                 # Основной файл приложения
├── go.mod                  # Модули Go
├── go.sum                  # Зависимости
├── Makefile               # Скрипты сборки
├── clinic.db              # База данных (создается автоматически)
└── docs/                  # Документация Swagger
    ├── docs.go
    ├── swagger.json
    └── swagger.yaml
```

## 🚀 Быстрый старт

1. **Клонируйте репозиторий**
2. **Запустите приложение**:
   ```bash
   make
   ```
3. **Откройте документацию**: http://localhost:8080/swagger/index.html
4. **Начните тестировать API** используя Swagger UI

## 📊 Тестовые данные

При первом запуске автоматически создаются тестовые данные:
- 5 пациентов
- 4 врача разных специализаций
- Медицинские приемы
- Результаты анализов
- Записи медицинского анамнеза

## 🔍 Примеры запросов

### Создание пациента
```bash
curl -X POST http://localhost:8080/patients \
  -H "Content-Type: application/json" \
  -d '{
    "full_name": "Иванов Петр Сидорович",
    "birth_date": "1980-05-15T00:00:00Z",
    "gender": "male",
    "phone": "+79990000000",
    "email": "ivanov@example.com"
  }'
```

### Получение приемов с фильтрацией
```bash
# Приемы конкретного пациента
curl "http://localhost:8080/appointments?patient_id=1"

# Приемы конкретного врача  
curl "http://localhost:8080/appointments?doctor_id=2"
```

## 🐛 Решение проблем

### Ошибка порта
Если порт 8080 занят, измените порт в функции `main()`:
```go
router.Run(":8081") // или другой свободный порт
```

### Проблемы с базой данных
Удалите файл `clinic.db` и перезапустите приложение:
```bash
make clean
make
```

## 👥 Разработка

### Добавление новых эндпоинтов
1. Добавьте обработчик в `main.go`
2. Добавьте Swagger аннотации
3. Обновите документацию: `swag init`
4. Протестируйте через Swagger UI

### Структура обработчика
```go
// @Summary Описание
// @Tags tag-name
// @Accept json
// @Produce json
// @Param param-name path/query int true "Описание параметра"
// @Success 200 {object} Model
// @Router /endpoint [get]
func handler(c *gin.Context) {
    // логика обработчика
}
```

## 📄 Лицензия

[Указать лицензию проекта]

---

**Примечание**: Это демонстрационное приложение. Для продакшн использования рекомендуется:
- Добавить аутентификацию и авторизацию
- Использовать PostgreSQL или другую production БД
- Добавить логирование
- Настроить окружения (dev/staging/prod)
