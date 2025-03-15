package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/smtp"
	"runtime/debug"
	"spb-remont/models"
	"strconv"
)

// Объявление констант
const (
	employeePath = "/employee" // Убираем слеш в конце
	successPath  = "/success"
)

// Объявление переменных в глобальной области с указанием типов
var employees []models.Employee
var director models.Director
var services []models.Service
var reviews []models.Review
var projects []models.Project

// errorHandler функция для обработки ошибок
func errorHandler(w http.ResponseWriter, err error, message string) {
	log.Printf("%s: %v\nStacktrace:\n%s", message, err, debug.Stack())
	http.Error(w, "Ошибка на сервере", http.StatusInternalServerError) // Не показываем детали ошибки клиенту
}

// Функция для отправки email
func sendEmail(to string, subject string, body string) error {
	from := "mukhammadjonbayitov@icloud.com" // Замените на свой email
	password := "20010001"                   // Замените на свой пароль

	// Настройки SMTP-сервера (Gmail)
	smtpServer := "smtp.gmail.com" // Замените, если используете другой сервер
	smtpPort := 587                // Замените, если используете другой порт

	// Создаем сообщение
	message := []byte(fmt.Sprintf("To: %s\r\n", to) +
		fmt.Sprintf("Subject: %s\r\n", subject) +
		"\r\n" +
		body)

	// Авторизация
	auth := smtp.PlainAuth("", from, password, smtpServer)

	// Отправка email
	err := smtp.SendMail(fmt.Sprintf("%s:%d", smtpServer, smtpPort), auth, from, []string{to}, message)
	if err != nil {
		log.Printf("Ошибка отправки email: %v", err)
		return err
	}
	return nil
}

// indexHandler обработчик для главной страницы
func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html", "templates/layout.html")
	if err != nil {
		errorHandler(w, err, "Error parsing templates in indexHandler")
		return
	}

	// Данные для главной страницы
	data := models.Page{
		Title:     "Главная",
		Employees: employees,
		Director:  director,
	}

	err = tmpl.ExecuteTemplate(w, "layout", data)
	if err != nil {
		errorHandler(w, err, "Error executing template in indexHandler")
		return
	}
}

// servicesHandler обработчик для страницы услуг
func servicesHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/services.html", "templates/layout.html")
	if err != nil {
		errorHandler(w, err, "Error parsing templates in servicesHandler")
		return
	}

	pageStr := r.URL.Query().Get("page")
	page := 1

	if pageStr != "" {
		var err error
		page, err = strconv.Atoi(pageStr)
		if err != nil {
			log.Printf("Error parsing page number in servicesHandler: %v", err)
			http.Error(w, "Неверный номер страницы", http.StatusBadRequest)
			return
		}
	}

	pageSize := 3
	startIndex := (page - 1) * pageSize
	endIndex := startIndex + pageSize
	if endIndex > len(services) {
		endIndex = len(services)
	}

	var prevPage int
	if page > 1 {
		prevPage = page - 1
	}
	var nextPage int
	if endIndex < len(services) {
		nextPage = page + 1
	}

	err = tmpl.ExecuteTemplate(w, "layout", models.Page{
		Title:    "Услуги",
		Services: services[startIndex:endIndex],
		PrevPage: prevPage,
		NextPage: nextPage,
	})
	if err != nil {
		errorHandler(w, err, "Error executing template in servicesHandler")
		return
	}
}

// portfolioHandler обработчик для страницы портфолио
func portfolioHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/portfolio.html", "templates/layout.html")
	if err != nil {
		errorHandler(w, err, "Error parsing templates in portfolioHandler")
		return
	}

	err = tmpl.ExecuteTemplate(w, "layout", models.Page{
		Title:    "Портфолио",
		Projects: projects,
	})
	if err != nil {
		errorHandler(w, err, "Error executing template in portfolioHandler")
		return
	}
}

// contactsHandler обработчик для страницы контактов
func contactsHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/contacts.html", "templates/layout.html")
	if err != nil {
		errorHandler(w, err, "Error parsing templates in contactsHandler")
		return
	}

	err = tmpl.ExecuteTemplate(w, "layout", models.Page{
		Title: "Контакты",
	})
	if err != nil {
		errorHandler(w, err, "Error executing template in contactsHandler")
		return
	}
}

// employeesHandler обработчик для страницы сотрудников
func employeesHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/employees.html", "templates/layout.html")
	if err != nil {
		errorHandler(w, err, "Error parsing templates in employeesHandler")
		return
	}

	err = tmpl.ExecuteTemplate(w, "layout", models.Page{
		Title:     "Сотрудники",
		Employees: employees,
	})
	if err != nil {
		errorHandler(w, err, "Error executing template in employeesHandler")
		return
	}
}

// employeeHandler обработчик для страницы сотрудника
func employeeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/employee.html", "templates/layout.html")
	if err != nil {
		errorHandler(w, err, "Error parsing templates in employeeHandler")
		return
	}

	// Получаем ID сотрудника из параметров запроса
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "Не указан ID сотрудника", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("Error parsing employee ID in employeeHandler: %v", err)
		http.Error(w, "Неверный ID сотрудника", http.StatusBadRequest)
		return
	}

	// Ищем сотрудника по ID
	var selectedEmployee models.Employee
	for _, employee := range employees {
		if employee.ID == id {
			selectedEmployee = employee
			break
		}
	}

	// Если сотрудник не найден, возвращаем 404
	if selectedEmployee.ID == 0 {
		http.NotFound(w, r)
		return
	}

	// Выполняем шаблон
	err = tmpl.ExecuteTemplate(w, "layout", models.Page{
		Title:    "Информация о сотруднике",
		Employee: selectedEmployee,
	})
	if err != nil {
		errorHandler(w, err, "Error executing template in employeeHandler")
		return
	}
}

// reviewsHandler обработчик для страницы отзывов
func reviewsHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/reviews.html", "templates/layout.html")
	if err != nil {
		errorHandler(w, err, "Error parsing templates in reviewsHandler")
		return
	}

	err = tmpl.ExecuteTemplate(w, "layout", models.Page{
		Title:   "Отзывы",
		Reviews: reviews,
	})
	if err != nil {
		errorHandler(w, err, "Error executing template in reviewsHandler")
		return
	}
}

// directorHandler обработчик для страницы директора
func directorHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/director.html", "templates/layout.html")
	if err != nil {
		errorHandler(w, err, "Error parsing templates in directorHandler")
		return
	}

	err = tmpl.ExecuteTemplate(w, "layout", models.Page{
		Title:    "Директор",
		Director: director,
	})
	if err != nil {
		errorHandler(w, err, "Error executing template in directorHandler")
		return
	}
}

// successHandler обработчик для страницы успеха
func successHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/success.html", "templates/layout.html")
	if err != nil {
		errorHandler(w, err, "Error parsing templates in successHandler")
		return
	}

	// Парсим форму
	err = r.ParseForm()
	if err != nil {
		log.Printf("Error parsing form in successHandler: %v", err)
		http.Error(w, "Ошибка обработки формы", http.StatusBadRequest)
		return
	}

	// Получаем данные из формы
	name := r.Form.Get("name")
	email := r.Form.Get("email")
	message := r.Form.Get("message")

	// Отправляем email
	subject := "Сообщение с вашего сайта"
	body := fmt.Sprintf("Имя: %s\nEmail: %s\nСообщение: %s", name, email, message)
	err = sendEmail(email, subject, body)
	if err != nil {
		log.Printf("Ошибка отправки email: %v", err)
		http.Error(w, "Ошибка при отправке сообщения", http.StatusInternalServerError)
		return
	}

	// Логируем данные (в production лучше использовать более надежный способ логирования)
	fmt.Printf("Name: %s, Email: %s, Message: %s \n", name, email, message)

	// Выполняем шаблон
	err = tmpl.ExecuteTemplate(w, "layout", models.Page{
		Title:   "Сообщение отправлено!",
		Name:    name,
		Email:   email,
		Message: message,
	})
	if err != nil {
		errorHandler(w, err, "Error executing template in successHandler")
		return
	}

	// Перенаправляем пользователя на страницу подтверждения (например, на главную страницу)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func main() {
	// Инициализация значений переменных в функции main
	employees = []models.Employee{
		{ID: 1, FirstName: "Юрий", LastName: "Юрийов", Position: "Строитель", Description: "Опытный строитель", Image: "/static/images/yuriy.jpg"},
		{ID: 2, FirstName: "Наташа", LastName: "Карпович", Position: "Менеджер", Description: "Квалифицированный менеджер", Image: "/static/images/natasha.jpg"},
	}
	director = models.Director{
		ID: 1, FirstName: "Мухаммаджон", LastName: "Байитов", Position: "Генеральный директор", Description: "Генеральный директор компании", Image: "/static/images/mukhammadjon.jpg"}
	services = []models.Service{
		{ID: 1, Name: "Ремонт ванной", Description: "Ремонт ванной комнаты под ключ", Price: 50000},
		{ID: 2, Name: "Электромонтаж", Description: "Полный электромонтаж квартиры", Price: 30000},
		{ID: 3, Name: "Малярные работы", Description: "Покраска стен и потолков", Price: 20000},
	}
	reviews = []models.Review{
		{ID: 1, Author: "Иван", Text: "Отличная работа!"},
		{ID: 2, Author: "Петр", Text: "Все понравилось."},
	}
	projects = []models.Project{
		{ID: 1, Name: "Квартира на Ленина", Description: "Ремонт двухкомнатной квартиры",
			Images: []string{"/static/images/img1.jpg", "/static/images/img2.jpg"},
			Videos: []string{"/static/videos/video1.mp4", "/static/videos/video2.mp4"}},
	}

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/services", servicesHandler)
	http.HandleFunc("/portfolio", portfolioHandler)
	http.HandleFunc("/contacts", contactsHandler)
	http.HandleFunc("/employees", employeesHandler)
	http.HandleFunc("/employee", employeeHandler) // Изменили путь
	http.HandleFunc("/reviews", reviewsHandler)
	http.HandleFunc("/director", directorHandler)
	http.HandleFunc("/success", successHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
