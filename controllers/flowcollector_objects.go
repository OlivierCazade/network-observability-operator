package controllers

import (
	"embed"
	"fmt"
	"path/filepath"

	"github.com/netobserv/flowlogs-pipeline/pkg/confgen"
	"github.com/netobserv/network-observability-operator/controllers/constants"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	downstreamLabelKey       = "openshift.io/cluster-monitoring"
	downstreamLabelValue     = "true"
	roleSuffix               = "-metrics-reader"
	monitoringServiceAccount = "prometheus-k8s"
	monitoringNamespace      = "openshift-monitoring"
)

func buildNamespace(ns string, isDownstream bool) *corev1.Namespace {
	labels := map[string]string{}
	if isDownstream {
		labels[downstreamLabelKey] = downstreamLabelValue
	}
	return &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:   ns,
			Labels: labels,
		},
	}
}

func buildRoleMonitoringReader(ns string) *rbacv1.ClusterRole {
	cr := rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name: constants.OperatorName + roleSuffix,
		},
		Rules: []rbacv1.PolicyRule{
			{APIGroups: []string{""},
				Verbs:     []string{"get", "list", "watch"},
				Resources: []string{"pods", "services", "endpoints"},
			},
			{
				NonResourceURLs: []string{"/metrics"},
				Verbs:           []string{"get"},
			},
		},
	}
	return &cr
}

func buildRoleBindingMonitoringReader(ns string) *rbacv1.ClusterRoleBinding {
	return &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      constants.OperatorName + roleSuffix,
			Namespace: ns,
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     constants.OperatorName + roleSuffix,
		},
		Subjects: []rbacv1.Subject{{
			Kind:      "ServiceAccount",
			Name:      monitoringServiceAccount,
			Namespace: monitoringNamespace,
		}},
	}
}

//go:embed health_dashboard
var healthDashboardEmbed embed.FS

const (
	dashboardConfigDir          = "health_dashboard"
	healthDashboardName         = "NetobservHealth"
	healthDashboardTitle        = "Netobserv health metrics"
	healthDashboardTags         = "['netobserv-health']"
	healthDashboardCMName       = "grafana-dashboard-netobserv-health"
	healthDashboardCMNamespace  = "openshift-config-managed"
	healthDashboardCMAnnotation = "console.openshift.io/dashboard"
	healthDashboardCMFile       = "netobserv-health-metrics.json"
)

func buildHealthDashboard() (*corev1.ConfigMap, error) {
	entries, err := healthDashboardEmbed.ReadDir(dashboardConfigDir)
	if err != nil {
		return nil, fmt.Errorf("failed to access metrics_definitions directory: %w", err)
	}

	cg := confgen.NewConfGen(&confgen.Options{})

	config := confgen.Config{
		Visualization: confgen.ConfigVisualization{
			Grafana: confgen.ConfigVisualizationGrafana{
				Dashboards: []confgen.ConfigVisualizationGrafanaDashboard{
					{
						Name:          healthDashboardName,
						Title:         healthDashboardTitle,
						TimeFrom:      "now",
						Tags:          healthDashboardTags,
						SchemaVersion: "16",
					},
				},
			},
		},
	}
	cg.SetConfig(&config)

	for _, entry := range entries {
		fileName := entry.Name()
		srcPath := filepath.Join(dashboardConfigDir, fileName)

		input, err := healthDashboardEmbed.ReadFile(srcPath)
		if err != nil {
			return nil, fmt.Errorf("error reading metrics file %s; %w", srcPath, err)
		}
		err = cg.ParseDefinition(fileName, input)
		if err != nil {
			return nil, fmt.Errorf("error parsing metrics file %s; %w", srcPath, err)
		}
	}

	jsonStr, _ := cg.GenerateGrafanaJson()
	configMap := corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      healthDashboardCMName,
			Namespace: healthDashboardCMNamespace,
			Labels: map[string]string{
				healthDashboardCMAnnotation: "true",
			},
		},
		Data: map[string]string{
			healthDashboardCMFile: jsonStr,
		},
	}
	return &configMap, nil
}
