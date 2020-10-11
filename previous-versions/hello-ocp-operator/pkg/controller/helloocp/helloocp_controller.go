package helloocp

import (
	"context"

	"github.com/go-logr/logr"
	//routev1 "github.com/openshift/api/route/v1"
	helloocpv1alpha1 "github.com/thisisdavidbell/hello-ocp/hello-ocp-operator/pkg/apis/helloocp/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_helloocp")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new Helloocp Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileHelloocp{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("helloocp-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource Helloocp
	err = c.Watch(&source.Kind{Type: &helloocpv1alpha1.Helloocp{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner Helloocp
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &helloocpv1alpha1.Helloocp{},
	})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &corev1.Service{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &helloocpv1alpha1.Helloocp{},
	})
	if err != nil {
		return err
	}

	/*	err = c.Watch(&source.Kind{Type: &routev1.Route{}}, &handler.EnqueueRequestForOwner{
			IsController: true,
			OwnerType:    &helloocpv1alpha1.Helloocp{},
		})
		if err != nil {
			return err
		}
	*/
	return nil
}

// blank assignment to verify that ReconcileHelloocp implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileHelloocp{}

// ReconcileHelloocp reconciles a Helloocp object
type ReconcileHelloocp struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a Helloocp object and makes changes based on the state read
// and what is in the Helloocp.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileHelloocp) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling Helloocp")

	// Fetch the Helloocp instance
	instance := &helloocpv1alpha1.Helloocp{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	result, err := reconcilePod(instance, r, reqLogger)
	if err != nil || result.Requeue == true {
		// reconcile requested requeue or errored, so requeue
		return result, err
	} // else all good carry on

	result, err = reconcileService(instance, r, reqLogger)
	if err != nil || result.Requeue == true {
		// reconcile requested requeue or errored, so requeue
		return result, err
	} // else all good carry on

	// if we got here all went well. Do not requeue
	return reconcile.Result{}, nil
}

// newPodForCR returns a busybox pod with the same name/namespace as the cr
func newPodForCR(cr *helloocpv1alpha1.Helloocp) *corev1.Pod {
	labels := map[string]string{
		"app": cr.Name,
	}
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name + "-pod",
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:    "hello-ocp",
					Image:   "somerandomhostnamem/drb/hello-ocp:v0.0.1",
					Command: []string{"./hello-ocp"},
					Env: []corev1.EnvVar{
						{
							Name:  "HELLONAME",
							Value: cr.Spec.HelloName,
						},
					},
				},
			},
		},
	}
}

func reconcilePod(instance *helloocpv1alpha1.Helloocp, r *ReconcileHelloocp, reqLogger logr.Logger) (reconcile.Result, error) {

	// Define a new Pod object
	pod := newPodForCR(instance)

	// Set Helloocp instance as the owner and controller
	if err := controllerutil.SetControllerReference(instance, pod, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	// Check if this Pod already exists
	found := &corev1.Pod{}
	err := r.client.Get(context.TODO(), types.NamespacedName{Name: pod.Name, Namespace: pod.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new Pod", "Pod.Namespace", pod.Namespace, "Pod.Name", pod.Name)
		err = r.client.Create(context.TODO(), pod)
		if err != nil {
			return reconcile.Result{}, err
		}

		// Pod created successfully - don't requeue
		reqLogger.Info("Pod successfully created - dont requeue", "Pod.Namespace", pod.Namespace, "Pod.Name", pod.Name)
		return reconcile.Result{}, nil
	} else if err != nil {
		return reconcile.Result{}, err
	}

	// Orig: Pod already exists - don't requeue
	// Changed: for now, always update the pod (this is even if not changed - just wip)

	//WIP TODO - hack - loop through correctly
	if found.Spec.Containers[0].Env[0].Value != pod.Spec.Containers[0].Env[0].Value {
		// Clearly a silly approach - have to delete as cannot change env and would be meaningless to anyway.
		reqLogger.Info("Always delete pod as envs dont match.", "Pod.Namespace", pod.Namespace, "Pod.Name", pod.Name)
		reqLogger.Info(pod.Spec.Containers[0].Env[0].Value, "Pod.Namespace", pod.Namespace, "Pod.Name", pod.Name)

		err = r.client.Delete(context.TODO(), found)
		if err != nil {
			reqLogger.Info("pod Delete errored", "Pod.Namespace", pod.Namespace, "Pod.Name", pod.Name)
			return reconcile.Result{}, err
		}
	}
	reqLogger.Info("pod  successfully - dont requeue", "Pod.Namespace", pod.Namespace, "Pod.Name", pod.Name)

	// pod all good  - dont requeue.
	return reconcile.Result{}, nil
}

//apiVersion: v1
//kind: Service
//metadata:
//  name: example
//  namespace: drb-hello
//spec:
//  selector:
//    app: example-helloocp-3
//  ports:
//    - protocol: TCP
//      port: 8080
//      targetPort: 8080

// newPodForCR returns a busybox pod with the same name/namespace as the cr
func newServiceForCR(cr *helloocpv1alpha1.Helloocp) *corev1.Service {
	selector := map[string]string{
		"app": cr.Name,
	}

	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name + "-service",
			Namespace: cr.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Selector: selector,
			Ports: []corev1.ServicePort{
				{
					Protocol:   corev1.ProtocolTCP,
					Port:       int32(8080),
					TargetPort: intstr.FromInt(8080),
				},
			},
		},
	}
}

func reconcileService(instance *helloocpv1alpha1.Helloocp, r *ReconcileHelloocp, reqLogger logr.Logger) (reconcile.Result, error) {

	// Define a new Pod object
	pod := newServiceForCR(instance)

	// Set Helloocp instance as the owner and controller
	if err := controllerutil.SetControllerReference(instance, pod, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	// Check if this Pod already exists
	found := &corev1.Pod{}
	err := r.client.Get(context.TODO(), types.NamespacedName{Name: pod.Name, Namespace: pod.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new Pod", "Pod.Namespace", pod.Namespace, "Pod.Name", pod.Name)
		err = r.client.Create(context.TODO(), pod)
		if err != nil {
			return reconcile.Result{}, err
		}

		// Pod created successfully - don't requeue
		reqLogger.Info("Pod successfully created - dont requeue", "Pod.Namespace", pod.Namespace, "Pod.Name", pod.Name)
		return reconcile.Result{}, nil
	} else if err != nil {
		return reconcile.Result{}, err
	}

	// Orig: Pod already exists - don't requeue
	return reconcile.Result{}, nil
}
