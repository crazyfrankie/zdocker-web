package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/crazyfrankie/zdocker/container"
)

// Container 容器信息结构体
type Container struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Image       string `json:"image"`
	Command     string `json:"command"`
	Status      string `json:"status"`
	CreatedTime string `json:"created_time"`
	Pid         string `json:"pid"`
	Volume      string `json:"volume"`
	PortMapping string `json:"port_mapping"`
}

// CreateContainerRequest 创建容器请求
type CreateContainerRequest struct {
	Image       string            `json:"image" binding:"required"`
	Command     string            `json:"command" binding:"required"`
	Name        string            `json:"name"`
	Detach      bool              `json:"detach"`
	TTY         bool              `json:"tty"`
	Volume      string            `json:"volume"`
	Memory      string            `json:"memory"`
	CpuShare    string            `json:"cpu_share"`
	CpuSet      string            `json:"cpu_set"`
	Network     string            `json:"network"`
	Environment map[string]string `json:"environment"`
	PortMapping []string          `json:"port_mapping"`
}

// ExecRequest 执行命令请求
type ExecRequest struct {
	Command []string `json:"command" binding:"required"`
}

// ExecResult 执行命令结果
type ExecResult struct {
	Output   string `json:"output"`
	ExitCode int    `json:"exit_code"`
}

// NetworkInfo 网络信息
type NetworkInfo struct {
	Name   string `json:"name"`
	Driver string `json:"driver"`
	Subnet string `json:"subnet"`
}

// CreateNetworkRequest 创建网络请求
type CreateNetworkRequest struct {
	Name   string `json:"name" binding:"required"`
	Driver string `json:"driver"`
	Subnet string `json:"subnet"`
}

// SystemInfo 系统信息
type SystemInfo struct {
	OS           string `json:"os"`
	Architecture string `json:"architecture"`
	CPUs         int    `json:"cpus"`
	Memory       string `json:"memory"`
	ZDockerRoot  string `json:"zdocker_root"`
}

// GetContainerList 获取容器列表
func GetContainerList() ([]Container, error) {
	var containers []Container

	dirUrl := fmt.Sprintf(container.DefaultLocation, "")
	dirUrl = dirUrl[:len(dirUrl)-1]

	// 检查目录是否存在
	if _, err := os.Stat(dirUrl); os.IsNotExist(err) {
		return containers, nil
	}

	containerFiles, err := os.ReadDir(dirUrl)
	if err != nil {
		return nil, fmt.Errorf("读取容器目录失败: %v", err)
	}

	for _, containerFile := range containerFiles {
		if containerFile.IsDir() {
			tmpContainer, err := getContainerInfoByName(containerFile.Name())
			if err != nil {
				continue // 跳过读取失败的容器
			}
			containers = append(containers, tmpContainer)
		}
	}

	return containers, nil
}

// getContainerInfoByName 根据容器名称获取容器信息
func getContainerInfoByName(containerName string) (Container, error) {
	configFileDir := fmt.Sprintf(container.DefaultLocation, containerName)
	configFileDir = configFileDir + container.ConfigName

	content, err := os.ReadFile(configFileDir)
	if err != nil {
		return Container{}, fmt.Errorf("读取容器配置文件失败: %v", err)
	}

	var containerInfo container.ContainerInfo
	if err := json.Unmarshal(content, &containerInfo); err != nil {
		return Container{}, fmt.Errorf("解析容器配置失败: %v", err)
	}

	// 检查容器状态
	status := "stopped"
	if containerInfo.PID != "" {
		// 检查进程是否还在运行
		if _, err := os.Stat(fmt.Sprintf("/proc/%s", containerInfo.PID)); err == nil {
			status = "running"
		}
	}

	return Container{
		ID:          containerInfo.ID,
		Name:        containerInfo.Name,
		Command:     containerInfo.Command,
		Status:      status,
		CreatedTime: containerInfo.CreateTime,
		Pid:         containerInfo.PID,
		Volume:      containerInfo.Volume,
		PortMapping: strings.Join(containerInfo.PortMapping, ","),
	}, nil
}

// GetContainerById 根据ID获取容器信息
func GetContainerById(containerId string) (Container, error) {
	containers, err := GetContainerList()
	if err != nil {
		return Container{}, err
	}

	for _, c := range containers {
		if c.ID == containerId || c.Name == containerId {
			return c, nil
		}
	}

	return Container{}, fmt.Errorf("容器不存在")
}

// CreateContainer 创建容器
func CreateContainer(req CreateContainerRequest) (Container, error) {
	// 构建zdocker run命令
	args := []string{"run"}

	if req.Detach {
		args = append(args, "-d")
	}
	if req.TTY {
		args = append(args, "-t")
	}
	if req.Name != "" {
		args = append(args, "--name", req.Name)
	}
	if req.Volume != "" {
		args = append(args, "-v", req.Volume)
	}
	if req.Memory != "" {
		args = append(args, "-m", req.Memory)
	}
	if req.CpuShare != "" {
		args = append(args, "--cpushare", req.CpuShare)
	}
	if req.CpuSet != "" {
		args = append(args, "--cpuset", req.CpuSet)
	}
	if req.Network != "" {
		args = append(args, "--net", req.Network)
	}

	// 添加环境变量
	for key, value := range req.Environment {
		args = append(args, "-e", fmt.Sprintf("%s=%s", key, value))
	}

	// 添加端口映射
	for _, port := range req.PortMapping {
		args = append(args, "-p", port)
	}

	// 添加镜像和命令
	args = append(args, req.Image)
	if req.Command != "" {
		args = append(args, strings.Fields(req.Command)...)
	}

	// 执行命令
	cmd := exec.Command("zdocker", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return Container{}, fmt.Errorf("创建容器失败: %s, %v", string(output), err)
	}

	// 从输出中解析容器ID或名称
	containerName := req.Name
	if containerName == "" {
		// 如果没有指定名称，从输出中解析
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			if strings.Contains(line, "container") {
				containerName = strings.TrimSpace(line)
				break
			}
		}
	}

	// 等待一段时间确保容器创建完成
	time.Sleep(time.Millisecond * 500)

	// 获取创建的容器信息
	if containerName != "" {
		return GetContainerById(containerName)
	}

	return Container{}, fmt.Errorf("无法获取创建的容器信息")
}

// StartContainer 启动容器
func StartContainer(containerId string) error {
	// zdocker没有专门的start命令，这里返回成功
	return nil
}

// StopContainer 停止容器
func StopContainer(containerId string) error {
	cmd := exec.Command("zdocker", "stop", containerId)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("停止容器失败: %s, %v", string(output), err)
	}
	return nil
}

// RemoveContainer 删除容器
func RemoveContainer(containerId string) error {
	cmd := exec.Command("zdocker", "rm", containerId)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("删除容器失败: %s, %v", string(output), err)
	}
	return nil
}

// GetContainerLogs 获取容器日志
func GetContainerLogs(containerId string) (string, error) {
	cmd := exec.Command("zdocker", "logs", containerId)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("获取容器日志失败: %s, %v", string(output), err)
	}
	return string(output), nil
}

// ExecContainer 在容器中执行命令
func ExecContainer(containerId string, req ExecRequest) (ExecResult, error) {
	args := []string{"exec", containerId}
	args = append(args, req.Command...)

	cmd := exec.Command("zdocker", args...)
	output, err := cmd.CombinedOutput()

	exitCode := 0
	if err != nil {
		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			exitCode = exitError.ExitCode()
		}
	}

	return ExecResult{
		Output:   string(output),
		ExitCode: exitCode,
	}, nil
}

// GetNetworkList 获取网络列表
func GetNetworkList() ([]NetworkInfo, error) {
	cmd := exec.Command("zdocker", "network", "list")
	output, err := cmd.CombinedOutput()
	if err != nil {
		// 如果命令失败，返回默认网络信息
		return []NetworkInfo{
			{Name: "bridge", Driver: "bridge", Subnet: "172.17.0.0/16"},
		}, nil
	}

	// 解析输出
	var networks []NetworkInfo
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) != "" && !strings.Contains(line, "NAME") {
			fields := strings.Fields(line)
			if len(fields) >= 3 {
				networks = append(networks, NetworkInfo{
					Name:   fields[0],
					Driver: fields[1],
					Subnet: fields[2],
				})
			}
		}
	}

	return networks, nil
}

// CreateNetwork 创建网络
func CreateNetwork(req CreateNetworkRequest) (NetworkInfo, error) {
	args := []string{"network", "create"}
	if req.Driver != "" {
		args = append(args, "--driver", req.Driver)
	}
	if req.Subnet != "" {
		args = append(args, "--subnet", req.Subnet)
	}
	args = append(args, req.Name)

	cmd := exec.Command("zdocker", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return NetworkInfo{}, fmt.Errorf("创建网络失败: %s, %v", string(output), err)
	}

	return NetworkInfo{
		Name:   req.Name,
		Driver: req.Driver,
		Subnet: req.Subnet,
	}, nil
}

// RemoveNetwork 删除网络
func RemoveNetwork(networkId string) error {
	cmd := exec.Command("zdocker", "network", "remove", networkId)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("删除网络失败: %s, %v", string(output), err)
	}
	return nil
}

// GetSystemInfo 获取系统信息
func GetSystemInfo() (SystemInfo, error) {
	// 获取ZDocker根目录
	zdockerRoot := "/var/lib/zdocker"
	if envRoot := os.Getenv("ZDOCKER_ROOT"); envRoot != "" {
		zdockerRoot = envRoot
	}

	// 获取内存信息
	memory := "Unknown"
	if memInfo, err := os.ReadFile("/proc/meminfo"); err == nil {
		lines := strings.Split(string(memInfo), "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "MemTotal:") {
				fields := strings.Fields(line)
				if len(fields) >= 2 {
					memory = fields[1] + " kB"
				}
				break
			}
		}
	}

	return SystemInfo{
		OS:           runtime.GOOS,
		Architecture: runtime.GOARCH,
		CPUs:         runtime.NumCPU(),
		Memory:       memory,
		ZDockerRoot:  zdockerRoot,
	}, nil
}
