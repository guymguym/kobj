package server

import (
	"time"

	kobjv1 "github.com/kobj-io/kobj/pkg/apis/kobj/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	metav1beta1 "k8s.io/apimachinery/pkg/apis/meta/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/duration"
	"k8s.io/kubernetes/pkg/printers"
)

func KobjTableHandlers(h printers.PrintHandler) {
	columns := []metav1beta1.TableColumnDefinition{
		{Name: "Namespace", Type: "string", Format: "namespace", Description: metav1.ObjectMeta{}.SwaggerDoc()["namespace"]},
		{Name: "Name", Type: "string", Format: "name", Description: metav1.ObjectMeta{}.SwaggerDoc()["name"]},
		{Name: "Size", Type: "int", Description: "The object size in bytes"},
		{Name: "Age", Type: "string", Description: metav1.ObjectMeta{}.SwaggerDoc()["creationTimestamp"]},
	}
	h.TableHandler(columns, KobjPrinter)
	h.TableHandler(columns, KobjListPrinter)

}

func KobjPrinter(ko *kobjv1.Kobj, options printers.GenerateOptions) ([]metav1beta1.TableRow, error) {
	row := metav1beta1.TableRow{
		Object: runtime.RawExtension{Object: ko},
	}
	row.Cells = append(row.Cells, ko.Namespace, ko.Name, len(ko.Data), translateTimestampSince(ko.CreationTimestamp))
	return []metav1beta1.TableRow{row}, nil
}

func KobjListPrinter(list *kobjv1.KobjList, options printers.GenerateOptions) ([]metav1beta1.TableRow, error) {
	rows := make([]metav1beta1.TableRow, 0, len(list.Items))
	for i := range list.Items {
		r, err := KobjPrinter(&list.Items[i], options)
		if err != nil {
			return nil, err
		}
		rows = append(rows, r...)
	}
	return rows, nil
}

// translateTimestampSince returns the elapsed time since timestamp in
// human-readable approximation.
func translateTimestampSince(timestamp metav1.Time) string {
	if timestamp.IsZero() {
		return "<unknown>"
	}

	return duration.HumanDuration(time.Since(timestamp.Time))
}
