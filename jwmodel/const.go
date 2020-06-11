package jwmodel

const (
	// 主机
	JwHost = "http://jwgl.ouc.edu.cn"
	// 数据页
	DataTable = JwHost + "/taglib/DataTable.jsp"
	// 登录
	Login = JwHost + "/cas/logon.action"
	// 验证码
	Captcha = JwHost + "/cas/genValidateCode"
	// 检查会话有效性
	OnlineMessage = JwHost + "/online/message"
	// 获取所有学科的专业的列表
	SpecialityList = JwHost + "/frame/droplist/getDropLists.action"
	// 获取指定年级的专业信息
	GradeSpeciality = JwHost + "/jw/common/getStuGradeSpeciatyInfo.action"
	// 获取二维形式的课表
	HtmlTimeTable = JwHost + "/student/wsxk.zx.bykb.jsp"
	// 获取学生选课币使用情况
	Coin = JwHost + "/jw/common/getSelectLessonPointsInfo.action"
	// 检测是否可以重修某一门课
	RepeatCourse = JwHost + "/jw/common/isRepeatCourse.action"
	// 选课
	SaveCourse = JwHost + "/jw/common/saveElectiveCourse.action"
	// 选课前检验是否具有资格的接口
	ConfirmSelectable = JwHost + "/jw/common/isSelectableSkbjdm.action"
	// 删除已选课程
	CancelCourse = JwHost + "/jw/common/cancelElectiveCourse.action"
	// 加密参数接口
	KeyTimeParam = JwHost + "/custom/js/SetKingoEncypt.jsp"
	// 获取成绩接口
	Score = JwHost + "/student/xscj.stuckcj_data.jsp"
)
