package npcore

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/lmxia/gaia/pkg/apis/apps/v1alpha1"
	ncsnp "github.com/lmxia/gaia/pkg/networkfilter/model"
	"github.com/lmxia/gaia/pkg/networkfilter/nputil"
)

func SetRbsAndNetworksRequirment() ([]*v1alpha1.ResourceBinding, *v1alpha1.NetworkRequirement) {

	var rbs []*v1alpha1.ResourceBinding
	var rb0 = v1alpha1.ResourceBinding{
		Spec: v1alpha1.ResourceBindingSpec{
			AppID: "0",
			RbApps: []*v1alpha1.ResourceBindingApps{
				0: {
					ClusterName: "Domain1",
					Replicas: map[string]int32{
						"a": 2,
						"b": 0,
					},
				},
				1: {
					ClusterName: "Domain4",
					Replicas: map[string]int32{
						"a": 0,
						"b": 1,
					},
				},
				2: {
					ClusterName: "Domain3",
					Replicas: map[string]int32{
						"a": 0,
						"c": 2,
					},
				},
			},
		},
	}
	rbs = append(rbs, &rb0)
	var rb1 = v1alpha1.ResourceBinding{
		Spec: v1alpha1.ResourceBindingSpec{
			AppID: "0",
			RbApps: []*v1alpha1.ResourceBindingApps{
				0: {
					ClusterName: "Domain1",
					Replicas: map[string]int32{
						"b": 2,
						"c": 1,
					},
				},
				1: {
					ClusterName: "Domain4",
					Replicas: map[string]int32{
						"b": 1,
						"c": 1,
					},
				},
				/*2: {
					ClusterName: "Domain3",
					Replicas: map[string]int32{
						"a": 0,
						"c": 2,
					},
				},*/
			},
		},
	}
	rbs = append(rbs, &rb1)

	var networkReq = v1alpha1.NetworkRequirement{
		Spec: v1alpha1.NetworkRequirementSpec{
			NetworkCommunication: []v1alpha1.NetworkCommunication{
				0: {
					Name:   "a",
					SelfID: []string{"sca1", "sca2"},
					InterSCNID: []v1alpha1.InterSCNID{
						0: {
							Source: v1alpha1.Direction{
								Id: "sca1",
							},
							Destination: v1alpha1.Direction{
								Id: "scb1",
							},
							Sla: v1alpha1.AppSlaAttr{
								Delay:     10000,
								Lost:      10000,
								Jitter:    1000,
								Bandwidth: 100,
							},
						},
						1: {
							Source: v1alpha1.Direction{
								Id: "sca2",
							},
							Destination: v1alpha1.Direction{
								Id: "scc1",
							},
							Sla: v1alpha1.AppSlaAttr{
								Delay:     10000,
								Lost:      10000,
								Jitter:    1000,
								Bandwidth: 100,
							},
						},
					},
				},
				1: {
					Name:   "b",
					SelfID: []string{"scb1"},
					InterSCNID: []v1alpha1.InterSCNID{
						0: {
							Source: v1alpha1.Direction{
								Id: "scb1",
							},
							Destination: v1alpha1.Direction{
								Id: "scc1",
							},
							Sla: v1alpha1.AppSlaAttr{
								Delay:     10000,
								Lost:      10000,
								Jitter:    1000,
								Bandwidth: 100,
							},
						},
					},
				},
				2: {
					Name:       "c",
					SelfID:     []string{"scc1"},
					InterSCNID: []v1alpha1.InterSCNID{},
				},
			},
		},
	}
	return rbs, &networkReq
}

func BuildNetworkDomainEdge() map[string][]byte {
	nputil.TraceInfoBegin("------------------------------------------------------")

	var domainTopoCacheArry []*ncsnp.DomainTopoCacheNotify
	var domainTopoMsg = make(map[string][]byte)

	domainTopoCache1 := new(ncsnp.DomainTopoCacheNotify)
	domainTopoCache1.LocalDomainId = 1
	domainTopoCache1.LocalDomainName = "Domain1"
	domainTopoCache1.LocalNodeSN = "1-1"
	domainVLink12 := new(ncsnp.DomainVLink)
	domainVLink12.LocalDomainName = "Domain1"
	domainVLink12.LocalDomainId = 1
	domainVLink12.RemoteDomainName = "Domain2"
	domainVLink12.RemoteDomainId = 2
	domainVLink12.LocalNodeSN = "Node12"
	domainVLink12.RemoteNodeSN = "Node21"
	domainVLink12.AttachDomainId = 1012
	domainVLink12.AttachDomainName = "Fabric12"
	vLinkSlaAttr12 := new(ncsnp.VLinkSla)
	vLinkSlaAttr12.Delay = 1
	vLinkSlaAttr12.Bandwidth = 10000
	vLinkSlaAttr12.FreeBandwidth = 10000
	domainVLink12.VLinkSlaAttr = vLinkSlaAttr12
	domainTopoCache1.DomainVLinkArray = append(domainTopoCache1.DomainVLinkArray, domainVLink12)
	domainVLink13 := new(ncsnp.DomainVLink)
	domainVLink13.LocalDomainName = "Domain1"
	domainVLink13.LocalDomainId = 1
	domainVLink13.RemoteDomainName = "Domain3"
	domainVLink13.RemoteDomainId = 3
	domainVLink13.LocalNodeSN = "Node13"
	domainVLink13.RemoteNodeSN = "Node31"
	domainVLink13.AttachDomainId = 1013
	domainVLink12.AttachDomainName = "Fabric13"
	vLinkSlaAttr13 := new(ncsnp.VLinkSla)
	vLinkSlaAttr13.Delay = 1
	vLinkSlaAttr13.Bandwidth = 10000
	vLinkSlaAttr13.FreeBandwidth = 10000
	domainVLink13.VLinkSlaAttr = vLinkSlaAttr13
	domainTopoCache1.DomainVLinkArray = append(domainTopoCache1.DomainVLinkArray, domainVLink13)
	domainTopoCacheArry = append(domainTopoCacheArry, domainTopoCache1)
	content, err := proto.Marshal(domainTopoCache1)
	if err != nil {
		nputil.TraceErrorWithStack(err)
		return nil
	}
	domainTopoMsg[domainTopoCache1.LocalDomainName] = content

	domainTopoCache2 := new(ncsnp.DomainTopoCacheNotify)
	domainTopoCache2.LocalDomainId = 2
	domainTopoCache2.LocalDomainName = "Domain2"
	domainTopoCache2.LocalNodeSN = "2-1"
	domainVLink23 := new(ncsnp.DomainVLink)
	domainVLink23.LocalDomainName = "Domain2"
	domainVLink23.LocalDomainId = 2
	domainVLink23.RemoteDomainName = "Domain3"
	domainVLink23.RemoteDomainId = 3
	domainVLink23.LocalNodeSN = "Node23"
	domainVLink23.RemoteNodeSN = "Node32"
	domainVLink23.AttachDomainId = 1023
	domainVLink12.AttachDomainName = "Fabric23"
	vLinkSlaAttr23 := new(ncsnp.VLinkSla)
	vLinkSlaAttr23.Delay = 2
	vLinkSlaAttr23.Bandwidth = 10000
	vLinkSlaAttr23.FreeBandwidth = 10000
	domainVLink23.VLinkSlaAttr = vLinkSlaAttr23
	domainTopoCache2.DomainVLinkArray = append(domainTopoCache2.DomainVLinkArray, domainVLink23)

	domainVLink21 := new(ncsnp.DomainVLink)
	domainVLink21.LocalDomainName = "Domain2"
	domainVLink21.LocalDomainId = 2
	domainVLink21.RemoteDomainName = "Domain1"
	domainVLink21.RemoteDomainId = 1
	domainVLink21.LocalNodeSN = "Node21"
	domainVLink21.RemoteNodeSN = "Node12"
	domainVLink21.AttachDomainId = 1012
	domainVLink12.AttachDomainName = "Fabric12"
	vLinkSlaAttr21 := new(ncsnp.VLinkSla)
	vLinkSlaAttr21.Delay = 1
	vLinkSlaAttr21.Bandwidth = 10000
	vLinkSlaAttr21.FreeBandwidth = 10000
	domainVLink21.VLinkSlaAttr = vLinkSlaAttr21
	domainTopoCache2.DomainVLinkArray = append(domainTopoCache2.DomainVLinkArray, domainVLink21)
	domainTopoCacheArry = append(domainTopoCacheArry, domainTopoCache2)
	content, err = proto.Marshal(domainTopoCache2)
	if err != nil {
		nputil.TraceErrorWithStack(err)
		return nil
	}
	domainTopoMsg[domainTopoCache2.LocalDomainName] = content

	domainTopoCache3 := new(ncsnp.DomainTopoCacheNotify)
	domainTopoCache3.LocalDomainId = 3
	domainTopoCache3.LocalDomainName = "Domain3"
	domainTopoCache3.LocalNodeSN = "3-1"
	domainVLink34 := new(ncsnp.DomainVLink)
	domainVLink34.LocalDomainName = "Domain3"
	domainVLink34.LocalDomainId = 3
	domainVLink34.RemoteDomainName = "Domain4"
	domainVLink34.RemoteDomainId = 4
	domainVLink34.LocalNodeSN = "Node34"
	domainVLink34.RemoteNodeSN = "Node43"
	domainVLink34.AttachDomainId = 1034
	domainVLink12.AttachDomainName = "Fabric34"
	vLinkSlaAttr34 := new(ncsnp.VLinkSla)
	vLinkSlaAttr34.Delay = 3
	vLinkSlaAttr34.Bandwidth = 10000
	vLinkSlaAttr34.FreeBandwidth = 10000
	domainVLink34.VLinkSlaAttr = vLinkSlaAttr34
	domainTopoCache3.DomainVLinkArray = append(domainTopoCache3.DomainVLinkArray, domainVLink34)
	domainVLink32 := new(ncsnp.DomainVLink)
	domainVLink32.LocalDomainName = "Domain3"
	domainVLink32.LocalDomainId = 3
	domainVLink32.RemoteDomainName = "Domain2"
	domainVLink32.RemoteDomainId = 2
	domainVLink32.LocalNodeSN = "Node32"
	domainVLink32.RemoteNodeSN = "Node23"
	domainVLink32.AttachDomainId = 1023
	domainVLink12.AttachDomainName = "Fabric23"
	vLinkSlaAttr32 := new(ncsnp.VLinkSla)
	vLinkSlaAttr32.Delay = 2
	vLinkSlaAttr32.Bandwidth = 10000
	vLinkSlaAttr32.FreeBandwidth = 10000
	domainVLink32.VLinkSlaAttr = vLinkSlaAttr32
	domainTopoCache3.DomainVLinkArray = append(domainTopoCache3.DomainVLinkArray, domainVLink32)
	domainTopoCacheArry = append(domainTopoCacheArry, domainTopoCache3)
	content, err = proto.Marshal(domainTopoCache3)
	if err != nil {
		nputil.TraceErrorWithStack(err)
		return nil
	}
	domainTopoMsg[domainTopoCache3.LocalDomainName] = content

	domainTopoCache4 := new(ncsnp.DomainTopoCacheNotify)
	domainTopoCache4.LocalDomainId = 4
	domainTopoCache4.LocalDomainName = "Domain4"
	domainTopoCache4.LocalNodeSN = "Node43"
	domainVLink43 := new(ncsnp.DomainVLink)
	domainVLink43.LocalDomainName = "Domain4"
	domainVLink43.LocalDomainId = 4
	domainVLink43.RemoteDomainName = "Domain3"
	domainVLink43.RemoteDomainId = 3
	domainVLink43.LocalNodeSN = "Node43"
	domainVLink43.RemoteNodeSN = "Node34"
	domainVLink43.AttachDomainId = 1034
	domainVLink12.AttachDomainName = "Fabric34"
	vLinkSlaAttr43 := new(ncsnp.VLinkSla)
	vLinkSlaAttr43.Delay = 3
	vLinkSlaAttr43.Bandwidth = 10000
	vLinkSlaAttr43.FreeBandwidth = 10000
	domainVLink43.VLinkSlaAttr = vLinkSlaAttr43
	domainTopoCache4.DomainVLinkArray = append(domainTopoCache4.DomainVLinkArray, domainVLink43)

	domainVLink42 := new(ncsnp.DomainVLink)
	domainVLink42.LocalDomainName = "Domain4"
	domainVLink42.LocalDomainId = 4
	domainVLink42.RemoteDomainName = "Domain2"
	domainVLink42.RemoteDomainId = 2
	domainVLink42.LocalNodeSN = "Node42"
	domainVLink42.RemoteNodeSN = "Node24"
	domainVLink42.AttachDomainId = 1024
	domainVLink12.AttachDomainName = "Fabric24"
	vLinkSlaAttr42 := new(ncsnp.VLinkSla)
	vLinkSlaAttr42.Delay = 2
	vLinkSlaAttr42.Bandwidth = 10000
	vLinkSlaAttr42.FreeBandwidth = 10000
	domainVLink42.VLinkSlaAttr = vLinkSlaAttr42
	domainTopoCache4.DomainVLinkArray = append(domainTopoCache4.DomainVLinkArray, domainVLink42)
	domainTopoCacheArry = append(domainTopoCacheArry, domainTopoCache4)
	content, err = proto.Marshal(domainTopoCache4)
	if err != nil {
		nputil.TraceErrorWithStack(err)
		return nil
	}
	domainTopoMsg[domainTopoCache4.LocalDomainName] = content

	fmt.Printf("Len of domainTopoCache is (+%d)\n", len(domainTopoCacheArry))
	for i, domainTopoCache := range domainTopoCacheArry {
		fmt.Printf("domainTopoCacheArry[%d] is (%+v)\n", i, domainTopoCache)
	}

	nputil.TraceInfoEnd("------------------------------------------------------")
	return domainTopoMsg
}