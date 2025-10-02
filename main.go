package main

import (
	"fmt"
	"net/http"
	"time"

	_ "demeda/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// @title Demo Medical Database API
// @version 1.0
// @description REST API для медицинской клиники с управлением пациентами, врачами, приемами и анамнезом
// @termsOfService http://swagger.io/terms/

// @host localhost:8080
// @BasePath /
// @schemes http

// ErrorResponse представляет стандартный ответ об ошибке
type ErrorResponse struct {
	Error string `json:"error"`
}

// Patient представляет пациента клиники
// @Description Информация о пациенте
type Patient struct {
	ID             uint             `gorm:"primaryKey" json:"id"`
	CreatedAt      time.Time        `json:"created_at"`
	FullName       string           `gorm:"not null" json:"full_name"`
	BirthDate      time.Time        `gorm:"not null" json:"birth_date"`
	Gender         string           `gorm:"not null;check:gender IN ('male','female')" json:"gender"`
	Phone          string           `json:"phone"`
	Email          string           `json:"email"`
	Appointments   []Appointment    `json:"appointments,omitempty"`
	MedicalHistory []MedicalHistory `json:"medical_history,omitempty"`
}

// Doctor представляет врача клиники
// @Description Информация о враче
type Doctor struct {
	ID             uint          `gorm:"primaryKey" json:"id"`
	CreatedAt      time.Time     `json:"created_at"`
	FullName       string        `gorm:"not null" json:"full_name"`
	Specialization string        `gorm:"not null" json:"specialization"`
	Phone          string        `json:"phone"`
	Email          string        `json:"email"`
	Appointments   []Appointment `json:"appointments,omitempty"`
}

// Appointment представляет медицинский прием
// @Description Информация о медицинском приеме
type Appointment struct {
	ID           uint          `gorm:"primaryKey" json:"id"`
	CreatedAt    time.Time     `json:"created_at"`
	PatientID    uint          `gorm:"not null" json:"patient_id"`
	DoctorID     uint          `gorm:"not null" json:"doctor_id"`
	Date         time.Time     `gorm:"not null" json:"date"`
	Diagnosis    string        `json:"diagnosis"`
	Treatment    string        `json:"treatment"`
	Notes        string        `json:"notes"`
	Patient      Patient       `gorm:"foreignKey:PatientID" json:"patient,omitempty"`
	Doctor       Doctor        `gorm:"foreignKey:DoctorID" json:"doctor,omitempty"`
	MedicalTests []MedicalTest `json:"medical_tests,omitempty"`
}

// MedicalTest представляет медицинский тест
// @Description Результаты медицинских тестов
type MedicalTest struct {
	ID             uint        `gorm:"primaryKey" json:"id"`
	CreatedAt      time.Time   `json:"created_at"`
	AppointmentID  uint        `gorm:"not null" json:"appointment_id"`
	Name           string      `gorm:"not null" json:"name"`
	Result         string      `json:"result"`
	Unit           string      `json:"unit"`
	ReferenceRange string      `json:"reference_range"`
	Appointment    Appointment `gorm:"foreignKey:AppointmentID" json:"appointment,omitempty"`
}

// MedicalHistory представляет запись медицинского анамнеза
// @Description Медицинский анамнез пациента
type MedicalHistory struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	PatientID   uint      `gorm:"not null" json:"patient_id"`
	HistoryType string    `gorm:"not null" json:"history_type"`
	Description string    `gorm:"not null" json:"description"`
	StartDate   time.Time `json:"start_date"`
	Severity    string    `json:"severity"`
	Status      string    `json:"status"`
	Notes       string    `json:"notes"`
	Patient     Patient   `gorm:"foreignKey:PatientID" json:"patient,omitempty"`
}

// DTO для создания/обновления записей
type CreatePatientRequest struct {
	FullName  string    `json:"full_name" binding:"required"`
	BirthDate time.Time `json:"birth_date" binding:"required"`
	Gender    string    `json:"gender" binding:"required"`
	Phone     string    `json:"phone"`
	Email     string    `json:"email"`
}

type CreateAppointmentRequest struct {
	PatientID uint      `json:"patient_id" binding:"required"`
	DoctorID  uint      `json:"doctor_id" binding:"required"`
	Date      time.Time `json:"date" binding:"required"`
	Diagnosis string    `json:"diagnosis"`
	Treatment string    `json:"treatment"`
	Notes     string    `json:"notes"`
}

type CreateMedicalHistoryRequest struct {
	PatientID   uint      `json:"patient_id" binding:"required"`
	HistoryType string    `json:"history_type" binding:"required"`
	Description string    `json:"description" binding:"required"`
	StartDate   time.Time `json:"start_date"`
	Severity    string    `json:"severity"`
	Status      string    `json:"status"`
	Notes       string    `json:"notes"`
}

var db *gorm.DB

func main() {
	var err error
	db, err = gorm.Open(sqlite.Open("clinic.db"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}

	// Автоматическое создание таблиц
	err = db.AutoMigrate(&Patient{}, &Doctor{}, &Appointment{}, &MedicalTest{}, &MedicalHistory{})
	if err != nil {
		panic("Database migration failed")
	}

	// Генерация тестовых данных
	seedDatabase(db)

	// Настройка роутера
	router := gin.Default()

	// CORS middleware
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Группа маршрутов для пациентов
	patients := router.Group("/patients")
	{
		patients.GET("", getPatients)
		patients.GET("/:id", getPatient)
		patients.POST("", createPatient)
		patients.PUT("/:id", updatePatient)
		patients.DELETE("/:id", deletePatient)
		patients.GET("/:id/appointments", getPatientAppointments)
		patients.GET("/:id/medical-history", getPatientMedicalHistory)
	}

	// Группа маршрутов для врачей
	doctors := router.Group("/doctors")
	{
		doctors.GET("", getDoctors)
		doctors.GET("/:id", getDoctor)
		doctors.GET("/:id/appointments", getDoctorAppointments)
	}

	// Группа маршрутов для приемов
	appointments := router.Group("/appointments")
	{
		appointments.GET("", getAppointments)
		appointments.GET("/:id", getAppointment)
		appointments.POST("", createAppointment)
		appointments.PUT("/:id", updateAppointment)
		appointments.DELETE("/:id", deleteAppointment)
		appointments.GET("/:id/tests", getAppointmentTests)
	}

	// Группа маршрутов для анамнеза
	medicalHistory := router.Group("/medical_history")
	{
		medicalHistory.GET("", getMedicalHistory)
		medicalHistory.POST("", createMedicalHistory)
		medicalHistory.DELETE("/:id", deleteMedicalHistory)
	}

	// Запуск сервера
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	fmt.Println("Сервер запущен на http://localhost:8080")
	fmt.Println("Swagger документация доступна на http://localhost:8080/swagger/index.html")
	router.Run(":8080")
}

// Обработчики для пациентов

// GetPatients godoc
// @Summary Получить список пациентов
// @Description Получить список всех пациентов
// @Tags patients
// @Accept json
// @Produce json
// @Success 200 {array} Patient
// @Failure 500 {object} ErrorResponse
// @Router /patients [get]
func getPatients(c *gin.Context) {
	var patients []Patient
	if err := db.Find(&patients).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, patients)
}

// GetPatient godoc
// @Summary Получить пациента по ID
// @Description Получить подробную информацию о пациенте включая анамнез и приемы
// @Tags patients
// @Accept json
// @Produce json
// @Param id path int true "ID пациента"
// @Success 200 {object} Patient
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /patients/{id} [get]
func getPatient(c *gin.Context) {
	id := c.Param("id")
	var patient Patient
	if err := db.Preload("MedicalHistory").Preload("Appointments").Preload("Appointments.Doctor").First(&patient, id).Error; err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Patient not found"})
		return
	}
	c.JSON(http.StatusOK, patient)
}

// CreatePatient godoc
// @Summary Создать нового пациента
// @Description Создать запись нового пациента в системе
// @Tags patients
// @Accept json
// @Produce json
// @Param patient body CreatePatientRequest true "Данные пациента"
// @Success 201 {object} Patient
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /patients [post]
func createPatient(c *gin.Context) {
	var req CreatePatientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	patient := Patient{
		FullName:  req.FullName,
		BirthDate: req.BirthDate,
		Gender:    req.Gender,
		Phone:     req.Phone,
		Email:     req.Email,
	}

	if err := db.Create(&patient).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, patient)
}

// UpdatePatient godoc
// @Summary Обновить данные пациента
// @Description Обновить информацию о существующем пациенте
// @Tags patients
// @Accept json
// @Produce json
// @Param id path int true "ID пациента"
// @Param patient body CreatePatientRequest true "Обновленные данные пациента"
// @Success 200 {object} Patient
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /patients/{id} [put]
func updatePatient(c *gin.Context) {
	id := c.Param("id")
	var patient Patient
	if err := db.First(&patient, id).Error; err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Patient not found"})
		return
	}

	var req CreatePatientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	patient.FullName = req.FullName
	patient.BirthDate = req.BirthDate
	patient.Gender = req.Gender
	patient.Phone = req.Phone
	patient.Email = req.Email

	if err := db.Save(&patient).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, patient)
}

// DeletePatient godoc
// @Summary Удалить пациента
// @Description Удалить запись пациента из системы
// @Tags patients
// @Accept json
// @Produce json
// @Param id path int true "ID пациента"
// @Success 200 {object} string
// @Failure 500 {object} ErrorResponse
// @Router /patients/{id} [delete]
func deletePatient(c *gin.Context) {
	id := c.Param("id")
	if err := db.Delete(&Patient{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, "Patient deleted")
}

// GetPatientAppointments godoc
// @Summary Получить приемы пациента
// @Description Получить список всех приемов конкретного пациента
// @Tags patients
// @Accept json
// @Produce json
// @Param id path int true "ID пациента"
// @Success 200 {array} Appointment
// @Failure 500 {object} ErrorResponse
// @Router /patients/{id}/appointments [get]
func getPatientAppointments(c *gin.Context) {
	id := c.Param("id")
	var appointments []Appointment
	if err := db.Preload("Doctor").Where("patient_id = ?", id).Find(&appointments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, appointments)
}

// GetPatientMedicalHistory godoc
// @Summary Получить анамнез пациента
// @Description Получить медицинский анамнез конкретного пациента
// @Tags patients
// @Accept json
// @Produce json
// @Param id path int true "ID пациента"
// @Success 200 {array} MedicalHistory
// @Failure 500 {object} ErrorResponse
// @Router /patients/{id}/medical-history [get]
func getPatientMedicalHistory(c *gin.Context) {
	id := c.Param("id")
	var history []MedicalHistory
	if err := db.Where("patient_id = ?", id).Find(&history).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, history)
}

// Обработчики для врачей

// GetDoctors godoc
// @Summary Получить список врачей
// @Description Получить список всех врачей клиники
// @Tags doctors
// @Accept json
// @Produce json
// @Success 200 {array} Doctor
// @Failure 500 {object} ErrorResponse
// @Router /doctors [get]
func getDoctors(c *gin.Context) {
	var doctors []Doctor
	if err := db.Find(&doctors).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, doctors)
}

// GetDoctor godoc
// @Summary Получить врача по ID
// @Description Получить подробную информацию о враче
// @Tags doctors
// @Accept json
// @Produce json
// @Param id path int true "ID врача"
// @Success 200 {object} Doctor
// @Failure 404 {object} ErrorResponse
// @Router /doctors/{id} [get]
func getDoctor(c *gin.Context) {
	id := c.Param("id")
	var doctor Doctor
	if err := db.First(&doctor, id).Error; err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Doctor not found"})
		return
	}
	c.JSON(http.StatusOK, doctor)
}

// GetDoctorAppointments godoc
// @Summary Получить приемы врача
// @Description Получить список всех приемов конкретного врача
// @Tags doctors
// @Accept json
// @Produce json
// @Param id path int true "ID врача"
// @Success 200 {array} Appointment
// @Failure 500 {object} ErrorResponse
// @Router /doctors/{id}/appointments [get]
func getDoctorAppointments(c *gin.Context) {
	id := c.Param("id")
	var appointments []Appointment
	if err := db.Preload("Patient").Where("doctor_id = ?", id).Find(&appointments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, appointments)
}

// Обработчики для приемов

// GetAppointments godoc
// @Summary Получить список приемов
// @Description Получить список медицинских приемов с возможностью фильтрации
// @Tags appointments
// @Accept json
// @Produce json
// @Param patient_id query int false "Фильтр по ID пациента"
// @Param doctor_id query int false "Фильтр по ID врача"
// @Success 200 {array} Appointment
// @Failure 500 {object} ErrorResponse
// @Router /appointments [get]
func getAppointments(c *gin.Context) {
	var appointments []Appointment
	query := db.Preload("Patient").Preload("Doctor")

	// Фильтрация по patient_id если указана
	if patientID := c.Query("patient_id"); patientID != "" {
		query = query.Where("patient_id = ?", patientID)
	}

	// Фильтрация по doctor_id если указана
	if doctorID := c.Query("doctor_id"); doctorID != "" {
		query = query.Where("doctor_id = ?", doctorID)
	}

	if err := query.Find(&appointments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, appointments)
}

// GetAppointment godoc
// @Summary Получить прием по ID
// @Description Получить подробную информацию о медицинском приеме
// @Tags appointments
// @Accept json
// @Produce json
// @Param id path int true "ID приема"
// @Success 200 {object} Appointment
// @Failure 404 {object} ErrorResponse
// @Router /appointments/{id} [get]
func getAppointment(c *gin.Context) {
	id := c.Param("id")
	var appointment Appointment
	if err := db.Preload("Patient").Preload("Doctor").Preload("MedicalTests").First(&appointment, id).Error; err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Appointment not found"})
		return
	}
	c.JSON(http.StatusOK, appointment)
}

// CreateAppointment godoc
// @Summary Создать новый прием
// @Description Создать запись о медицинском приеме
// @Tags appointments
// @Accept json
// @Produce json
// @Param appointment body CreateAppointmentRequest true "Данные приема"
// @Success 201 {object} Appointment
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /appointments [post]
func createAppointment(c *gin.Context) {
	var req CreateAppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	appointment := Appointment{
		PatientID: req.PatientID,
		DoctorID:  req.DoctorID,
		Date:      req.Date,
		Diagnosis: req.Diagnosis,
		Treatment: req.Treatment,
		Notes:     req.Notes,
	}

	if err := db.Create(&appointment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, appointment)
}

// UpdateAppointment godoc
// @Summary Обновить данные приема
// @Description Обновить информацию о медицинском приеме
// @Tags appointments
// @Accept json
// @Produce json
// @Param id path int true "ID приема"
// @Param appointment body CreateAppointmentRequest true "Обновленные данные приема"
// @Success 200 {object} Appointment
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /appointments/{id} [put]
func updateAppointment(c *gin.Context) {
	id := c.Param("id")
	var appointment Appointment
	if err := db.First(&appointment, id).Error; err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Appointment not found"})
		return
	}

	var req CreateAppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	appointment.PatientID = req.PatientID
	appointment.DoctorID = req.DoctorID
	appointment.Date = req.Date
	appointment.Diagnosis = req.Diagnosis
	appointment.Treatment = req.Treatment
	appointment.Notes = req.Notes

	if err := db.Save(&appointment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, appointment)
}

// DeleteAppointment godoc
// @Summary Удалить прием
// @Description Удалить запись о медицинском приеме
// @Tags appointments
// @Accept json
// @Produce json
// @Param id path int true "ID приема"
// @Success 200 {object} string
// @Failure 500 {object} ErrorResponse
// @Router /appointments/{id} [delete]
func deleteAppointment(c *gin.Context) {
	id := c.Param("id")
	if err := db.Delete(&Appointment{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, "Appointment deleted")
}

// GetAppointmentTests godoc
// @Summary Получить тесты приема
// @Description Получить медицинские тесты конкретного приема
// @Tags appointments
// @Accept json
// @Produce json
// @Param id path int true "ID приема"
// @Success 200 {array} MedicalTest
// @Failure 500 {object} ErrorResponse
// @Router /appointments/{id}/tests [get]
func getAppointmentTests(c *gin.Context) {
	id := c.Param("id")
	var tests []MedicalTest
	if err := db.Where("appointment_id = ?", id).Find(&tests).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, tests)
}

// Обработчики для анамнеза

// GetMedicalHistory godoc
// @Summary Получить анамнез
// @Description Получить записи медицинского анамнеза с возможностью фильтрации
// @Tags medical-history
// @Accept json
// @Produce json
// @Param patient_id query int false "Фильтр по ID пациента"
// @Param type query string false "Фильтр по типу анамнеза"
// @Success 200 {array} MedicalHistory
// @Failure 500 {object} ErrorResponse
// @Router /medical-history [get]
func getMedicalHistory(c *gin.Context) {
	var history []MedicalHistory
	query := db.Preload("Patient")

	if patientID := c.Query("patient_id"); patientID != "" {
		query = query.Where("patient_id = ?", patientID)
	}

	if historyType := c.Query("type"); historyType != "" {
		query = query.Where("history_type = ?", historyType)
	}

	if err := query.Find(&history).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, history)
}

// CreateMedicalHistory godoc
// @Summary Создать запись анамнеза
// @Description Создать новую запись в медицинском анамнезе пациента
// @Tags medical-history
// @Accept json
// @Produce json
// @Param history body CreateMedicalHistoryRequest true "Данные анамнеза"
// @Success 201 {object} MedicalHistory
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /medical-history [post]
func createMedicalHistory(c *gin.Context) {
	var req CreateMedicalHistoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	history := MedicalHistory{
		PatientID:   req.PatientID,
		HistoryType: req.HistoryType,
		Description: req.Description,
		StartDate:   req.StartDate,
		Severity:    req.Severity,
		Status:      req.Status,
		Notes:       req.Notes,
	}

	if err := db.Create(&history).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, history)
}

// DeleteMedicalHistory godoc
// @Summary Удалить запись анамнеза
// @Description Удалить запись из медицинского анамнеза
// @Tags medical-history
// @Accept json
// @Produce json
// @Param id path int true "ID записи анамнеза"
// @Success 200 {object} string
// @Failure 500 {object} ErrorResponse
// @Router /medical-history/{id} [delete]
func deleteMedicalHistory(c *gin.Context) {
	id := c.Param("id")
	if err := db.Delete(&MedicalHistory{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, "Medical history record deleted")
}

func seedDatabase(db *gorm.DB) {
	// Очистка существующих данных
	db.Exec("DELETE FROM medical_histories")
	db.Exec("DELETE FROM medical_tests")
	db.Exec("DELETE FROM appointments")
	db.Exec("DELETE FROM patients")
	db.Exec("DELETE FROM doctors")

	// Генерация пациентов
	patients := []Patient{
		{FullName: "Иванов Иван Иванович", BirthDate: time.Date(1985, 5, 15, 0, 0, 0, 0, time.UTC), Gender: "male", Phone: "+79990000001", Email: "ivanov@mail.ru"},
		{FullName: "Петрова Мария Сергеевна", BirthDate: time.Date(1990, 8, 22, 0, 0, 0, 0, time.UTC), Gender: "female", Phone: "+79990000002", Email: "petrova@mail.ru"},
		{FullName: "Сидоров Алексей Владимирович", BirthDate: time.Date(1978, 3, 10, 0, 0, 0, 0, time.UTC), Gender: "male", Phone: "+79990000003", Email: "sidorov@mail.ru"},
		{FullName: "Кузнецова Елена Викторовна", BirthDate: time.Date(1982, 11, 5, 0, 0, 0, 0, time.UTC), Gender: "female", Phone: "+79990000004", Email: "kuznetsova@mail.ru"},
		{FullName: "Смирнов Дмитрий Петрович", BirthDate: time.Date(1995, 7, 30, 0, 0, 0, 0, time.UTC), Gender: "male", Phone: "+79990000005", Email: "smirnov@mail.ru"},
	}
	db.Create(&patients)

	// Генерация врачей
	doctors := []Doctor{
		{FullName: "Прохоров Андрей Васильевич", Specialization: "Кардиолог", Phone: "+79991111111", Email: "prokhorov@clinic.ru"},
		{FullName: "Громова Ольга Игоревна", Specialization: "Невролог", Phone: "+79991111112", Email: "gromova@clinic.ru"},
		{FullName: "Белов Станислав Михайлович", Specialization: "Терапевт", Phone: "+79991111113", Email: "belov@clinic.ru"},
		{FullName: "Ковальчук Анна Денисовна", Specialization: "Офтальмолог", Phone: "+79991111114", Email: "kovalchuk@clinic.ru"},
	}
	db.Create(&doctors)

	// Генерация приемов
	appointments := []Appointment{
		{PatientID: 1, DoctorID: 1, Date: time.Now().Add(-24 * time.Hour), Diagnosis: "Гипертония", Treatment: "Контроль давления, лизиноприл 10 мг 1 раз в день", Notes: "Жалобы на головные боли"},
		{PatientID: 2, DoctorID: 2, Date: time.Now().Add(-12 * time.Hour), Diagnosis: "Мигрень", Treatment: "Ибупрофен при болях, режим сна", Notes: "Рекомендован отдых"},
		{PatientID: 3, DoctorID: 3, Date: time.Now().Add(-6 * time.Hour), Diagnosis: "ОРВИ", Treatment: "Обильное питье, парацетамол", Notes: "Температура 37.8"},
		{PatientID: 4, DoctorID: 4, Date: time.Now().Add(-3 * time.Hour), Diagnosis: "Конъюнктивит", Treatment: "Глазные капли Офтальмоферон", Notes: "Назначен повторный прием через 5 дней"},
		{PatientID: 5, DoctorID: 1, Date: time.Now(), Diagnosis: "Аритмия", Treatment: "Холтеровское мониторирование", Notes: "Направлен на дополнительное обследование"},
	}
	db.Create(&appointments)

	// Генерация медицинских тестов
	medicalTests := []MedicalTest{
		{AppointmentID: 1, Name: "Артериальное давление", Result: "140/90", Unit: "мм рт.ст.", ReferenceRange: "120/80"},
		{AppointmentID: 1, Name: "Холестерин", Result: "5.2", Unit: "ммоль/л", ReferenceRange: "3.5-5.2"},
		{AppointmentID: 2, Name: "МРТ головного мозга", Result: "Без патологий", Unit: "-", ReferenceRange: "-"},
		{AppointmentID: 3, Name: "Температура тела", Result: "37.8", Unit: "°C", ReferenceRange: "36.6"},
		{AppointmentID: 4, Name: "Острота зрения", Result: "0.8", Unit: "усл.ед.", ReferenceRange: "1.0"},
		{AppointmentID: 5, Name: "ЭКГ", Result: "Мерцательная аритмия", Unit: "-", ReferenceRange: "Синусовый ритм"},
	}
	db.Create(&medicalTests)

	// Генерация анамнеза
	medicalHistories := []MedicalHistory{
		// Аллергии
		{PatientID: 1, HistoryType: "allergy", Description: "Аллергия на пенициллин", StartDate: time.Date(2005, 1, 1, 0, 0, 0, 0, time.UTC), Severity: "severe", Status: "active", Notes: "Анафилактический шок при приеме"},
		{PatientID: 2, HistoryType: "allergy", Description: "Сезонная аллергия на пыльцу", StartDate: time.Date(2010, 1, 1, 0, 0, 0, 0, time.UTC), Severity: "moderate", Status: "active", Notes: "Обострение весной"},

		// Хронические заболевания
		{PatientID: 1, HistoryType: "chronic", Description: "Артериальная гипертензия", StartDate: time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC), Severity: "moderate", Status: "chronic", Notes: "Постоянный прием препаратов"},
		{PatientID: 3, HistoryType: "chronic", Description: "Сахарный диабет 2 типа", StartDate: time.Date(2018, 1, 1, 0, 0, 0, 0, time.UTC), Severity: "mild", Status: "chronic", Notes: "Контроль диеты"},
		{PatientID: 4, HistoryType: "chronic", Description: "Бронхиальная астма", StartDate: time.Date(2012, 1, 1, 0, 0, 0, 0, time.UTC), Severity: "mild", Status: "chronic", Notes: "Ингалятор по необходимости"},

		// Перенесенные операции
		{PatientID: 2, HistoryType: "surgery", Description: "Аппендэктомия", StartDate: time.Date(2015, 6, 15, 0, 0, 0, 0, time.UTC), Severity: "moderate", Status: "resolved", Notes: "Восстановление прошло без осложнений"},
		{PatientID: 5, HistoryType: "surgery", Description: "Артроскопия коленного сустава", StartDate: time.Date(2020, 3, 10, 0, 0, 0, 0, time.UTC), Severity: "moderate", Status: "resolved", Notes: "Спортивная травма"},

		// Семейный анамнез
		{PatientID: 1, HistoryType: "family", Description: "Инфаркт миокарда у отца в 55 лет", StartDate: time.Date(2010, 1, 1, 0, 0, 0, 0, time.UTC), Severity: "severe", Status: "active", Notes: "Наследственная предрасположенность"},
		{PatientID: 3, HistoryType: "family", Description: "Онкологические заболевания у родственников", StartDate: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), Severity: "moderate", Status: "active", Notes: "Бабушка - рак молочной железы"},

		// Вредные привычки
		{PatientID: 3, HistoryType: "habit", Description: "Курение", StartDate: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), Severity: "moderate", Status: "active", Notes: "10 сигарет в день, 20 лет стажа"},
		{PatientID: 5, HistoryType: "habit", Description: "Злоупотребление алкоголем", StartDate: time.Date(2018, 1, 1, 0, 0, 0, 0, time.UTC), Severity: "mild", Status: "resolved", Notes: "Воздержание 2 года"},
	}
	db.Create(&medicalHistories)
}
