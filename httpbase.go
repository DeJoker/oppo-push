package oppopush

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"fmt"
	"os"
	"io"
	"bytes"
	"mime/multipart"
	"net/textproto"
)

var (
	client = &http.Client{
		Timeout : time.Second * 60,
	}
)

const (
	RetryTimes       = 2
)

func doPost(url string, form url.Values) ([]byte, error) {
	req, err := http.NewRequest("POST", url, strings.NewReader(form.Encode()))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	tryTime := 0
tryAgain:
	resp, err := client.Do(req)
	if err != nil {
		tryTime++
		if tryTime < RetryTimes {
			goto tryAgain
		}
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("HTTP status code:"+strconv.Itoa(resp.StatusCode))
	}
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	str := string(result)
	str, err = strconv.Unquote(str)
	if err != nil {
		str = string(result)
	}
	return []byte(str), nil
}

func doGet(url string, params string) ([]byte, error) {
	req, err := http.NewRequest("GET", url+params, nil)
	if err != nil {
		panic(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("HTTP status code:"+strconv.Itoa(resp.StatusCode))
	}
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return result, nil
}


func commonUploadRealFile(filename, url string, params map[string]string) string {
	bufReader := new(bytes.Buffer)
	mpWriter := multipart.NewWriter(bufReader)

	//添加参数
	for key, value := range params {
		mpWriter.WriteField(key, value)
	}
	//最后写入文件
	constructFormFile(mpWriter, filename)
	mpWriter.Close()

	req, _ := http.NewRequest("POST", url, bufReader)

	req.Header.Add("content-type", mpWriter.FormDataContentType())

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return err.Error()
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	return UrlDecode(body)
}



func constructFormFile(writer *multipart.Writer, filename string) error {
	f, _ := os.Open(filename)
	defer f.Close()

	buffer := new(bytes.Buffer)
	_, err := io.Copy(buffer, f)
	if err != nil {
		return fmt.Errorf("copying f %v", err)
	}

	//推断Content-Type
	contentType := http.DetectContentType(buffer.Bytes())

	waitToWriteContent, _ := createFormFile(writer, contentType, "file", filename)
	//写入文件内容
	waitToWriteContent.Write(buffer.Bytes())

	return nil
}

func createFormFile(w *multipart.Writer, contentType, fieldname, filename string) (io.Writer, error) {
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
			escapeQuotes(fieldname), escapeQuotes(filename)))
	h.Set("Content-Type", contentType)
	return w.CreatePart(h)
}


var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")
func escapeQuotes(s string) string {
	return quoteEscaper.Replace(s)
}

//Json对& < > 会默认unicode编码，escape，当前未找到合适解决方案
func UrlDecode(body []byte) string {
	content := string(body)
	content = strings.Replace(content, "\\u003c", "<", -1)
	content = strings.Replace(content, "\\u003e", ">", -1)
	content = strings.Replace(content, "\\u0026", "&", -1)

	return content
}




