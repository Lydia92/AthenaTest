package spilt

import (
	"crypto/md5"
	"encoding/hex"
	"regexp"
	"strings"
)

func Split(str string) map[string]string {
	var grants string
	var database string
	var table string
	var user string
	var hosts string
	var option string
	//var results string
	res := make(map[string]string)
	if !strings.Contains(str, "PROXY ON") || !strings.Contains(str, "USAGE") {
		srr := strings.Split(str, " ON ")[0]
		grants = strings.Split(srr, "GRANT")[1]
		//fmt.Println(grants)
		srr1 := strings.Split(str, " ON ")[1]
		database = strings.Split(srr1, "TO")[0]
		table = strings.Split(database, ".")[0]
		database = strings.Split(database, ".")[0]

		str1 := strings.Split(srr1, "TO")[1]
		user = strings.Split(str1, "@")[0]
		hosts = strings.Split(str1, "@")[1]
		hosts = strings.Split(hosts, "' ")[0]
		hosts = strings.Split(hosts, "'")[1]
		//hosts= strings.Replace(hosts, " ", "", -1)

		//fmt.Println(len(strings.Split(hosts,"WITH")))
		if strings.HasSuffix(str, "GRANT OPTION") {
			option = "1"
		} else {
			option = "0"
		}
		//results=fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t%s",grants,database,table,user,hosts,option)
		res["grants"] = grants
		res["database"] = database
		res["table"] = table
		res["user"] = user
		res["hosts"] = hosts
		res["option"] = option

		//fmt.Println(res)
	}
	return res
}
func MD5(text string) string {
	ctx := md5.New()
	ctx.Write([]byte(text))
	return hex.EncodeToString(ctx.Sum(nil))
}

func SplitSQL(sql string) []string {
	var result []string
	reg := regexp.MustCompile(`(from|join)\s\S+`)
	res := reg.FindAllString(sql, -1)
	for _, v := range res {
		if strings.Contains(v, "from") {
			s := strings.Split(v, "from")[1]
			if strings.Contains(s, ",") {
				sptring := strings.Split(s, ",")
				for _, v1 := range sptring {
					result = append(result, v1)
				}
			} else {
				result = append(result, s)
			}

		} else if strings.Contains(v, "join") {
			s := strings.Split(v, "join")[1]
			result = append(result, s)
		}
	}
	return result
}
