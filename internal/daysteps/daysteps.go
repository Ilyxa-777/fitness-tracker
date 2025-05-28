package daysteps

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

// parsePackage парсит строку данных активности и возвращает количество шагов и продолжительность активности.
func parsePackage(data string) (int, time.Duration, error) {
	// TODO: реализовать функцию
	ds := strings.Split(data, ",")

	if len(ds) != 2 {
		return 0, 0, fmt.Errorf("invalid activity data format")
	}

	steps, err := strconv.Atoi(ds[0])
	if err != nil {
		return 0, 0, err
	}

	if steps <= 0 {
		return 0, 0, fmt.Errorf("number of steps cannot be negative")
	}

	duration, err := time.ParseDuration(ds[1])
	if err != nil {
		return 0, 0, err
	}

	if duration <= 0 {
		return 0, 0, fmt.Errorf("time cannot be negative")
	}

	return steps, duration, nil
}

// DayActionInfo выводит информацию о количестве шагов, дистанции и сожженных калориях.
func DayActionInfo(data string, weight, height float64) string {
	steps, duration, err := parsePackage(data)
	if err != nil {
		log.Printf("ошибка: %v", err)
		return ""
	}

	distanseInMeters := float64(steps) * stepLength
	distanseInKm := distanseInMeters / mInKm
	spentCalories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		log.Printf("ошибка: %v", err)
		return ""
	}

	return fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n", steps, distanseInKm, spentCalories)
}
