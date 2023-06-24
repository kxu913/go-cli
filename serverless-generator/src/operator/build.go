package operator

import (
	"bufio"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
)

var (
	dockerRegistryUserID = ""
	Output               = GetEnvWithDefault("output", "D:\\tmp\\wsl")
	AliAuthConfig        = types.AuthConfig{
		Username:      "your account",
		Password:      "your pwd",
		ServerAddress: "https://registry.cn-shenzhen.aliyuncs.com",
	}
	format        = "20060102150405"
	AliRepository = "registry.cn-shenzhen.aliyuncs.com/kevin-demo/cdi"

	TxAuthConfig = types.AuthConfig{
		Username:      "your account",
		Password:      "your pwd",
		ServerAddress: "https://ccr.ccs.tencentyun.com/minicloud/cdi",
	}

	TxRepository = "ccr.ccs.tencentyun.com/minicloud/cdi"
)

type ErrorLine struct {
	Error       string      `json:"error"`
	ErrorDetail ErrorDetail `json:"errorDetail"`
}

type ErrorDetail struct {
	Message string `json:"message"`
}

func GetCloudProviderConfig(cloudProvider string) (types.AuthConfig, string) {
	if strings.EqualFold(cloudProvider, "ali") {
		return AliAuthConfig, AliRepository

	} else {
		return TxAuthConfig, TxRepository
	}
}

func getPort(project string) int {
	b, err := ioutil.ReadFile(fmt.Sprintf("%s/%s/Dockerfile", Output, project))
	if err != nil {
		panic(err)
	}
	lineBytes := string(b)
	lines := strings.Split(lineBytes, "\n")
	for _, l := range lines {
		portLine := strings.Split(l, "EXPOSE")
		if len(portLine) > 1 {
			p, err := strconv.Atoi(strings.TrimSpace(portLine[1]))
			if err != nil {
				panic(err)
			}
			return p

		}
	}
	return 0
}

func BuildImage(project string, cloudProvider string) (string, int) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*120)
	defer cancel()
	tar, err := archive.TarWithOptions(fmt.Sprintf("%s/%s/", Output, project), &archive.TarOptions{})
	if err != nil {
		panic(err)
	}

	_, repository := GetCloudProviderConfig(cloudProvider)
	remoteTag := fmt.Sprintf("%s:%s", repository, time.Now().Format(format))
	opts := types.ImageBuildOptions{
		Dockerfile: "Dockerfile",
		Tags:       []string{project + ":latest", remoteTag},
		Remove:     true,
	}
	res, err := cli.ImageBuild(ctx, tar, opts)
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	err = print(res.Body)
	if err != nil {
		panic(err)
	}
	// pushImageToRemote(cli, aliAuthConfig, aliTag)
	// pushImageToRemote(cli, txAuthConfig, txTag)
	return remoteTag, getPort(project)

}

func PushImageToRemote(cloudProvider string, remoteTag string) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	config, _ := GetCloudProviderConfig(cloudProvider)
	authConfigBytes, _ := json.Marshal(config)
	authConfigEncoded := base64.URLEncoding.EncodeToString(authConfigBytes)
	repositoryOpts := types.ImagePushOptions{RegistryAuth: authConfigEncoded}
	rd, err := cli.ImagePush(context.Background(), remoteTag, repositoryOpts)
	if err != nil {
		panic(err)
	}
	err = print(rd)
	if err != nil {
		panic(err)
	}
}

func print(rd io.Reader) error {
	var lastLine string

	scanner := bufio.NewScanner(rd)
	for scanner.Scan() {
		lastLine = scanner.Text()
		fmt.Println(scanner.Text())
	}

	errLine := &ErrorLine{}
	json.Unmarshal([]byte(lastLine), errLine)
	if errLine.Error != "" {
		return errors.New(errLine.Error)
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}
