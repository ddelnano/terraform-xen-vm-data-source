## Terraform Xen Vm Data source

While working with the [Xen Orchestra terraform provider](), it's often helpful to be able to retrieve the IP addresses of an existing VM(s). This is not possible at the moment until [#112](https://github.com/terra-farm/terraform-provider-xenorchestra/issues/112) is implemented, however, Terraform has an `external` data source that allows using an external program to augment existing terraform attributes and data. This repo is a proof of concept on how to access guest tools IP address information from Xen VMs until [#112](https://github.com/terra-farm/terraform-provider-xenorchestra/issues/112) is implemented inside the provider.

### Argument Reference
The terraform `external` data source requires providing the `program` and for this specific script `query` is not an optional field.

- `query` - A map whose keys are a resource identifier (`xenorchestra_vm.vm`) and whose values are the Xen UUID of the given VM. Please the examples below for more details.

### Attributes Reference
The following attributes are exported:
- `result` - A map of resource name to object that contains `ip_address`, `ipv4_address` and `ipv6_address`. An example can be seen below in the Testing the script section.

### Using this script
1.Compile this program `go build -o main`
2. Move the `main` binary into a directory that is accessible from your Terraform code
3. Create a `external` data source using this program and providing the correct inputs

```
resource "xenorchestra_vm" "vm" {
 ....
 ....
}

data "external" "xen_vm_ips" {
  program = ["${path.module}/path/to/main"]
  query = {
    # The query argument expects a VM resource name as a key and
    # the VM's UUID as the value
    "xenorchestra_vm.vm": xenorchestra_vm.vm.id
  }
}
```
4. Export the `XAPI_HOST`, `XAPI_USERNAME` and `XAPI_PASSWORD` to your current shell
5. Run `terraform plan` and `terraform apply`


### Testing the script

The terraform external program protocol is documented in more detail [here](https://registry.terraform.io/providers/hashicorp/external/latest/docs/data-sources/data_source#external-program-protocol).

In order to test the script you will need to pass the "query" arguments to the program's stdin. Below is an example of a successful test of the script

```
$ echo '{"query": {"xenorchestra_vm.vm": "fd882c8e-e344-22e2-bd9c-ecb3cddae519"}}' | go run main.go
{"xenorchestra_vm.vm":{"ip_address":"172.16.210.6","ipv4_address":"172.16.210.6","ipv6_address":"2a01:240:ab08:4:1c07:34ff:fee2:a5d1"}}

```

