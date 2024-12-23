package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
	"takeout/common"
	"takeout/common/e"
	"takeout/global"
	"takeout/internal/api/request"
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
// @Param data body request.EmployeePageQueryDTO true "查询员工请求信息"
// @Success 200 {object} common.Result{data:common.PageResult} "success"
// @Failure 400 {object} common.Result{} "fail"
// @Router /admin/page/ [get]
func (ec *EmployeeController) PageQuery(ctx *gin.Context) {
	var (
		code                 = e.SUCCESS
		err                  error
		employeePageQueryDTO request.EmployeePageQueryDTO
	)
	err = ctx.ShouldBindWith(&employeePageQueryDTO, binding.JSON)
	if err != nil {
		code = e.ERROR
		global.Log.Error("AddEmployee  invalid params err:", err.Error())
		e.Send(ctx, code)
		return
	}
	//pageResult, err := ec.service.PageQuery(ctx, employeePageQueryDTO)

}

func (ec *EmployeeController) GetById(ctx *gin.Context)        {}
func (ec *EmployeeController) OnOrOff(ctx *gin.Context)        {}
func (ec *EmployeeController) EditPassword(ctx *gin.Context)   {}
func (ec *EmployeeController) UpdateEmployee(ctx *gin.Context) {}
