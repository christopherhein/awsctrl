/*
Copyright © 2019 AWS Controller authors

Licensed under the Apache License, Version 2.0 (the &#34;License&#34;);
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an &#34;AS IS&#34; BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	"fmt"
	"strings"

	metav1alpha1 "go.awsctrl.io/manager/apis/meta/v1alpha1"
	controllerutils "go.awsctrl.io/manager/controllers/utils"
	cfnencoder "go.awsctrl.io/manager/encoding/cloudformation"

	"github.com/awslabs/goformation/v4/cloudformation"
	"github.com/awslabs/goformation/v4/cloudformation/certificatemanager"
	"k8s.io/client-go/dynamic"
)

// GetNotificationARNs is an autogenerated deepcopy function, will return notifications for stack
func (in *Certificate) GetNotificationARNs() []string {
	notifcations := []string{}
	for _, notifarn := range in.Spec.NotificationARNs {
		notifcations = append(notifcations, *notifarn)
	}
	return notifcations
}

// GetTemplate will return the JSON version of the CFN to use.
func (in *Certificate) GetTemplate(client dynamic.Interface) (string, error) {
	if client == nil {
		return "", fmt.Errorf("k8s client not loaded for template")
	}
	template := cloudformation.NewTemplate()

	template.Description = "AWS Controller - certificatemanager.Certificate (ac-{TODO})"

	template.Outputs = map[string]interface{}{
		"ResourceRef": map[string]interface{}{
			"Value": cloudformation.Ref("Certificate"),
			"Export": map[string]interface{}{
				"Name": in.Name + "Ref",
			},
		},
	}

	certificatemanagerCertificate := &certificatemanager.Certificate{}

	if in.Spec.DomainName != "" {
		certificatemanagerCertificate.DomainName = in.Spec.DomainName
	}

	certificatemanagerCertificateDomainValidationOptions := []certificatemanager.Certificate_DomainValidationOption{}

	for _, item := range in.Spec.DomainValidationOptions {
		certificatemanagerCertificateDomainValidationOption := certificatemanager.Certificate_DomainValidationOption{}

		if item.DomainName != "" {
			certificatemanagerCertificateDomainValidationOption.DomainName = item.DomainName
		}

		if item.ValidationDomain != "" {
			certificatemanagerCertificateDomainValidationOption.ValidationDomain = item.ValidationDomain
		}

	}

	if len(certificatemanagerCertificateDomainValidationOptions) > 0 {
		certificatemanagerCertificate.DomainValidationOptions = certificatemanagerCertificateDomainValidationOptions
	}
	if len(in.Spec.SubjectAlternativeNames) > 0 {
		certificatemanagerCertificate.SubjectAlternativeNames = in.Spec.SubjectAlternativeNames
	}

	// TODO(christopherhein): implement tags this could be easy now that I have the mechanims of nested objects
	if in.Spec.ValidationMethod != "" {
		certificatemanagerCertificate.ValidationMethod = in.Spec.ValidationMethod
	}

	template.Resources = map[string]cloudformation.Resource{
		"Certificate": certificatemanagerCertificate,
	}

	// json, err := template.JSONWithOptions(&intrinsics.ProcessorOptions{NoEvaluateConditions: true})
	json, err := template.JSON()
	if err != nil {
		return "", err
	}

	return string(json), nil
}

// GetStackID will return stackID
func (in *Certificate) GetStackID() string {
	return in.Status.StackID
}

// GenerateStackName will generate a StackName
func (in *Certificate) GenerateStackName() string {
	return strings.Join([]string{"certificatemanager", "certificate", in.GetName(), in.GetNamespace()}, "-")
}

// GetStackName will return stackName
func (in *Certificate) GetStackName() string {
	return in.Spec.StackName
}

// GetTemplateVersionLabel will return the stack template version
func (in *Certificate) GetTemplateVersionLabel() (value string, ok bool) {
	value, ok = in.Labels[controllerutils.StackTemplateVersionLabel]
	return
}

// GetParameters will return CFN Parameters
func (in *Certificate) GetParameters() map[string]string {
	params := map[string]string{}
	cfnencoder.MarshalTypes(params, in.Spec, "Parameter")
	return params
}

// GetCloudFormationMeta will return CFN meta object
func (in *Certificate) GetCloudFormationMeta() metav1alpha1.CloudFormationMeta {
	return in.Spec.CloudFormationMeta
}

// GetStatus will return the CFN Status
func (in *Certificate) GetStatus() metav1alpha1.ConditionStatus {
	return in.Status.Status
}

// SetStackID will put a stackID
func (in *Certificate) SetStackID(input string) {
	in.Status.StackID = input
	return
}

// SetStackName will return stackName
func (in *Certificate) SetStackName(input string) {
	in.Spec.StackName = input
	return
}

// SetTemplateVersionLabel will set the template version label
func (in *Certificate) SetTemplateVersionLabel() {
	if len(in.Labels) == 0 {
		in.Labels = map[string]string{}
	}

	in.Labels[controllerutils.StackTemplateVersionLabel] = controllerutils.ComputeHash(in.Spec)
}

// TemplateVersionChanged will return bool if template has changed
func (in *Certificate) TemplateVersionChanged() bool {
	// Ignore bool since it will still record changed
	label, _ := in.GetTemplateVersionLabel()
	return label != controllerutils.ComputeHash(in.Spec)
}

// SetStatus will set status for object
func (in *Certificate) SetStatus(status *metav1alpha1.StatusMeta) {
	in.Status.StatusMeta = *status
}
