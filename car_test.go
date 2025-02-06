package main

import (
	_ "github.com/Seven11Eleven/golangkbtu2025/practices"
	"reflect"
	"sort"
	"testing"
)

func TestNewCar(t *testing.T) {
	car := NewCar("Toyota", "Camry", 2020)
	if car.Brand != "Toyota" || car.Model != "Camry" || car.Year != 2020 {
		t.Errorf("NewCar() failed. Got: %v", car)
	}
}

func TestCarMethods(t *testing.T) {
	car := &Car{Brand: "Ford", Model: "Mustang", Year: 1969}

	startMsg := car.StartEngine()
	if startMsg != "Engine of Ford Mustang started." {
		t.Errorf("Expected 'Engine of Ford Mustang started.', got: %s", startMsg)
	}

	infoMsg := car.Info()
	expectedInfo := "Brand: Ford, Model: Mustang, Year: 1969"
	if infoMsg != expectedInfo {
		t.Errorf("Expected %s, got: %s", expectedInfo, infoMsg)
	}
}

func TestGroupCarsByBrand(t *testing.T) {
	cars := []Car{
		*NewCar("Toyota", "Camry", 2020),
		*NewCar("Toyota", "Corolla", 2018),
		*NewCar("Ford", "Mustang", 1969),
	}

	grouped := GroupCarsByBrand(cars)

	if len(grouped) != 2 {
		t.Errorf("Expected 2 brands, got %d", len(grouped))
	}

	// Проверим, что внутри Toyota ровно 2 машины
	toyotaCars, ok := grouped["Toyota"]
	if !ok {
		t.Error("Expected to have key 'Toyota' in map")
	} else {
		if len(toyotaCars) != 2 {
			t.Errorf("Expected 2 Toyota cars, got %d", len(toyotaCars))
		}
	}

	// Проверим, что внутри Ford ровно 1 машина
	fordCars, ok := grouped["Ford"]
	if !ok {
		t.Error("Expected to have key 'Ford' in map")
	} else {
		if len(fordCars) != 1 {
			t.Errorf("Expected 1 Ford car, got %d", len(fordCars))
		}
	}
}

func TestSortCars(t *testing.T) {
	cars := []Car{
		{Brand: "BMW", Model: "X5", Year: 2022},
		{Brand: "Ford", Model: "Mustang", Year: 1969},
		{Brand: "Toyota", Model: "Camry", Year: 2020},
	}

	// Сортируем по году
	sort.Slice(cars, func(i, j int) bool {
		return cars[i].Year < cars[j].Year
	})

	got := []int{cars[0].Year, cars[1].Year, cars[2].Year}
	want := []int{1969, 2020, 2022}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("SortCars failed. Got %v, wanted %v", got, want)
	}
}
