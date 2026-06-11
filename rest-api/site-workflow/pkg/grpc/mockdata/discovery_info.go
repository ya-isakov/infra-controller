// SPDX-FileCopyrightText: Copyright (c) 2026 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package mockdata

import wflows "github.com/NVIDIA/infra-controller/rest-api/workflow-schema/schema/site-agent/workflows/v1"

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

// MachineDiscoveryInfo returns hardware discovery metadata for mocked test machines.
func MachineDiscoveryInfo() *wflows.DiscoveryInfo {
	machineArch := wflows.CpuArchitecture_X86_64

	return &wflows.DiscoveryInfo{
		NetworkInterfaces: []*wflows.NetworkInterface{
			{
				MacAddress: "58:A2:E1:5B:D1:B0",
				PciProperties: pciProperties(
					"Mellanox Technologies",
					"MT43244 BlueField-3 integrated ConnectX-7 network controller",
					"/devices/pci0000:00/0000:00:01.3/0000:01:00.0/net/enp1s0np0",
					0,
					"MT43244 BlueField-3 integrated ConnectX-7 network controller",
					"0000:01:00.0",
				),
				Lldp: &wflows.NetworkInterfaceLldp{
					PortId:           "Eth1/1",
					SwitchId:         strPtr("00:01:00:00:02:00"),
					SwitchSystemName: "leaf-0",
				},
			},
			{
				MacAddress: "6C:B3:11:8D:8F:70",
				PciProperties: pciProperties(
					"Intel Corporation",
					"I350 Gigabit Network Connection",
					"/devices/pci0000:a0/0000:a0:01.3/0000:a3:00.0/net/ens11f0",
					0,
					"I350 Gigabit Network Connection",
					"0000:a3:00.0",
				),
				Lldp: &wflows.NetworkInterfaceLldp{
					PortId:           "Eth1/2",
					SwitchId:         strPtr("00:01:00:00:02:00"),
					SwitchSystemName: "leaf-0",
				},
			},
			{
				MacAddress: "6C:B3:11:8D:8F:71",
				PciProperties: pciProperties(
					"Intel Corporation",
					"I350 Gigabit Network Connection",
					"/devices/pci0000:a0/0000:a0:01.3/0000:a3:00.1/net/ens11f1",
					0,
					"I350 Gigabit Network Connection",
					"0000:a3:00.1",
				),
				Lldp: &wflows.NetworkInterfaceLldp{
					PortId:           "Eth1/11",
					SwitchId:         strPtr("00:01:00:00:02:00"),
					SwitchSystemName: "leaf-0",
				},
			},
		},
		BlockDevices: []*wflows.BlockDevice{
			{Model: "SAMSUNG MZQL23T8HCLS-00A07", Revision: "GDC5A02Q", Serial: "S64HNT0Y201196", DeviceType: "disk"},
			{Model: "SAMSUNG MZQL23T8HCLS-00A07", Revision: "GDC5A02Q", Serial: "S64HNT0Y201195", DeviceType: "disk"},
			{Model: "SAMSUNG MZQL23T8HCLS-00A07", Revision: "GDC5A02Q", Serial: "S64HNT0Y202243", DeviceType: "disk"},
			{Model: "SAMSUNG MZQL23T8HCLS-00A07", Revision: "GDC5A02Q", Serial: "S64HNT0Y201188", DeviceType: "disk"},
			{Model: "SAMSUNG MZQL21T9HCJR-00A07", Revision: "GDC5A02Q", Serial: "S64GNN0WB19128", DeviceType: "disk"},
			{Model: "SAMSUNG MZQL21T9HCJR-00A07", Revision: "GDC5A02Q", Serial: "S64GNN0WB19130", DeviceType: "disk"},
			{Model: "SAMSUNG MZQL2960HCJR-00A07", Revision: "GDC5A02Q", Serial: "S64FNN0XB34201", DeviceType: "disk"},
		},
		MachineType: "x86_64",
		MachineArch: &machineArch,
		NvmeDevices: []*wflows.NvmeDevice{
			{Model: "SAMSUNG MZQL23T8HCLS-00A07", FirmwareRev: "GDC5A02Q", Serial: "S64HNT0Y201196"},
			{Model: "SAMSUNG MZQL23T8HCLS-00A07", FirmwareRev: "GDC5A02Q", Serial: "S64HNT0Y201195"},
			{Model: "SAMSUNG MZQL23T8HCLS-00A07", FirmwareRev: "GDC5A02Q", Serial: "S64HNT0Y202243"},
			{Model: "SAMSUNG MZQL23T8HCLS-00A07", FirmwareRev: "GDC5A02Q", Serial: "S64HNT0Y201188"},
			{Model: "SAMSUNG MZQL21T9HCJR-00A07", FirmwareRev: "GDC5A02Q", Serial: "S64GNN0WB19128"},
			{Model: "SAMSUNG MZQL21T9HCJR-00A07", FirmwareRev: "GDC5A02Q", Serial: "S64GNN0WB19130"},
			{Model: "SAMSUNG MZQL2960HCJR-00A07", FirmwareRev: "GDC5A02Q", Serial: "S64FNN0XB34201"},
		},
		DmiData: &wflows.DmiData{
			BoardName:     "MZG3-GU0-000",
			BoardVersion:  "03000300",
			BiosVersion:   "R23_F20",
			ProductSerial: "DPG5NS621A0001",
			BoardSerial:   "PK1N6300085",
			ChassisSerial: "2451R26302R1.0U1057",
			BiosDate:      "04/01/2026",
			ProductName:   "R263-ZG0-AAL2-000",
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
				Guid: "0000000000000000",
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
				Guid: "0000000000000001",
			},
		},
		Gpus: []*wflows.Gpu{
			{
				Name:           "NVIDIA A30",
				Serial:         "1651922012475",
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

// EnsureMachineDiscoveryInfo backfills hardware discovery metadata when absent.
func EnsureMachineDiscoveryInfo(machine *wflows.Machine) {
	if machine == nil || machine.DiscoveryInfo != nil {
		return
	}
	machine.DiscoveryInfo = MachineDiscoveryInfo()
}
