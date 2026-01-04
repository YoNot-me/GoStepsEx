package spentcalories

import (
	"errors"
	"log"
	"strconv"
	"strings"
	"time"
	"fmt"
	// "reflect" без надобности?
	// "github.com/golang/protobuf/ptypes/duration" зачем он? Я ознакомился с пакетом, но хватает пакета time
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

func parseTraining(data string) (int, string, time.Duration, error) {
	sliceString := strings.Split(data, ",")
	if len(sliceString) < 3 || len(sliceString) > 3 {
		return 0,"",0, errors.New("parseTraining, data -> slice < 3 || slice > 3")
	}

	steps, err := strconv.Atoi(sliceString[0])
	if err != nil {
		return 0,"",0, errors.New("parseTraining, slice -> steps err")
	}

	if steps <= 0 {
		return 0,"",0, errors.New("parseTraining, steps empty")
	}

	time, err := time.ParseDuration(sliceString[2])
	if err != nil {
		return 0,"",0, errors.New("time.ParseDuration, slice -> time err")
	} else if time <= 0 {
		return 0,"",0, errors.New("time.ParseDuration, slice -> time 0")
	}

	return steps, sliceString[1], time, nil
}

func distance(steps int, height float64) float64 {
	return (float64(steps)*(height*stepLengthCoefficient))/mInKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}
	distance := distance(steps, height)
	return distance / duration.Hours()
	
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, trainingType, time, err := parseTraining(data)
	if err != nil {
		log.Println(err)
	} else if weight <= 0 {
		return "", errors.New("TrainingInfo, weight <= 0")
	} else if height <= 0 {
		return "", errors.New("TrainingInfo, height <= 0")
	}

	walkDistance := distance(steps, height)
	meanSpeed := meanSpeed(steps, height, time)


	switch trainingType {
		case "Ходьба":	
			calories, err := WalkingSpentCalories(steps, weight, height, time)
			if err != nil {
				return "", errors.New("TrainingInfo, swtich -> Ходьба err")
			}

			stringReturn := fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n", trainingType, time.Hours(), walkDistance, meanSpeed, calories)
			return stringReturn, nil

		case "Бег":
			calories, err := RunningSpentCalories(steps, weight, height, time)
			if err != nil {
				return "", errors.New("TrainingInfo, swtich -> Бег err")
			}

			stringReturn := fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n", trainingType, time.Hours(), walkDistance, meanSpeed, calories)
			return stringReturn, nil

		default:
			return "", errors.New("неизвестный тип тренировки")

	}
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, errors.New("RunningSpentCalories, steps <= 0")
	} else if weight <= 0 {
		return 0, errors.New("RunningSpentCalories, weight <= 0")
	} else if height <= 0 {
		return 0, errors.New("RunningSpentCalories, height <= 0")
	} else if duration <= 0 {
		return 0, errors.New("RunningSpentCalories, duration <= 0")
	}

	averagSpeed := meanSpeed(steps, height, duration)
	minutesDuration := duration.Minutes()
	return  (weight * averagSpeed * minutesDuration) / minInH, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, errors.New("RunningSpentCalories, steps <= 0")
	} else if weight <= 0 {
		return 0, errors.New("RunningSpentCalories, weight <= 0")
	} else if height <= 0 {
		return 0, errors.New("RunningSpentCalories, height <= 0")
	} else if duration <= 0 {
		return 0, errors.New("RunningSpentCalories, duration <= 0")
	}

	averagSpeed := meanSpeed(steps, height, duration)
	minutesDuration := duration.Minutes()
	calories :=  (weight * averagSpeed * minutesDuration) / minInH
	return calories * walkingCaloriesCoefficient, nil
}
