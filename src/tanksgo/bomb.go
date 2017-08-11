package main

type Bomb struct {
	Point		Point
	Direction 	int
	Active 		bool
}

func (bomb *Bomb) init(tank Tank)  {
	bomb.Active = true
	bomb.Direction = tank.Direction
	switch tank.Direction {
	case RIGHT:
		bomb.Point = Point{X: tank.Borders.Points[RIGHTUP].X+1, Y: tank.Borders.Points[RIGHTUP].Y+tankWidth/2 }
	case LEFT:
		bomb.Point = Point{X: tank.Borders.Points[LEFTUP].X-1, Y: tank.Borders.Points[LEFTUP].Y+tankWidth/2}
	case UP:
		bomb.Point = Point{X: tank.Borders.Points[LEFTUP].X+tankWidth/2, Y: tank.Borders.Points[LEFTUP].Y-1 }
	case DOWN:
		bomb.Point = Point{X: tank.Borders.Points[LEFTDOWN].X+tankWidth/2, Y: tank.Borders.Points[LEFTDOWN].Y+1}
	}
}

func (bomb *Bomb) move() {
	switch bomb.Direction {
	case RIGHT:
		bomb.Point.X++
	case LEFT:
		bomb.Point.X--
	case UP:
		bomb.Point.Y--
	case DOWN:
		bomb.Point.Y++
	}
	if bomb.Point.X < 1 || bomb.Point.X >= mapWidth-nameTableWidth || bomb.Point.Y < 1 || bomb.Point.Y >= mapHeight-1 {
		bomb.Active = false
	}
}