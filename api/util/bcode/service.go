package bcode

// service: 10000~10099
var (
	//ErrPortNotFound -
	ErrPortNotFound = newByMessage(404, 10001, "service port not found")
	//ErrServiceMonitorNotFound -
	ErrServiceMonitorNotFound = newByMessage(404, 10101, "service monitor not found")
	//ErrServiceMonitorNameExist -
	ErrServiceMonitorNameExist = newByMessage(400, 10102, "service monitor name is exist")
)
