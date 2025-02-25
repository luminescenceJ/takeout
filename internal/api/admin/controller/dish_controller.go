package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"takeout/common"
	"takeout/common/e"
	"takeout/global"
	"takeout/internal/api/admin/request"
	"takeout/internal/api/admin/response"
	"takeout/internal/service"
)

type DishController struct {
	service service.IDishService
}

func NewDishController(service service.IDishService) *DishController {
	return &DishController{service: service}
}

// AddDish 新增菜品数据接口
// @Tags dish
// @Security JWTAuth
// @Produce json
// @Param data body request.DishDTO true "新增信息"
// @Success 200 {object} common.Result "success"
// @Failure 400 {object} common.Result "Invalid request payload"
// @Failure 500 {object} common.Result "Internal Server Failure"
// @Router /admin/dish [post]
func (c DishController) AddDish(ctx *gin.Context) {
	var (
		code    = e.SUCCESS
		err     error
		dishDTO request.DishDTO
	)
	// 记录请求参数
	global.Log.Debug("Received request to add dish", "params", dishDTO)

	// 参数绑定
	if err = ctx.Bind(&dishDTO); err != nil {
		global.Log.Warn("Failed to bind DishDTO", "error", err)
		ctx.JSON(http.StatusBadRequest, common.Result{
			Code: e.ERROR,
			Msg:  "Invalid request data",
		})
		return
	}

	// 调用服务层方法
	if err = c.service.AddDishWithFlavors(ctx, dishDTO); err != nil {
		code = e.ERROR
		global.Log.Error("Failed to add dish with flavors", "error", err.Error(), "params", dishDTO)
		ctx.JSON(http.StatusInternalServerError, common.Result{
			Code: code,
			Msg:  "Failed to add dish",
		})
		return
	}

	// 记录成功信息
	global.Log.Info("Successfully added new dish", "dish", dishDTO)

	ctx.JSON(http.StatusOK, common.Result{
		Code: code,
		Msg:  e.GetMsg(code),
	})
}

// PageQuery @PageQuery 菜品分页查询
// @Tags dish
// @Security JWTAuth
// @Produce json
// @Param data body request.DishPageQueryDTO true "查询信息"
// @Success 200 {object} common.Result{Data=common.PageResult} "success"
// @Failure 400 {object} common.Result "Invalid request payload"
// @Failure 500 {object} common.Result "Internal Server Failure"
// @Router /admin/dish/page [get]
func (c DishController) PageQuery(ctx *gin.Context) {
	var (
		code    = e.SUCCESS
		err     error
		pageDTO request.DishPageQueryDTO
		pageRES *common.PageResult
	)

	// 记录分页查询请求参数
	global.Log.Debug("Received request to query dish page", "params", pageDTO)

	// 参数绑定
	if err = ctx.ShouldBind(&pageDTO); err != nil {
		global.Log.Warn("Failed to bind DishPageQueryDTO", "error", err)
		ctx.JSON(http.StatusBadRequest, common.Result{})
		return
	}

	// 调用服务层分页查询
	if pageRES, err = c.service.PageQuery(ctx, &pageDTO); err != nil {
		code = e.ERROR
		global.Log.Error("PageQuery failed", "error", err.Error(), "params", pageDTO)
		ctx.JSON(http.StatusInternalServerError, common.Result{})
		return
	}

	// 记录成功的分页查询结果
	global.Log.Info("Successfully fetched paginated dishes", "total", pageRES.Total, "page", pageDTO.Page)

	ctx.JSON(http.StatusOK, common.Result{
		Code: code,
		Data: pageRES,
		Msg:  e.GetMsg(code),
	})
}

// GetById @GetById 根据id查询菜品信息
// @Tags dish
// @Security JWTAuth
// @Produce json
// @Param id path string true "菜品id"
// @Success 200 {object} common.Result "success"
// @Failure 500 {object} common.Result "Internal Server Failure"
// @Router /admin/dish/{id} [get]
func (c DishController) GetById(ctx *gin.Context) {
	var (
		code   = e.SUCCESS
		err    error
		id     uint64
		dishVO response.DishVo
	)
	id, _ = strconv.ParseUint(ctx.Param("id"), 10, 64)

	// 记录查询菜品的 ID 参数
	global.Log.Debug("Received request to get dish by ID", "id", id)

	// 调用服务层根据 ID 获取菜品
	if dishVO, err = c.service.GetByIdWithFlavors(ctx, id); err != nil {
		code = e.ERROR
		global.Log.Warn("GetById failed", "id", id, "error", err.Error())
		ctx.JSON(http.StatusInternalServerError, common.Result{})
		return
	}

	// 记录查询成功的结果
	global.Log.Info("Successfully fetched dish", "id", id, "dish", dishVO)

	ctx.JSON(http.StatusOK, common.Result{
		Code: code,
		Data: dishVO,
		Msg:  e.GetMsg(code),
	})
}

// List @List 根据分类id查询菜品信息
// @Tags dish
// @Security JWTAuth
// @Produce json
// @Param categoryId query string true "分类id"
// @Success 200 {object} common.Result{Data=response.DishListVo} "success"
// @Failure 400 {object} common.Result "Invalid request payload"
// @Failure 500 {object} common.Result "Internal Server Failure"
// @Router /admin/dish/list [get]
func (c DishController) List(ctx *gin.Context) {
	var (
		code       = e.SUCCESS
		err        error
		categoryId uint64
		dishList   []response.DishListVo
	)

	// 记录查询的分类 ID
	categoryId, err = strconv.ParseUint(ctx.Query("categoryId"), 10, 64)
	global.Log.Debug("Received request to list dishes by category", "categoryId", categoryId)

	// 调用服务层查询菜品列表
	if dishList, err = c.service.List(ctx, categoryId); err != nil {
		code = e.ERROR
		global.Log.Warn("List failed", "categoryId", categoryId, "error", err.Error())
		ctx.JSON(http.StatusInternalServerError, common.Result{})
		return
	}

	// 记录查询成功的结果
	global.Log.Info("Successfully fetched dish list", "categoryId", categoryId, "total", len(dishList))

	ctx.JSON(http.StatusOK, common.Result{
		Code: code,
		Data: dishList,
		Msg:  e.GetMsg(code),
	})
}

// Delete @Delete 删除菜品信息
// @Tags dish
// @Security JWTAuth
// @Produce json
// @Param ids query string true "删除id集合"
// @Success 200 {object} common.Result{} "success"
// @Failure 500 {object} common.Result "Internal Server Failure"
// @Router /admin/dish [delete]
func (c DishController) Delete(ctx *gin.Context) {
	var (
		code = e.SUCCESS
		err  error
		ids  string
	)

	// 记录删除请求的 ID 集合
	ids = ctx.Query("ids")
	global.Log.Debug("Received request to delete dishes", "ids", ids)

	// 调用服务层删除菜品
	if err = c.service.Delete(ctx, ids); err != nil {
		code = e.ERROR
		global.Log.Warn("Delete failed", "ids", ids, "error", err.Error())
		ctx.JSON(http.StatusInternalServerError, common.Result{
			Code: code,
			Msg:  e.GetMsg(code),
		})
		return
	}

	// 记录成功删除的信息
	global.Log.Info("Successfully deleted dishes", "ids", ids)

	ctx.JSON(http.StatusOK, common.Result{
		Code: code,
		Msg:  e.GetMsg(code),
	})
}

// Update @Update 修改菜品信息
// @Tags dish
// @Security JWTAuth
// @Produce json
// @Param dto body request.DishUpdateDTO true "修改菜品的信息"
// @Success 200 {object} common.Result{} "success"
// @Failure 400 {object} common.Result "Invalid request payload"
// @Failure 500 {object} common.Result "Internal Server Faliure"
// @Router /admin/dish [put]
func (c DishController) Update(ctx *gin.Context) {
	var (
		code      = e.SUCCESS
		err       error
		updateDTO request.DishUpdateDTO
	)
	if err = ctx.ShouldBindJSON(&updateDTO); err != nil {
		global.Log.Debug("param Update failed")
		ctx.JSON(http.StatusBadRequest, common.Result{})
		return
	}

	if err = c.service.Update(ctx, updateDTO); err != nil {
		code = e.ERROR
		global.Log.Warn("AddDish failed", err.Error())
		ctx.JSON(http.StatusInternalServerError, common.Result{
			Code: code,
			Msg:  e.GetMsg(code),
		})
		return
	}

	ctx.JSON(http.StatusOK, common.Result{
		Code: code,
		Msg:  e.GetMsg(code),
	})
}

// OnOrClose @OnOrClose 菜品启售或禁售
// @Tags dish
// @Security JWTAuth
// @Produce json
// @Param id query string true "id"
// @Param status path string true "status"
// @Success 200 {object} common.Result{} "success"
// @Failure 500 {object} common.Result "Internal Server Faliure"
// @Router /admin/dish/status/{status} [post]
func (c DishController) OnOrClose(ctx *gin.Context) {
	var (
		code   = e.SUCCESS
		err    error
		id     uint64
		status int
	)
	id, _ = strconv.ParseUint(ctx.Query("id"), 10, 64)
	status, _ = strconv.Atoi(ctx.Param("status"))

	if err = c.service.OnOrClose(ctx, id, status); err != nil {
		code = e.ERROR
		global.Log.Warn("OnOrClose failed: ", err.Error())
		ctx.JSON(http.StatusInternalServerError, common.Result{
			Code: code,
			Msg:  e.GetMsg(code),
		})
		return
	}

	ctx.JSON(http.StatusOK, common.Result{
		Code: code,
		Msg:  e.GetMsg(code),
	})
}

//package controller
//
//import (
//	"github.com/gin-gonic/gin"
//	"net/http"
//	"strconv"
//	"takeout/common"
//	"takeout/common/e"
//	"takeout/global"
//	"takeout/internal/api/admin/request"
//	"takeout/internal/api/admin/response"
//	"takeout/internal/service"
//)
//
//type DishController struct {
//	service service.IDishService
//}
//
//// AddDish @AddDish 新增菜品数据接口
//// @Tags dish
//// @Security JWTAuth
//// @Produce json
//// @Param data body request.DishDTO true "新增信息"
//// @Success 200 {object} common.Result "success"
//// @Failure 400 {object} common.Result "Invalid request payload"
//// @Failure 500 {object} common.Result "Internal Server Faliure"
//// @Router /admin/dish [post]
//func (c DishController) AddDish(ctx *gin.Context) {
//	var (
//		code    = e.SUCCESS
//		err     error
//		dishDTO request.DishDTO
//	)
//	if err = ctx.Bind(&dishDTO); err != nil {
//		global.Log.Debug("param DishDTO failed")
//		ctx.JSON(http.StatusBadRequest, common.Result{})
//		return
//	}
//
//	if err = c.service.AddDishWithFlavors(ctx, dishDTO); err != nil {
//		code = e.ERROR
//		global.Log.Warn("AddDish failed", err.Error())
//	}
//
//	ctx.JSON(http.StatusOK, common.Result{
//		Code: code,
//		Msg:  e.GetMsg(code),
//	})
//}
//
//// PageQuery @PageQuery 菜品分页查询
//// @Tags dish
//// @Security JWTAuth
//// @Produce json
//// @Param data body request.DishPageQueryDTO true "新增信息"
//// @Success 200 {object} common.Result{Data=common.PageResult} "success"
//// @Failure 400 {object} common.Result "Invalid request payload"
//// @Failure 500 {object} common.Result "Internal Server Faliure"
//// @Router /admin/dish/page [get]
//func (c DishController) PageQuery(ctx *gin.Context) {
//	var (
//		code    = e.SUCCESS
//		err     error
//		pageDTO request.DishPageQueryDTO
//		pageRES *common.PageResult
//	)
//
//	if err = ctx.ShouldBind(&pageDTO); err != nil {
//		global.Log.Debug("param DishPageQueryDTO failed")
//		ctx.JSON(http.StatusBadRequest, common.Result{})
//		return
//	}
//
//	if pageRES, err = c.service.PageQuery(ctx, &pageDTO); err != nil {
//		code = e.ERROR
//		global.Log.Warn("PageQuery failed", err.Error())
//		ctx.JSON(http.StatusInternalServerError, common.Result{})
//		return
//	}
//
//	ctx.JSON(http.StatusOK, common.Result{
//		Code: code,
//		Data: pageRES,
//		Msg:  e.GetMsg(code),
//	})
//}
//
//// GetById @GetById 根据id查询菜品信息
//// @Tags dish
//// @Security JWTAuth
//// @Produce json
//// @Param id path string true "菜品id"
//// @Success 200 {object} common.Result "success"
//// @Failure 500 {object} common.Result "Internal Server Faliure"
//// @Router /admin/dish/{id} [get]
//func (c DishController) GetById(ctx *gin.Context) {
//	var (
//		code   = e.SUCCESS
//		err    error
//		id     uint64
//		dishVO response.DishVo
//	)
//	id, _ = strconv.ParseUint(ctx.Param("id"), 10, 64)
//	if dishVO, err = c.service.GetByIdWithFlavors(ctx, id); err != nil {
//		code = e.ERROR
//		global.Log.Warn("GetById failed : ", err.Error())
//		ctx.JSON(http.StatusInternalServerError, common.Result{})
//		return
//	}
//
//	ctx.JSON(http.StatusOK, common.Result{
//		Code: code,
//		Data: dishVO,
//		Msg:  e.GetMsg(code),
//	})
//}
//
//// List @List 根据分类id查询菜品信息
//// @Tags dish
//// @Security JWTAuth
//// @Produce json
//// @Param categoryId query string true "分类id"
//// @Success 200 {object} common.Result{Data=response.DishListVo} "success"
//// @Failure 400 {object} common.Result "Invalid request payload"
//// @Failure 500 {object} common.Result "Internal Server Faliure"
//// @Router /admin/dish/list [get]
//func (c DishController) List(ctx *gin.Context) {
//	var (
//		code       = e.SUCCESS
//		err        error
//		categoryId uint64
//		dishList   []response.DishListVo
//	)
//
//	categoryId, err = strconv.ParseUint(ctx.Query("categoryId"), 10, 64)
//
//	if dishList, err = c.service.List(ctx, categoryId); err != nil {
//		code = e.ERROR
//		global.Log.Warn("List failed : ", err.Error())
//		ctx.JSON(http.StatusInternalServerError, common.Result{})
//		return
//	}
//
//	ctx.JSON(http.StatusOK, common.Result{
//		Code: code,
//		Data: dishList,
//		Msg:  e.GetMsg(code),
//	})
//}
//
//// Update @Update 修改菜品信息
//// @Tags dish
//// @Security JWTAuth
//// @Produce json
//// @Param dto body request.DishUpdateDTO true "修改菜品的信息"
//// @Success 200 {object} common.Result{} "success"
//// @Failure 400 {object} common.Result "Invalid request payload"
//// @Failure 500 {object} common.Result "Internal Server Faliure"
//// @Router /admin/dish [put]
//func (c DishController) Update(ctx *gin.Context) {
//	var (
//		code      = e.SUCCESS
//		err       error
//		updateDTO request.DishUpdateDTO
//	)
//	if err = ctx.ShouldBindJSON(&updateDTO); err != nil {
//		global.Log.Debug("param Update failed")
//		ctx.JSON(http.StatusBadRequest, common.Result{})
//		return
//	}
//
//	if err = c.service.Update(ctx, updateDTO); err != nil {
//		code = e.ERROR
//		global.Log.Warn("AddDish failed", err.Error())
//		ctx.JSON(http.StatusInternalServerError, common.Result{
//			Code: code,
//			Msg:  e.GetMsg(code),
//		})
//		return
//	}
//
//	ctx.JSON(http.StatusOK, common.Result{
//		Code: code,
//		Msg:  e.GetMsg(code),
//	})
//}
//
//// OnOrClose @OnOrClose 菜品启售或禁售
//// @Tags dish
//// @Security JWTAuth
//// @Produce json
//// @Param id query string true "id"
//// @Param status path string true "status"
//// @Success 200 {object} common.Result{} "success"
//// @Failure 500 {object} common.Result "Internal Server Faliure"
//// @Router /admin/dish/status/{status} [post]
//func (c DishController) OnOrClose(ctx *gin.Context) {
//	var (
//		code   = e.SUCCESS
//		err    error
//		id     uint64
//		status int
//	)
//	id, _ = strconv.ParseUint(ctx.Query("id"), 10, 64)
//	status, _ = strconv.Atoi(ctx.Param("status"))
//
//	if err = c.service.OnOrClose(ctx, id, status); err != nil {
//		code = e.ERROR
//		global.Log.Warn("OnOrClose failed: ", err.Error())
//		ctx.JSON(http.StatusInternalServerError, common.Result{
//			Code: code,
//			Msg:  e.GetMsg(code),
//		})
//		return
//	}
//
//	ctx.JSON(http.StatusOK, common.Result{
//		Code: code,
//		Msg:  e.GetMsg(code),
//	})
//}
//
//// Delete @Delete 删除菜品信息
//// @Tags dish
//// @Security JWTAuth
//// @Produce json
//// @Param ids query string true "删除id集合"
//// @Success 200 {object} common.Result{} "success"
//// @Failure 500 {object} common.Result "Internal Server Faliure"
//// @Router /admin/dish [delete]
//func (c DishController) Delete(ctx *gin.Context) {
//	var (
//		code = e.SUCCESS
//		err  error
//		ids  string
//	)
//	ids = ctx.Query("ids")
//
//	if err = c.service.Delete(ctx, ids); err != nil {
//		code = e.ERROR
//		global.Log.Warn("Delete failed: ", err.Error())
//		ctx.JSON(http.StatusInternalServerError, common.Result{
//			Code: code,
//			Msg:  e.GetMsg(code),
//		})
//		return
//	}
//
//	ctx.JSON(http.StatusOK, common.Result{
//		Code: code,
//		Msg:  e.GetMsg(code),
//	})
//}
//
//func NewDishController(service service.IDishService) *DishController {
//	return &DishController{service: service}
//}
