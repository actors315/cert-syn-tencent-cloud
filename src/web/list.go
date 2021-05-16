package web

import (
	"fmt"
	"html/template"
	"net/http"
	"qcloud-tools/src/db"
	"qcloud-tools/src/tools"
	"time"
)

type IssueShow struct {
	Id            uint64
	DnsApi        string
	CdnType       string `default:"cdn"`
	MainDomain    string
	ExtraDomain   string
	LastIssueTime string
}

type List struct {
	Item []IssueShow
}

func GetList(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "text/html")

	rootPath := tools.GetRootPath()
	templatePath := fmt.Sprintf("%s/web/list.html", rootPath)
	tpl, err := template.ParseFiles(templatePath)
	if nil != err {
		fmt.Fprint(writer, "<div>Error ~~</div>")
		return
	}

	sqlStr := "SELECT id,dns_api,type,main_domain,extra_domain,last_issue_time FROM issue_info order by id desc"
	rows, err := db.QcloudToolDb.Query(sqlStr)
	if err != nil {
		fmt.Fprint(writer, "<div>Error Query~~</div>")
		return
	}
	defer rows.Close()

	var list List
	for rows.Next() {
		var issue IssueShow
		var lastIssueTime int64
		err = rows.Scan(
			&issue.Id,
			&issue.DnsApi,
			&issue.CdnType,
			&issue.MainDomain,
			&issue.ExtraDomain,
			&lastIssueTime)
		if nil != err {
			fmt.Println("scan row error:", err)
			continue
		}
		if lastIssueTime > 0 {
			issue.LastIssueTime = time.Unix(lastIssueTime, 0).Format("2006-01-02")
		}
		list.Item = append(list.Item, issue)
	}

	_ = tpl.Execute(writer, list)
}
