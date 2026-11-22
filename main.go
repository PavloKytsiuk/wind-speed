package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/PavloKytsiuk/wind-speed/analyzer"
	"github.com/compilyator/VNTU_GoLang/labs/Lab2Data/wind"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Введіть місто для генерації даних швидкості вітру (наприклад, Вінниця): ")
	if !scanner.Scan() {
		fmt.Println("\nПомилка читання міста.")
		return
	}
	city := strings.TrimSpace(scanner.Text())
	if city == "" {
		fmt.Println("Місто не може бути порожнім. Будь ласка, введіть дійсне місто.")
		return
	}

	fmt.Print("Введіть значення seed для генерації даних: ")
	if !scanner.Scan() {
		fmt.Println("\nПомилка читання seed.")
		return
	}
	seed, err := strconv.ParseInt(strings.TrimSpace(scanner.Text()), 10, 64)
	if err != nil {
		fmt.Println("Невірне значення seed. Будь ласка, введіть ціле число.")
		return
	}

	fmt.Print("Введіть кількість даних (від 5 до 100): ")
	if !scanner.Scan() {
		fmt.Println("\nПомилка читання кількості.")
		return
	}
	count, err := strconv.Atoi(strings.TrimSpace(scanner.Text()))
	if err != nil || count < analyzer.MinCount || count > analyzer.MaxCount {
		fmt.Printf("Помилка: кількість даних повинна бути від %d до %d.\n", analyzer.MinCount, analyzer.MaxCount)
		return
	}

	windSpeed, err := wind.GenerateWithSeed(count, city, seed)
	if err != nil {
		fmt.Println("Помилка при генерації даних швидкості вітру:", err)
		return
	}

	series, err := analyzer.NewWindSeries(city, windSpeed)
	if err != nil {
		fmt.Println("Помилка валідації даних:", err)
		return
	}

	report := series.Analyze()
	printReport(report, series.Speeds())
}

func printReport(report analyzer.WindReport, speeds []float64) {
	fmt.Printf("Мінімум: %.2f км/год\n", report.Min)
	fmt.Printf("Максимум: %.2f км/год\n", report.Max)
	fmt.Printf("Сума: %.2f км/год\n", report.Sum)
	fmt.Printf("Середнє: %.2f км/год\n", report.Average)
	fmt.Printf("Місто: %s\n", report.City)
	fmt.Printf("Отриманий зріз швидкостей (км/год): %v\n\n", speeds)

	fmt.Println("Категорії швидкості вітру:")
	for _, cat := range report.Categories {
		fmt.Printf("- %s (%s): %d\n", cat.Name, cat.RangeLabel, cat.Count)
	}

	fmt.Println("\nПовтори:")
	if len(report.Repeats) == 0 {
		fmt.Println("Повторів не знайдено.")
		return
	}
	for _, rep := range report.Repeats {
		fmt.Printf("Значення %.2f зустрічається %d разів\n", rep.Value, rep.Count)
	}
}
