package service

import (
	"allsum_oa/model"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/jinzhu/gorm"
	"strconv"
	"strings"
)

//表单模板相关
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

//审批单模板相关
func GetApprocvaltplList(prefix string, params ...string) (atpls []*model.Approvaltpl, e error) {
	db := model.NewOrm().Table(prefix + "." + model.Approvaltpl{}.TableName()).Order("no desc")
	atpls = []*model.Approvaltpl{}
	if len(params) != 0 {
		e = db.Where("name like ?", "%"+params[0]+"%").Find(&atpls).Error
	} else {
		e = db.Find(&atpls).Error
	}
	return
}

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

func GetMatchGroupsOfRole(prefix string, rid int) (groups []*model.Group, e error) {
	sql := fmt.Sprintf(
		`SELECT *
		FROM "%s"."group"
		WHERE id IN (SELECT t2.group_id
					 FROM "%s".user_role AS t1
					   INNER JOIN "%s".user_group AS t2
						 ON t1.user_id = t2.user_id
					 WHERE t1.role_id = ?);`, prefix, prefix, prefix)
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

func AddApproval(prefix string, a *model.Approval, atplNo string) (e error) {
	db := model.NewOrm()
	user := &model.User{}
	group := &model.Group{}
	role := &model.Role{}
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
	//获取并写入每一步流程
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
	//写入审批单
	a.CurrentFlow = aflows[0].Id
	e = tx.Table(prefix + "." + a.TableName()).Create(a).Error
	if e != nil {
		tx.Rollback()
		return
	}
	//写入表单内容
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
	//根据模板和发起人信息找出需要进行的审批流程
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
	//倒序遍历
	for i := len(realFlows) - 1; i >= 0; i-- {
		v := realFlows[i]
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
	//在指定组织寻找角色
	db := model.NewOrm()
	users = []*model.User{}
	if atplFlow.GroupId != 0 {
		sql := fmt.Sprintf(
			`SELECT *
			FROM "%s".allsum_user
			WHERE id IN (SELECT t1.user_id
						 FROM "%s".user_group AS t1
						   INNER JOIN "%s".user_role AS t2
							 ON t1.user_id = t2.user_id
						 WHERE t2.role_id = ? AND t1.group_id = ?);`, prefix, prefix, prefix)
		e = db.Raw(sql, atplFlow.RoleId, atplFlow.GroupId).Scan(&users).Error
		return
	}
	//在发起人所在组织树路径上寻找符合条件的角色,先找上级,再找下级
	sql := fmt.Sprintf(
		`SELECT *
		FROM "%s".allsum_user
		WHERE id IN
			  (SELECT t1.user_id
			   FROM "%s".user_group AS t1
				 INNER JOIN "%s".user_role AS t2
				   ON t1.user_id = t2.user_id
			   WHERE t2.role_id =? AND t1.group_id IN (?) )`, prefix, prefix, prefix)

	me := &model.Group{}
	e = db.Table(prefix+"."+me.TableName()).First(me, "id=?", a.GroupId).Error
	if e != nil {
		return
	}
	//找上级,找到就返回
	fathers := []int{}
	for _, v := range strings.Split(me.Path, "-") {
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
	e = db.Table(prefix+"."+me.TableName()).Where("path like ?", me.Path+"-%").Pluck("id", &children).Error
	if e != nil {
		return
	}
	e = db.Raw(sql, atplFlow.RoleId, children).Scan(&users).Error
	return
}

func Approve(prefix string, a *model.Approval, af *model.ApproveFlow) (e error) {
	db := model.NewOrm()
	tx := db.Begin()
	//审批，修改一步审批流程状态
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
	} else {
		nextFlow := &model.ApproveFlow{}
		e = db.Table(prefix+"."+nextFlow.TableName()).Order("id").Limit(1).
			Where("approval_no=? and id>?", a.No, a.CurrentFlow).Find(nextFlow).Error
		if e == gorm.ErrRecordNotFound {
			//流程走完,整个审批单通过
			a.Status = model.ApprovalStatAccessed
			e = tx.Table(prefix + "." + a.TableName()).Model(a).Updates(a).Error
			if e != nil {
				tx.Rollback()
				return
			}
			go newMsgToCreator(prefix, a)
		} else if e == nil {
			//继续下一步
			a.CurrentFlow = nextFlow.Id
			e = tx.Table(prefix + "." + a.TableName()).Model(a).Updates(a).Error
			if e != nil {
				tx.Rollback()
				return
			}
			go newMsgToApprovers(prefix, nextFlow.MatchUsers, a)
		} else {
			tx.Rollback()
			return errors.New("审批失败")
		}
	}
	return tx.Commit().Error
}

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

func GetTodoApprovalsToMe(prefix string, uid int, params ...string) (alist []*model.Approval, e error) {
	db := model.NewOrm()
	alist = []*model.Approval{}
	sql := fmt.Sprintf(`select * from "%s".approval as t1 inner join "%s".approve_flow as t2
		on t1.current_flow = t2.id
		where t1.status=%d and t2.status=%d and t2.match_users like '%%-%d-%%' `,
		prefix, prefix, model.ApprovalStatWaiting, model.ApprovalStatWaiting, uid)

	if len(params) != 0 && len(params[0]) != 0 {
		sql += fmt.Sprintf(`and t2.ctime>='%s' `, params[0])
	}
	sql += `order by t2.ctime desc`
	e = db.Raw(sql).Scan(&alist).Error
	return
}

func GetFinishedApprovalsToMe(prefix string, uid int, params ...string) (alist []*model.Approval, e error) {
	db := model.NewOrm()
	alist = []*model.Approval{}
	//approve_flow 中user_id有值表示已经审批过
	sql := fmt.Sprintf(`select * from "%s".approval as t1 inner join "%s".approve_flow as t2
		on t1.no = t2.approval_no where t2.user_id=%d `, prefix, prefix, uid)
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
