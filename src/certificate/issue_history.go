package certificate

import (
	"fmt"
	"qcloud-tools/src/db"
	"strings"
	"time"
)

type IssueHistory struct {
	IssueDomain string
	PublicKey   string
	PrivateKey  string
	CreatedAt   uint
}

func GetLatestValidRecord(domain string) (history IssueHistory) {

	sqlStr := "SELECT issue_domain,public_key,private_key,created_at FROM issue_history WHERE issue_domain in (?) AND created_at > ? ORDER BY id DESC LIMIT 1"
	now := time.Now().Unix()

	var arr []string
	arr = append(arr, domain)

	index := strings.Index(domain, ".")

	arr = append(arr, "*"+domain[index:])

	domain = strings.Join(arr, "','")

	rows, err := db.QcloudToolDb.Query(sqlStr, "'"+domain+"'", now-86400*62)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(
			history.IssueDomain,
			history.PublicKey,
			history.PrivateKey,
			history.CreatedAt)

		if err != nil {
			fmt.Println(err)
		}
	}


	return
}

func (history IssueHistory) Add() {
	sql := "INSERT INTO issue_history (issue_domain,public_key,private_key,created_at) values (?, ?, ?, ?)"
	_, _ = db.QcloudToolDb.Insert(sql,
		history.IssueDomain,
		history.PublicKey,
		history.PrivateKey,
		history.CreatedAt)
}