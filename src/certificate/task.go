package certificate

import (
	"context"
	"fmt"
	"qcloud-tools/src/db"
	"strings"
	"time"
)

func TickerSchedule(ctx context.Context) {

	ticker := time.NewTicker(time.Duration(86400) * time.Second)
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Ticker Schedule done:", ctx.Err())
			return
		case <-ticker.C:
			go checkUpdate()
		default:
			time.Sleep(time.Duration(300) * time.Second)
		}
	}
}

func checkUpdate() {

	var err error
	fields := "id,secret_id,secret_key,type,cdn_domain,issue_id"
	sqlStr := fmt.Sprintf("SELECT %s FROM issue_sync WHERE last_issue_time < ? AND last_check_time < ?", fields)
	now := time.Now().Unix()

	rows, err := db.QcloudToolDb.Query(sqlStr, now-62*86400, now-86400)
	if err != nil {
		return
	}
	defer rows.Close()

	var rowIdArr []interface{}

	for rows.Next() {
		var issue IssueSync
		err = rows.Scan(
			&issue.Id,
			&issue.SecretId,
			&issue.SecretKey,
			&issue.CdnType,
			&issue.CdnDomain,
			&issue.IssueId)
		if nil != err {
			fmt.Println("scan row error:", err)
			continue
		}
		rowIdArr = append(rowIdArr, issue.Id)

		issue.IssueCert()
	}

	if nil != rowIdArr {
		sqlStr := fmt.Sprintf("UPDATE issue_sync SET last_check_time = %d WHERE id IN (%s?)", now, strings.Repeat("?, ", len(rowIdArr)-1))
		_, _ = db.QcloudToolDb.Update(sqlStr, rowIdArr...)
	}

	fmt.Println("check update done")
}
