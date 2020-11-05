package handler

import (
	api_model "github.com/goodrain/rainbond/api/model"
	"github.com/goodrain/rainbond/api/util/bcode"
	"github.com/goodrain/rainbond/db"
	dbmodel "github.com/goodrain/rainbond/db/model"
	"github.com/jinzhu/gorm"
)

//UpdateServiceMonitor update service monitor
func (s *ServiceAction) UpdateServiceMonitor(tenantID, serviceID, name string, update api_model.UpdateServiceMonitorRequestStruct) (*dbmodel.TenantServiceMonitor, error) {
	sm, err := db.GetManager().TenantServiceMonitorDao().GetByName(serviceID, name)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, bcode.ErrServiceMonitorNotFound
		}
		return nil, err
	}
	_, err = db.GetManager().TenantServicesPortDao().GetPort(serviceID, update.Port)
	if err != nil {
		return nil, bcode.ErrPortNotFound
	}
	sm.ServiceShowName = update.ServiceShowName
	sm.Port = update.Port
	sm.Path = update.Path
	sm.Interval = update.Interval
	return sm, db.GetManager().TenantServiceMonitorDao().UpdateModel(sm)
}

//DeleteServiceMonitor delete
func (s *ServiceAction) DeleteServiceMonitor(tenantID, serviceID, name string) (*dbmodel.TenantServiceMonitor, error) {
	sm, err := db.GetManager().TenantServiceMonitorDao().GetByName(serviceID, name)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, bcode.ErrServiceMonitorNotFound
		}
		return nil, err
	}
	return sm, db.GetManager().TenantServiceMonitorDao().DeleteServiceMonitor(sm)
}

//AddServiceMonitor add service monitor
func (s *ServiceAction) AddServiceMonitor(tenantID, serviceID string, add api_model.AddServiceMonitorRequestStruct) (*dbmodel.TenantServiceMonitor, error) {
	_, err := db.GetManager().TenantServicesPortDao().GetPort(serviceID, add.Port)
	if err != nil {
		return nil, bcode.ErrPortNotFound
	}
	sm := dbmodel.TenantServiceMonitor{
		Name:            add.Name,
		TenantID:        tenantID,
		ServiceID:       serviceID,
		ServiceShowName: add.ServiceShowName,
		Port:            add.Port,
		Path:            add.Path,
		Interval:        add.Interval,
	}
	return &sm, db.GetManager().TenantServiceMonitorDao().AddModel(&sm)
}
