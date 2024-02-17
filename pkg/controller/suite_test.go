/*
Copyright 2024 Jason Kölker.

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

package controller

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	testifySuite "github.com/stretchr/testify/suite"

	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	networkv1 "github.com/jkoelker/natpmp-controller/api/v1"
)

// Suite is the base test suite for the controller tests.
type Suite struct {
	testifySuite.Suite

	cfg       *rest.Config
	k8sClient client.Client
	testEnv   *envtest.Environment
}

// SetupSuite sets up the test suite.
func (suite *Suite) SetupSuite() {
	logf.SetLogger(zap.New(zap.UseDevMode(true)))

	suite.testEnv = &envtest.Environment{
		CRDDirectoryPaths:     []string{filepath.Join("..", "..", "config", "crd", "bases")},
		ErrorIfCRDPathMissing: true,

		// The BinaryAssetsDirectory is only required if you want to run the tests directly
		// without call the makefile target test. If not informed it will look for the
		// default path defined in controller-runtime which is /usr/local/kubebuilder/.
		// Note that you must have the required binaries setup under the bin directory to perform
		// the tests directly. When we run make test it will be setup and used automatically.
		BinaryAssetsDirectory: filepath.Join("..", "..", "bin", "k8s",
			fmt.Sprintf("1.28.3-%s-%s", runtime.GOOS, runtime.GOARCH)),
	}

	cfg, err := suite.testEnv.Start()
	suite.Require().NoError(err)
	suite.Require().NotNil(cfg)
	suite.cfg = cfg

	err = networkv1.AddToScheme(scheme.Scheme)
	suite.Require().NoError(err)

	//+kubebuilder:scaffold:scheme

	k8sClient, err := client.New(cfg, client.Options{Scheme: scheme.Scheme})
	suite.Require().NoError(err)
	suite.Require().NotNil(k8sClient)
	suite.k8sClient = k8sClient
}

// TearDownSuite tears down the test suite.
func (suite *Suite) TearDownSuite() {
	suite.Require().NoError(suite.testEnv.Stop())
}

// TestControllers runs the controller tests.
func TestSuiteControllers(t *testing.T) {
	testifySuite.Run(t, new(Suite))
}
