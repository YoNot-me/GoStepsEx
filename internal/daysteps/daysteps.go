package daysteps

import (
	"time"
	"strconv"
	"strings"
	"errors"
	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
	"fmt"
	"log"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

func parsePackage(data string) (int, time.Duration, error) {
	sliceString := strings.Split(data, ",")
	if len(sliceString) <= 1 {
		return 0,0,errors.New("parsePackage, data -> slice < 2")
	} else if len(sliceString) > 2 {
		return 0,0,errors.New("parsePackage, data -> slice > 2")
	}

 	steps, err := strconv.Atoi(sliceString[0])
	if err != nil {
		return 0,0, errors.New("strconv.Atoi(sliceString[0])")
	} else if steps <= 0 {
		return 0,0, errors.New("parsePackage, data -> step == 0")
	}

	time, err := time.ParseDuration(sliceString[1])
	if err != nil {
		return 0,0, errors.New("time.ParseDuration(sliceString[1])")
	} else if time <= 0 {
		return 0,0, errors.New("time.ParseDuration(sliceString[1]) - > time <= 0")
	}

	return steps, time, nil
}

func DayActionInfo(data string, weight, height float64) string {
	steps, time, err := parsePackage(data)
	if err != nil {
		log.Println(err)
		return ""
	}
	if steps <= 0 || time <= 0 {
		log.Println(err)
		return 	""
	}
	

	stepLength := (float64(steps)*stepLength)/mInKm
	calories,err := spentcalories.WalkingSpentCalories(steps, weight, height, time)
	if err != nil {
		return ""
	}
	
	returnVal := fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n", steps, stepLength, calories)
	return returnVal
}