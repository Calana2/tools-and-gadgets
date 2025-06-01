// GOOS=windows GOARCH=amd64 go build -o ovm.exe main.go

package main
import (
	"fmt"
	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"
)

// key types
const (
	REG_SZ = iota
	REG_MULTI_SZ
	REG_BINARY
	REG_DWORD
)

// root keys
const (
	HKEY_CLASSES_ROOT_STR     = "HKEY_CLASSES_ROOT"
	HKEY_CURRENT_USER_STR     = "HKEY_CURRENT_USER"
	HKEY_LOCAL_MACHINE_STR    = "HKEY_LOCAL_MACHINE"
	HKEY_USERS_STR            = "HKEY_USERS"
	HKEY_CURRENT_CONFIG_STR   = "HKEY_CURRENT_CONFIG"
	HKEY_PERFORMANCE_DATA_STR = "HKEY_PERFORMANCE_DATA"
)

// key validator
type K struct {
	root      registry.Key
	path      string
	valueName string
	regType   int
	cmpstr    string
}

/* Retrieve root registry key in string format */
func rootKeyToString(rk registry.Key) string {
	switch rk {
	case registry.Key(windows.HKEY_CLASSES_ROOT):
		return HKEY_CLASSES_ROOT_STR
	case registry.Key(windows.HKEY_CURRENT_USER):
		return HKEY_CURRENT_USER_STR
	case registry.Key(windows.HKEY_LOCAL_MACHINE):
		return HKEY_LOCAL_MACHINE_STR
	case registry.Key(windows.HKEY_USERS):
		return HKEY_USERS_STR
	case registry.Key(windows.HKEY_CURRENT_CONFIG):
		return HKEY_CURRENT_CONFIG_STR
	case registry.Key(windows.HKEY_PERFORMANCE_DATA):
		return HKEY_PERFORMANCE_DATA_STR
	default:
		return "UNKNOWN_ROOT_KEY"
	}
}

/* Reads a key from registry and returns the value and a formatted string */
func keyReader(rk registry.Key, path string, name string, valueType int) (interface{}, string, error) {
	key, err := registry.OpenKey(rk, path, registry.QUERY_VALUE)
	if err != nil {
		return nil, "", fmt.Errorf("Error opening key: %v", err)
	}
	defer key.Close()
	var value interface{}
	switch valueType {
	case REG_SZ:
		value, _, err = key.GetStringValue(name)
	case REG_MULTI_SZ:
		value, _, err = key.GetStringsValue(name)
	case REG_BINARY:
		value, _, err = key.GetBinaryValue(name)
	case REG_DWORD:
		value, _, err = key.GetIntegerValue(name)
	default:
		//fmt.Printf("Not supported value: %d\n",valueType)
	}
	if err != nil {
		//fmt.Printf("Error reading the value: %v\n",err)
	}
	fullPath := fmt.Sprintf("%s\\%s\\%s", rootKeyToString(rk), path, name)
	if value != nil {
		return value, fullPath, nil
	}
	return nil, fullPath, nil
}

func detectVM(usefulKeys []*K) {
	// Keys related to VM
	var keysFound []string
	RegValuePath := []string{"HARDWARE\\DEVICEMAP\\Scsi\\Scsi Port 0\\Scsi Bus 0\\Target Id 0\\Logical Unit Id 0",
		"HARDWARE\\DEVICEMAP\\Scsi\\Scsi Port 1\\Scsi Bus 0\\Target Id 0\\Logical Unit Id 0",
		"HARDWARE\\DEVICEMAP\\Scsi\\Scsi Port 2\\Scsi Bus 0\\Target Id 0\\Logical Unit Id 0",
		"SOFTWARE\\VMware, Inc.\\VMware Tools",
		"HARDWARE\\Description\\System",
		"SOFTWARE\\Oracle\\VirtualBox Guest Additions",
		"SYSTEM\\ControlSet001\\Services\\Disk\\Enum",
		"HARDWARE\\ACPI\\DSDT\\VBOX__",
		"HARDWARE\\ACPI\\FADT\\VBOX__",
		"HARDWARE\\ACPI\\RSDT\\VBOX__",
		"SYSTEM\\ControlSet001\\Services\\VBoxGuest",
		"SYSTEM\\ControlSet001\\Services\\VBoxMouse",
		"SYSTEM\\ControlSet001\\Services\\VBoxService",
		"SYSTEM\\ControlSet001\\Services\\VBoxSF",
		"SYSTEM\\ControlSet001\\Services\\VBoxVideo",
	}
	for _, rvp := range RegValuePath {
		_, path, err := keyReader(registry.LOCAL_MACHINE, rvp, "", REG_SZ)
		if err != nil {
			continue
		}
		keysFound = append(keysFound, path)
	}
	// Useful Keys
	var uKeysFound []string
	if usefulKeys != nil {
		for _, k := range usefulKeys {
			v, path, err := keyReader(k.root, k.path, k.valueName, k.regType)
			if err != nil {
				continue
			}
			uKeysFound = append(uKeysFound, fmt.Sprintf("%s: %v\n", path, v))
		}
	}
	// Dump
  if keysFound == nil && uKeysFound == nil {
    fmt.Println("No Virtual Environment Artefacts detected.")
    return
  }
	fmt.Println("\n\n====================== VM Keys found ======================")
	for _, k := range keysFound {
		fmt.Println(k)
	}
	fmt.Println("====================== xx xxxx xxxxx ======================")
	fmt.Println("\n\n====================== VM Modified Keys found ======================")
	for _, k := range uKeysFound {
		fmt.Println(k)
	}
	fmt.Println("====================== xx xxxx xxxxx ======================")
}

func main() {
	usefulKeys := []*K{
		&K{root: registry.LOCAL_MACHINE, path: "HARDWARE\\DEVICEMAP\\Scsi\\Scsi Port 0\\Scsi Bus 0\\Target Id 0\\Logical Unit Id 0", valueName: "Identifier", regType: REG_SZ, cmpstr: "VMware"},
		&K{root: registry.LOCAL_MACHINE, path: "HARDWARE\\DEVICEMAP\\Scsi\\Scsi Port 0\\Scsi Bus 0\\Target Id 0\\Logical Unit Id 0", valueName: "Identifier", regType: REG_SZ, cmpstr: "VBOX"},
		&K{root: registry.LOCAL_MACHINE, path: "HARDWARE\\DEVICEMAP\\Scsi\\Scsi Port 0\\Scsi Bus 0\\Target Id 0\\Logical Unit Id 0", valueName: "Identifier", regType: REG_SZ, cmpstr: "QEMU"},
		&K{root: registry.LOCAL_MACHINE, path: "HARDWARE\\DEVICEMAP\\Scsi\\Scsi Port 1\\Scsi Bus 0\\Target Id 0\\Logical Unit Id 0", valueName: "Identifier", regType: REG_SZ, cmpstr: "Vmware"},
		&K{root: registry.LOCAL_MACHINE, path: "HARDWARE\\DEVICEMAP\\Scsi\\Scsi Port 2\\Scsi Bus 0\\Target Id 0\\Logical Unit Id 0", valueName: "Identifier", regType: REG_SZ, cmpstr: "Vmware"},

		&K{root: registry.LOCAL_MACHINE, path: "HARDWARE\\Description\\System", valueName: "SystemBiosVersion", regType: REG_MULTI_SZ, cmpstr: "VBOX"},
		&K{root: registry.LOCAL_MACHINE, path: "HARDWARE\\Description\\System", valueName: "VideoBiosVersion", regType: REG_MULTI_SZ, cmpstr: "VIRTUALBOX"},
		&K{root: registry.LOCAL_MACHINE, path: "HARDWARE\\Description\\System", valueName: "VideoBiosVersion", regType: REG_MULTI_SZ, cmpstr: "QEMU"},
		&K{root: registry.LOCAL_MACHINE, path: "HARDWARE\\Description\\System", valueName: "SystemBiosDate", regType: REG_SZ, cmpstr: "06/23/99"},
	}
  fmt.Println("====================== Windows Sandbox Detector ======================")
	detectVM(usefulKeys)
}
