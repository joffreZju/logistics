package service

import (
	"allsum_account/model"
)

//隐藏文档，启用文档
func SetDocStatus(st, tp, did int) (err error) {
	if st == model.DocUsing {
		err = model.SetDocHide(tp)
		if err != nil {
			return err
		}
	}

	d := &model.Document{
		Status: st,
		Id:     did,
	}
	err = model.UpdateDocument(d)
	//err = model.UpdateDocument(d, "Status")
	return
}

func NewDocument(d *model.Document) (err error) {
	model.SetDocHide(d.DocType)
	err = model.CreateDocument(d)
	return
}

func NewFile(f *model.File) (err error) {
	err = model.CreateFile(f)
	return
}
