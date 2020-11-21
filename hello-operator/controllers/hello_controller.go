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

package controllers

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/prometheus/common/log"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	hellogroupv1alpha1 "github.com/thisisdavidbell/hello-ocp/1-create-operator/api/v1alpha1"
)

// HelloReconciler reconciles a Hello object
type HelloReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=hellogroup.my.domain,resources=hellos,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=hellogroup.my.domain,resources=hellos/status,verbs=get;update;patch

// Reconcile - perform the reconcile to manage Hello kind
func (r *HelloReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	_ = context.Background()
	_ = r.Log.WithValues("hello", req.NamespacedName)

	// Fetch the Hello instance
	instance := &hellogroupv1alpha1.Hello{}
	err := r.Get(context.TODO(), req.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return ctrl.Result{}, err
	}

	// Check if the deployment already exists, if not create a new one
	result, err := r.reconcileDeployment(instance)
	if err != nil || result.Requeue == true {
		// reconcile requested requeue or errored, so requeue
		return result, err
	} // else all good carry on

	return ctrl.Result{}, nil
}

func (r *HelloReconciler) reconcileDeployment(instance *hellogroupv1alpha1.Hello) (ctrl.Result, error) {

	// Define a new deployment
	dep := r.newDeploymentForHello(instance)

	// Set Hello instance as the owner and controller
	ctrl.SetControllerReference(instance, dep, r.Scheme)

	found := &appsv1.Deployment{}
	err := r.Get(context.TODO(), types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		log.Info("Creating a new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
		err = r.Create(context.TODO(), dep)
		if err != nil {
			log.Error(err, "Failed to create new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
			return ctrl.Result{}, err
		}
		// Deployment created successfully - return and requeue
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		log.Error(err, "Failed to get Deployment")
		return ctrl.Result{}, err
	}

	// deployment already exists

	return ctrl.Result{}, nil
}

func (r *HelloReconciler) newDeploymentForHello(instance *hellogroupv1alpha1.Hello) *appsv1.Deployment {
	labels := map[string]string{
		"app": instance.Name,
	}
	var replicas int32 = 1

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.Name,
			Namespace: instance.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:    "hello",
							Image:   "DOCKERHOSTNAME/DOCKERNAMESPACE/hello:v0.0.1",
							Command: []string{"./hello"},
						},
					},
				},
			},
		},
	}
	return dep
}

// SetupWithManager - a method
func (r *HelloReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&hellogroupv1alpha1.Hello{}).
		Complete(r)
}
