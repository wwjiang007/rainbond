// Copyright (C) 2014-2018 Goodrain Co., Ltd.
// RAINBOND, Application Management Platform

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version. For any non-GPL usage of Rainbond,
// one or multiple Commercial Licenses authorized by Goodrain Co., Ltd.
// must be obtained first.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package server

import (
	"context"
	"fmt"
	"github.com/goodrain/rainbond/discover.v2"
	"github.com/goodrain/rainbond/node/initiate"
	"github.com/goodrain/rainbond/util/constants"
	"k8s.io/client-go/kubernetes"
	"os"
	"os/signal"
	"syscall"

	"github.com/goodrain/rainbond/cmd/node/option"
	eventLog "github.com/goodrain/rainbond/event"
	"github.com/goodrain/rainbond/node/api"
	"github.com/goodrain/rainbond/node/api/controller"
	"github.com/goodrain/rainbond/node/core/store"
	"github.com/goodrain/rainbond/node/kubecache"
	"github.com/goodrain/rainbond/node/masterserver"
	"github.com/goodrain/rainbond/node/nodem"
	"github.com/goodrain/rainbond/node/nodem/docker"
	"github.com/goodrain/rainbond/node/nodem/envoy"
	etcdutil "github.com/goodrain/rainbond/util/etcd"
	k8sutil "github.com/goodrain/rainbond/util/k8s"

	"github.com/sirupsen/logrus"
)

//Run start run
func Run(cfg *option.Conf) error {
	var stoped = make(chan struct{})
	stopfunc := func() error {
		close(stoped)
		return nil
	}
	startfunc := func() error {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		etcdClientArgs := &etcdutil.ClientArgs{
			Endpoints:   cfg.EtcdEndpoints,
			CaFile:      cfg.EtcdCaFile,
			CertFile:    cfg.EtcdCertFile,
			KeyFile:     cfg.EtcdKeyFile,
			DialTimeout: cfg.EtcdDialTimeout,
		}
		if err := cfg.ParseClient(ctx, etcdClientArgs); err != nil {
			return fmt.Errorf("config parse error:%s", err.Error())
		}

		config, err := k8sutil.NewRestConfig(cfg.K8SConfPath)
		if err != nil {
			return err
		}
		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			return err
		}

		k8sDiscover := discover.NewK8sDiscover(ctx, clientset, cfg)
		defer k8sDiscover.Stop()

		nodemanager, err := nodem.NewNodeManager(ctx, cfg)
		if err != nil {
			return fmt.Errorf("create node manager failed: %s", err)
		}
		if err := nodemanager.InitStart(); err != nil {
			return err
		}

		err = eventLog.NewManager(eventLog.EventConfig{
			EventLogServers: cfg.EventLogServer,
			DiscoverArgs:    etcdClientArgs,
		})
		if err != nil {
			logrus.Errorf("error creating eventlog manager")
			return nil
		}
		defer eventLog.CloseManager()
		logrus.Debug("create and start event log client success")

		kubecli, err := kubecache.NewKubeClient(cfg, clientset)
		if err != nil {
			return err
		}
		defer kubecli.Stop()

		if cfg.ImageRepositoryHost == constants.DefImageRepository {
			hostManager, err := initiate.NewHostManager(cfg, k8sDiscover)
			if err != nil {
				return fmt.Errorf("create new host manager: %v", err)
			}
			hostManager.Start()
		}

		logrus.Debugf("rbd-namespace=%s; rbd-docker-secret=%s", os.Getenv("RBD_NAMESPACE"), os.Getenv("RBD_DOCKER_SECRET"))
		// sync docker inscure registries cert info into all rainbond node
		if err = docker.SyncDockerCertFromSecret(clientset, os.Getenv("RBD_NAMESPACE"), os.Getenv("RBD_DOCKER_SECRET")); err != nil { // TODO fanyangyang namespace secretname
			return fmt.Errorf("sync docker cert from secret error: %s", err.Error())
		}

		// init etcd client
		if err = store.NewClient(ctx, cfg, etcdClientArgs); err != nil {
			return fmt.Errorf("Connect to ETCD %s failed: %s", cfg.EtcdEndpoints, err)
		}
		errChan := make(chan error, 3)
		if err := nodemanager.Start(errChan); err != nil {
			return fmt.Errorf("start node manager failed: %s", err)
		}
		defer nodemanager.Stop()
		logrus.Debug("create and start node manager moudle success")

		//master服务在node服务之后启动
		var ms *masterserver.MasterServer
		if cfg.RunMode == "master" {
			ms, err = masterserver.NewMasterServer(nodemanager.GetCurrentNode(), kubecli)
			if err != nil {
				logrus.Errorf(err.Error())
				return err
			}
			ms.Cluster.UpdateNode(nodemanager.GetCurrentNode())
			if err := ms.Start(errChan); err != nil {
				logrus.Errorf(err.Error())
				return err
			}
			defer ms.Stop(nil)
			logrus.Debug("create and start master server moudle success")
		}
		//create api manager
		apiManager := api.NewManager(*cfg, nodemanager.GetCurrentNode(), ms, kubecli)
		if err := apiManager.Start(errChan); err != nil {
			return err
		}
		if err := nodemanager.AddAPIManager(apiManager); err != nil {
			return err
		}
		defer apiManager.Stop()

		//create service mesh controller
		grpcserver, err := envoy.CreateDiscoverServerManager(clientset, *cfg)
		if err != nil {
			return err
		}
		if err := grpcserver.Start(errChan); err != nil {
			return err
		}
		defer grpcserver.Stop()

		logrus.Debug("create and start api server moudle success")

		defer controller.Exist(nil)
		//step finally: listen Signal
		term := make(chan os.Signal)
		signal.Notify(term, os.Interrupt, syscall.SIGTERM)
		select {
		case <-stoped:
			logrus.Infof("windows service stoped..")
		case <-term:
			logrus.Warn("Received SIGTERM, exiting gracefully...")
		case err := <-errChan:
			logrus.Errorf("Received a error %s, exiting gracefully...", err.Error())
		}
		logrus.Info("See you next time!")
		return nil
	}
	err := initService(cfg, startfunc, stopfunc)
	if err != nil {
		return err
	}
	return nil
}
