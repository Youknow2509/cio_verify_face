package health

import (
	"context"
	"fmt"
	"os"
	"sync"
	"syscall"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/process"
	domainHealth "github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/domain/health"
	"github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/domain/model"
	"github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/global"
	wsCore "github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/interfaces/ws/core"
)

// Impl health check
type HealthCheck struct {
}

// CheckDownstreamServices implements health.IHealthCheck.
func (h *HealthCheck) CheckDownstreamServices(ctx context.Context) *model.ComponentCheck {
	return nil
}

// CheckSystemResource implements health.IHealthCheck.
func (h *HealthCheck) CheckSystemResource(ctx context.Context) *model.ComponentCheck {
	var wg sync.WaitGroup
	// Sử dụng WaitGroup để chờ tất cả các goroutine hoàn thành.
	wg.Add(3)

	// Các biến để lưu kết quả từ các goroutine.
	// Sử dụng con trỏ để có thể phân biệt giữa giá trị rỗng và kết quả thực sự.
	var cpuCheck, memCheck, fdCheck *model.ComponentCheck

	// 1. Chạy kiểm tra CPU trong một goroutine riêng.
	go func() {
		defer wg.Done() // Báo cho WaitGroup biết goroutine này đã xong.
		cpuCheck = getCpuUsage(ctx)
	}()

	// 2. Chạy kiểm tra Memory trong một goroutine riêng.
	go func() {
		defer wg.Done()
		memCheck = getMemoryUsage(ctx)
	}()

	// 3. Chạy kiểm tra File Descriptors trong một goroutine riêng.
	go func() {
		defer wg.Done()
		fdCheck = getFileDescriptorUsage(ctx)
	}()

	// Chờ tất cả các kiểm tra song song hoàn tất.
	wg.Wait()

	// Tổng hợp trạng thái cuối cùng từ các kiểm tra con.
	// Đây là một quy tắc nghiệp vụ quan trọng: trạng thái tổng thể sẽ là trạng thái "tệ nhất"
	// trong số các thành phần con.
	overallStatus := aggregateStatus(cpuCheck.Status, memCheck.Status, fdCheck.Status)

	return &model.ComponentCheck{
		Status: overallStatus,
		Details: map[string]interface{}{
			"cpu":             cpuCheck,
			"memory":          memCheck,
			"fileDescriptors": fdCheck,
		},
	}
}

// CheckWebSocketServer implements health.IHealthCheck.
func (h *HealthCheck) CheckWebSocketServer(ctx context.Context) *model.ComponentCheck {
	// 1. Kiểm tra số lượng kết nối WebSocket đang hoạt động.
	activeConnections := wsCore.GetHub().NumClients()
	// 2. Lấy số lượng kết nối tối đa được phép.
	maxConnections := global.ServerWsSetting.MaxConnectionSystem
	// 3. Tính toán và quyết định trạng thái dựa trên ngưỡng
	usagePercent := (float64(activeConnections) / float64(maxConnections)) * 100.0
	status := model.StatusUp
	if usagePercent >= global.ServerSetting.OutOfServiceThreshold {
		status = model.StatusOutOfService // Quá tải, không nên nhận thêm kết nối
	} else if usagePercent >= global.ServerSetting.DegradedThreshold {
		status = model.StatusDegraded // Gần quá tải, cần theo dõi
	}
	return &model.ComponentCheck{
		Status: status,
		Details: map[string]interface{}{
			"activeConnections": activeConnections,
			"maxConnections":    maxConnections,
			"connectionUsage":   fmt.Sprintf("%.2f%%", usagePercent),
		},
	}
}

// New Health check
func NewHealthCheck() domainHealth.IHealthCheck {
	return &HealthCheck{}
}

// =================================
// 	Helper check
// =================================
/**
 * getCpuUsage lấy thông tin sử dụng CPU.
 */
func getCpuUsage(ctx context.Context) *model.ComponentCheck {
	// Sử dụng gopsutil để lấy phần trăm CPU sử dụng.
	// Tham số đầu tiên `0` nghĩa là so sánh với lần gọi trước, `false` là tính cho tất cả các core.
	percent, err := cpu.PercentWithContext(ctx, 0, false)
	if err != nil || len(percent) == 0 {
		return &model.ComponentCheck{
			Status:  model.StatusDown,
			Details: map[string]interface{}{"error": "could not get cpu usage"},
		}
	}

	p := percent[0]
	status := model.StatusUp
	// Ngưỡng cảnh báo: trên 90% là quá tải, trên 75% là suy giảm.
	// Các con số này nên được cấu hình từ bên ngoài thay vì hard-code.
	if p > 90.0 {
		status = model.StatusOutOfService
	} else if p > 75.0 {
		status = model.StatusDegraded
	}

	return &model.ComponentCheck{
		Status: status,
		Details: map[string]interface{}{
			"usage": fmt.Sprintf("%.2f%%", p),
		},
	}
}

/**
 * getMemoryUsage lấy thông tin sử dụng bộ nhớ.
 */
func getMemoryUsage(ctx context.Context) *model.ComponentCheck {
	vm, err := mem.VirtualMemoryWithContext(ctx)
	if err != nil {
		return &model.ComponentCheck{
			Status:  model.StatusDown,
			Details: map[string]interface{}{"error": "could not get memory usage"},
		}
	}

	status := model.StatusUp
	// Ngưỡng cảnh báo: trên 90% là quá tải, trên 80% là suy giảm.
	if vm.UsedPercent > 90.0 {
		status = model.StatusOutOfService
	} else if vm.UsedPercent > 80.0 {
		status = model.StatusDegraded
	}

	return &model.ComponentCheck{
		Status: status,
		Details: map[string]interface{}{
			"used":  fmt.Sprintf("%d MB", vm.Used/1024/1024),
			"total": fmt.Sprintf("%d MB", vm.Total/1024/1024),
			"usage": fmt.Sprintf("%.2f%%", vm.UsedPercent),
		},
	}
}

/**
 * getFileDescriptorUsage lấy thông tin sử dụng File Descriptors.
 * Đây là chỉ số SỐNG CÒN đối với một WebSocket Gateway.
 */
func getFileDescriptorUsage(ctx context.Context) *model.ComponentCheck {
	// Lấy process hiện tại.
	p, err := process.NewProcessWithContext(ctx, int32(os.Getpid()))
	if err != nil {
		return &model.ComponentCheck{
			Status:  model.StatusDown,
			Details: map[string]interface{}{"error": "could not get current process"},
		}
	}

	// Lấy số lượng FDs đang được sử dụng bởi process này.
	numFDs, err := p.NumFDsWithContext(ctx)
	if err != nil {
		return &model.ComponentCheck{
			Status:  model.StatusDown,
			Details: map[string]interface{}{"error": "could not get number of file descriptors"},
		}
	}

	// Lấy giới hạn FDs của hệ điều hành cho process (ulimit -n).
	// Lưu ý: syscall là đặc thù cho từng HĐH. Đoạn mã này hoạt động trên Linux/macOS.
	var rlim syscall.Rlimit
	err = syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rlim)
	if err != nil {
		return &model.ComponentCheck{
			Status:  model.StatusDown,
			Details: map[string]interface{}{"error": "could not get file descriptor limit"},
		}
	}

	limit := rlim.Cur
	usage := float64(numFDs) / float64(limit)
	status := model.StatusUp
	// Ngưỡng cho FDs nên chặt chẽ hơn vì chạm đỉnh là sập ngay.
	if usage > 0.95 {
		status = model.StatusOutOfService
	} else if usage > 0.85 {
		status = model.StatusDegraded
	}

	return &model.ComponentCheck{
		Status: status,
		Details: map[string]interface{}{
			"used":  numFDs,
			"limit": limit,
			"usage": fmt.Sprintf("%.2f%%", usage*100),
		},
	}
}

/**
 * getFileDescriptorUsage lấy thông tin sử dụng File Descriptors.
 * Đây là chỉ số SỐNG CÒN đối với một WebSocket Gateway.
 */
func aggregateStatus(statuses ...model.StatusType) model.StatusType {
	finalStatus := model.StatusUp
	for _, s := range statuses {
		if s == model.StatusDown {
			return model.StatusDown
		}
		if s == model.StatusOutOfService {
			finalStatus = model.StatusOutOfService
		}
		if s == model.StatusDegraded && finalStatus != model.StatusOutOfService {
			finalStatus = model.StatusDegraded
		}
	}
	return finalStatus
}
