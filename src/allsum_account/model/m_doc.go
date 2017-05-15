package model

import (
	"time"
)

const (
	//doctype
	DocUsing = 1
	DocHide  = 2
)

type File struct {
	FileNo     string    `gorm:"primary_key;size:50;not null"`
	Uid        int       `gorm:"not null" json:",omitempty"`
	Name       string    `gorm:"not null;size:100" json:"name,omitempty"`
	Mime       string    `gorm:"size:250"`
	Size       int       `json:"size,omitempty"`
	Md5        string    `gorm:"not null;size:50" json:"md5,omitempty"`
	CreateTime time.Time `gorm:"default:current_timestamp" json:",omitempty"`
	Data       string    `gorm:"type:bytea" json:"-"`
}

type Document struct {
	Id       int    `gorm:"primary_key;auto_increment"`
	DocType  int    `gorm:"not null"`
	Uploader int    `gorm:"not null" json:",omitempty"`
	FileNo   string `gorm:"not null" json:",omitempty"`
	Desc     string `json:",omitempty"`
	Status   int    `json:",omitempty"`
}

//插入文件
func CreateFile(f *File) (err error) {
	err = NewOrm().Create(f).Error
	return
}

//删除文件
func DeleteFile(no string) (err error) {
	f := &File{
		FileNo: no,
	}
	err = NewOrm().Delete(f).Error
	return
}

//获取文件按文件编号
func GetFile(id string) (f *File, err error) {
	f = &File{FileNo: id}
	err = NewOrm().Find(f).Error
	return
}

//检查文件是否存在
func CheckFileExist(id string) bool {
	var c int
	NewOrm(ReadOnly).Find(&File{FileNo: id}).Count(&c)
	if c != 1 {
		return false
	}
	return true
}

//获取文件按用户/// 此处以后可能有分页
func GetFilesByUser(userid int) (fs []*File, err error) {
	err = NewOrm(ReadOnly).Table("File").Where("Uid = ?", userid).Find(&fs).Error
	return
}

//-------------------------------Document

//创建文档
func CreateDocument(doc *Document) (err error) {
	err = NewOrm().Create(doc).Error
	return
}

//更新文档
func UpdateDocument(doc *Document) (err error) {
	err = NewOrm().Model(doc).Updates(doc).Error
	return
}

//获取文档
func GetDocument(id int) (doc *Document, err error) {
	doc = &Document{Id: id}
	err = NewOrm(ReadOnly).First(doc, id).Error
	return
}

func SetDocHide(tp int) (err error) {
	err = NewOrm().Table("Document").Where("DocType = ?", tp).Update("Status", DocHide).Error

	return
}

//根据文档类型获取当前有效文档
func GetDocByType(tp int) (doc *Document, err error) {
	doc = &Document{}
	err = NewOrm().Table("Document").Where("DocType = ? and Status = ?", tp, DocUsing).First(doc).Error
	return
}

func GetDocListByType(tp int) (docs []*Document, err error) {
	err = NewOrm().Table("Document").Where("DocType = ?", tp).Find(&docs).Error
	return
}

/*
//获取文档列表
func ListDocument(limit int, mark int) (docs []*Document, err error) {
	_, err = NewOrm(ReadOnly).QueryTable("Document").Limit(limit, mark).All(&docs)
	return
}
//根据userid获取文档
func GetDocumentByUserId(userid int64) (docs []*Document, err error) {
	_, err = NewOrm(ReadOnly).QueryTable("Document").Filter("Uploader", userid).All(&docs)
	return
}

//删除文档
func DeleteDocument(doc *Document) (err error) {
	_, err = orm.NewOrm().Delete(doc)
	return
}
*/
