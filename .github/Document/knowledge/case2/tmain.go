package main

import (
	"fmt"
	"sync"
	"time"
)

type Worker struct {
	id       int
	taskCh   chan int
	stopCh   chan struct{}
	disabled bool
}

func (w *Worker) Start(wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()

		// ตัวแปร channel ที่อาจเป็น nil (disabled)
		var activeTaskCh chan int = w.taskCh

		for {
			select {
			case task, ok := <-activeTaskCh:
				if !ok {
					fmt.Printf("Worker %d: task channel closed\n", w.id)
					return
				}
				fmt.Printf("Worker %d: กำลังทำงาน %d\n", w.id, task)
				time.Sleep(500 * time.Millisecond)

			case <-w.stopCh:
				// ปิด worker แบบถาวร
				fmt.Printf("Worker %d: หยุดทำงาน\n", w.id)
				return

			case <-time.After(2 * time.Second):
				// ถ้า activeTaskCh == nil จะมาที่นี่ทุก 2 วินาที
				if activeTaskCh == nil {
					fmt.Printf("Worker %d: ถูก disable อยู่ กำลังรอ enable...\n", w.id)
				} else {
					fmt.Printf("Worker %d: idle\n", w.id)
				}
			}
		}
	}()
}

func (w *Worker) Enable() {
	w.disabled = false
	// เปลี่ยนเป็นใช้ task channel จริง
	// จริงๆต้องมีวิธีเข้าถึงตัวแปร activeTaskCh ข้างใน goroutine
	// ซึ่งต้องใช้ channel สั่ง หรือ redesign
	// วิธีง่าย: ส่ง signal ไปอีก channel เพื่อบอกให้เปลี่ยน activeTaskCh
}

func main() {
	taskCh := make(chan int, 10)
	stopCh := make(chan struct{})

	// สร้าง worker
	go func() {
		// ตัวแปร active channel สามารถเปลี่ยนค่า nil/non-nil ได้
		var activeCh chan int = taskCh // เริ่มต้น active

		for {
			select {
			case task, ok := <-activeCh:
				if !ok {
					return
				}
				fmt.Printf("Worker: ทำงาน %d\n", task)

			case <-stopCh:
				return

			case <-time.After(3 * time.Second):
				// toggle activeCh ทุก 3 วินาที (disable/enable)
				if activeCh == nil {
					activeCh = taskCh
					fmt.Println("🟢 Enable worker")
				} else {
					activeCh = nil
					fmt.Println("🔴 Disable worker (ไม่รับงานใหม่)")
				}
			}
		}
	}()

	// ส่งงาน
	go func() {
		for i := 1; i <= 20; i++ {
			taskCh <- i
			time.Sleep(500 * time.Millisecond)
		}
		close(taskCh)
	}()

	time.Sleep(15 * time.Second)
	close(stopCh)
}
