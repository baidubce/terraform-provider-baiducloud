package bec

import (
	"fmt"
	"log"

	"github.com/baidubce/bce-sdk-go/services/bec"
	"github.com/baidubce/bce-sdk-go/services/bec/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func ResourceVMInstance() *schema.Resource {
	return &schema.Resource{
		Description: "Use this resource to manage BEC VM Instance. \n\n" +
			"More information can be found in the [Developer Guide](https://cloud.baidu.com/doc/BEC/s/jknpo0evo). \n\n",

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Create: resourceVMInstanceCreate,
		Read:   resourceVMInstanceRead,
		Update: resourceVMInstanceUpdate,
		Delete: resourceVMInstanceDelete,

		Schema: map[string]*schema.Schema{
			"service_id": {
				Type:        schema.TypeString,
				Description: "ID of the vm instance group.",
				Required:    true,
			},
			"vm_name": {
				Type:        schema.TypeString,
				Description: "Name of the vm instance. If empty, system will assign one.",
				Optional:    true,
				Computed:    true,
			},
			"host_name": {
				Type:        schema.TypeString,
				Description: "Host name of the vm instance. If empty, system will assign one.",
				Optional:    true,
				Computed:    true,
			},
			"region_id": {
				Type:        schema.TypeString,
				Description: "Node ID, composed of lowercase letters of [`country`-`city`-`isp`]. Can be obtained through data source `baiducloud_bec_nodes`.",
				Required:    true,
				ForceNew:    true,
			},
			"spec": {
				Type:        schema.TypeString,
				Description: "Specification family.",
				Optional:    true,
			},
			"cpu": {
				Type:         schema.TypeInt,
				Description:  "CPU core count of the vm instance. At least 1 core.",
				Required:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},
			"memory": {
				Type:         schema.TypeInt,
				Description:  "Memory size (GB) of the vm instance. At least 1 GB.",
				Required:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},
			"image_type": {
				Type:         schema.TypeString,
				Description:  "Valid values: `bec`(public image or bec custom image), `bcc`(bcc custom image)",
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"bcc", "bec"}, false),
			},
			"image_id": {
				Type:        schema.TypeString,
				Description: "ID of the image.",
				Required:    true,
				ForceNew:    true,
			},
			"system_volume": {
				Type:        schema.TypeList,
				Description: "System volume config of the vm instance.",
				Required:    true,
				MaxItems:    1,
				Elem:        VolumeConfigSchema(),
			},
			"data_volume": {
				Type:        schema.TypeList,
				Description: "Data volume config of the vm instance.",
				Optional:    true,
				Elem:        VolumeConfigSchema(),
			},
			"need_public_ip": {
				Type:        schema.TypeBool,
				Description: "Whether to open public network. Defaults to `false`.",
				Optional:    true,
				Default:     false,
				ForceNew:    true,
			},
			"need_ipv6_public_ip": {
				Type:             schema.TypeBool,
				Description:      "Whether to open IPv6 public network. Defaults to `false`.",
				Optional:         true,
				Default:          false,
				DiffSuppressFunc: needPublicIPDiffSuppress,
			},
			"bandwidth": {
				Type:             schema.TypeInt,
				Description:      "Public network bandwidth size (Mbps).",
				Optional:         true,
				DiffSuppressFunc: needPublicIPDiffSuppress,
			},
			"dns_config": {
				Type:        schema.TypeList,
				Description: "DNS config of the vm instance.",
				Required:    true,
				MaxItems:    1,
				Elem:        DNSConfigSchema(),
			},
			"network_config": {
				Type:        schema.TypeList,
				Description: "Network config of the vm instance. If not set, system will use default network config.",
				Optional:    true,
				Elem:        NetworkConfigListReadSchema(),
			},
			"key_config": {
				Type:        schema.TypeList,
				Description: "Password or keypair config of the vm instance.",
				Required:    true,
				MaxItems:    1,
				Elem:        KeyConfigSchema(),
			},
			// computed
			"status": {
				Type: schema.TypeString,
				Description: "Status of the vm instance. Possible values: `CREATING`, `RUNNING`, `STOPPING`, `STOPPED`, `RESTARTING`, " +
					"`REINSTALLING`, `STARTING`, `IMAGING`, `FAILED`, `UNKNOWN`",
				Computed: true,
			},
			"internal_ip": {
				Type:        schema.TypeString,
				Description: "Local network IPv4 address of the vm instance.",
				Computed:    true,
			},
			"public_ip": {
				Type:        schema.TypeString,
				Description: "Public network IPv4 address of the vm instance.",
				Computed:    true,
			},
			"ipv6_public_ip": {
				Type:        schema.TypeString,
				Description: "Public network IPv6 address of the vm instance.",
				Computed:    true,
			},
			"create_time": {
				Type:        schema.TypeString,
				Description: "Creation time of the vm instance.",
				Computed:    true,
			},
		},
		CustomizeDiff: volumeConfigCustomizeDiff,
	}
}

func VolumeConfigSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "Name of the disk.",
				Required:    true,
			},
			"size_in_gb": {
				Type:        schema.TypeInt,
				Description: "Size (GB) of the disk.",
				Required:    true,
			},
			"volume_type": {
				Type:         schema.TypeString,
				Description:  "Type of the disk. Valid values: `NVME`(SSD), `SATA`(HDD).",
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"NVME", "SATA"}, false),
			},
			"pvc_name": {
				Type:        schema.TypeString,
				Description: "PVC name of the disk.",
				Computed:    true,
			},
		},
	}
}

func DNSConfigSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"dns_type": {
				Type: schema.TypeString,
				Description: "DNS type. Valid values: `NONE`(no DNS config), `DEFAULT`(114.114.114.114 for domestic nodes, 8.8.8.8 for overseas nodes), " +
					"`LOCAL`(local dns of node), `CUSTOMIZE`",
				Required: true,
			},
			"dns_address": {
				Type:        schema.TypeList,
				Description: "Custom DNS address.",
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func KeyConfigSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:         schema.TypeString,
				Description:  "Valid values: `bccKeyPair`, `password`",
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"bccKeyPair", "password"}, false),
			},
			"bcc_key_pair_id_list": {
				Type:        schema.TypeList,
				Description: "Key pair ID list.",
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				ConflictsWith: []string{"key_config.0.admin_pass"},
			},
			"admin_pass": {
				Type: schema.TypeString,
				Description: "Length of the password is limited to 8 to 32 characters. Letters, numbers and symbols must exist at the same time, " +
					"and the symbols are limited to `!@#$%^+*()`",
				Optional:      true,
				Sensitive:     true,
				ConflictsWith: []string{"key_config.0.bcc_key_pair_id_list"},
			},
		},
	}
}

func resourceVMInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*connectivity.BaiduClient)
	serviceID := d.Get("service_id").(string)

	raw, err := conn.WithBECClient(func(client *bec.Client) (interface{}, error) {
		return client.CreateVmServiceInstance(serviceID, buildCreationArgs(d))
	})
	log.Printf("[DEBUG] Create BEC VM Instance result: %+v", raw)
	if err != nil {
		return fmt.Errorf("error creating BEC VM Instance: %w", err)
	}
	response := raw.(*api.CreateVmServiceResult)
	if !response.Result || len(response.Details.Instances) == 0 {
		return fmt.Errorf("error creating BEC VM Instance: %+v", raw)
	}

	instance := response.Details.Instances[0]
	d.SetId(instance.VmId)

	if _, err = waitVMInstanceAvailable(conn, d.Id()); err != nil {
		return fmt.Errorf("error waiting BEC VM Instance (%s) becoming available: %w", d.Id(), err)
	}
	return resourceVMInstanceRead(d, meta)
}

func resourceVMInstanceRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*connectivity.BaiduClient)

	raw, err := conn.WithBECClient(func(client *bec.Client) (interface{}, error) {
		return client.GetVirtualMachine(d.Id())
	})
	log.Printf("[DEBUG] Read BEC VM Instance (%s) result: %+v", d.Id(), raw)
	if err != nil {
		return fmt.Errorf("error reading BEC VM Instance (%s): %w", d.Id(), err)
	}

	detail := raw.(*api.VmInstanceDetailsVo)
	if err := d.Set("service_id", detail.ServiceId); err != nil {
		return fmt.Errorf("error setting service_id: %w", err)
	}
	if err := d.Set("vm_name", detail.VmName); err != nil {
		return fmt.Errorf("error setting vm_name: %w", err)
	}
	if err := d.Set("host_name", detail.Hostname); err != nil {
		return fmt.Errorf("error setting host_name: %w", err)
	}
	if err := d.Set("region_id", detail.RegionId); err != nil {
		return fmt.Errorf("error setting region_id: %w", err)
	}
	if err := d.Set("spec", detail.Spec); err != nil {
		return fmt.Errorf("error setting spec: %w", err)
	}
	if err := d.Set("cpu", detail.Cpu); err != nil {
		return fmt.Errorf("error setting cpu: %w", err)
	}
	if err := d.Set("memory", detail.Mem); err != nil {
		return fmt.Errorf("error setting memory: %w", err)
	}
	if err := d.Set("image_type", flattenImageType(detail.OsImage.ImageType)); err != nil {
		return fmt.Errorf("error setting image_type: %w", err)
	}
	if err := d.Set("image_id", detail.OsImage.ImageId); err != nil {
		return fmt.Errorf("error setting image_id: %w", err)
	}
	if err := d.Set("system_volume", flattenSystemVolume(detail.SystemVolume)); err != nil {
		return fmt.Errorf("error setting system_volume: %w", err)
	}
	if err := d.Set("data_volume", flattenDataVolumes(detail.DataVolumeList)); err != nil {
		return fmt.Errorf("error setting data_volume: %w", err)
	}
	if err := d.Set("need_public_ip", detail.NeedPublicIp); err != nil {
		return fmt.Errorf("error setting need_public_ip: %w", err)
	}
	if detail.NeedPublicIp {
		if err := d.Set("need_ipv6_public_ip", detail.NeedIpv6PublicIp); err != nil {
			return fmt.Errorf("error setting need_ipv6_public_ip: %w", err)
		}
		if err := d.Set("bandwidth", flattenBandwidth(detail.Bandwidth)); err != nil {
			return fmt.Errorf("error setting bandwidth: %w", err)
		}
	}

	if err := d.Set("dns_config", flattenDNSConfig(detail.Dns)); err != nil {
		return fmt.Errorf("error setting spec: %w", err)
	}

	adminPass := d.Get("key_config.0.admin_pass").(string)
	if err := d.Set("key_config", flattenKeyConfig(detail.BccKeyPairList, adminPass)); err != nil {
		return fmt.Errorf("error setting key_config: %w", err)
	}

	if err := d.Set("status", detail.Status); err != nil {
		return fmt.Errorf("error setting status: %w", err)
	}
	if err := d.Set("internal_ip", detail.InternalIp); err != nil {
		return fmt.Errorf("error setting internal_ip: %w", err)
	}
	if err := d.Set("public_ip", detail.PublicIp); err != nil {
		return fmt.Errorf("error setting public_ip: %w", err)
	}
	if err := d.Set("ipv6_public_ip", detail.Ipv6PublicIp); err != nil {
		return fmt.Errorf("error setting ipv6_public_ip: %w", err)
	}
	if err := d.Set("create_time", detail.CreateTime); err != nil {
		return fmt.Errorf("error setting create_time: %w", err)
	}

	return nil
}

func resourceVMInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*connectivity.BaiduClient)

	if err := updateVMName(d, conn); err != nil {
		return fmt.Errorf("error updating BEC VM Instance (%s) vmName: %w", d.Id(), err)
	}
	if err := updateHostName(d, conn); err != nil {
		return fmt.Errorf("error updating BEC VM Instance (%s) hostname: %w", d.Id(), err)
	}
	if err := updateResource(d, conn); err != nil {
		return fmt.Errorf("error updating BEC VM Instance (%s) resource: %w", d.Id(), err)
	}
	if err := updatePassword(d, conn); err != nil {
		return fmt.Errorf("error updating BEC VM Instance (%s) password: %w", d.Id(), err)
	}
	return resourceVMInstanceRead(d, meta)
}

func resourceVMInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*connectivity.BaiduClient)

	raw, err := conn.WithBECClient(func(client *bec.Client) (interface{}, error) {
		return client.DeleteVmInstance(d.Id())
	})
	log.Printf("[DEBUG] Delete BEC VM Instance (%s) result: %+v", d.Id(), raw)

	if err != nil {
		return fmt.Errorf("error delete BEC VM Instance (%s): %w", d.Id(), err)
	}

	response := raw.(*api.ActionInfoVo)
	if !response.Result {
		return fmt.Errorf("error delete BEC VM Instance (%s): %+v", d.Id(), raw)
	}
	return nil
}

func buildCreationArgs(d *schema.ResourceData) *api.CreateVmServiceArgs {
	args := &api.CreateVmServiceArgs{
		VmName:            d.Get("vm_name").(string),
		Hostname:          d.Get("host_name").(string),
		DeployInstances:   expandDeployInstances(d.Get("region_id").(string)),
		Spec:              d.Get("spec").(string),
		Cpu:               d.Get("cpu").(int),
		Memory:            d.Get("memory").(int),
		ImageType:         api.ImageType(d.Get("image_type").(string)),
		ImageId:           d.Get("image_id").(string),
		SystemVolume:      expandSystemVolume(d.Get("system_volume").([]interface{})),
		DataVolumeList:    expandDataVolumes(d.Get("data_volume").([]interface{})),
		NeedPublicIp:      d.Get("need_public_ip").(bool),
		DnsConfig:         expandDNSConfig(d.Get("dns_config").([]interface{})),
		KeyConfig:         expandKeyConfig(d.Get("key_config").([]interface{})),
		NetworkConfigList: expandNetworkConfigList(d.Get("network_config").([]interface{})),
	}
	if args.NeedPublicIp {
		args.NeedIpv6PublicIp = d.Get("need_ipv6_public_ip").(bool)
		args.Bandwidth = d.Get("bandwidth").(int)
	}
	return args
}

func updateVMName(d *schema.ResourceData, conn *connectivity.BaiduClient) error {
	if d.HasChange("vm_name") {
		args := &api.UpdateVmInstanceArgs{
			Type:   "vmName",
			VmName: d.Get("vm_name").(string),
		}
		log.Printf("[DEBUG] Update BEC VM Instance (%s) vmName: %+v", d.Id(), args)
		return updateVMInstance(conn, d.Id(), args)
	}
	return nil
}

func updateHostName(d *schema.ResourceData, conn *connectivity.BaiduClient) error {
	if d.HasChange("host_name") {
		args := &api.UpdateVmInstanceArgs{
			Type:     "hostname",
			Hostname: d.Get("host_name").(string),
		}
		log.Printf("[DEBUG] Update BEC VM Instance (%s) hostname: %+v", d.Id(), args)
		return updateVMInstanceAndWait(conn, d.Id(), args)
	}
	return nil
}

func updateResource(d *schema.ResourceData, conn *connectivity.BaiduClient) error {
	if d.HasChanges("cpu", "memory", "data_volume", "need_ipv6_public_ip", "bandwidth", "dns_config") {
		args := &api.UpdateVmInstanceArgs{
			Type: "resource",
		}

		if d.HasChange("cpu") {
			log.Printf("[DEBUG] Update BEC VM cpu")
			args.Cpu = d.Get("cpu").(int)
		}
		if d.HasChange("memory") {
			log.Printf("[DEBUG] Update BEC VM memory")
			args.Memory = d.Get("memory").(int)
		}
		if d.HasChanges("data_volume") {
			log.Printf("[DEBUG] Update BEC VM data_volume")
			args.DataVolumeList = expandDataVolumes(d.Get("data_volume").([]interface{}))
		}
		if d.HasChange("need_ipv6_public_ip") {
			log.Printf("[DEBUG] Update BEC VM need_ipv6_public_ip")
			args.NeedIpv6PublicIp = d.Get("need_ipv6_public_ip").(bool)
		}
		if d.HasChange("bandwidth") {
			log.Printf("[DEBUG] Update BEC VM bandwidth")
			args.Bandwidth = d.Get("bandwidth").(int)
			args.Cpu = d.Get("cpu").(int)
		}
		if d.HasChange("dns_config") {
			log.Printf("[DEBUG] Update BEC VM dns_config")
			args.DnsConfig = expandDNSConfig(d.Get("dns_config").([]interface{}))
		}
		log.Printf("[DEBUG] Update BEC VM Instance (%s) resource: %+v", d.Id(), args)
		return updateVMInstanceAndWait(conn, d.Id(), args)
	}
	return nil
}

func updatePassword(d *schema.ResourceData, conn *connectivity.BaiduClient) error {
	if d.HasChange("key_config") {
		args := &api.UpdateVmInstanceArgs{
			Type:      "password",
			KeyConfig: expandKeyConfig(d.Get("key_config").([]interface{})),
		}
		log.Printf("[DEBUG] Update BEC VM Instance (%s) password: %+v", d.Id(), args)
		return updateVMInstanceAndWait(conn, d.Id(), args)
	}
	return nil
}

func updateVMInstance(conn *connectivity.BaiduClient, vmID string, args *api.UpdateVmInstanceArgs) error {
	raw, err := conn.WithBECClient(func(client *bec.Client) (interface{}, error) {
		return client.UpdateVmInstance(vmID, args)
	})
	if err != nil {
		return err
	}
	response := raw.(*api.UpdateVmDeploymentResult)
	if !response.Result {
		return fmt.Errorf("error updating BEC VM Instance: %+v", raw)
	}
	return nil
}

func updateVMInstanceAndWait(conn *connectivity.BaiduClient, vmID string, args *api.UpdateVmInstanceArgs) error {
	if err := updateVMInstance(conn, vmID, args); err != nil {
		return err
	}
	if _, err := waitVMInstanceAvailable(conn, vmID); err != nil {
		return err
	}
	return nil
}
