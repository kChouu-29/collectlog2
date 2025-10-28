package main

import (
	"collectlogupdate/controller" // <--- ĐÃ SỬA: Bổ sung import controller
	"collectlogupdate/server"
)

func main() {
	// <--- ĐÃ SỬA: Khởi động worker pool với 4 worker (số lượng hợp lý cho dev/test)
	// Bạn có thể tăng số lượng này (ví dụ: 1.5 đến 2 lần số Core CPU) cho Production.
	controller.StartWorkerPool(4)

	server.StartServer()
}
