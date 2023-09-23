package callbacks

import (
	"water-tank-api/controllers"
	"water-tank-api/infra/web"

	iris "github.com/kataras/iris/v12"
)

func GetAll(ctx iris.Context) {
	controller := web.ExternalController()
	groupName := ctx.Params().Get("group")

	response, err := controller.Get(groupName)

	if err != nil {
		switch response.Code {
		case controllers.WaterTankNotFound:
			ctx.StatusCode(iris.StatusNotFound)
		case controllers.WaterTankInternalServerError:
			ctx.StatusCode(iris.StatusInternalServerError)
		}
		ctx.JSON(response)
		return
	}

	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(response)
}
