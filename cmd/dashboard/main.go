/*


Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"context"
	"flag"
	"os"

	"github.com/eddycharly/kloops/apis/config/v1alpha1"
	_ "github.com/eddycharly/kloops/pkg/chatbot/pluginimports"
	"github.com/eddycharly/kloops/pkg/dashboard/server"
	"github.com/eddycharly/kloops/pkg/utils"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	_ = clientgoscheme.AddToScheme(scheme)
	_ = v1alpha1.AddToScheme(scheme)
}

func main() {
	var namespace string
	var metricsAddr string
	var enableLeaderElection bool
	flag.StringVar(&namespace, "namespace", "default", "The namespace to watch.")
	flag.StringVar(&metricsAddr, "metrics-addr", ":8080", "The address the metric endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "enable-leader-election", false, "Enable leader election for controller manager. Enabling this will ensure there is only one active controller manager.")
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseDevMode(true)))

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:             scheme,
		Namespace:          namespace,
		MetricsBindAddress: metricsAddr,
		Port:               9443,
		LeaderElection:     enableLeaderElection,
		LeaderElectionID:   "1e9e8f6c.kloops.io",
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	// Define all broadcasters/channels
	// Keep broadcaster channels open indefinitely
	var resourcesChannel = make(chan utils.SocketData)

	var resourcesBroadcaster = utils.NewBroadcaster(resourcesChannel)

	cache := mgr.GetCache()
	informer, err := cache.GetInformer(context.TODO(), &v1alpha1.RepoConfig{})
	if err != nil {
		setupLog.Error(err, "unable to start controller")
		os.Exit(1)
	}

	utils.NewController(
		resourcesChannel,
		informer,
		"created",
		"updated",
		"deleted",
	)

	stopCh := ctrl.SetupSignalHandler()

	server := server.NewServer(namespace, mgr.GetConfig(), mgr.GetClient(), resourcesBroadcaster, ctrl.Log)

	go func() {
		if err := server.Start("", 8090); err != nil {
			setupLog.Error(err, "problem starting dashboard")
			os.Exit(1)
		}
	}()

	setupLog.Info("starting manager")
	if err := mgr.Start(stopCh); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}
