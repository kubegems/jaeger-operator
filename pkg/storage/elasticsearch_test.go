package storage

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	esv1 "github.com/openshift/elasticsearch-operator/apis/logging/v1"

	v1 "github.com/jaegertracing/jaeger-operator/apis/v1"
)

func TestCreateElasticsearchCR(t *testing.T) {
	storageClassName := "floppydisk"
	genuuid1 := "myprojectfoo"
	genuuidmaster1 := "myprojectfoomaster"
	genuuid2 := "myprojectfoobar"
	genuuidmaster2 := "myprojectfoobarmaster"
	genuuid3 := "mytolerableprojecttolerations"

	toleration := corev1.Toleration{
		Key:      "special",
		Operator: "Equal",
		Value:    "false",
		Effect:   "NoSchedule",
	}
	tolerations := []corev1.Toleration{toleration}

	tests := []struct {
		name      string
		namespace string
		jEsSpec   v1.ElasticsearchSpec
		esSpec    esv1.ElasticsearchSpec
	}{
		{
			name:      "foo",
			namespace: "myproject",
			jEsSpec: v1.ElasticsearchSpec{
				Name:             "elasticsearch",
				NodeCount:        2,
				RedundancyPolicy: esv1.FullRedundancy,
				Storage: esv1.ElasticsearchStorageSpec{
					StorageClassName: &storageClassName,
				},
			},
			esSpec: esv1.ElasticsearchSpec{
				ManagementState:  esv1.ManagementStateManaged,
				RedundancyPolicy: esv1.FullRedundancy,
				Spec:             esv1.ElasticsearchNodeSpec{},
				Nodes: []esv1.ElasticsearchNode{
					{
						NodeCount: 2,
						Storage:   esv1.ElasticsearchStorageSpec{StorageClassName: &storageClassName},
						Roles:     []esv1.ElasticsearchNodeRole{esv1.ElasticsearchRoleMaster, esv1.ElasticsearchRoleClient, esv1.ElasticsearchRoleData},
						GenUUID:   &genuuid1,
					},
				},
			},
		},
		{
			name:      "foo",
			namespace: "myproject",
			jEsSpec: v1.ElasticsearchSpec{
				Name:             "elasticsearch",
				NodeCount:        5,
				RedundancyPolicy: esv1.FullRedundancy,
				Storage: esv1.ElasticsearchStorageSpec{
					StorageClassName: &storageClassName,
				},
			},
			esSpec: esv1.ElasticsearchSpec{
				ManagementState:  esv1.ManagementStateManaged,
				RedundancyPolicy: esv1.FullRedundancy,
				Spec:             esv1.ElasticsearchNodeSpec{},
				Nodes: []esv1.ElasticsearchNode{
					{
						NodeCount: 3,
						Storage:   esv1.ElasticsearchStorageSpec{StorageClassName: &storageClassName},
						Roles:     []esv1.ElasticsearchNodeRole{esv1.ElasticsearchRoleMaster, esv1.ElasticsearchRoleClient, esv1.ElasticsearchRoleData},
						GenUUID:   &genuuidmaster1,
					},
					{
						NodeCount: 2,
						Storage:   esv1.ElasticsearchStorageSpec{StorageClassName: &storageClassName},
						Roles:     []esv1.ElasticsearchNodeRole{esv1.ElasticsearchRoleClient, esv1.ElasticsearchRoleData},
						GenUUID:   &genuuid1,
					},
				},
			},
		},
		{
			name:      "foo-ba%r",
			namespace: "myproje&ct",
			jEsSpec: v1.ElasticsearchSpec{
				Name:             "elasticsearch",
				NodeCount:        5,
				RedundancyPolicy: esv1.FullRedundancy,
				Storage: esv1.ElasticsearchStorageSpec{
					StorageClassName: &storageClassName,
				},
			},
			esSpec: esv1.ElasticsearchSpec{
				ManagementState:  esv1.ManagementStateManaged,
				RedundancyPolicy: esv1.FullRedundancy,
				Spec:             esv1.ElasticsearchNodeSpec{},
				Nodes: []esv1.ElasticsearchNode{
					{
						NodeCount: 3,
						Storage:   esv1.ElasticsearchStorageSpec{StorageClassName: &storageClassName},
						Roles:     []esv1.ElasticsearchNodeRole{esv1.ElasticsearchRoleMaster, esv1.ElasticsearchRoleClient, esv1.ElasticsearchRoleData},
						GenUUID:   &genuuidmaster2,
					},
					{
						NodeCount: 2,
						Storage:   esv1.ElasticsearchStorageSpec{StorageClassName: &storageClassName},
						Roles:     []esv1.ElasticsearchNodeRole{esv1.ElasticsearchRoleClient, esv1.ElasticsearchRoleData},
						GenUUID:   &genuuid2,
					},
				},
			},
		},
		{
			name:      "tolerations",
			namespace: "mytolerableproject",
			jEsSpec: v1.ElasticsearchSpec{
				Name:             "elasticsearch",
				NodeCount:        2,
				RedundancyPolicy: esv1.FullRedundancy,
				Tolerations:      tolerations,
				Storage: esv1.ElasticsearchStorageSpec{
					StorageClassName: &storageClassName,
				},
			},
			esSpec: esv1.ElasticsearchSpec{
				ManagementState:  esv1.ManagementStateManaged,
				RedundancyPolicy: esv1.FullRedundancy,
				Spec: esv1.ElasticsearchNodeSpec{
					Tolerations: tolerations,
				},
				Nodes: []esv1.ElasticsearchNode{
					{
						NodeCount: 2,
						Storage:   esv1.ElasticsearchStorageSpec{StorageClassName: &storageClassName},
						Roles:     []esv1.ElasticsearchNodeRole{esv1.ElasticsearchRoleMaster, esv1.ElasticsearchRoleClient, esv1.ElasticsearchRoleData},
						GenUUID:   &genuuid3,
					},
				},
			},
		},
	}
	for _, test := range tests {
		j := v1.NewJaeger(types.NamespacedName{Name: test.name, Namespace: test.namespace})
		j.Spec.Storage.Elasticsearch = test.jEsSpec
		es := &ElasticsearchDeployment{Jaeger: j}
		cr := es.Elasticsearch()
		assert.Equal(t, test.namespace, cr.Namespace)
		assert.Equal(t, "elasticsearch", cr.Name)
		trueVar := true
		assert.Equal(t, []metav1.OwnerReference{{Name: test.name, Controller: &trueVar}}, cr.OwnerReferences)
		assert.Equal(t, cr.Spec, test.esSpec)
	}
}

func TestInject(t *testing.T) {
	tests := []struct {
		pod      *corev1.PodSpec
		expected *corev1.PodSpec
		es       v1.ElasticsearchSpec
	}{
		{
			es: v1.ElasticsearchSpec{
				Name: "elasticsearch",
			},
			pod: &corev1.PodSpec{
				Containers: []corev1.Container{{
					Args:         []string{"foo"},
					VolumeMounts: []corev1.VolumeMount{{Name: "lol"}},
				}},
			},
			expected: &corev1.PodSpec{
				Containers: []corev1.Container{{
					Args: []string{
						"foo",
						"--es.server-urls=https://elasticsearch:9200",
						"--es.tls.enabled=true",
						"--es.tls.ca=" + caPath,
						"--es.tls.cert=" + certPath,
						"--es.tls.key=" + keyPath,
						"--es.timeout=15s",
						"--es.num-shards=0",
						"--es.num-replicas=1",
					},
					VolumeMounts: []corev1.VolumeMount{
						{Name: "lol"},
						{Name: volumeName, ReadOnly: true, MountPath: volumeMountPath},
					},
				}},
				Volumes: []corev1.Volume{{Name: "certs", VolumeSource: corev1.VolumeSource{
					Secret: &corev1.SecretVolumeSource{
						SecretName: "jaeger-elasticsearch"}}},
				}},
		},
		{
			es: v1.ElasticsearchSpec{Name: "elasticsearch"},
			pod: &corev1.PodSpec{
				Containers: []corev1.Container{{
					Args: []string{"--es.num-shards=15", "--es.num-replicas=55", "--es.timeout=99s"},
				}},
			},
			expected: &corev1.PodSpec{
				Containers: []corev1.Container{{
					Args: []string{
						"--es.num-shards=15",
						"--es.num-replicas=55",
						"--es.timeout=99s",
						"--es.server-urls=https://elasticsearch:9200",
						"--es.tls.enabled=true",
						"--es.tls.ca=" + caPath,
						"--es.tls.cert=" + certPath,
						"--es.tls.key=" + keyPath,
					},
					VolumeMounts: []corev1.VolumeMount{
						{Name: volumeName, ReadOnly: true, MountPath: volumeMountPath},
					},
				}},
				Volumes: []corev1.Volume{{Name: "certs", VolumeSource: corev1.VolumeSource{
					Secret: &corev1.SecretVolumeSource{
						SecretName: "jaeger-elasticsearch"}}},
				}},
		},
		{
			pod: &corev1.PodSpec{Containers: []corev1.Container{{}}},
			es: v1.ElasticsearchSpec{
				Name:             "my-es",
				NodeCount:        15,
				RedundancyPolicy: esv1.FullRedundancy,
			},
			expected: &corev1.PodSpec{
				Containers: []corev1.Container{{
					Args: []string{
						"--es.server-urls=https://my-es:9200",
						"--es.tls.enabled=true",
						"--es.tls.ca=" + caPath,
						"--es.tls.cert=" + certPath,
						"--es.tls.key=" + keyPath,
						"--es.timeout=15s",
						"--es.num-shards=15",
						"--es.num-replicas=14",
					},
					VolumeMounts: []corev1.VolumeMount{
						{Name: volumeName, ReadOnly: true, MountPath: volumeMountPath},
					},
				}},
				Volumes: []corev1.Volume{{Name: "certs", VolumeSource: corev1.VolumeSource{
					Secret: &corev1.SecretVolumeSource{
						SecretName: "jaeger-my-es"}}},
				}},
		},
		{
			es: v1.ElasticsearchSpec{
				Name:             "es-tenant2",
				NodeCount:        15,
				RedundancyPolicy: esv1.FullRedundancy,
			},
			pod: &corev1.PodSpec{
				Containers: []corev1.Container{
					{
						Args: []string{"--es-archive.enabled=true"},
					},
				},
			},
			expected: &corev1.PodSpec{
				Containers: []corev1.Container{{
					Args: []string{
						"--es-archive.enabled=true",
						"--es.server-urls=https://es-tenant2:9200",
						"--es.tls.enabled=true",
						"--es.tls.ca=" + caPath,
						"--es.tls.cert=" + certPath,
						"--es.tls.key=" + keyPath,
						"--es.timeout=15s",
						"--es.num-shards=15",
						"--es.num-replicas=14",
						"--es-archive.server-urls=https://es-tenant2:9200",
						"--es-archive.tls.enabled=true",
						"--es-archive.tls.ca=" + caPath,
						"--es-archive.tls.cert=" + certPath,
						"--es-archive.tls.key=" + keyPath,
						"--es-archive.timeout=15s",
						"--es-archive.num-shards=15",
						"--es-archive.num-replicas=14",
					},
					VolumeMounts: []corev1.VolumeMount{
						{Name: volumeName, ReadOnly: true, MountPath: volumeMountPath},
					},
				}},
				Volumes: []corev1.Volume{{Name: "certs", VolumeSource: corev1.VolumeSource{
					Secret: &corev1.SecretVolumeSource{
						SecretName: "jaeger-es-tenant2"}}},
				}},
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			es := &ElasticsearchDeployment{Jaeger: v1.NewJaeger(types.NamespacedName{Name: "hoo"})}
			es.Jaeger.Spec.Storage.Elasticsearch = test.es
			es.InjectStorageConfiguration(test.pod)
			assert.Equal(t, test.expected, test.pod)
		})
	}

}

func TestCalculateReplicaShards(t *testing.T) {
	tests := []struct {
		dataNodes int
		redType   esv1.RedundancyPolicyType
		shards    int
	}{
		{redType: esv1.ZeroRedundancy, dataNodes: 1, shards: 0},
		{redType: esv1.ZeroRedundancy, dataNodes: 1, shards: 0},
		{redType: esv1.SingleRedundancy, dataNodes: 1, shards: 1},
		{redType: esv1.SingleRedundancy, dataNodes: 20, shards: 1},
		{redType: esv1.MultipleRedundancy, dataNodes: 1, shards: 0},
		{redType: esv1.MultipleRedundancy, dataNodes: 20, shards: 9},
		{redType: esv1.FullRedundancy, dataNodes: 1, shards: 0},
		{redType: esv1.FullRedundancy, dataNodes: 20, shards: 19},
	}
	for _, test := range tests {
		assert.Equal(t, test.shards, calculateReplicaShards(test.redType, test.dataNodes))
	}
}
