package main

type Tank struct {
	Borders   Rectangle
	Direction int
}

func (tank *Tank) moveTo(direction int, crash bool) {
	switch direction {
	case RIGHT:
		for i := range tank.Borders.Points {
			if crash {
				tank.Borders.Points[i].X--
			} else {
				tank.Borders.Points[i].X++
			}
		}
	case LEFT:
		for i, _ := range tank.Borders.Points {
			if crash {
				tank.Borders.Points[i].X++
			} else {
				tank.Borders.Points[i].X--
			}
		}
	case UP:
		for i, _ := range tank.Borders.Points {
			if crash {
				tank.Borders.Points[i].Y++
			} else {
				tank.Borders.Points[i].Y--
			}
		}
	case DOWN:
		for i, _ := range tank.Borders.Points {
			if crash {
				tank.Borders.Points[i].Y--
			} else {
				tank.Borders.Points[i].Y++
			}
		}
	}
}

