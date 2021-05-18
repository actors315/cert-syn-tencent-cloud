package certificate

import (
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"qcloud-tools/src/db"
	"qcloud-tools/src/tools"
	"regexp"
	"time"
)

type Issue struct {
	SecretId      string
	SecretKey     string
	AppIdName     string
	AppIdValue    string
	AppKeyName    string
	AppKeyValue   string
	DnsApi        string
	CdnType       string `default:"cdn"`
	MainDomain    string
	ExtraDomain   string
}

func (issue *Issue) GenerateScript() (string, error) {

	rootPath := tools.GetRootPath()

	fileName := fmt.Sprintf("%s/shell/%s.sh", rootPath, issue.MainDomain)
	f, err := os.Create(fileName)
	defer f.Close()

	if nil != err {
		fmt.Printf("创建文件失败: %s \n", fileName)
		return "", err
	}

	templatePath := fmt.Sprintf("%s/config/issue-template.tpl", rootPath)
	tpl, err := template.ParseFiles(templatePath)
	if nil != err {
		fmt.Printf("模板文件不存在: %s \n", templatePath)
		return "", err
	}

	if err := tpl.Execute(f, issue); err != nil {
		fmt.Printf("生成脚本失败：%s \n", err)
		return "", err
	}

	if err = f.Chmod(0777); err != nil {
		fmt.Printf("更改文件权限失败：%s \n", err)
		return fileName, err
	}

	return fileName, nil
}

func (issue *Issue) IssueCertByScript(rowId uint64) {
	fileName, err := issue.GenerateScript()
	if err != nil {
		return
	}

	command := exec.Command(fileName)
	stdout, _ := command.Output()
	command.Wait()

	var privateKeyPath, publicKeyPath string

	content := string(stdout)

	fmt.Println("issue result", content)

	privateKeyRegexp, _ := regexp.Compile(`Your cert key is in  (.*\/\.acme\.sh\/.*[\S])`)
	publicKeyRegexp, _ := regexp.Compile(`And the full chain certs is there:  (.*\/\.acme\.sh\/.*[\S])`)

	var regexpResult []string
	regexpResult = privateKeyRegexp.FindStringSubmatch(content)
	if nil != regexpResult {
		privateKeyPath = regexpResult[1]
	}
	regexpResult = publicKeyRegexp.FindStringSubmatch(content)
	if nil != regexpResult {
		publicKeyPath = regexpResult[1]
	}

	if "" == privateKeyPath || "" == publicKeyPath {
		fmt.Printf("update certificate failed,private %s, public %s \n", privateKeyPath, publicKeyPath)
		return
	}

	// 更新证书到 cdn 或者 ecdn
	var syncInstance ISync
	sync := Sync{
		SecretId:       issue.SecretId,
		SecretKey:      issue.SecretKey,
		Domain:         issue.MainDomain,
		privateKeyPath: privateKeyPath,
		publicKeyPath:  publicKeyPath,
	}

	switch issue.CdnType {
	case "ecdn":
		syncInstance = EcdnSync{sync}
	default:
		syncInstance = CdnSync{sync}
	}

	result := syncInstance.UpdateCredential()

	// 更新数据库信息
	if result && rowId > 0 {
		now := time.Now().Unix()
		sqlStr := "UPDATE issue_info SET last_issue_time = ? WHERE id = ?"
		_, _ = db.QcloudToolDb.Update(sqlStr, now, rowId)
	}
}
