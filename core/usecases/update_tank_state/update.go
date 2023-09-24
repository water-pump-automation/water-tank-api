package update_tank_state

import (
	stack "water-tank-api/core/entity/error_stack"
	data "water-tank-api/core/entity/water_tank"
	get_tank "water-tank-api/core/usecases/get_tank"
)

type UpdateWaterTank struct {
	tank       data.WaterTankData
	getUsecase *get_tank.GetWaterTank
}

func NewWaterTankUpdate(tank data.WaterTankData) *UpdateWaterTank {
	return &UpdateWaterTank{
		tank:       tank,
		getUsecase: get_tank.NewGetWaterTank(tank),
	}
}

func (conn *UpdateWaterTank) Update(tank string, currentLevel data.Capacity) (err stack.ErrorStack) {
	var tankState *get_tank.WaterTankState

	tankState, err = conn.getUsecase.Get(tank)

	if err.HasError() {
		if entity := err.EntityError(); entity != nil {
			err.Append(WaterTankErrorServerError(entity.Error()))
			return
		}

		err.Append(WaterTankErrorNotFound(err.UsecaseError().Error()))
		return
	}

	if currentLevel > tankState.Tank.MaximumCapacity {
		err.Append(WaterTankCurrentWaterLevelBiggerThanMax)
		return
	}

	if currentLevel < 0 {
		err.Append(WaterTankCurrentWaterLevelSmallerThanZero)
		return
	}

	_, updateErr := conn.tank.UpdateWaterTankState(tank, currentLevel)

	if updateErr.HasError() {
		err.Append(WaterTankErrorServerError(updateErr.EntityError().Error()))
	}

	return
}
