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
	Id   string
	Name string
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

var DnsApiList = []DnsApi{
	{
		Id:   "dns_dp",
		Name: "dnspod",
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
		DnsApiList []DnsApi
	}{
		DnsApiList,
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

	switch form.DnsApi {
	case "dns_dp":
		issue.AppIdName = "DP_Id"
		issue.AppKeyName = "DP_Key"
	}

	sqlStr := `INSERT INTO issue_info (
secret_id,secret_key,dns_api,app_id,app_id_value,app_key,app_key_value,type,main_domain,extra_domain
) VALUES (?,?,?,?,?,?,?,?,?,?)`
	lastInsertId, err := db.QcloudToolDb.Insert(sqlStr,
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

	go issue.IssueCertByScript(uint64(lastInsertId))

	return
}
