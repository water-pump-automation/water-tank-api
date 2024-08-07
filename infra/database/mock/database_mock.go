package database_mock

import (
	"time"
	"water-tank-api/app/core/entity/access"
	stack "water-tank-api/app/core/entity/error_stack"
	"water-tank-api/app/core/entity/water_tank"
)

type WaterTankMockData struct {
	states map[string]map[string]*water_tank.WaterTank
}

var MockTimeNow = time.Now()

func NewWaterTankMockData() *WaterTankMockData {
	return &WaterTankMockData{
		states: map[string]map[string]*water_tank.WaterTank{
			"GROUP_1": {
				"TANK_1": {
					Name:              "TANK_1",
					Group:             "GROUP_1",
					MaximumCapacity:   100,
					TankState:         water_tank.Empty,
					CurrentWaterLevel: 0,
					Access:            "a",
					LastFullTime:      MockTimeNow,
				},
				"TANK_2": {
					Name:              "TANK_2",
					Group:             "GROUP_1",
					MaximumCapacity:   80,
					TankState:         water_tank.Filling,
					CurrentWaterLevel: 50,
					Access:            "b",
					LastFullTime:      MockTimeNow,
				},
				"TANK_3": {
					Name:              "TANK_3",
					Group:             "GROUP_1",
					MaximumCapacity:   120,
					TankState:         water_tank.Full,
					CurrentWaterLevel: 120,
					Access:            "c",
					LastFullTime:      MockTimeNow,
				},
			},
			"GROUP_2": {
				"TANK_1": {
					Name:              "TANK_1",
					Group:             "GROUP_2",
					MaximumCapacity:   100,
					TankState:         water_tank.Empty,
					CurrentWaterLevel: 0,
					Access:            "d",
					LastFullTime:      MockTimeNow,
				},
				"TANK_2": {
					Name:              "TANK_2",
					Group:             "GROUP_2",
					MaximumCapacity:   80,
					TankState:         water_tank.Full,
					CurrentWaterLevel: 80,
					Access:            "e",
					LastFullTime:      MockTimeNow,
				},
			},
			"GROUP_3": {
				"TANK_1": {
					Name:              "TANK_1",
					Group:             "GROUP_3",
					MaximumCapacity:   120,
					TankState:         water_tank.Filling,
					CurrentWaterLevel: 90,
					Access:            "f",
					LastFullTime:      MockTimeNow,
				},
			},
		},
	}
}

func (tank *WaterTankMockData) GetWaterTankState(group string, names ...string) (state *water_tank.WaterTank, err stack.ErrorStack) {
	state = tank.states[group][names[0]]
	return
}

func (tank *WaterTankMockData) GetTankGroupState(groups ...string) (state []*water_tank.WaterTank, err stack.ErrorStack) {
	if group, exists := tank.states[groups[0]]; exists {
		for _, tank := range group {
			state = append(state, tank)
		}
	}
	return
}

func (tank *WaterTankMockData) CreateWaterTank(name string, group string, accessToken access.AccessToken, capacity water_tank.Capacity) (err stack.ErrorStack) {
	if _, exists := tank.states[group]; !exists {
		tank.states[group] = map[string]*water_tank.WaterTank{
			name: {
				Name:              name,
				Group:             group,
				MaximumCapacity:   capacity,
				TankState:         water_tank.Empty,
				CurrentWaterLevel: 0,
				Access:            accessToken,
			},
		}
		return
	}

	tank.states[group][name] = &water_tank.WaterTank{
		Name:              name,
		Group:             group,
		MaximumCapacity:   capacity,
		TankState:         water_tank.Empty,
		CurrentWaterLevel: 0,
		Access:            accessToken,
	}

	return
}

func (tank *WaterTankMockData) UpdateWaterTankState(name string, group string, waterLevel water_tank.Capacity, levelState water_tank.State) (state *water_tank.WaterTank, err stack.ErrorStack) {
	if group, exists := tank.states[group]; exists {
		group[name].CurrentWaterLevel = waterLevel
		group[name].TankState = levelState
	}
	return
}
