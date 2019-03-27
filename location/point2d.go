package location

import (
	"math"
)

type Point2D struct {
	x, y float64
}

func NewPoint2D(x float64, y float64) *Point2D {
	return &Point2D{x: x, y: y}
}

func (p *Point2D) SetX(v float64) {
	p.x = v
}

func (p *Point2D) SetY(v float64) {
	p.y = v
}

func (p Point2D) GetX() float64 {
	return p.x
}

func (p Point2D) GetY() float64 {
	return p.y
}

func (p Point2D) DistanceTo(d Point2D) float64 {
	return math.Sqrt((p.x-d.x)*(p.x-d.x) + (p.y-d.y)*(p.y-d.y))
}

func (p Point2D) IsOnSameLine(ps ...Point2D) bool {
	if len(ps) == 1 || len(ps) == 0 {
		return true
	}
	if s, diff := p.IsTheSame(ps...); !s || diff < 3 {
		return true
	}
	type vector [2]float64
	var vectors = make(map[vector]struct{})
	for _, v := range ps {
		vectors[vector{p.x - v.x, p.y - v.y}] = struct{}{}
	}
	for k1 := range vectors {
		for k2 := range vectors {
			if math.Round(k1[0]*k2[1]-k1[1]*k2[0]) != 0 {
				return false
			}
		}
	}
	return true
}

func (p Point2D) IsTheSame(ps ...Point2D) (bool, int) {
	if len(ps)+1 == 1 {
		return true, 1
	}
	tmpMap := make(map[Point2D]struct{})
	for _, v := range ps {
		tmpMap[v] = struct{}{}
	}
	tmpMap[p] = struct{}{}
	if len(tmpMap) == 1 {
		return true, 1
	}
	return false, len(tmpMap)
}
