package main
import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/vadxq/go-rest-starter/internal/apps/app/bootstrap"
	"github.com/vadxq/go-rest-starter/pkg/logger"
)

func main() {
	var appLogger logger.Logger = logger.Default()

	// สร้าง instance ของแอปพลิเคชัน
	application, err := bootstrap.New()
	if err != nil {
		appLogger.Error("สร้างแอปพลิเคชันล้มเหลว", "error", err)
		os.Exit(1)
	}
	appLogger = application.Logger()

	// เริ่มต้นเซิร์ฟเวอร์ HTTP
	serverErrCh := application.StartServer()

	// รอสัญญาณหรือข้อผิดพลาดจากเซิร์ฟเวอร์
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	select {
	case err := <-serverErrCh:
		appLogger.Error("ข้อผิดพลาดของเซิร์ฟเวอร์", "error", err)
	case sig := <-signalCh:
		appLogger.Info("ได้รับสัญญาณจากระบบ เริ่มการปิดอย่างนุ่มนวล", "signal", sig.String())
	}

	// ปิดแอปพลิเคชันอย่างนุ่มนวล
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := application.Shutdown(ctx); err != nil {
		appLogger.Error("ปิดแอปพลิเคชันล้มเหลว", "error", err)
		os.Exit(1)
	}
}
