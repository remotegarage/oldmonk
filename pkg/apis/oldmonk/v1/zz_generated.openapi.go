// +build !ignore_autogenerated

// This file was autogenerated by openapi-gen. Do not edit it manually!

package v1

import (
	spec "github.com/go-openapi/spec"
	common "k8s.io/kube-openapi/pkg/common"
)

func GetOpenAPIDefinitions(ref common.ReferenceCallback) map[string]common.OpenAPIDefinition {
	return map[string]common.OpenAPIDefinition{
		"github.com/remotegarage/oldmonk/pkg/apis/oldmonk/v1.ListOptions":         schema_pkg_apis_oldmonk_v1_ListOptions(ref),
		"github.com/remotegarage/oldmonk/pkg/apis/oldmonk/v1.QueueAutoScalerSpec": schema_pkg_apis_oldmonk_v1_QueueAutoScalerSpec(ref),
		"github.com/remotegarage/oldmonk/pkg/apis/oldmonk/v1.ScaleSpec":           schema_pkg_apis_oldmonk_v1_ScaleSpec(ref),
	}
}

func schema_pkg_apis_oldmonk_v1_ListOptions(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "ListOptions defines the desired state of Queue",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"uri": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"string"},
							Format: "",
						},
					},
					"region": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"string"},
							Format: "",
						},
					},
					"type": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"string"},
							Format: "",
						},
					},
					"queue": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"string"},
							Format: "",
						},
					},
					"vshost": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"string"},
							Format: "",
						},
					},
					"key": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"string"},
							Format: "",
						},
					},
					"exchange": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"string"},
							Format: "",
						},
					},
					"tube": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"string"},
							Format: "",
						},
					},
					"group": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"string"},
							Format: "",
						},
					},
					"topic": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"string"},
							Format: "",
						},
					},
				},
				Required: []string{"uri"},
			},
		},
	}
}

func schema_pkg_apis_oldmonk_v1_QueueAutoScalerSpec(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "AutoScalerSpec defines the desired state of QueueAutoScaler",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"type": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"string"},
							Format: "",
						},
					},
					"option": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("github.com/remotegarage/oldmonk/pkg/apis/oldmonk/v1.ListOptions"),
						},
					},
					"minPods": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"integer"},
							Format: "int32",
						},
					},
					"maxPods": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"integer"},
							Format: "int32",
						},
					},
					"scaleUp": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("github.com/remotegarage/oldmonk/pkg/apis/oldmonk/v1.ScaleSpec"),
						},
					},
					"scaleDown": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("github.com/remotegarage/oldmonk/pkg/apis/oldmonk/v1.ScaleSpec"),
						},
					},
					"secrets": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"string"},
							Format: "",
						},
					},
					"deployment": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"string"},
							Format: "",
						},
					},
					"appSpec": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("k8s.io/api/core/v1.Container"),
						},
					},
					"labels": {
						SchemaProps: spec.SchemaProps{
							Type: []string{"object"},
							AdditionalProperties: &spec.SchemaOrBool{
								Allows: true,
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Type:   []string{"string"},
										Format: "",
									},
								},
							},
						},
					},
					"strategy": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("k8s.io/api/apps/v1.DeploymentStrategy"),
						},
					},
					"volume": {
						SchemaProps: spec.SchemaProps{
							Type: []string{"array"},
							Items: &spec.SchemaOrArray{
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Ref: ref("k8s.io/api/core/v1.Volume"),
									},
								},
							},
						},
					},
					"autopilot": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"boolean"},
							Format: "",
						},
					},
				},
				Required: []string{"type", "option", "minPods", "maxPods", "scaleUp", "scaleDown", "secrets", "deployment", "autopilot"},
			},
		},
		Dependencies: []string{
			"github.com/remotegarage/oldmonk/pkg/apis/oldmonk/v1.ListOptions", "github.com/remotegarage/oldmonk/pkg/apis/oldmonk/v1.ScaleSpec", "k8s.io/api/apps/v1.DeploymentStrategy", "k8s.io/api/core/v1.Container", "k8s.io/api/core/v1.Volume"},
	}
}

func schema_pkg_apis_oldmonk_v1_ScaleSpec(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "ScaleSpec defines the desired state of Autoscaler",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"threshold": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"integer"},
							Format: "int32",
						},
					},
					"amount": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"integer"},
							Format: "int32",
						},
					},
				},
				Required: []string{"threshold", "amount"},
			},
		},
	}
}
