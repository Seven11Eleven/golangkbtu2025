package practices

import "fmt"

type Car struct {
	Brand string
	Model string
	Year  int
}

type Vehicle interface {
	StartEngine() string
	Info() string
}

func NewCar(brand, model string, year int) *Car {
	return &Car{
		Brand: brand,
		Model: model,
		Year:  year,
	}
}

func (c *Car) StartEngine() string {
	return fmt.Sprintf("Engine of %s %s started.", c.Brand, c.Model)
}

func (c *Car) Info() string {
	return fmt.Sprintf("Brand: %s, Model: %s, Year: %d", c.Brand, c.Model, c.Year)
}

func GroupCarsByBrand(cars []Car) map[string][]*Car {
	groupedCars := make(map[string][]*Car)
	for i := range cars {
		groupedCars[cars[i].Brand] = append(groupedCars[cars[i].Brand], &cars[i])
	}
	return groupedCars
}
