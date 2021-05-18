package web

import (
	"fmt"
	"html/template"
	"net/http"
	"qcloud-tools/src/certificate"
	"qcloud-tools/src/db"
	"qcloud-tools/src/tools"
)

const (
	GET  = "GET"
	POST = "POST"
)

type DnsApi struct {
	Id         string
	Name       string
	AppIdName  string
	AppKeyName string
}

type IssueForm struct {
	SecretId    string
	SecretKey   string
	CdnType     string
	DnsApi      string
	AppId       string
	AppValue    string
	MainDomain  string
	ExtraDomain string
}

type CdnType struct {
	Id   string
	Name string
}

var DnsApiList = map[string]DnsApi{
	"dns_dp": {
		Id:         "dns_dp",
		Name:       "dnspod",
		AppIdName:  "DP_Id",
		AppKeyName: "DP_Key",
	},
	"dns_cf": {
		Id:         "dns_cf",
		Name:       "Cloudflare",
		AppIdName:  "CF_Token",
		AppKeyName: "CF_Account_ID",
	},
	"dns_gd": {
		Id:         "dns_cf",
		Name:       "GoDaddy.com",
		AppIdName:  "GD_Key",
		AppKeyName: "GD_Secret",
	},
	"dns_aws": {
		Id:         "dns_aws",
		Name:       "Amazon Route53",
		AppIdName:  "AWS_ACCESS_KEY_ID",
		AppKeyName: "AWS_SECRET_ACCESS_KEY",
	},
	"dns_ali": {
		Id:         "dns_ali",
		Name:       "Aliyun",
		AppIdName:  "Ali_Key",
		AppKeyName: "Ali_Secret",
	},
}

var CdnTypeList = []CdnType{
	{
		Id:   "cdn",
		Name: "内容分发网络",
	},
	{
		Id:   "ecdn",
		Name: "全站加速网络",
	},
}

func AddDomain(writer http.ResponseWriter, request *http.Request) {

	if POST == request.Method {
		_ = request.ParseForm()

		var issueForm IssueForm
		issueForm.SecretId = request.Form.Get("secret_id")
		issueForm.SecretKey = request.Form.Get("secret_key")
		issueForm.CdnType = request.Form.Get("type")
		issueForm.DnsApi = request.Form.Get("dns_api")
		issueForm.AppId = request.Form.Get("app_id_value")
		issueForm.AppValue = request.Form.Get("app_key_value")
		issueForm.MainDomain = request.Form.Get("main_domain")
		issueForm.ExtraDomain = request.Form.Get("extra_domain")

		if err := issueForm.Add(); err != nil {
			fmt.Fprintf(writer, fmt.Sprintf(`{"code":1,"msg":%s}`, err))
		} else {
			fmt.Fprintf(writer, `{"code":0}`)
		}
		return
	}

	writer.Header().Set("Content-Type", "text/html")
	rootPath := tools.GetRootPath()
	templatePath := fmt.Sprintf("%s/web/add.html", rootPath)
	tpl, _ := template.ParseFiles(templatePath)

	var form = struct {
		DnsApiList  map[string]DnsApi
		CdnTypeList []CdnType
	}{
		DnsApiList,
		CdnTypeList,
	}

	_ = tpl.Execute(writer, form)
}

func (form IssueForm) Add() (err error) {

	issue := certificate.Issue{
		SecretId:    form.SecretId,
		SecretKey:   form.SecretKey,
		AppIdValue:  form.AppId,
		AppKeyValue: form.AppValue,
		DnsApi:      form.DnsApi,
		CdnType:     form.CdnType,
		MainDomain:  form.MainDomain,
		ExtraDomain: form.ExtraDomain,
	}

	dnsapi, ok := DnsApiList[form.DnsApi]
	if !ok {
		fmt.Printf("%s 不存在\n", form.DnsApi)
		return err
	}

	issue.AppIdName = dnsapi.AppIdName
	issue.AppKeyName = dnsapi.AppKeyName

	sqlStr := `INSERT INTO issue_info (
secret_id,secret_key,dns_api,app_id,app_id_value,app_key,app_key_value,type,main_domain,extra_domain
) VALUES (?,?,?,?,?,?,?,?,?,?)`
	_, err = db.QcloudToolDb.Insert(sqlStr,
		issue.SecretId,
		issue.SecretKey,
		issue.DnsApi,
		issue.AppIdName,
		issue.AppIdValue,
		issue.AppKeyName,
		issue.AppKeyValue,
		issue.CdnType,
		issue.MainDomain,
		issue.ExtraDomain)
	if nil != err {
		return err
	}

	return
}
