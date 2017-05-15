package wxisv

import (
	"encoding/xml"
	"fmt"
	"strings"
)

func (cli *Client) downLoad(cmd string, arg map[string]string) (resp []byte, err error) {
	signType := getSignType(cmd)
	odrInXml := cli.signedOrderRequestXmlString(signType, arg)
	resp, err = cli.doHttpPost(cli.GateWay+cmd, []byte(odrInXml), needSecure(cmd))
	if err != nil {
		return nil, err
	}

	if strings.Contains(string(resp), "xml") {
		var reply CommonReply
		err = xml.Unmarshal(resp, &reply)
		if err != nil {
			return
		}
		err = reply.GetError()
		return
	}

	fmt.Printf("%s\n", string(resp))

	return resp, nil
}

type ReverseReply struct {
	Appid  string `xml:"appid"`
	Mchid  string `xml:"mch_id"`
	Subid  string `xml:"sub_mch_id"`
	Recall string `xml:"recall"`
}

//撤销订单
//需要双向认证
func (p *Pay) ReverseOrder(subid, tradeno string) (*ReverseReply, error) {
	input := map[string]string{
		"out_trade_no": tradeno,
	}
	if len(subid) != 0 {
		input["sub_mch_id"] = subid
	}
	var resp struct {
		CommonReply
		ReverseReply
	}

	err := p.Client.sendCommand("secapi/pay/reverse", input, &resp)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &resp.ReverseReply, nil
}

type RefundApplyReply struct {
	Appid     string `xml:"appid"`
	Mchid     string `xml:"mch_id"`
	Subid     string `xml:"sub_mch_id"`
	TradeNo   string `xml:"out_trade_no"`
	TotalFee  int    `xml:"total_fee"` //订单金额
	RefundFee int    `xml:"refund_fee"`
	OpUserId  string `xml:"op_user_id"`
}

//申请退款
//需要双向认证
func (p *Pay) RefundApply(subid, tradeno string) (*ReverseReply, error) {
	input := map[string]string{
		"out_trade_no": tradeno,
		"op_user_id":   p.Client.MchId,
	}
	if len(subid) != 0 {
		input["sub_mch_id"] = subid
	}
	var resp struct {
		CommonReply
		ReverseReply
	}

	err := p.Client.sendCommand("secapi/pay/refund", input, &resp)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &resp.ReverseReply, nil
}

type RefundQueryApply struct {
	Appid        string `xml:"appid"`
	Mchid        string `xml:"mch_id"`
	Subid        string `xml:"sub_mch_id"`
	TradeNo      string `xml:"out_trade_no"`
	TotalFee     int    `xml:"total_fee"`
	CashFee      int    `xml:"cash_fee"`
	RefundCount  int    `xml:"refund_count"`
	RefundFee_0  int    `xml:"refund_fee_0"`
	RefundStatus string `xml:"refund_status_0"`
}

//查询退款
//this api not support partial refund
func (p *Pay) QueryRefund(subid, tradeno string) (*RefundQueryApply, error) {
	input := map[string]string{
		"out_trade_no": tradeno,
	}
	if len(subid) != 0 {
		input["sub_mch_id"] = subid
	}
	var resp struct {
		CommonReply
		RefundQueryApply
	}
	err := p.Client.sendCommand("pay/refundquery", input, &resp)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &resp.RefundQueryApply, nil
}

//下载账单
func (p *Pay) DownLoadBill(subid, date string) (string, error) {
	input := map[string]string{
		"bill_date": date,
		"bill_type": "ALL",
	}
	if len(subid) != 0 {
		input["sub_mch_id"] = subid
	}
	res, err := p.Client.downLoad("pay/downloadbill", input)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return string(res), err
}

func init() {
	setSignType("pay/refundquery", "MD5")
	setSignType("pay/downloadbill", "MD5")
	setSignType("secapi/pay/reverse", "MD5")
	setSignType("secapi/pay/refund", "MD5")
}
