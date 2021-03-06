package map2d

import (
	"encoding/json"
	"github.com/chfanghr/WTFCarProject/grid"
	"github.com/chfanghr/WTFCarProject/location"
)

const perGridSize = 0.01

type Map2D struct {
	Map struct {
		Size struct {
			X float64 `json:"x"`
			Y float64 `json:"y"`
		} `json:"size"`
		Barriers []struct {
			Required [3]location.Point2D `json:"required"`
			Optional []location.Point2D  `json:"optional"`
		} `json:"barriers"`
	} `json:"map"`
}

func NewMap2d(raw []byte) (*Map2D, error) {
	res := &Map2D{}
	err := json.Unmarshal(raw, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}
func (m *Map2D) isValid() bool {
	return m.Map.Size.X > 0 && m.Map.Size.Y > 0 &&
		func() bool {
			for _, b := range m.Map.Barriers {
				isOutOfMap := func(p location.Point2D) bool {
					return !(p.GetX() > m.Map.Size.X || p.GetY() > m.Map.Size.Y || p.GetX() < 0 || p.GetY() < 0)
				}
				for _, p := range b.Required {
					if !isOutOfMap(p) {
						return false
					}
				}
				for _, p := range b.Optional {
					if !isOutOfMap(p) {
						return false
					}
				}
			}
			return true
		}() &&
		func() bool {
			for _, b := range m.Map.Barriers {
				tmp := append(b.Optional, b.Required[1:]...)
				if b.Required[0].IsOnSameLine(tmp...) {
					return false
				}
			}
			return true
		}()
}
func (m *Map2D) toGrid() *grid.Grid {
	if !m.isValid() {
		return nil
	}
	gridXSize, gridYSize := int(m.Map.Size.X/perGridSize), int(m.Map.Size.Y/perGridSize)
	g := grid.NewGrid(gridXSize, gridYSize)
	var barriers [][][]float64
	for _, b := range m.Map.Barriers {
		barrier := grid2DsToFloat2DArray(pointsToGrid2DArray(append(b.Optional, b.Required[:]...)...)...)
		barriers = append(barriers, barrier)
	}
	for x := 0; x < gridXSize; x++ {
		for y := 0; y < gridYSize; y++ {
			for _, b := range barriers {
				if InPolygon([]float64{float64(x), float64(y)}, b) {
					g.SetBarrier(x, y, true)
				}
			}
		}
	}
	return g
}
