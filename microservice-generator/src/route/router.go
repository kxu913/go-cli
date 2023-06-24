package route

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"microservice-generator/model"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/labstack/echo/v4"
)

func GetEnvWithDefault(key string, defValue string) string {
	val, err := os.LookupEnv(key)
	if !err {
		return defValue
	}
	return val

}

var (
	basicRequestUrl   = GetEnvWithDefault("basic_host", "http://localhost:1323")
	dbRequestUrl      = GetEnvWithDefault("db_host", "http://localhost:1324")
	graphqlRequestUrl = GetEnvWithDefault("graphql_host", "http://localhost:9004")
	Output            = GetEnvWithDefault("output", "d:/tmp/wsl/")
	sampleDateTime    = "20230507164000"
)

func CreateService(c echo.Context) error {
	request := &model.ServiceRequest{}
	if err := c.Bind(request); err != nil {
		fmt.Println(err)
		return c.String(http.StatusNotAcceptable, "Invalid Json")
	}
	basicRequest := request.Basic
	fmt.Println(basicRequest)
	_basicRequest, _ := json.Marshal(basicRequest)
	http.Post(fmt.Sprintf("%s/cli/v1/init", basicRequestUrl), "application/json", bytes.NewBuffer(_basicRequest))
	dbRequest := model.SetDBRequestFromBasicRequest(request)
	fmt.Println(dbRequest)
	_dbRequest, _ := json.Marshal(dbRequest)
	http.Post(fmt.Sprintf("%s/cli/v1/db/%s", dbRequestUrl, dbRequest.Table), "application/json", bytes.NewBuffer(_dbRequest))

	if request.Graphql.SQL != "" {
		graphqlRequest := model.SetGraphqlRequestFromBasicRequest(request)
		fmt.Println(graphqlRequest)
		_graphqlRequest, _ := json.Marshal(graphqlRequest)
		http.Post(fmt.Sprintf("%s/graphql/v1/sql", graphqlRequestUrl), "application/json", bytes.NewBuffer(_graphqlRequest))
	}

	zipName := basicRequest.ProjectName + "-" + time.Now().Format(sampleDateTime) + ".zip"
	err := ZipSource(fmt.Sprintf("%s%s", Output, basicRequest.ProjectName), fmt.Sprintf("%s%s/%s", Output, basicRequest.ProjectName, zipName))
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusInternalServerError, "zip files error.")
	}

	return c.Attachment(fmt.Sprintf("%s%s/%s", Output, basicRequest.ProjectName, zipName), zipName)

}

func Group(e *echo.Echo) {
	t := e.Group("/it/v1")
	t.POST("/create", CreateService)

}

func ZipSource(source, target string) error {
	// 1. Create a ZIP file and zip.Writer
	f, err := os.Create(target)
	if err != nil {
		return err
	}
	defer f.Close()

	writer := zip.NewWriter(f)
	defer writer.Close()

	// 2. Go through all the files of the source
	return filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 3. Create a local file header
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		// set compression
		header.Method = zip.Deflate

		// 4. Set relative path of a file as the header name
		header.Name, err = filepath.Rel(filepath.Dir(source), path)
		if err != nil {
			return err
		}
		if info.IsDir() {
			header.Name += "/"
		}

		// 5. Create writer for the file header and save content of the file
		headerWriter, err := writer.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		_, err = io.Copy(headerWriter, f)
		return err
	})
}
