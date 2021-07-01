package certificate

import (
	"context"
	"fmt"
	"qcloud-tools/src/db"
	"strconv"
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

	rows, err := db.QcloudToolDb.Query(sqlStr, now-31*86400, now-86400)
	if err != nil {
		return
	}
	defer rows.Close()

	var rowIdArr []string

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
		rowIdArr = append(rowIdArr, strconv.FormatUint(issue.Id, 10))

		issue.IssueCert()
	}

	if nil != rowIdArr {
		sqlStr := "UPDATE issue_info SET last_check_time = ? WHERE id IN (?)"
		_, _ = db.QcloudToolDb.Update(sqlStr, now, strings.Join(rowIdArr, ","))
	}

	fmt.Println("check update done")
}
