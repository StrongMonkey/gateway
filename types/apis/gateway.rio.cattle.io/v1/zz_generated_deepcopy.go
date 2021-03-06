package v1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GatewayDestination) DeepCopyInto(out *GatewayDestination) {
	*out = *in
	out.Namespaced = in.Namespaced
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GatewayDestination.
func (in *GatewayDestination) DeepCopy() *GatewayDestination {
	if in == nil {
		return nil
	}
	out := new(GatewayDestination)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *GatewayDestination) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GatewayDestinationList) DeepCopyInto(out *GatewayDestinationList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]GatewayDestination, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GatewayDestinationList.
func (in *GatewayDestinationList) DeepCopy() *GatewayDestinationList {
	if in == nil {
		return nil
	}
	out := new(GatewayDestinationList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *GatewayDestinationList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GatewayDestinationSpec) DeepCopyInto(out *GatewayDestinationSpec) {
	*out = *in
	if in.MatchHeader != nil {
		in, out := &in.MatchHeader, &out.MatchHeader
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GatewayDestinationSpec.
func (in *GatewayDestinationSpec) DeepCopy() *GatewayDestinationSpec {
	if in == nil {
		return nil
	}
	out := new(GatewayDestinationSpec)
	in.DeepCopyInto(out)
	return out
}
