terraform {
    required_providers {
        hyperv = {
            source  = "taliesins/hyperv"
            version = "1.0.4"
        }
    }
}

/*provider "hyperv" {
    user        = "pm-admin"
    password    = "Ruebennase01"
    host        = "172.16.10.107"
    port        = 5985
    https       = false
    insecure    = true
    use_ntlm    = true
    script_path = "C:/Temp/terraform_%RAND%.cmd"
    timeout     = "40s"
}*/

provider "hyperv" {
    user        = "Administrator"
    password    = "Ruebennase01"
    host        = "172.16.14.27"
    port        = 5985
    https       = false
    insecure    = true
    use_ntlm    = true
    script_path = "C:/Temp/terraform_%RAND%.cmd"
    timeout     = "10000s"
}

/*resource "hyperv_dvd" "cp_dvd" {
    path        = "c:\\users\\administrator\\documents\\vms\\pm-vm-microk8s-test\\virtual hard disks\\test.iso"
    ip          = "172.16.14.84"
}*/

resource "hyperv_vhd" "base_vhdx" {
    path   = "c:\\users\\administrator\\documents\\vms\\parentdisks\\pm-vm-microk8s-test.vhdx"
    source = "http://172.16.14.28/repository/vm-images/vhdx/microk8s-controlplane.vhdx"
}

resource "hyperv_vhd" "win_test_vhdx" {
    path   = "c:\\users\\administrator\\documents\\vms\\pm-vm-microk8s-test\\virtual hard disks\\pm-vm-microk8s-test.vhdx"
    parent_path = hyperv_vhd.base_vhdx.path
    #source = "http://172.16.14.28/repository/vm-images/vhdx/microk8s-controlplane.vhdx"
    #parent_path = "http://172.16.14.28/repository/vm-images/vhdx/microk8s-controlplane.vhdx"
    vhd_type = "Differencing"
}

# Create a server
resource "hyperv_machine_instance" "win_test" {
    name                 = "pm-vm-microk8s-test"
    static_memory        = true
    path                 = "c:\\users\\administrator\\documents\\vms"
    processor_count      = 2
    memory_startup_bytes = 2294967296

    hard_disk_drives {
        controller_location = "0"
        controller_number   = "0"
        path                = hyperv_vhd.win_test_vhdx.path
    }

    # Create dvd drive
    /*dvd_drives {
        controller_number   = "0"
        controller_location = "1"
        path                = hyperv_dvd.cp_dvd.path
    }*/

    vm_firmware {
        enable_secure_boot = "Off"
        boot_order {
            boot_type           = "HardDiskDrive"
            controller_number   = "0"
            controller_location = "0"
        }
    }

    network_adaptors {
        name        = "wan"
        switch_name = "test-switch"
    }
}

resource "hyperv_vhd" "win_test2_vhdx" {
    path   = "c:\\users\\administrator\\documents\\vms\\pm-vm-microk8s-test2\\virtual hard disks\\pm-vm-microk8s-test2.vhdx"
    parent_path = hyperv_vhd.base_vhdx.path
    #source = "http://172.16.14.28/repository/vm-images/vhdx/microk8s-controlplane.vhdx"
    #parent_path = "http://172.16.14.28/repository/vm-images/vhdx/microk8s-controlplane.vhdx"
    vhd_type = "Differencing"
}

# Create a server
resource "hyperv_machine_instance" "win_test2" {
    name                 = "pm-vm-microk8s-test2"
    static_memory        = true
    path                 = "c:\\users\\administrator\\documents\\vms"
    processor_count      = 2
    memory_startup_bytes = 2294967296

    hard_disk_drives {
        controller_location = "0"
        controller_number   = "0"
        path                = hyperv_vhd.win_test2_vhdx.path
    }

    # Create dvd drive
    /*dvd_drives {
        controller_number   = "0"
        controller_location = "1"
        path                = hyperv_dvd.cp_dvd.path
    }*/

    vm_firmware {
        enable_secure_boot = "Off"
        boot_order {
            boot_type           = "HardDiskDrive"
            controller_number   = "0"
            controller_location = "0"
        }
    }

    network_adaptors {
        name        = "wan"
        switch_name = "test-switch"
    }
}
