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
	"context"
	"fmt"
	"net"
	"time"

	natpmp "github.com/jackpal/go-nat-pmp"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	networkv1 "github.com/jkoelker/natpmp-controller/api/v1"
)

const (
	three = 3
	four  = 4
)

func renewIn(lifetime int) time.Duration {
	return three * (time.Duration(lifetime) * time.Second) / four
}

// NatPMPReconciler reconciles a NatPMP object.
type NatPMPReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=network.natpmp.jkoelker.github.io,resources=natpmps,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=network.natpmp.jkoelker.github.io,resources=natpmps/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=network.natpmp.jkoelker.github.io,resources=natpmps/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (reconciler *NatPMPReconciler) Reconcile(
	ctx context.Context,
	req ctrl.Request,
) (ctrl.Result, error) {
	start := time.Now()

	var natpmpCR networkv1.NatPMP
	if err := reconciler.Get(ctx, req.NamespacedName, &natpmpCR); err != nil {
		if errors.IsNotFound(err) {
			Info(
				ctx,
				"NatPMP resource not found. Ignoring since object must be deleted",
				"namespace", req.NamespacedName.Namespace,
				"name", req.NamespacedName.Name,
			)

			return ctrl.Result{}, nil
		}

		return ctrl.Result{}, WrapError(ctx, err, "unable to fetch NatPMP")
	}

	gateway, protocol, errs := ValidateNatPMP(natpmpCR)
	if len(errs) > 0 {
		gvk := networkv1.GroupVersionKind()
		err := errors.NewInvalid(gvk.GroupKind(), natpmpCR.Name, errs)

		return ctrl.Result{}, WrapError(ctx, err, "invalid NatPMP")
	}

	client := natpmp.NewClient(gateway)

	external, err := client.GetExternalAddress()
	if err != nil {
		return ctrl.Result{}, WrapError(ctx, err, "unable to get external IP")
	}

	response, err := client.AddPortMapping(
		protocol,
		natpmpCR.Spec.InternalPort,
		natpmpCR.Spec.ExternalPort,
		natpmpCR.Spec.Lifetime,
	)
	if err != nil {
		return ctrl.Result{}, WrapError(ctx, err, "unable to add port mapping")
	}

	externalIP := net.IP(external.ExternalIPAddress[:])
	mappedLifetime := int(response.PortMappingLifetimeInSeconds)

	natpmpCR.Status.ExternalIP = externalIP.String()
	natpmpCR.Status.MappedExternalPort = int(response.MappedExternalPort)
	natpmpCR.Status.MappedInternalPort = int(response.InternalPort)
	natpmpCR.Status.MappedLifetime = mappedLifetime
	natpmpCR.Status.SecondsSinceStartOfEpoch = int(response.SecondsSinceStartOfEpoc)

	if err := reconciler.Status().Update(ctx, &natpmpCR); err != nil {
		return ctrl.Result{}, WrapError(ctx, err, "unable to update NatPMP status")
	}

	if err := reconciler.ApplyTemplates(ctx, natpmpCR); err != nil {
		return ctrl.Result{}, WrapError(ctx, err, "unable to apply templates")
	}

	// Renew the port mapping 3/4 of the way through the lifetime. Taking into
	// account the time it took to get here.
	renewAfter := renewIn(mappedLifetime) - time.Since(start)
	if renewAfter < 0 {
		renewAfter = 0
	}

	return ctrl.Result{
		RequeueAfter: renewAfter,
	}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (reconciler *NatPMPReconciler) SetupWithManager(mgr ctrl.Manager) error {
	err := ctrl.NewControllerManagedBy(mgr).
		For(&networkv1.NatPMP{}).
		Complete(reconciler)
	if err != nil {
		return fmt.Errorf("unable to complete controller: %w", err)
	}

	return nil
}

// ApplyTemplates applies the templates from the NatPMP CR to the cluster.
func (reconciler *NatPMPReconciler) ApplyTemplates(
	ctx context.Context,
	natpmpCR networkv1.NatPMP,
) error {
	objects, err := ProcessTemplates(natpmpCR)
	if err != nil {
		return WrapError(ctx, err, "unable to process templates")
	}

	for idx := range objects {
		object := objects[idx]

		if err := ctrl.SetControllerReference(&natpmpCR, object, reconciler.Scheme); err != nil {
			return WrapError(ctx, err, "unable to set controller reference")
		}

		opts := []client.PatchOption{
			client.ForceOwnership,
			client.FieldOwner("natpmp-controller"),
		}
		if err := reconciler.Patch(ctx, object, client.Apply, opts...); err != nil {
			return WrapError(ctx, err, "unable to apply templates")
		}
	}

	return nil
}
