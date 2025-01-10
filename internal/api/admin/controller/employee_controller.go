package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
	"strconv"
	"takeout/common"
	"takeout/common/e"
	"takeout/common/enum"
	"takeout/global"
	"takeout/internal/api/admin/request"
	"takeout/internal/model"
	"takeout/internal/service"
)

type EmployeeController struct {
	service service.IEmployeeService
}

func NewEmployeeController(employeeService service.IEmployeeService) *EmployeeController {
	return &EmployeeController{service: employeeService}
}

// @Login 员工登录接口
// @Tags Employee
// @Produce json
// @Param data body request.EmployeeLogin true "员工登录信息"
// @Success 200 {object} common.Result{data=response.EmployeeLogin} "success"
// @Failure 400 {object} common.Result "Invalid request payload"
// @Failure 401 {object} common.Result "wrong password or error on json web token generate"
// @Router /admin/employee/login [post]
func (ec *EmployeeController) Login(ctx *gin.Context) {
	code := e.SUCCESS
	employeeLogin := request.EmployeeLogin{}
	err := ctx.ShouldBindWith(&employeeLogin, binding.JSON) // 确保请求体是JSON格式
	if err != nil {
		code = e.ERROR
		global.Log.Debug("Invalid request payload:", err.Error())
		ctx.JSON(http.StatusBadRequest, common.Result{
			Code: code,
			Msg:  "Invalid request payload",
		})
		return
	}

	resp, err := ec.service.Login(ctx, employeeLogin)
	if err != nil {
		code = e.ERROR
		global.Log.Warn("EmployeeController login Error:", err.Error())
		ctx.JSON(http.StatusUnauthorized, common.Result{
			Code: code,
			Msg:  "wrong password or error on json web token generate",
		})
		return
	}

	ctx.JSON(http.StatusOK, common.Result{
		Code: code,
		Data: resp,
	})
}

// @Logout 员工退出接口
// @Security JWTAuth
// @Tags Employee
// @Produce json
// @Success 200 {object} common.Result{} "success"
// @Router /admin/employee/logout [post]
func (ec *EmployeeController) Logout(ctx *gin.Context) {
	code := e.SUCCESS
	var err error
	err = ec.service.Logout(ctx)
	if err != nil {
		code = e.ERROR
		global.Log.Warn("EmployeeController logout Error:", err.Error())
		ctx.JSON(http.StatusBadRequest, common.Result{
			Code: code,
			Msg:  err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, common.Result{
		Code: code,
	})
}

// @AddEmployee 注册员工接口
// @Security JWTAuth
// @Tags Employee
// @Produce json
// @Param data body request.EmployeeDTO true "新增员工信息"
// @Success 200 {object} common.Result{} "success"
// @Failure 500 {object} common.Result "Dupliciated Username"
// @Router /admin/employee/ [post]
func (ec *EmployeeController) AddEmployee(ctx *gin.Context) {
	var (
		code     = e.SUCCESS
		err      error
		employee request.EmployeeDTO
	)
	err = ctx.ShouldBindWith(&employee, binding.JSON)
	if err != nil {
		global.Log.Debug("AddEmployee Error:", err.Error())
		return
	}

	err = ec.service.CreateEmployee(ctx, employee)
	if err != nil {
		code = e.ERROR
		global.Log.Warn("AddEmployee  Error:", err.Error())
		ctx.JSON(http.StatusOK, common.Result{
			Code: code,
			Msg:  err.Error(),
		})
	}

	ctx.JSON(http.StatusOK, common.Result{
		Code: code,
	})
}

// @PageQuery 查询员工分页
// @Security JWTAuth
// @Tags Employee
// @Produce json
// @Param data query request.EmployeePageQueryDTO true "查询员工请求信息"
// @Success 200 {object} common.Result{data=common.PageResult} "success"
// @Failure 501 {object} common.Result{} "fail"
// @Router /admin/employee/page/ [get]
func (ec *EmployeeController) PageQuery(ctx *gin.Context) {
	var (
		code                 = e.SUCCESS
		err                  error
		employeePageQueryDTO request.EmployeePageQueryDTO
		pageResult           *common.PageResult
	)
	err = ctx.BindQuery(&employeePageQueryDTO)
	if err != nil {
		code = e.ERROR
		global.Log.Error("AddEmployee  invalid params err:", err.Error())
		ctx.JSON(http.StatusOK, common.Result{
			Code: code,
			Msg:  e.GetMsg(code),
		})
		return
	}
	// 分页查询
	pageResult, err = ec.service.PageQuery(ctx, employeePageQueryDTO)
	if err != nil {
		code = e.ERROR
		global.Log.Warn("PageQuery  Error:", err.Error())
		ctx.JSON(http.StatusNotImplemented, common.Result{
			Code: code,
			Msg:  err.Error(),
		})
	}
	ctx.JSON(http.StatusOK, common.Result{
		Code: code,
		Data: pageResult,
	})
}

// @OnOrOff 禁用员工账号
// @Security JWTAuth
// @Tags Employee
// @Produce json
// @Param status path string true "员工状态"
// @Param id query string true "查询员工请求信息"
// @Success 200 {object} common.Result{} "success"
// @Failure 501 {object} common.Result{} "fail"
// @Router /admin/employee/status/{status} [post]
func (ec *EmployeeController) OnOrOff(ctx *gin.Context) {
	var (
		code   = e.SUCCESS
		id     uint64
		err    error
		status int
	)
	id, err = strconv.ParseUint(ctx.Query("id"), 10, 64)
	if err != nil {
		code = e.ERROR
		global.Log.Error("OnOrOff invalid params err:", err.Error())
		ctx.JSON(http.StatusOK, common.Result{
			Code: code,
			Msg:  e.GetMsg(code),
		})
		return
	}
	status, _ = strconv.Atoi(ctx.Param("status"))
	err = ec.service.SetStatus(ctx, id, status)
	if err != nil {
		code = e.ERROR
		global.Log.Warn("OnOrOff Status  Error:", err.Error())
		ctx.JSON(http.StatusInternalServerError, common.Result{
			Code: code,
			Msg:  err.Error(),
		})
	}
	// 更新员工状态
	global.Log.Info("启用Or禁用员工状态：", ", id:", id, ", status:", status)
	ctx.JSON(http.StatusOK, common.Result{
		Code: code,
	})
}

// @GetById 根据id查询员工信息
// @Security JWTAuth
// @Tags Employee
// @Produce json
// @Param id path string true "员工id"
// @Success 200 {object} common.Result{} "success"
// @Failure 400 {object} common.Result{} "fail"
// @Router /admin/employee/{id} [get]
func (ec *EmployeeController) GetById(ctx *gin.Context) {
	var (
		code     = e.SUCCESS
		id       uint64
		err      error
		employee model.Employee
	)
	id, err = strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		code = e.ERROR
		global.Log.Error("GetById invalid params err:", err.Error())
		ctx.JSON(http.StatusBadRequest, common.Result{
			Code: code,
			Msg:  e.GetMsg(code),
		})
		return
	}
	employee, err = ec.service.GetById(ctx, id)
	if err != nil {
		code = e.ERROR
		global.Log.Warn("GetById Error:", err.Error())
		ctx.JSON(http.StatusInternalServerError, common.Result{
			Code: code,
			Msg:  err.Error(),
		})
	}
	// 更新员工状态
	global.Log.Info("根据id查询员工信息：", ", id:", id, ", employee:", employee)
	ctx.JSON(http.StatusOK, common.Result{
		Code: code,
		Data: employee,
	})
}

// @UpdateEmployee 更新员工信息
// @Security JWTAuth
// @Tags Employee
// @Produce json
// @Param employee body request.EmployeeDTO true "信息"
// @Success 200 {object} common.Result{} "success"
// @Failure 400 {object} common.Result{} "fail"
// @Failure 500 {object} common.Result{} "fail"
// @Router /admin/employee [put]
func (ec *EmployeeController) UpdateEmployee(ctx *gin.Context) {
	var (
		code        = e.SUCCESS
		employeeDTO request.EmployeeDTO
		err         error
	)

	err = ctx.ShouldBindWith(&employeeDTO, binding.JSON)
	if err != nil {
		code = e.ERROR
		global.Log.Error("UpdateEmployee Bind Params Error:", err.Error())
		ctx.JSON(http.StatusBadRequest, common.Result{
			Code: code,
			Msg:  e.GetMsg(code),
		})
	}
	err = ec.service.UpdateEmployee(ctx, employeeDTO)
	if err != nil {
		code = e.ERROR
		global.Log.Warn("UpdateEmployee Error:", err.Error())
		ctx.JSON(http.StatusInternalServerError, common.Result{
			Code: code,
			Msg:  e.GetMsg(code),
		})
	}
	ctx.JSON(http.StatusOK, common.Result{
		Code: code,
		Msg:  e.GetMsg(code),
	})
}

// @EditPassword 修改员工密码
// @Security JWTAuth
// @Tags Employee
// @Produce json
// @Param employee body request.EmployeeEditPassword true "id和新旧密码"
// @Success 200 {object} common.Result{} "success"
// @Failure 400 {object} common.Result{} "fail"
// @Failure 406 {object} common.Result{} "fail"
// @Router /admin/employee/editPassword [put]
func (ec *EmployeeController) EditPassword(ctx *gin.Context) {
	var (
		EditPasswordDto request.EmployeeEditPassword
		code            = e.SUCCESS
		err             error
	)
	err = ctx.ShouldBindWith(&EditPasswordDto, binding.JSON)
	if err != nil {
		code = e.ERROR
		global.Log.Error("EditPassword Bind Params Error:", err.Error())
		ctx.JSON(http.StatusBadRequest, common.Result{
			Code: code,
			Msg:  e.GetMsg(code),
		})
	}
	// 从上下文获取员工id
	if id, ok := ctx.Get(enum.CurrentId); ok {
		EditPasswordDto.EmpId = id.(uint64)
	}

	err = ec.service.EditPassword(ctx, EditPasswordDto)
	if err != nil {
		code = e.ERROR
		global.Log.Error("Wrong Password Error:", err.Error())
		ctx.JSON(http.StatusNotAcceptable, common.Result{
			Code: code,
			Msg:  e.GetMsg(code),
		})
	}

	ctx.JSON(http.StatusOK, common.Result{
		Code: code,
	})

}
