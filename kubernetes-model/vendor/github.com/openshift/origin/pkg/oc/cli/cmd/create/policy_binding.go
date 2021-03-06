/**
 * Copyright (C) 2015 Red Hat, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *         http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package create

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/kubernetes/pkg/kubectl/cmd/templates"
	cmdutil "k8s.io/kubernetes/pkg/kubectl/cmd/util"

	authorizationapi "github.com/openshift/origin/pkg/authorization/apis/authorization"
	authorizationclient "github.com/openshift/origin/pkg/authorization/generated/internalclientset/typed/authorization/internalversion"
	"github.com/openshift/origin/pkg/oc/cli/util/clientcmd"
)

const PolicyBindingRecommendedName = "policybinding"

var (
	policyBindingLong = templates.LongDesc(`Create a policy binding that references the policy in the targeted namespace.`)

	policyBindingExample = templates.Examples(`
		# Create a policy binding in namespace "foo" that references the policy in namespace "bar"
  	%[1]s bar -n foo`)
)

type CreatePolicyBindingOptions struct {
	BindingNamespace string
	PolicyNamespace  string

	BindingClient authorizationclient.PolicyBindingsGetter

	Mapper       meta.RESTMapper
	OutputFormat string
	Out          io.Writer
	Printer      ObjectPrinter
}

type ObjectPrinter func(runtime.Object, io.Writer) error

// NewCmdCreatePolicyBinding is a macro command to create a new policy binding.
func NewCmdCreatePolicyBinding(name, fullName string, f *clientcmd.Factory, out io.Writer) *cobra.Command {
	o := &CreatePolicyBindingOptions{Out: out}

	cmd := &cobra.Command{
		Use:     name + " TARGET_POLICY_NAMESPACE",
		Short:   "Create a policy binding that references the policy in the targeted namespace.",
		Long:    policyBindingLong,
		Example: fmt.Sprintf(policyBindingExample, fullName),
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(cmd, f, args))
			cmdutil.CheckErr(o.Validate())
			cmdutil.CheckErr(o.Run())
		},
		Deprecated: fmt.Sprintf("will not work against 3.7 or later servers. Use (Cluster)RoleBindings instead."),
	}
	cmdutil.AddOutputFlagsForMutation(cmd)
	return cmd
}

func (o *CreatePolicyBindingOptions) Complete(cmd *cobra.Command, f *clientcmd.Factory, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("exactly one argument (policy namespace) is supported, not: %v", args)
	}
	o.PolicyNamespace = args[0]

	namespace, _, err := f.DefaultNamespace()
	if err != nil {
		return err
	}
	o.BindingNamespace = namespace

	discovery, err := f.DiscoveryClient()
	if err != nil {
		return err
	}

	if err := clientcmd.LegacyPolicyResourceGate(discovery); err != nil {
		return err
	}
	client, err := f.OpenshiftInternalAuthorizationClient()
	if err != nil {
		return err
	}

	o.BindingClient = client.Authorization()

	o.Mapper, _ = f.Object()
	o.OutputFormat = cmdutil.GetFlagString(cmd, "output")

	o.Printer = func(obj runtime.Object, out io.Writer) error {
		return f.PrintObject(cmd, false, o.Mapper, obj, out)
	}

	return nil
}

func (o *CreatePolicyBindingOptions) Validate() error {
	if len(o.BindingNamespace) == 0 {
		return fmt.Errorf("destination namespace is required")
	}
	if len(o.PolicyNamespace) == 0 {
		return fmt.Errorf("referenced policy namespace is required")
	}
	if o.BindingClient == nil {
		return fmt.Errorf("BindingClient is required")
	}
	if o.Mapper == nil {
		return fmt.Errorf("Mapper is required")
	}
	if o.Out == nil {
		return fmt.Errorf("Out is required")
	}
	if o.Printer == nil {
		return fmt.Errorf("Printer is required")
	}

	return nil
}

func (o *CreatePolicyBindingOptions) Run() error {
	binding := &authorizationapi.PolicyBinding{}
	binding.PolicyRef.Namespace = o.PolicyNamespace
	binding.PolicyRef.Name = authorizationapi.PolicyName
	binding.Name = authorizationapi.GetPolicyBindingName(binding.PolicyRef.Namespace)

	actualBinding, err := o.BindingClient.PolicyBindings(o.BindingNamespace).Create(binding)
	if err != nil {
		return err
	}

	if useShortOutput := o.OutputFormat == "name"; useShortOutput || len(o.OutputFormat) == 0 {
		cmdutil.PrintSuccess(o.Mapper, useShortOutput, o.Out, "policybinding", actualBinding.Name, false, "created")
		return nil
	}

	return o.Printer(actualBinding, o.Out)
}
