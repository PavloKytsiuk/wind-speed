package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/compilyator/VNTU_GoLang/labs/Lab2Data/wind"
)

func GenerateWithSeed(count int, city string, seed int64) ([]float64, error) {
	return wind.GenerateWithSeed(count, city, seed)
}

func main() {
	var count int
	var city string
	var seed int64

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Введіть місто для генерації даних швидкості вітру (наприклад, Вінниця): ")
	if !scanner.Scan() {
		fmt.Println("\nПомилка читання міста.")
		return
	}
	city = strings.TrimSpace(scanner.Text())
	if city == "" {
		fmt.Println("Місто не може бути порожнім. Будь ласка, введіть дійсне місто.")
		return
	}

	fmt.Print("Введіть значення seed для генерації даних: ")
	if !scanner.Scan() {
		fmt.Println("\nПомилка читання seed.")
		return
	}
	seedValue, err := strconv.ParseInt(strings.TrimSpace(scanner.Text()), 10, 64)
	if err != nil {
		fmt.Println("Невірне значення seed. Будь ласка, введіть ціле число.")
		return
	}
	seed = seedValue

	fmt.Print("Введіть кількість даних (від 5 до 100): ")
	if !scanner.Scan() {
		fmt.Println("\nПомилка читання кількості.")
		return
	}
	countValue, err := strconv.Atoi(strings.TrimSpace(scanner.Text()))
	if err != nil || countValue < 5 || countValue > 100 {
		fmt.Println("Помилка: Кількість даних повинна бути від 5 до 100.")
		return
	}
	count = countValue

	windSpeed, err := GenerateWithSeed(count, city, seed)
	if err != nil {
		fmt.Println("Помилка при генерації даних швидкості вітру:", err)
		return
	}
	categoryCounts := make(map[string]int)
	v := 0
	if v < 0 || v > 150 {
		fmt.Println("Помилка: значення швидкості вітру повинно бути в межах від 0 до 150 км/год.")
		return
	}
	for _, v := range windSpeed {
		switch {
		case v <= 5.0:
			categoryCounts["Штиль"]++
		case v <= 30.0:
			categoryCounts["Слабкий"]++
		case v <= 60.0:
			categoryCounts["Помірний"]++
		default:
			categoryCounts["Сильний"]++
		}
	}
	min := windSpeed[0]
	max := windSpeed[0]
	sum := 0.0

	for _, v := range windSpeed {
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
		sum += v
	}
	average := sum / float64(count)

	fmt.Printf("Мінімум: %.2f км/год\n", min)
	fmt.Printf("Максимум: %.2f км/год\n", max)
	fmt.Printf("Сума: %.2f км/год\n", sum)
	fmt.Printf("Середнє: %.2f км/год\n", average)
	fmt.Printf("Місто: %s\n", city)
	fmt.Printf("Отриманий зріз швидкостей (км/год): %v\n\n", windSpeed)
	fmt.Println("Категорії швидкості вітру:")
	fmt.Printf("- Штиль (<= 5): %d\n", categoryCounts["Штиль"])
	fmt.Printf("- Слабкий (5.1 - 30): %d\n", categoryCounts["Слабкий"])
	fmt.Printf("- Помірний (30.1 - 60): %d\n", categoryCounts["Помірний"])
	fmt.Printf("- Сильний (> 60): %d\n\n", categoryCounts["Сильний"])
	fmt.Println("Повтори (Алгоритм вкладених циклів):")
	var printedO2 []float64
	foundO2 := false

	for i := 0; i < len(windSpeed); i++ {
		alreadyCounted := false
		for _, p := range printedO2 {
			if windSpeed[i] == p {
				alreadyCounted = true
				break
			}
		}
		if alreadyCounted {
			continue
		}

		repeats := 1
		for j := i + 1; j < len(windSpeed); j++ {
			if windSpeed[i] == windSpeed[j] {
				repeats++
			}
		}

		if repeats > 1 {
			fmt.Printf("Значення %.2f зустрічається %d разів\n", windSpeed[i], repeats)
			printedO2 = append(printedO2, windSpeed[i])
			foundO2 = true
		}
	}
	if !foundO2 {
		fmt.Println("Повторів не знайдено.")
	}
	fmt.Println("\nПовтори (Алгоритм через map):")
	freqMap := make(map[float64]int)
	var order []float64

	for _, v := range windSpeed {
		if freqMap[v] == 0 {
			order = append(order, v)
		}
		freqMap[v]++

	}
	foundMap := false
	for _, v := range order {
		if freqMap[v] > 1 {
			fmt.Printf("Значення %.2f зустрічається %d разів\n", v, freqMap[v])
			foundMap = true
		}
	}
	if !foundMap {
		fmt.Println("Повторів не знайдено.")
	}
}
