package helper

import (
	"bytes"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

type HttpMethod string
type ContentType string

const (
	MethodGet            HttpMethod  = "GET"
	MethodPost           HttpMethod  = "POST"
	ContentJSON          ContentType = "application/json"
	ContentFormUrlEncode ContentType = "application/x-www-form-urlencoded"
)

type ApiHelper struct {
	Token          string
	Type           string
	Language       string
	BaseUrl        string
	QueryParam     string
	Path           string
	Body           []byte
	ContentType    ContentType
	Method         HttpMethod
	File           string
	FileFieldName  string
	ContentTypeStr string
	err            error
	BodyBuffer     *bytes.Buffer
}

func NewApiHelper(path, token, baseUrl, apiType, language string) *ApiHelper {
	apiHelper := &ApiHelper{Token: token, Type: "Bot", BaseUrl: "https://www.kaiheila.cn", Language: "zh-CN"}

	if baseUrl != "" {
		apiHelper.BaseUrl = baseUrl
	}
	if apiType != "" {
		apiHelper.Type = apiType
	}
	if language != "" {
		apiHelper.Language = language
	}
	apiHelper.Path = path
	apiHelper.ContentType = ContentJSON
	apiHelper.Method = MethodGet

	return apiHelper
}
func (h *ApiHelper) SetQuery(values map[string]string) *ApiHelper {

	for k, v := range values {
		if h.QueryParam == "" {
			h.QueryParam = fmt.Sprintf("%s=%s", k, v)
		} else {
			h.QueryParam += fmt.Sprintf("&%s=%s", k, v)
		}
	}
	return h

}

func (h *ApiHelper) SetBody(body []byte) *ApiHelper {
	h.Body = body
	return h
}

func (h *ApiHelper) SetContentType(contentType ContentType) *ApiHelper {
	h.ContentType = contentType
	return h
}

func (h *ApiHelper) SetUploadFile(filePath string) *ApiHelper {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	fileName := filepath.Base(filePath)
	// this step is very important
	fileWriter, err := bodyWriter.CreateFormFile("file", fileName)
	if err != nil {
		log.Error(err)
		h.err = err
		return h
	}

	// open file handle
	fh, err := os.Open(filePath)
	defer fh.Close()
	if err != nil {
		log.Error(err)
		h.err = err
		return h
	}

	//iocopy
	_, err = io.Copy(fileWriter, fh)
	if err != nil {
		log.Error(err)
		h.err = err
		return h
	}
	bodyWriter.Close()
	h.BodyBuffer = bodyBuf
	h.ContentTypeStr = bodyWriter.FormDataContentType()
	return h
}
func (h *ApiHelper) Get() ([]byte, error) {
	h.Method = MethodGet
	return h.Send()
}
func (h *ApiHelper) Post() ([]byte, error) {
	h.Method = MethodPost
	return h.Send()
}

func (h *ApiHelper) getReqPath() string {
	reqPath := ""
	if strings.HasPrefix(h.Path, "/") || strings.HasSuffix(h.BaseUrl, "/") {
		reqPath = h.BaseUrl + h.Path
	} else {
		reqPath = h.BaseUrl + "/" + h.Path
	}
	if h.QueryParam != "" {
		reqPath += "?" + h.QueryParam
	}
	return reqPath
}
func (h *ApiHelper) setHeader(req *http.Request) {
	if h.ContentTypeStr != "" {
		req.Header.Set("Content-Type", h.ContentTypeStr)
	} else {
		req.Header.Set("Content-Type", string(h.ContentType))
	}
	req.Header.Set("Authorization", fmt.Sprintf("%s %s", h.Type, h.Token))
	req.Header.Set("Accept-Language", h.Language)
}
func (h *ApiHelper) Send() ([]byte, error) {
	if h.err != nil {
		return nil, h.err
	}
	client := &http.Client{}
	reqPath := h.getReqPath()
	var req *http.Request
	var err error
	if h.Body != nil {
		req, err = http.NewRequest(string(h.Method), reqPath, bytes.NewBuffer(h.Body))
	} else if h.BodyBuffer != nil {
		req, err = http.NewRequest(string(h.Method), reqPath, h.BodyBuffer)
	} else {
		req, err = http.NewRequest(string(h.Method), reqPath, nil)
	}
	if err != nil {
		return nil, err
	}
	h.setHeader(req)
	printRequestAsCurl(req)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		var data []byte
		if resp.Body != nil {
			data, _ = io.ReadAll(resp.Body)
			log.WithField("statusCode", resp.StatusCode).WithField("data", string(data)).Error("http error", reqPath)
		} else {
			log.WithField("statusCode", resp.StatusCode).WithField("data", nil).Error("http error", reqPath)
		}

		return nil, errors.New("http error")
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (h *ApiHelper) String() string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("Path:%s BaseUrl:%s Method:%s Query:%s ", h.Path, h.BaseUrl, h.Method, h.QueryParam))
	if len(h.Body) > 0 {
		sb.WriteString(fmt.Sprintf("Body:%s", string(h.Body)))
	}
	return sb.String()
}

func printRequestAsCurl(req *http.Request) {
	method := req.Method
	urlStr := req.URL.String()
	u, _ := url.Parse(urlStr)

	// 构建 curl 命令
	curlCmd := "curl -X " + method + " "
	if u.User != nil {
		curlCmd += "-u " + u.User.String() + " "
	}
	curlCmd += "'" + urlStr + "'"

	// 处理请求头
	for key, values := range req.Header {
		for _, value := range values {
			curlCmd += " -H '" + key + ": " + value + "'"
		}
	}

	// 处理请求体（如果有）
	if req.Body != nil {
		// 这里假设请求体是字符串，实际可能需要更复杂的处理
		body, _ := req.GetBody()
		b, _ := io.ReadAll(body)
		curlCmd += " --data '" + string(b) + "'"
	}

	fmt.Println(curlCmd)
}

//func (h *ApiHelper) PrintAsCurl() string {
//	sb := strings.Builder{}
//	//curl -X POST 'http://127.0.0.1/notify_users'    -H 'Content-Type: application/json'    -d '{"author_id":1978522,"only_visible_id":0,"guild_id":469731,"channel_id":2487355,"mention_here":0, "mention_all":0, "mention_roles":[], "mention_target_ids":[]}'
//	reqPath := h.getReqPath()
//
//	curlStr := `curl -X POST '%s'    -H '%s'   -d '%s'`
//
//}
