package util

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func Trim(s string) string {
	s = strings.Replace(s, "<br />", "", -1)
	s = strings.Replace(s, "<br>", "", -1)
	s = strings.Replace(s, "&quot;", "\"", -1)
	s = strings.Replace(s, "&#039;", "'", -1)
	return s
}

func ReadReqBody(req *http.Request) map[string]string {
	body, _ := ioutil.ReadAll(req.Body)
	// Info.Println("body = ", string(body))
	m := make(map[string]string)

	content := strings.Split(string(body), "&")
	for _, c := range content {
		p := strings.Split(c, "=")
		if len(p) == 2 {
			m[p[0]] = p[1]
		}
	}
	return m
}

func AuthenticateRequest(req *http.Request) bool {
	pass := false
	if strings.Contains(req.URL.Path, "/scd/callback") {
		return true
	}

	/*for _, c := range req.Cookies() { // Simple authentication, should use OAuth 2.0
		if c.Name == "JSESSIONID" {
			pass = true
			break
		}
	}*/

	basicAuthPrefix := "Basic "
	user := []byte("foo")
	passwd := []byte("bar")

	auth := req.Header.Get("Authorization")
	if strings.HasPrefix(auth, basicAuthPrefix) {
		payload, err := base64.StdEncoding.DecodeString(
			auth[len(basicAuthPrefix):],
		)
		if err == nil {
			pair := bytes.SplitN(payload, []byte(":"), 2)
			if len(pair) == 2 && bytes.Equal(pair[0], user) &&
				bytes.Equal(pair[1], passwd) {

				return true
			}
		}
	}
	return pass
}

func RunStmt(stmt string) string {
	stmt, _ = url.QueryUnescape(stmt)
	Info.Println("Running SQL:", stmt)

	db := GetDB()
	defer db.Close()

	var msg string

	stmtTmp := strings.ToLower(stmt)
	if strings.HasPrefix(stmtTmp, "update") || strings.HasPrefix(stmtTmp, "insert") || strings.HasPrefix(stmtTmp, "delete") {
		tx, err := db.Begin()

		if err != nil {
			return err.Error()
		}
		stmt, err := tx.Prepare(stmt)
		if err != nil {
			return err.Error()
		}

		res, err := stmt.Exec()
		if err != nil {
			return err.Error()
		}

		err = tx.Commit()
		if err != nil {
			return err.Error()
		}

		stmt.Close()

		affect, err := res.RowsAffected()
		msg = fmt.Sprintf("%d %s", affect, "rows affected")

		if err != nil {
			return err.Error()
		}
	} else if strings.HasPrefix(stmtTmp, "select") {
		rows, err := db.Query(stmt)
		if err != nil {
			return err.Error()
		}

		cols, err := rows.Columns()
		if err != nil {
			return err.Error()
		}

		rawResult := make([][]byte, len(cols))
		rowResult := make([]string, len(cols))

		var result []string

		dest := make([]interface{}, len(cols)) // A temporary interface{} slice
		for i, _ := range rawResult {
			dest[i] = &rawResult[i] // Put pointers to each string in the interface slice
		}

		for rows.Next() {
			//pointers := make([]interface{}, len(columnNames))
			/*for i, colName := range columnNames {
			    pointers[i] = reflect.ValueOf(obj).FieldByName(strings.Title(colName)).Addr().Interface()

			}*/
			//err = rows.Scan(pointers...)
			err = rows.Scan(dest...)
			if err != nil {
				return err.Error()
			}

			for i, raw := range rawResult {
				if raw == nil {
					rowResult[i] = ""
				} else {
					rowResult[i] = string(raw)
				}
			}
			result = append(result, strings.Join(rowResult, " | "))
		}

		msg = strings.Join(cols, " | ") + "<br/>" + strings.Join(result, "<br/>")
	}

	return msg
}

func CompressData(res http.ResponseWriter, data []byte) error {
	res.Header().Set("Content-Encoding", "gzip")
	writer, err := gzip.NewWriterLevel(res, gzip.BestCompression)
	if err != nil {
		Error.Println("fail to compress data", err)
		return err
	}

	defer writer.Close()
	writer.Write(data)

	return nil
}
