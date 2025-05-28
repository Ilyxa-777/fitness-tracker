package spentcalories

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

// parseTraining парсит строку данных активности и возвращает количество шагов, тип активности и продолжительность активности.
func parseTraining(data string) (int, string, time.Duration, error) {
	ds := strings.Split(data, ",")
	if len(ds) != 3 {
		return 0, "", 0, fmt.Errorf("invalid data format")
	}

	activityType := ds[1]

	if !(activityType == "Бег" || activityType == "Ходьба") {
		return 0, activityType, 0, fmt.Errorf("unknown training type")
	}

	steps, err := strconv.Atoi(ds[0])
	if err != nil {
		return 0, activityType, 0, err
	}

	duration, err := time.ParseDuration(ds[2])
	if err != nil {
		return 0, activityType, 0, err
	}

	if steps <= 0 {
		return 0, activityType, 0, fmt.Errorf("number of steps cannot be negative")
	}

	if duration <= 0 {
		return 0, activityType, 0, fmt.Errorf("activity duration cannot be negative")
	}

	return steps, activityType, duration, nil

}

// distance рассчитывает дистанцию, пройденную за тренировку.
func distance(steps int, height float64) float64 {
	stepLength := stepLengthCoefficient * height
	distance := float64(steps) * stepLength
	distanceInKm := distance / mInKm
	return distanceInKm
}

// meanSpeed рассчитывает среднюю скорость, с которой была пройдена дистанция.
func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}

	distance := distance(steps, height)
	meanSpeed := distance / duration.Hours()
	return meanSpeed
}

// возвращает данные тренировки
func TrainingInfo(data string, weight, height float64) (string, error) {
	// TODO: установить точку остановки на следующей строке
	steps, activityType, duration, err := parseTraining(data)
	if err != nil {
		return "", err
	}

	var calories float64
	var errCal error

	switch activityType {
	case "Бег":
		calories, errCal = RunningSpentCalories(steps, weight, height, duration)
	case "Ходьба":
		calories, errCal = WalkingSpentCalories(steps, weight, height, duration)
	default:
		return "", fmt.Errorf("unknown training type")
	}

	if errCal != nil {
		return "", errCal
	}

	distance := distance(steps, height)
	speed := meanSpeed(steps, height, duration)

	return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
		activityType, duration.Hours(), distance, speed, calories), nil
}

// RunningSpentCalories рассчитывает количество сожженных калорий при беге.
func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if duration <= 0 {
		return 0, fmt.Errorf("activity duration cannot be negative or zero")
	}

	if steps <= 0 {
		return 0, fmt.Errorf("number of steps cannot be negative or zero")
	}

	if weight <= 0 {
		return 0, fmt.Errorf("weight cannot be negative or zero")
	}

	if height <= 0 {
		return 0, fmt.Errorf("height cannot be negative or zero")
	}

	meanSpeed := meanSpeed(steps, height, duration)
	durationInMinutes := duration.Minutes()

	calories := (weight * meanSpeed * durationInMinutes) / minInH

	return calories, nil
}

// WalkingSpentCalories рассчитывает количество сожженных калорий при ходьбе.
func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	runningCalories, err := RunningSpentCalories(steps, weight, height, duration)
	if err != nil {
		return .0, err
	}

	walkingCalories := runningCalories * walkingCaloriesCoefficient

	return walkingCalories, nil
}
