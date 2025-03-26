package logger

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

func SetupLogger(env string) *logrus.Logger {
	log := logrus.New()

	if env == "prod" {
		// Создание папки logs, если её нет
		if _, err := os.Stat("logs"); os.IsNotExist(err) {
			err := os.Mkdir("logs", 0755) // Создаём папку logs с правами 0755
			if err != nil {
				log.Fatal(err)
			}
		}

		// Формирование имени файла по дате (например, 2025-03-19.log)
		date := time.Now().Format("2006-01-02") // Получаем текущую дату в формате YYYY-MM-DD
		logFilePath := "logs/" + date + ".log"

		// Открытие файла для записи (создание нового или добавление в существующий)
		logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
		if err != nil {
			log.Fatal(err)
		}
		log.SetOutput(logFile) // В продакшн выводим в файл
	} else {
		log.SetOutput(os.Stdout) // Выводим только в консоль
	}

	// Настройка логирования в зависимости от окружения
	switch env {
	case "prod":
		log.SetLevel(logrus.WarnLevel)            // В продакшене логируем только предупреждения и ошибки
		log.SetFormatter(&logrus.JSONFormatter{}) // В продакшене выводим логи в формате JSON
	case "test":
		log.SetLevel(logrus.InfoLevel) // В тестах логируем информационные сообщения и ошибки
		log.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,                  // Добавляем полные метки времени
			TimestampFormat: "2006-01-02 15:04:05", // Настройка формата времени
		})
	case "dev":
		log.SetLevel(logrus.DebugLevel) // В разработке логируем отладочную информацию
		log.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,                  // Полные метки времени для разработчика
			TimestampFormat: "2006-01-02 15:04:05", // Настройка формата времени
		})
	default:
		log.SetLevel(logrus.InfoLevel) // По умолчанию логируем только информационные сообщения и ошибки
		log.SetFormatter(&logrus.TextFormatter{
			TimestampFormat: "2006-01-02 15:04:05", // Формат времени по умолчанию
		})
	}

	// Добавление дополнительной информации, например, о версии приложения или хосте
	log.WithFields(logrus.Fields{
		"app":     "ia-online-golang",
		"version": "1.0.0",
		"host":    "localhost",
	})

	return log
}
