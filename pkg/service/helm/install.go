package helm

import (
	appsv1 "k8s.io/api/apps/v1"
	appsv1beta1 "k8s.io/api/apps/v1beta1"
	appsv1beta2 "k8s.io/api/apps/v1beta2"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
	extv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/kubernetes/pkg/apis/apps"
)

// Scheme is the default instance of runtime.Scheme to which types in the Kubernetes API are already registered.
// NOTE: If you are copying this file to start a new api group, STOP! Copy the
// extensions group instead. This Scheme is special and should appear ONLY in
// the api group, unless you really know what you're doing.
// TODO(lavalamp): make the above error impossible.
var Scheme = runtime.NewScheme()

// Codecs provides access to encoding and decoding for the scheme
var Codecs = serializer.NewCodecFactory(Scheme)

// ParameterCodec handles versioning of objects that are converted to query parameters.
var ParameterCodec = runtime.NewParameterCodec(Scheme)

func init() {
	ExtensionInstall(Scheme)
	AppInstall(Scheme)
}

func ExtensionInstall(scheme *runtime.Scheme) {
	utilruntime.Must(apiextensions.AddToScheme(scheme))
	utilruntime.Must(extv1beta1.AddToScheme(scheme))
	utilruntime.Must(scheme.SetVersionPriority(extv1beta1.SchemeGroupVersion))
}

func AppInstall(scheme *runtime.Scheme) {
	utilruntime.Must(apps.AddToScheme(scheme))
	utilruntime.Must(appsv1beta1.AddToScheme(scheme))
	utilruntime.Must(appsv1beta2.AddToScheme(scheme))
	utilruntime.Must(appsv1.AddToScheme(scheme))
	utilruntime.Must(scheme.SetVersionPriority(appsv1.SchemeGroupVersion, appsv1beta2.SchemeGroupVersion, appsv1beta1.SchemeGroupVersion))

}
