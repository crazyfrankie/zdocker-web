package service

import (
	"errors"
	"fmt"
	"github.com/bytedance/sonic"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"syscall"
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
	dirUrl := fmt.Sprintf(container.DefaultLocation, "")
	dirUrl = dirUrl[:len(dirUrl)-1]

	files, err := os.ReadDir(dirUrl)
	if err != nil {
		return nil, err
	}

	containers := make([]Container, 0, len(files))
	for _, f := range files {
		info, err := getContainerInfo(f)
		if err != nil {
			continue
		}

		// Check if container process is still running and update status if needed
		if info.Status == container.RUNNING && info.PID != "" {
			if !isProcessRunning(info.PID) {
				if err := updateContainerStatusToExit(info.Name); err != nil {
				} else {
					info.Status = container.EXIT
					info.PID = ""
				}
			}
		}

		containers = append(containers, Container{
			ID:          info.ID,
			Name:        info.Name,
			Command:     info.Command,
			Status:      info.Status,
			CreatedTime: info.CreateTime,
			Pid:         info.PID,
			Volume:      info.Volume,
		})
	}

	return containers, nil
}

// getContainerInfoByName 根据容器名称获取容器信息
func getContainerInfo(file os.DirEntry) (*container.ContainerInfo, error) {
	var info container.ContainerInfo
	fileName := file.Name()
	cfgDir := fmt.Sprintf(container.DefaultLocation, fileName)
	cfgName := cfgDir + container.ConfigName
	data, err := os.ReadFile(cfgName)
	if err != nil {
		return nil, err
	}
	err = sonic.Unmarshal(data, &info)
	if err != nil {
		return nil, err
	}

	return &info, nil
}

// isProcessRunning checks if a process with the given PID is still running
func isProcessRunning(pidStr string) bool {
	if pidStr == "" {
		return false
	}

	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		return false
	}

	// Send signal 0 to check if process exists without affecting it
	if err := syscall.Kill(pid, 0); err != nil {
		// ESRCH means "No such process"
		if errors.Is(err, syscall.ESRCH) {
			return false
		}
		// Other errors might mean the process exists, but we don't have permission
		// In this case, we assume the process is running
		return true
	}
	return true
}

// updateContainerStatusToExit updates container status to EXIT and clears PID
func updateContainerStatusToExit(containerName string) error {
	dirUrl := fmt.Sprintf(container.DefaultLocation, containerName)
	cfgFile := dirUrl + container.ConfigName

	// Read current container info
	content, err := os.ReadFile(cfgFile)
	if err != nil {
		return fmt.Errorf("read container config error: %v", err)
	}

	var info container.ContainerInfo
	if err := sonic.Unmarshal(content, &info); err != nil {
		return fmt.Errorf("unmarshal container info error: %v", err)
	}

	// Update status and clear PID
	info.Status = container.EXIT
	info.PID = ""

	// Write back to file
	newContent, err := sonic.Marshal(info)
	if err != nil {
		return fmt.Errorf("marshal container info error: %v", err)
	}

	if err := os.WriteFile(cfgFile, newContent, 0622); err != nil {
		return fmt.Errorf("write container config error: %v", err)
	}

	return nil
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
func StopContainer(containerName string) error {
	cmd := exec.Command("zdocker", "stop", containerName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("停止容器失败: %s, %v", string(output), err)
	}
	return nil
}

// RemoveContainer 删除容器
func RemoveContainer(containerName string) error {
	cmd := exec.Command("zdocker", "rm", containerName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("删除容器失败: %s, %v", string(output), err)
	}
	return nil
}

// GetContainerLogs 获取容器日志
func GetContainerLogs(containerName string) (string, error) {
	cmd := exec.Command("zdocker", "logs", containerName)
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
