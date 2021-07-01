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
	fields := "id,secret_id,secret_key,app_id,app_id_value,app_key,app_key_value,dns_api,type,cdn_domain,main_domain,extra_domain FROM issue_info"
	sqlStr := fmt.Sprintf("SELECT %s WHERE last_issue_time < ? AND last_check_time < ?", fields)
	now := time.Now().Unix()

	rows, err := db.QcloudToolDb.Query(sqlStr, now-31*86400, now-86400)
	if err != nil {
		return
	}
	defer rows.Close()

	var rowIdArr []string

	for rows.Next() {
		var issue Issue
		var rowId uint64
		err = rows.Scan(
			&rowId,
			&issue.SecretId,
			&issue.SecretKey,
			&issue.AppIdName,
			&issue.AppIdValue,
			&issue.AppKeyName,
			&issue.AppKeyValue,
			&issue.DnsApi,
			&issue.CdnType,
			&issue.CdnDomain,
			&issue.MainDomain,
			&issue.ExtraDomain)
		if nil != err {
			fmt.Println("scan row error:", err)
			continue
		}
		rowIdArr = append(rowIdArr, strconv.FormatUint(rowId, 10))

		issue.IssueCert(rowId)
	}

	if nil != rowIdArr {
		sqlStr := "UPDATE issue_info SET last_check_time = ? WHERE id IN (?)"
		_, _ = db.QcloudToolDb.Update(sqlStr, now, strings.Join(rowIdArr, ","))
	}

	fmt.Println("check update done")
}
