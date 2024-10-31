/*
Copyright 2024.

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
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	crv1 "undercloud.kr/nova-custom-operator/api/v1"

	"time"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// VirtualMachineHostReconciler reconciles a VirtualMachineHost object
type VirtualMachineHostReconciler struct {
	client.Client
	Scheme      *runtime.Scheme
	prev_status string // added
}

// instance status struct
type Node struct {
	Time  string `json:"time"`
	State string `json:"state"`
}

// +kubebuilder:rbac:groups=cr.undercloud.kr,resources=virtualmachinehosts,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=cr.undercloud.kr,resources=virtualmachinehosts/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=cr.undercloud.kr,resources=virtualmachinehosts/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the VirtualMachineHost object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.19.0/pkg/reconcile
func (r *VirtualMachineHostReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger = log.FromContext(ctx)

	url := "http://127.0.0.1:8000/v1/nodes"
	resp, err := http.Get(url)
	if err != nil {
		logger.info("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.info("Request failed with status code: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.info("Failed to read response body: %v", err)
	}

	var node Node // for single instance
	if err := json.Unmarshal(body, &node); err != nil {
		logger.info("Failed to unmarshal JSON response: %v", err)
	}

	now_status := node.State

	if now_status != r.prev_status {
		logger.info("Successful to find converted status into: %v", now_status)
	}

	// something logic (ex. re-inspection, re-scrap, and so on)

	r.prev_status = now_status

	return ctrl.Result{RequeueAfter: time.Second}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *VirtualMachineHostReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&crv1.VirtualMachineHost{}).
		Named("virtualmachinehost").
		Complete(r)
}
