package service

import (
	"allsum_oa/model"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/jinzhu/gorm"
	"gopkg.in/gomail.v2"
	"strconv"
	"strings"
)

//表单模板相关，CURD
func GetFormtplList(prefix string) (ftpls []*model.Formtpl, e error) {
	ftpls = []*model.Formtpl{}
	e = model.NewOrm().Table(prefix + "." + model.Formtpl{}.TableName()).Order("no desc").Find(&ftpls).Error
	return
}

func AddFormtpl(prefix string, ftpl *model.Formtpl) (e error) {
	e = model.NewOrm().Table(prefix + "." + ftpl.TableName()).Create(ftpl).Error
	return
}

func UpdateFormtpl(prefix string, ftpl *model.Formtpl) (e error) {
	tx := model.NewOrm().Begin()
	c := tx.Table(prefix + "." + ftpl.TableName()).
		Model(ftpl).Updates(ftpl).RowsAffected
	if c != 1 {
		tx.Rollback()
		e = errors.New("wrong formtpl no")
		return
	}
	return tx.Commit().Error
}

func ControlFormtpl(prefix, no string, status int) (e error) {
	count := 0
	if status == model.TplDisabled {
		//如果有审批单模板用到了这个表单模板，那么不能禁用
		e = model.NewOrm().Table(prefix+"."+model.Approvaltpl{}.TableName()).
			Where("formtpl_no=?", no).Count(&count).Error
		if e != nil {
			return
		} else if count != 0 {
			return errors.New("some approvaltpl are using this formtpl")
		}
	}
	tx := model.NewOrm().Begin()
	c := tx.Table(prefix+"."+model.Formtpl{}.TableName()).
		Model(&model.Formtpl{No: no}).Update("status", status).RowsAffected
	if c != 1 {
		tx.Rollback()
		e = errors.New("wrong formtpl no")
		return
	}
	return tx.Commit().Error
}

func DelFormtpl(prefix, no string) (e error) {
	//如果有审批单模板用到了这个表单模板，那么不能删除
	count := 0
	e = model.NewOrm().Table(prefix+"."+model.Approvaltpl{}.TableName()).
		Where("formtpl_no=?", no).Count(&count).Error
	if e != nil {
		return
	} else if count != 0 {
		return errors.New("some approvaltpl are using this formtpl")
	}
	tx := model.NewOrm().Begin()
	c := tx.Table(prefix + "." + model.Formtpl{}.TableName()).
		Delete(&model.Formtpl{No: no}).RowsAffected
	if c != 1 {
		tx.Rollback()
		e = errors.New("wrong formtpl no")
		return
	}
	return tx.Commit().Error
}

//审批单模板相关，可变参数params用户根据名字搜索审批单模板
func GetApprovaltplList(prefix string, params ...string) (atpls []*model.Approvaltpl, e error) {
	db := model.NewOrm().Table(prefix + "." + model.Approvaltpl{}.TableName()).Order("no desc")
	atpls = []*model.Approvaltpl{}
	if len(params) != 0 {
		e = db.Where("name like ?", "%"+params[0]+"%").Find(&atpls).Error
	} else {
		e = db.Find(&atpls).Error
	}
	return
}

//获取审批单模板的详细信息，包括表单内容，以及设定的每一步流程
func GetApprovaltplDetail(prefix, atplno string) (atpl *model.Approvaltpl, e error) {
	db := model.NewOrm()
	atpl = &model.Approvaltpl{}
	e = db.Table(prefix+"."+model.Approvaltpl{}.TableName()).
		First(atpl, "no=?", atplno).Error
	if e != nil {
		return
	}
	atpl.FormtplContent = new(model.Formtpl)
	e = db.Table(prefix+"."+model.Formtpl{}.TableName()).
		First(atpl.FormtplContent, "no=?", atpl.FormtplNo).Error
	if e != nil {
		return
	}
	e = db.Table(prefix+"."+model.ApprovaltplFlow{}.TableName()).Order("id").
		Find(&atpl.FlowContent, "approvaltpl_no=?", atpl.No).Error
	return
}

//获取角色可以匹配到的组织（存在用户既是这个角色也在这个组织）
func GetMatchGroupsOfRole(prefix string, rid int) (groups []*model.Group, e error) {
	sql := fmt.Sprintf(
		`SELECT *
		FROM "%s"."%s"
		WHERE id IN (SELECT t2.group_id
					 FROM "%s"."%s" AS t1
					   INNER JOIN "%s"."%s" AS t2
						 ON t1.user_id = t2.user_id
					 WHERE t1.role_id = ?);`,
		prefix, model.Group{}.TableName(), prefix, model.UserRole{}.TableName(), prefix, model.UserGroup{}.TableName())
	groups = []*model.Group{}
	e = model.NewOrm().Raw(sql, rid).Scan(&groups).Error
	return
}

func AddApprovaltpl(prefix string, atpl *model.Approvaltpl) (e error) {
	tx := model.NewOrm().Begin()
	e = tx.Table(prefix + "." + atpl.TableName()).Create(atpl).Error
	if e != nil {
		tx.Rollback()
		return
	}
	for _, v := range atpl.FlowContent {
		e = tx.Table(prefix + "." + v.TableName()).Create(v).Error
		if e != nil {
			tx.Rollback()
			return
		}
	}
	return tx.Commit().Error
}

func UpdateApprovaltpl(prefix string, atpl *model.Approvaltpl) (e error) {
	tx := model.NewOrm().Begin()
	c := tx.Table(prefix + "." + atpl.TableName()).
		Model(atpl).Updates(atpl).RowsAffected
	if c != 1 {
		tx.Rollback()
		e = errors.New("wrong approvaltpl no")
		return
	}
	e = tx.Table(prefix+"."+model.ApprovaltplFlow{}.TableName()).
		Delete(&model.ApprovaltplFlow{}, "approvaltpl_no=?", atpl.No).Error
	if e != nil {
		tx.Rollback()
		return
	}
	for _, v := range atpl.FlowContent {
		e = tx.Table(prefix + "." + v.TableName()).Create(v).Error
		if e != nil {
			tx.Rollback()
			return
		}
	}
	return tx.Commit().Error
}

func ControlApprovaltpl(prefix, no string, status int) (e error) {
	tx := model.NewOrm().Begin()
	c := tx.Table(prefix+"."+model.Approvaltpl{}.TableName()).
		Model(&model.Approvaltpl{No: no}).Update("status", status).RowsAffected
	if c != 1 {
		tx.Rollback()
		e = errors.New("wrong approvaltpl no")
		return
	}
	return tx.Commit().Error
}

func DelApprovaltpl(prefix, no string) (e error) {
	tx := model.NewOrm().Begin()
	c := tx.Table(prefix + "." + model.Approvaltpl{}.TableName()).
		Delete(&model.Approvaltpl{No: no}).RowsAffected
	if c != 1 {
		tx.Rollback()
		e = errors.New("wrong approvaltpl no")
		return
	}
	e = tx.Table(prefix+"."+model.ApprovaltplFlow{}.TableName()).
		Delete(&model.ApprovaltplFlow{}, "approvaltpl_no=?", no).Error
	if e != nil {
		tx.Rollback()
		return
	}
	return tx.Commit().Error
}

//审批流相关
//发起人撤销审批流
func CancelApproval(prefix, no string) (e error) {
	db := model.NewOrm().Table(prefix + "." + model.Approval{}.TableName())
	a := &model.Approval{}
	e = db.First(a, "no=?", no).Error
	if e != nil {
		return
	}
	if a.Status != model.ApprovalStatWaiting {
		e = errors.New("审批单已经完成")
		return
	}
	tx := db.Begin()
	c := tx.Model(a).Update("status", model.ApprovalStatCanceled).RowsAffected
	if c != 1 {
		tx.Rollback()
		e = errors.New("审批单编号错误")
		return
	}
	return tx.Commit().Error
}

//发起一个审批流
func AddApproval(prefix string, a *model.Approval, atplNo string) (e error) {
	db := model.NewOrm()
	atpl := &model.Approvaltpl{No: atplNo}
	user := &model.User{}
	group := &model.Group{}
	role := &model.Role{}
	e = db.Table(prefix+"."+atpl.TableName()).First(atpl, atpl).Error
	if e != nil {
		return
	}
	e = db.Table(prefix+"."+user.TableName()).First(user, a.UserId).Error
	if e != nil {
		return
	}
	e = db.Table(prefix+"."+group.TableName()).First(group, a.GroupId).Error
	if e != nil {
		return
	}
	e = db.Table(prefix+"."+role.TableName()).First(role, a.RoleId).Error
	if e != nil {
		return
	}
	a.UserName = user.UserName
	a.GroupName = group.Name
	a.RoleName = role.Name
	//根据发起用户的信息以及这个审批单模板设定的规则，获取并写入每一步应该走的流程
	var aflows []*model.ApproveFlow
	aflows, e = getApprovalFlows(prefix, a, atplNo)
	if e != nil {
		return
	}
	tx := model.NewOrm().Begin()
	for _, v := range aflows {
		e = tx.Table(prefix + "." + v.TableName()).Create(v).Error
		if e != nil {
			tx.Rollback()
			return
		}
	}
	//写入审批单的基本信息
	a.CurrentFlow = aflows[0].Id
	a.EmailMsg = atpl.EmailMsg
	e = tx.Table(prefix + "." + a.TableName()).Create(a).Error
	if e != nil {
		tx.Rollback()
		return
	}
	//写入审批单对应的表单内容
	e = tx.Table(prefix + "." + a.FormContent.TableName()).Create(a.FormContent).Error
	if e != nil {
		tx.Rollback()
		return
	}
	go newMsgToApprovers(prefix, aflows[0].MatchUsers, a)
	return tx.Commit().Error
}

//根据模板以及发起人信息，获取所有要走的flows
func getApprovalFlows(prefix string, a *model.Approval, atplNo string) (aflows []*model.ApproveFlow, e error) {
	aflows = []*model.ApproveFlow{}
	db := model.NewOrm()
	atplFlows := []*model.ApprovaltplFlow{}
	e = db.Table(prefix + "." + model.ApprovaltplFlow{}.TableName()).Order("id").
		Where(&model.ApprovaltplFlow{ApprovaltplNo: atplNo}).Find(&atplFlows).Error
	if e != nil {
		return
	}
	//根据模板和发起人信息找出需要进行的审批流程realFlows
	realFlows := []*model.ApprovaltplFlow{}
	myLocation := -1
	//倒序遍历
	for i := len(atplFlows) - 1; i >= 0; i-- {
		f := atplFlows[i]
		if f.RoleId == a.RoleId && (f.GroupId == 0 || f.GroupId == a.GroupId) {
			//发起人所在位置
			myLocation = i
			realFlows = append(realFlows, f)
		} else if i > myLocation {
			//发起人后面的位置
			realFlows = append(realFlows, f)
		} else if i < myLocation && f.Necessary == model.FlowNecessaryYes {
			//发起人前面的位置，但是必审
			realFlows = append(realFlows, f)
		}
	}
	if len(realFlows) == 0 {
		e = errors.New("发起失败：没有符合条件的审批人!")
		return
	}
	//倒序遍历realFlows
	for i := len(realFlows) - 1; i >= 0; i-- {
		v := realFlows[i]
		//获取当前一步的符合条件的审批人
		var users []*model.User
		users, e = getMatchUsersOfFlow(prefix, a, v)
		if e != nil || len(users) == 0 {
			beego.Error(e)
			e = errors.New("发起失败：没有符合条件的审批人!")
			return
		}
		matchUsers := "-"
		for _, u := range users {
			matchUsers += fmt.Sprintf("%d-", u.Id)
		}
		role := &model.Role{Id: v.RoleId}
		e = db.Table(prefix + "." + role.TableName()).Where(role).Find(role).Error
		if e != nil {
			return
		}
		//构造一步流程准备返回
		af := &model.ApproveFlow{
			ApprovalNo: a.No,
			MatchUsers: matchUsers,
			RoleId:     role.Id,
			RoleName:   role.Name,
			Status:     model.ApprovalStatWaiting,
		}
		aflows = append(aflows, af)
	}
	return aflows, nil
}

//获取一步流程对应的用户
func getMatchUsersOfFlow(prefix string, a *model.Approval, atplFlow *model.ApprovaltplFlow) (users []*model.User, e error) {
	//通过确定组织和确定角色寻找符合条件的审批人
	db := model.NewOrm()
	users = []*model.User{}
	if atplFlow.GroupId != 0 {
		sql := fmt.Sprintf(
			`SELECT *
			FROM "%s"."%s"
			WHERE id IN (SELECT t1.user_id
						 FROM "%s"."%s" AS t1
						   INNER JOIN "%s"."%s" AS t2
							 ON t1.user_id = t2.user_id
						 WHERE t2.role_id = ? AND t1.group_id = ?);`,
			prefix, model.User{}.TableName(), prefix, model.UserGroup{}.TableName(), prefix, model.UserRole{}.TableName())
		e = db.Raw(sql, atplFlow.RoleId, atplFlow.GroupId).Scan(&users).Error
		return
	}
	//在发起人所在组织树路径上寻找符合条件的角色,先找上级,再找下级
	sql := fmt.Sprintf(
		`SELECT *
		FROM "%s"."%s"
		WHERE id IN
			  (SELECT t1.user_id
			   FROM "%s"."%s" AS t1
				 INNER JOIN "%s"."%s" AS t2
				   ON t1.user_id = t2.user_id
			   WHERE t2.role_id =? AND t1.group_id IN (?) )`,
		prefix, model.User{}.TableName(), prefix, model.UserGroup{}.TableName(), prefix, model.UserRole{}.TableName())

	//找到发起人所在组织
	myGroup := &model.Group{}
	e = db.Table(prefix+"."+myGroup.TableName()).First(myGroup, "id=?", a.GroupId).Error
	if e != nil {
		return
	}
	//找上级,找到就返回
	fathers := []int{}
	for _, v := range strings.Split(myGroup.Path, "-") {
		pid, e := strconv.Atoi(v)
		if e != nil {
			return nil, e
		}
		fathers = append(fathers, pid)
	}
	e = db.Raw(sql, atplFlow.RoleId, fathers).Scan(&users).Error
	if len(users) != 0 {
		return
	}
	//找下级
	children := []int{}
	e = db.Table(prefix+"."+myGroup.TableName()).Where("path like ?", myGroup.Path+"-%").Pluck("id", &children).Error
	if e != nil {
		return
	}
	e = db.Raw(sql, atplFlow.RoleId, children).Scan(&users).Error
	return
}

//审批操作
func Approve(prefix string, a *model.Approval, af *model.ApproveFlow) (e error) {
	db := model.NewOrm()
	tx := db.Begin()
	//审批，修改当前一步审批流程状态
	count := tx.Table(prefix+"."+af.TableName()).
		Where("id=? and status=?", af.Id, model.ApprovalStatWaiting).Updates(af).RowsAffected
	if count != 1 {
		tx.Rollback()
		return errors.New("审批失败")
	}
	if af.Status == model.ApprovalStatNotAccessed {
		//整个审批单未通过
		a.Status = model.ApprovalStatNotAccessed
		count = tx.Table(prefix + "." + a.TableName()).Model(a).Updates(a).RowsAffected
		if count != 1 {
			tx.Rollback()
			return errors.New("审批失败")
		}
		go newMsgToCreator(prefix, a)
		go sendEmailToCreator(prefix, a, af, nil)
	} else {
		//当前流程通过,查找下一步流程
		nextFlow := &model.ApproveFlow{}
		e = db.Table(prefix+"."+nextFlow.TableName()).Order("id").Limit(1).
			Where("approval_no=? and id>?", a.No, a.CurrentFlow).Find(nextFlow).Error
		if e == gorm.ErrRecordNotFound {
			//没有下一步,整个审批单通过
			a.Status = model.ApprovalStatAccessed
			e = tx.Table(prefix + "." + a.TableName()).Model(a).Updates(a).Error
			if e != nil {
				tx.Rollback()
				return
			}
			go newMsgToCreator(prefix, a)
			go sendEmailToCreator(prefix, a, af, nil)
		} else if e == nil {
			//继续下一步
			a.CurrentFlow = nextFlow.Id
			e = tx.Table(prefix + "." + a.TableName()).Model(a).Updates(a).Error
			if e != nil {
				tx.Rollback()
				return
			}
			go newMsgToApprovers(prefix, nextFlow.MatchUsers, a)
			go sendEmailToCreator(prefix, a, af, nextFlow)
		} else {
			tx.Rollback()
			return errors.New("审批失败")
		}
	}
	return tx.Commit().Error
}

//给发起人发送邮件通知
func sendEmailToCreator(prefix string, approval *model.Approval, currentFlow, nextFlow *model.ApproveFlow) {
	db := model.NewOrm()
	u := &model.User{Id: approval.UserId}
	e := db.Table(prefix+"."+u.TableName()).Find(u, u).Error
	if e != nil {
		beego.Error(e)
		return
	}
	if approval.EmailMsg == model.EmailMsgNo || len(u.Mail) == 0 {
		return
	}
	subject := fmt.Sprintf("OA系统通知:%s", approval.Name)

	body := fmt.Sprintf(`你的<b>%s</b>有最新状态,请到<a href="http://oa.allsum.cn">壹算科技OA</a>查看<br><br>`, approval.Name)
	body += fmt.Sprintf(`审批人：<b>%s</b><br>`, currentFlow.UserName)

	if currentFlow.Status == model.ApprovalStatAccessed {
		body += fmt.Sprint(`审批意见：<b>通过</b><br>`)
	} else {
		body += fmt.Sprint(`审批意见：<b>拒绝</b><br>`)
	}
	if len(currentFlow.Comment) != 0 {
		body += fmt.Sprintf(`审批备注：<b>%s</b>`, currentFlow.Comment)
	}
	if nextFlow == nil {
		body += fmt.Sprint(`下一步审批：审批单已结束，没有下一步<br>`)
	} else {
		body += fmt.Sprintf(`下一步审批：%s`, nextFlow.RoleName)
	}
	sendEmail([]string{u.Mail}, subject, body)
}

//读取配置发送邮件
func sendEmail(targets []string, subject, body string) {
	if len(targets) == 0 {
		return
	}
	smtpHost := beego.AppConfig.String("emailAccount::smtp")
	smtpPort, _ := beego.AppConfig.Int("emailAccount::port")
	from := beego.AppConfig.String("emailAccount::from")
	password := beego.AppConfig.String("emailAccount::password")

	m := gomail.NewMessage()
	m.SetAddressHeader("From", from, "")

	tos := []string{}
	for _, v := range targets {
		tos = append(tos, m.FormatAddress(v, ""))
	}
	m.SetHeader("To", tos...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	dial := gomail.NewPlainDialer(smtpHost, smtpPort, from, password)
	if e := dial.DialAndSend(m); e != nil {
		beego.Error("发送邮件失败:", e)
	} else {
		beego.Info("发送邮件成功:", targets)
	}
}

//给发起人App推送
func newMsgToCreator(company string, a *model.Approval) {
	msg := &model.Message{
		CompanyNo: company,
		UserId:    a.UserId,
		MsgType:   model.MsgTypeApprove,
		Content: model.JsonMap{
			"ApprovalNo": a.No,
		},
	}
	if a.Status == model.ApprovalStatAccessed {
		msg.Title = "你的" + a.Name + "审批单已通过"
	} else if a.Status == model.ApprovalStatNotAccessed {
		msg.Title = "你的" + a.Name + "审批单被拒绝"
	} else if a.Status == model.ApprovalStatStop {
		msg.Title = "你的" + a.Name + "审批单无法继续流转，请咨询管理员或客服"
	}
	e := SaveAndSendMsg(msg)
	if e != nil {
		beego.Error(e)
	}
}

//给审批人App推送
func newMsgToApprovers(company, matchUsers string, a *model.Approval) {
	title := "来自$" + a.UserName + "$的审批消息"
	users := strings.Split(strings.Trim(matchUsers, "-"), "-")
	for _, v := range users {
		uid, e := strconv.Atoi(v)
		if e != nil {
			beego.Error(e)
			continue
		}
		msg := &model.Message{
			CompanyNo: company,
			UserId:    uid,
			MsgType:   model.MsgTypeApprove,
			Title:     title,
			Content: model.JsonMap{
				"ApprovalNo": a.No,
			},
		}
		go SaveAndSendMsg(msg)
	}
}

func GetApproval(prefix, no string) (a *model.Approval, e error) {
	a = new(model.Approval)
	e = model.NewOrm().Table(prefix+"."+a.TableName()).First(a, "no=?", no).Error
	return
}

func GetApproveFlowById(prefix string, id int) (af *model.ApproveFlow, e error) {
	af = new(model.ApproveFlow)
	e = model.NewOrm().Table(prefix+"."+af.TableName()).First(af, "id=?", id).Error
	return
}

//根据开始时间和审批单状态，获取我发起的审批单
func GetApprovalsFromMe(prefix string, uid int, beginTime, condition string) (alist []*model.Approval, e error) {
	alist = []*model.Approval{}
	db := model.NewOrm().Table(prefix+"."+model.Approval{}.TableName()).
		Order("status, ctime desc").Where("user_id=?", uid)
	if len(beginTime) != 0 {
		db = db.Where("ctime>=?", beginTime)
	}
	if condition == model.GetApprovalApproving {
		db = db.Where("status=?", model.ApprovalStatWaiting)
	} else if condition == model.GetApprovalFinished {
		db = db.Where("status<>?", model.ApprovalStatWaiting)
	}
	e = db.Find(&alist).Error
	return
}

//获取需要我审批的审批单，可变参数过滤时间
func GetTodoApprovalsToMe(prefix string, uid int, params ...string) (alist []*model.Approval, e error) {
	db := model.NewOrm()
	alist = []*model.Approval{}
	sql := fmt.Sprintf(`select * from "%s"."%s" as t1 inner join "%s"."%s" as t2
		on t1.current_flow = t2.id
		where t1.status=%d and t2.status=%d and t2.match_users like '%%-%d-%%' `,
		prefix, model.Approval{}.TableName(), prefix, model.ApproveFlow{}.TableName(), model.ApprovalStatWaiting, model.ApprovalStatWaiting, uid)

	if len(params) != 0 && len(params[0]) != 0 {
		sql += fmt.Sprintf(`and t2.ctime>='%s' `, params[0])
	}
	sql += `order by t2.ctime desc`
	e = db.Raw(sql).Scan(&alist).Error
	return
}

//获取穷我审批过的审批单，可变参数过滤时间
func GetFinishedApprovalsToMe(prefix string, uid int, params ...string) (alist []*model.Approval, e error) {
	db := model.NewOrm()
	alist = []*model.Approval{}
	//approve_flow 中user_id有值表示已经审批过
	sql := fmt.Sprintf(`select * from "%s"."%s" as t1 inner join "%s"."%s" as t2
		on t1.no = t2.approval_no where t2.user_id=%d `,
		prefix, model.Approval{}.TableName(), prefix, model.ApproveFlow{}.TableName(), uid)
	if len(params) != 0 && len(params[0]) != 0 {
		sql += fmt.Sprintf(`and t2.ctime>='%s' `, params[0])
	}
	sql += `order by t2.ctime desc`
	e = db.Raw(sql).Scan(&alist).Error
	return
}

func GetApprovalDetail(prefix, no string) (a *model.Approval, e error) {
	a = new(model.Approval)
	db := model.NewOrm()
	e = db.Table(prefix+"."+a.TableName()).First(a, "no=?", no).Error
	if e != nil {
		return
	}
	a.FormContent = new(model.Form)
	e = db.Table(prefix+"."+model.Form{}.TableName()).
		First(a.FormContent, "no=?", a.FormNo).Error
	if e != nil {
		return
	}
	e = db.Table(prefix + "." + model.ApproveFlow{}.TableName()).Order("id").
		Where(&model.ApproveFlow{ApprovalNo: a.No}).Find(&a.ApproveFLows).Error
	return
}
