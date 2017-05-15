package wxisv

import (
	"fmt"
	"testing"
)

func TestQueryOrder2(t *testing.T) {
	orderNo := "30814"
	reply, err := pay.QueryOrder("", orderNo)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v\n", reply)
}

//先查询订单状态，再撤销
func TestReverse(t *testing.T) {
	reply3, err := pay.ReverseOrder("1428339002", "sandbox_test_31")
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%+v\n", reply3)
}

func TestRefundQuery(t *testing.T) {
	subid := "1900000109"
	orderNo := "30814"
	rep, err := pay.QueryRefund(subid, orderNo)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v\n", rep)
}

func TestRefundApply(t *testing.T) {
	reply3, err := pay.RefundApply("1428339002", "sandbox_test_31")
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%+v\n", reply3)
}

func TestDownBill(t *testing.T) {
	rep, err := pay.DownLoadBill("", "20170101")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("+%v", rep)
}
