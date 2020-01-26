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
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	metav1alpha1 "go.awsctrl.io/manager/apis/meta/v1alpha1"
	controllerutils "go.awsctrl.io/manager/controllers/utils"
	cfnencoder "go.awsctrl.io/manager/encoding/cloudformation"

	"github.com/awslabs/goformation/v4/cloudformation"
	"github.com/awslabs/goformation/v4/cloudformation/iam"
	"k8s.io/client-go/dynamic"
)

// GetNotificationARNs is an autogenerated deepcopy function, will return notifications for stack
func (in *User) GetNotificationARNs() []string {
	notifcations := []string{}
	for _, notifarn := range in.Spec.NotificationARNs {
		notifcations = append(notifcations, *notifarn)
	}
	return notifcations
}

// GetTemplate will return the JSON version of the CFN to use.
func (in *User) GetTemplate(client dynamic.Interface) (string, error) {
	if client == nil {
		return "", fmt.Errorf("k8s client not loaded for template")
	}
	template := cloudformation.NewTemplate()

	template.Description = "AWS Controller - iam.User (ac-{TODO})"

	template.Outputs = map[string]interface{}{
		"ResourceRef": map[string]interface{}{
			"Value": cloudformation.Ref("User"),
			"Export": map[string]interface{}{
				"Name": in.Name + "Ref",
			},
		},
		"Arn": map[string]interface{}{
			"Value":  cloudformation.GetAtt("User", "Arn"),
			"Export": map[string]interface{}{"Name": in.Name + "Arn"},
		},
	}

	iamUser := &iam.User{}

	if !reflect.DeepEqual(in.Spec.LoginProfile, User_LoginProfile{}) {
		iamUserLoginProfile := iam.User_LoginProfile{}

		if in.Spec.LoginProfile.PasswordResetRequired || !in.Spec.LoginProfile.PasswordResetRequired {
			iamUserLoginProfile.PasswordResetRequired = in.Spec.LoginProfile.PasswordResetRequired
		}

		if in.Spec.LoginProfile.Password != "" {
			iamUserLoginProfile.Password = in.Spec.LoginProfile.Password
		}

		iamUser.LoginProfile = &iamUserLoginProfile
	}

	if len(in.Spec.ManagedPolicyRefs) > 0 {
		iamUserManagedPolicyRefs := []string{}

		for _, item := range in.Spec.ManagedPolicyRefs {
			iamUserManagedPolicyRefsItem := item.DeepCopy()

			if iamUserManagedPolicyRefsItem.ObjectRef.Namespace == "" {
				iamUserManagedPolicyRefsItem.ObjectRef.Namespace = in.Namespace
			}

		}

		iamUser.ManagedPolicyArns = iamUserManagedPolicyRefs
	}

	if in.Spec.Path != "" {
		iamUser.Path = in.Spec.Path
	}

	if in.Spec.PermissionsBoundary != "" {
		iamUser.PermissionsBoundary = in.Spec.PermissionsBoundary
	}

	iamUserPolicies := []iam.User_Policy{}

	for _, item := range in.Spec.Policies {
		iamUserPolicy := iam.User_Policy{}

		if item.PolicyDocument != "" {
			iamUserPolicyJSON := make(map[string]interface{})
			err := json.Unmarshal([]byte(item.PolicyDocument), &iamUserPolicyJSON)
			if err != nil {
				return "", err
			}
			iamUserPolicy.PolicyDocument = iamUserPolicyJSON
		}

		if item.PolicyName != "" {
			iamUserPolicy.PolicyName = item.PolicyName
		}

	}

	if len(iamUserPolicies) > 0 {
		iamUser.Policies = iamUserPolicies
	}
	// TODO(christopherhein): implement tags this could be easy now that I have the mechanims of nested objects
	// TODO(christopherhein) move these to a defaulter
	if in.Spec.UserName == "" {
		iamUser.UserName = in.Name
	}

	if in.Spec.UserName != "" {
		iamUser.UserName = in.Spec.UserName
	}

	if len(in.Spec.Groups) > 0 {
		iamUser.Groups = in.Spec.Groups
	}

	template.Resources = map[string]cloudformation.Resource{
		"User": iamUser,
	}

	// json, err := template.JSONWithOptions(&intrinsics.ProcessorOptions{NoEvaluateConditions: true})
	json, err := template.JSON()
	if err != nil {
		return "", err
	}

	return string(json), nil
}

// GetStackID will return stackID
func (in *User) GetStackID() string {
	return in.Status.StackID
}

// GenerateStackName will generate a StackName
func (in *User) GenerateStackName() string {
	return strings.Join([]string{"iam", "user", in.GetName(), in.GetNamespace()}, "-")
}

// GetStackName will return stackName
func (in *User) GetStackName() string {
	return in.Spec.StackName
}

// GetTemplateVersionLabel will return the stack template version
func (in *User) GetTemplateVersionLabel() (value string, ok bool) {
	value, ok = in.Labels[controllerutils.StackTemplateVersionLabel]
	return
}

// GetParameters will return CFN Parameters
func (in *User) GetParameters() map[string]string {
	params := map[string]string{}
	cfnencoder.MarshalTypes(params, in.Spec, "Parameter")
	return params
}

// GetCloudFormationMeta will return CFN meta object
func (in *User) GetCloudFormationMeta() metav1alpha1.CloudFormationMeta {
	return in.Spec.CloudFormationMeta
}

// GetStatus will return the CFN Status
func (in *User) GetStatus() metav1alpha1.ConditionStatus {
	return in.Status.Status
}

// SetStackID will put a stackID
func (in *User) SetStackID(input string) {
	in.Status.StackID = input
	return
}

// SetStackName will return stackName
func (in *User) SetStackName(input string) {
	in.Spec.StackName = input
	return
}

// SetTemplateVersionLabel will set the template version label
func (in *User) SetTemplateVersionLabel() {
	if len(in.Labels) == 0 {
		in.Labels = map[string]string{}
	}

	in.Labels[controllerutils.StackTemplateVersionLabel] = controllerutils.ComputeHash(in.Spec)
}

// TemplateVersionChanged will return bool if template has changed
func (in *User) TemplateVersionChanged() bool {
	// Ignore bool since it will still record changed
	label, _ := in.GetTemplateVersionLabel()
	return label != controllerutils.ComputeHash(in.Spec)
}

// SetStatus will set status for object
func (in *User) SetStatus(status *metav1alpha1.StatusMeta) {
	in.Status.StatusMeta = *status
}
