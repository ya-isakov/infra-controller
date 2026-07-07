// SPDX-FileCopyrightText: Copyright (c) 2026 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package mockdata

import (
	"fmt"
	"hash/fnv"
	"strconv"
	"strings"

	wflows "github.com/NVIDIA/infra-controller/rest-api/workflow-schema/schema/site-agent/workflows/v1"
)

const (
	// MockHostCount is the number of preloaded mock hosts (IDs 0..MockHostCount-1).
	MockHostCount = 6

	// mockMachineUUIDPrefix is the fixed UUID prefix for preloaded mock hosts.
	// The host index is encoded in the last 12 hex digits of the UUID.
	mockMachineUUIDPrefix = "00000000-0000-4000-8000-"
)

func strPtr(s string) *string {
	return &s
}

func pciProperties(vendor, device, path string, numaNode int32, description, slot string) *wflows.PciDeviceProperties {
	return &wflows.PciDeviceProperties{
		Vendor:      vendor,
		Device:      device,
		Path:        path,
		NumaNode:    numaNode,
		Description: strPtr(description),
		Slot:        strPtr(slot),
	}
}

func memoryDevice(sizeMB uint32, memType string) *wflows.MemoryDevice {
	size := sizeMB
	return &wflows.MemoryDevice{
		SizeMb:  &size,
		MemType: strPtr(memType),
	}
}

// MockMachineID returns a stable UUID for a mock host index.
// Host ID is encoded in the last 12 hex digits, e.g. host 3 → …-000000000003.
func MockMachineID(hostID int) string {
	if hostID < 0 {
		hostID = 0
	}
	hostID = hostID % MockHostCount
	return fmt.Sprintf("%s%012x", mockMachineUUIDPrefix, hostID)
}

// HostIDFromMachineID maps a machine ID to a mock host index in [0, MockHostCount).
func HostIDFromMachineID(machineID string) int {
	machineID = strings.ToLower(strings.TrimSpace(machineID))
	if strings.HasPrefix(machineID, mockMachineUUIDPrefix) {
		suffix := machineID[len(mockMachineUUIDPrefix):]
		if n, err := strconv.ParseUint(suffix, 16, 64); err == nil && int(n) < MockHostCount {
			return int(n)
		}
	}

	h := fnv.New32a()
	_, _ = h.Write([]byte(machineID))
	return int(h.Sum32() % MockHostCount)
}

// MockHostname returns a host-specific hostname for mocked machines.
func MockHostname(hostID int) string {
	if hostID < 0 {
		hostID = 0
	}
	hostID = hostID % MockHostCount
	return fmt.Sprintf("mock-host-%d.nico.nvidia.com", hostID)
}

// ufoSimMAC mirrors the MAC formula ufo-simulator's create-vm.yml uses for VM NICs:
// vm_base_mac:00:(host_id*vm_port_count + ifaceOffset).
func ufoSimMAC(hostID, ifaceOffset int) string {
	return fmt.Sprintf("52:54:00:12:00:%02X", hostID*4+ifaceOffset)
}

// mockMAC generates a stable, vendor-scoped MAC for interfaces that ufo-simulator
// doesn't assign a predefined MAC to, such as the BlueField mgmt/OOB port.
func mockMAC(prefix string, lastOctet byte, hostID int) string {
	return fmt.Sprintf("%s:%02X", prefix, lastOctet+byte(hostID))
}

var mockSwitches = []struct {
	id   string
	name string
}{
	{"00:01:00:00:01:00", "leaf-0"},
	{"00:01:00:00:01:01", "leaf-1"},
}

// mockDataSwitchPort returns the LLDP port ID a verity switch would report for a mock
// host's Nth Intel data-NIC link (linkIndex 0 or 1) into a leaf switch, mirroring
// ufo-simulator's leaf-0/leaf-1 swp(2*host_id+linkIndex+1) topology translated to
// verity's "eth1/N" port naming (see ufo-simulator/ansible/vars/verity.yml).
func mockDataSwitchPort(hostID, linkIndex int) string {
	return fmt.Sprintf("eth1/%d", 2*hostID+linkIndex+1)
}

func mockIBGUID(base uint64, hostID int) string {
	return fmt.Sprintf("%016X", base+uint64(hostID))
}

// MachineDiscoveryInfo returns hardware discovery metadata for mock host 0.
func MachineDiscoveryInfo() *wflows.DiscoveryInfo {
	return MachineDiscoveryInfoForHost(0)
}

// MachineDiscoveryInfoForHost returns host-specific hardware discovery metadata.
func MachineDiscoveryInfoForHost(hostID int) *wflows.DiscoveryInfo {
	if hostID < 0 {
		hostID = 0
	}
	hostID = hostID % MockHostCount

	machineArch := wflows.CpuArchitecture_X86_64
	diskSerialBase := 201196 + hostID

	return &wflows.DiscoveryInfo{
		NetworkInterfaces: []*wflows.NetworkInterface{
			{
				MacAddress: mockMAC("58:A2:E1:5B:D1", 0xB0, hostID),
				PciProperties: pciProperties(
					"Mellanox Technologies",
					"MT43244 BlueField-3 integrated ConnectX-7 network controller",
					"/devices/pci0000:00/0000:00:01.3/0000:01:00.0/net/eth0",
					0,
					"MT43244 BlueField-3 integrated ConnectX-7 network controller",
					"0000:01:00.0",
				),
			},
			{
				MacAddress: ufoSimMAC(hostID, 0),
				PciProperties: pciProperties(
					"Intel Corporation",
					"I350 Gigabit Network Connection",
					"/devices/pci0000:a0/0000:a0:01.3/0000:a3:00.0/net/eth1",
					0,
					"I350 Gigabit Network Connection",
					"0000:a3:00.0",
				),
				Lldp: &wflows.NetworkInterfaceLldp{
					PortId:           "ifname=" + mockDataSwitchPort(hostID, 0),
					SwitchId:         strPtr("mac=" + mockSwitches[0].id),
					SwitchSystemName: mockSwitches[0].name,
				},
			},
			{
				MacAddress: ufoSimMAC(hostID, 1),
				PciProperties: pciProperties(
					"Intel Corporation",
					"I350 Gigabit Network Connection",
					"/devices/pci0000:a0/0000:a0:01.3/0000:a3:00.1/net/eth2",
					0,
					"I350 Gigabit Network Connection",
					"0000:a3:00.1",
				),
				Lldp: &wflows.NetworkInterfaceLldp{
					PortId:           "ifname=" + mockDataSwitchPort(hostID, 1),
					SwitchId:         strPtr("mac=" + mockSwitches[0].id),
					SwitchSystemName: mockSwitches[0].name,
				},
			},
			{
				MacAddress: ufoSimMAC(hostID, 2),
				PciProperties: pciProperties(
					"Intel Corporation",
					"I350 Gigabit Network Connection",
					"/devices/pci0000:a0/0000:a0:01.3/0000:a3:00.2/net/eth3",
					0,
					"I350 Gigabit Network Connection",
					"0000:a3:00.2",
				),
				Lldp: &wflows.NetworkInterfaceLldp{
					PortId:           "ifname=" + mockDataSwitchPort(hostID, 0),
					SwitchId:         strPtr("mac=" + mockSwitches[1].id),
					SwitchSystemName: mockSwitches[1].name,
				},
			},
			{
				MacAddress: ufoSimMAC(hostID, 3),
				PciProperties: pciProperties(
					"Intel Corporation",
					"I350 Gigabit Network Connection",
					"/devices/pci0000:a0/0000:a0:01.3/0000:a3:00.3/net/eth4",
					0,
					"I350 Gigabit Network Connection",
					"0000:a3:00.3",
				),
				Lldp: &wflows.NetworkInterfaceLldp{
					PortId:           "ifname=" + mockDataSwitchPort(hostID, 1),
					SwitchId:         strPtr("mac=" + mockSwitches[1].id),
					SwitchSystemName: mockSwitches[1].name,
				},
			},
		},
		BlockDevices: []*wflows.BlockDevice{
			{Model: "SAMSUNG MZQL23T8HCLS-00A07", Revision: "GDC5A02Q", Serial: fmt.Sprintf("S64HNT0Y2%06d", diskSerialBase), DeviceType: "disk"},
			{Model: "SAMSUNG MZQL23T8HCLS-00A07", Revision: "GDC5A02Q", Serial: fmt.Sprintf("S64HNT0Y2%06d", diskSerialBase+1), DeviceType: "disk"},
			{Model: "SAMSUNG MZQL23T8HCLS-00A07", Revision: "GDC5A02Q", Serial: fmt.Sprintf("S64HNT0Y2%06d", diskSerialBase+2), DeviceType: "disk"},
			{Model: "SAMSUNG MZQL23T8HCLS-00A07", Revision: "GDC5A02Q", Serial: fmt.Sprintf("S64HNT0Y2%06d", diskSerialBase+3), DeviceType: "disk"},
			{Model: "SAMSUNG MZQL21T9HCJR-00A07", Revision: "GDC5A02Q", Serial: fmt.Sprintf("S64GNN0W%06d", 19128+hostID), DeviceType: "disk"},
			{Model: "SAMSUNG MZQL21T9HCJR-00A07", Revision: "GDC5A02Q", Serial: fmt.Sprintf("S64GNN0W%06d", 19130+hostID), DeviceType: "disk"},
			{Model: "SAMSUNG MZQL2960HCJR-00A07", Revision: "GDC5A02Q", Serial: fmt.Sprintf("S64FNN0X%06d", 34201+hostID), DeviceType: "disk"},
		},
		MachineType: "x86_64",
		MachineArch: &machineArch,
		NvmeDevices: []*wflows.NvmeDevice{
			{Model: "SAMSUNG MZQL23T8HCLS-00A07", FirmwareRev: "GDC5A02Q", Serial: fmt.Sprintf("S64HNT0Y2%06d", diskSerialBase)},
			{Model: "SAMSUNG MZQL23T8HCLS-00A07", FirmwareRev: "GDC5A02Q", Serial: fmt.Sprintf("S64HNT0Y2%06d", diskSerialBase+1)},
			{Model: "SAMSUNG MZQL23T8HCLS-00A07", FirmwareRev: "GDC5A02Q", Serial: fmt.Sprintf("S64HNT0Y2%06d", diskSerialBase+2)},
			{Model: "SAMSUNG MZQL23T8HCLS-00A07", FirmwareRev: "GDC5A02Q", Serial: fmt.Sprintf("S64HNT0Y2%06d", diskSerialBase+3)},
			{Model: "SAMSUNG MZQL21T9HCJR-00A07", FirmwareRev: "GDC5A02Q", Serial: fmt.Sprintf("S64GNN0W%06d", 19128+hostID)},
			{Model: "SAMSUNG MZQL21T9HCJR-00A07", FirmwareRev: "GDC5A02Q", Serial: fmt.Sprintf("S64GNN0W%06d", 19130+hostID)},
			{Model: "SAMSUNG MZQL2960HCJR-00A07", FirmwareRev: "GDC5A02Q", Serial: fmt.Sprintf("S64FNN0X%06d", 34201+hostID)},
		},
		DmiData: &wflows.DmiData{
			BoardName:     fmt.Sprintf("MZG3-GU0-00%d", hostID),
			BoardVersion:  fmt.Sprintf("0300030%d", hostID),
			BiosVersion:   "R23_F20",
			ProductSerial: fmt.Sprintf("DPG5NS621A%04d", hostID),
			BoardSerial:   fmt.Sprintf("PK1N6300%03d", 85+hostID),
			ChassisSerial: fmt.Sprintf("2451R26302R1.0U10%02d", 57+hostID),
			BiosDate:      "04/01/2026",
			ProductName:   fmt.Sprintf("R263-ZG0-AAL2-00%d", hostID),
			SysVendor:     "Giga Computing",
		},
		InfinibandInterfaces: []*wflows.InfinibandInterface{
			{
				PciProperties: pciProperties(
					"Mellanox Technologies",
					"MT28800 Family [ConnectX-5 Ex]",
					"/devices/pci0000:e0/0000:e0:01.4/0000:e1:00.0/infiniband/ibp225s0f0",
					0,
					"MT28800 Family [ConnectX-5 Ex]",
					"0000:e1:00.0",
				),
				Guid: mockIBGUID(0, hostID),
			},
			{
				PciProperties: pciProperties(
					"Mellanox Technologies",
					"MT28800 Family [ConnectX-5 Ex]",
					"/devices/pci0000:e0/0000:e0:01.4/0000:e1:00.1/infiniband/ibp225s0f1",
					0,
					"MT28800 Family [ConnectX-5 Ex]",
					"0000:e1:00.1",
				),
				Guid: mockIBGUID(1, hostID),
			},
		},
		Gpus: []*wflows.Gpu{
			{
				Name:           "NVIDIA A30",
				Serial:         fmt.Sprintf("1651922012%03d", 475+hostID),
				DriverVersion:  "580.126.16",
				VbiosVersion:   "92.00.66.00.04",
				InforomVersion: "1001.0205.00.02",
				TotalMemory:    "24576 MiB",
				Frequency:      "930 MHz",
				PciBusId:       "00000000:C1:00.0",
			},
		},
		MemoryDevices: []*wflows.MemoryDevice{
			memoryDevice(32768, "DDR5"),
			memoryDevice(32768, "DDR5"),
			memoryDevice(32768, "DDR5"),
			memoryDevice(32768, "DDR5"),
			memoryDevice(32768, "DDR5"),
			memoryDevice(32768, "DDR5"),
			memoryDevice(32768, "DDR5"),
			memoryDevice(32768, "DDR5"),
		},
		CpuInfo: []*wflows.CpuInfo{
			{
				Model:   "AMD EPYC 9115 16-Core Processor",
				Vendor:  "AuthenticAMD",
				Sockets: 1,
				Cores:   16,
				Threads: 16,
			},
		},
	}
}

// EnsureMachineDiscoveryInfo backfills host-specific hardware discovery metadata when absent.
func EnsureMachineDiscoveryInfo(machine *wflows.Machine) {
	if machine == nil || machine.DiscoveryInfo != nil {
		return
	}
	hostID := 0
	if machine.Id != nil {
		hostID = HostIDFromMachineID(machine.Id.GetId())
	}
	machine.DiscoveryInfo = MachineDiscoveryInfoForHost(hostID)
}
