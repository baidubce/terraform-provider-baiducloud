package baiducloud

import (
	"fmt"
	"strconv"

	"github.com/baidubce/bce-sdk-go/services/cfc"
	"github.com/baidubce/bce-sdk-go/services/cfc/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

type CFCService struct {
	client *connectivity.BaiduClient
}

func (c *CFCService) CFCUpdateFunctionCode(functionName string, d *schema.ResourceData) error {
	updateArgs := &api.UpdateFunctionCodeArgs{
		FunctionName: functionName,
	}

	fileName, fileNameOk := d.GetOk("code_file_name")
	fileDir, fileDirOk := d.GetOk("code_file_dir")
	if fileNameOk {
		zipFile, err := loadFileContent(fileName.(string))
		if err != nil {
			return err
		}

		updateArgs.ZipFile = zipFile
	} else if fileDirOk {
		zipFile, err := zipFileDir(fileDir.(string))
		if err != nil {
			return err
		}

		updateArgs.ZipFile = zipFile
	} else {
		bucket, bucketOk := d.GetOk("code_bos_bucket")
		object, objectOk := d.GetOk("code_bos_object")

		if !bucketOk || !objectOk {
			return fmt.Errorf("code_bos_bucket and code_bos_object must all be set while using bos code source")
		}

		updateArgs.BosBucket = bucket.(string)
		updateArgs.BosObject = object.(string)
	}

	action := "Update Function " + functionName + " code"
	raw, err := c.client.WithCFCClient(func(client *cfc.Client) (i interface{}, e error) {
		return client.UpdateFunctionCode(updateArgs)
	})
	addDebug(action, raw)

	if err != nil {
		return WrapError(err)
	}

	return nil
}

func (c *CFCService) CFCUpdateFunctionConfigure(updateArgs *api.UpdateFunctionConfigurationArgs) error {
	action := "Update Function " + updateArgs.FunctionName + " configure"
	raw, err := c.client.WithCFCClient(func(client *cfc.Client) (i interface{}, e error) {
		return client.UpdateFunctionConfiguration(updateArgs)
	})
	addDebug(action, raw)

	if err != nil {
		return WrapError(err)
	}

	return nil
}

func (c *CFCService) CFCSetReservedConcurrent(functionName string, newConcurrent int) error {
	setArgs := &api.ReservedConcurrentExecutionsArgs{
		FunctionName:                 functionName,
		ReservedConcurrentExecutions: newConcurrent,
	}
	action := "Update Function " + functionName + " reserved concurrent executions to " + strconv.Itoa(newConcurrent)

	_, err := c.client.WithCFCClient(func(client *cfc.Client) (i interface{}, e error) {
		return nil, client.SetReservedConcurrentExecutions(setArgs)
	})
	addDebug(action, setArgs)

	if err != nil {
		return WrapError(err)
	}

	return nil
}

func (c *CFCService) CFCDeleteReservedConcurrent(functionName string) error {
	action := "Delete Function " + functionName + " reserved concurrent executions"

	_, err := c.client.WithCFCClient(func(client *cfc.Client) (i interface{}, e error) {
		return nil, client.DeleteReservedConcurrentExecutions(
			&api.DeleteReservedConcurrentExecutionsArgs{
				FunctionName: functionName,
			})
	})
	addDebug(action, functionName)

	if err != nil {
		return WrapError(err)
	}

	return nil
}

func (c *CFCService) CFCGetVersionsByFunction(functionName, functionVersion string) (*api.Function, error) {
	action := "Get version " + functionVersion + " by function " + functionName

	args := &api.ListVersionsByFunctionArgs{
		FunctionName: functionName,
	}
	for {
		raw, err := c.client.WithCFCClient(func(client *cfc.Client) (i interface{}, e error) {
			return client.ListVersionsByFunction(args)
		})

		if err != nil {
			return nil, WrapError(err)
		}

		addDebug(action, raw)
		response := raw.(*api.ListVersionsByFunctionResult)
		for _, f := range response.Versions {
			if f.Version == functionVersion {
				return f, nil
			}
		}

		if len(response.NextMarker) > 0 {
			args.Marker, _ = strconv.Atoi(response.NextMarker)
		} else {
			return nil, WrapError(Error("Function %s not exit version %s", functionName, functionVersion))
		}
	}
}

func (c *CFCService) CFCListAllVersionsByFunction(functionName string) ([]*api.Function, error) {
	action := "List all versions by function " + functionName

	result := make([]*api.Function, 0)

	args := &api.ListVersionsByFunctionArgs{
		FunctionName: functionName,
	}
	for {
		raw, err := c.client.WithCFCClient(func(client *cfc.Client) (i interface{}, e error) {
			return client.ListVersionsByFunction(args)
		})

		if err != nil {
			return nil, WrapError(err)
		}

		addDebug(action, raw)
		response := raw.(*api.ListVersionsByFunctionResult)
		result = append(result, response.Versions...)

		if len(response.NextMarker) > 0 {
			args.Marker, _ = strconv.Atoi(response.NextMarker)
		} else {
			return result, nil
		}
	}
}

func (c *CFCService) CFCDeleteFunctionVersion(functionName, functionVersion string) error {
	action := "Delete Function " + functionName + " version " + functionVersion

	_, err := c.client.WithCFCClient(func(client *cfc.Client) (i interface{}, e error) {
		return nil, client.DeleteFunction(&api.DeleteFunctionArgs{
			FunctionName: functionName,
			Qualifier:    functionVersion,
		})
	})
	addDebug(action, functionName)

	if err != nil {
		return WrapError(err)
	}

	return nil
}

func (c *CFCService) CFCGetTriggerByFunction(functionBrn, relationId string) (*api.RelationInfo, error) {
	action := "Get Function " + functionBrn + " trigger " + relationId

	raw, err := c.client.WithCFCClient(func(client *cfc.Client) (i interface{}, e error) {
		return client.ListTriggers(&api.ListTriggersArgs{FunctionBrn: functionBrn})
	})

	addDebug(action, raw)

	if err != nil {
		return nil, WrapError(err)
	}

	response, _ := raw.(*api.ListTriggersResult)
	for _, relation := range response.Relation {
		if relation.RelationId == relationId {
			return relation, nil
		}
	}

	return nil, Error(ResourceNotFound)
}

func (c *CFCService) ListAllFunctions() ([]*api.Function, error) {
	result := make([]*api.Function, 0)

	args := &api.ListFunctionsArgs{}
	for {
		raw, err := c.client.WithCFCClient(func(client *cfc.Client) (i interface{}, e error) {
			return client.ListFunctions(args)
		})

		if err != nil {
			return nil, err
		}

		response := raw.(*api.ListFunctionsResult)
		result = append(result, response.Functions...)

		if len(response.NextMarker) > 0 {
			args.Marker, _ = strconv.Atoi(response.NextMarker)
		} else {
			return result, nil
		}
	}
}

func (c *CFCService) faltternCFCFunctionVpcConfigToMap(config *api.VpcConfig) map[string]interface{} {
	if config == nil || len(config.SubnetIds) == 0 || len(config.SecurityGroupIds) == 0 {
		return nil
	}

	result := map[string]interface{}{}
	result["subnet_ids"] = schema.NewSet(schema.HashString, flattenStringListToInterface(config.SubnetIds))
	result["security_group_ids"] = schema.NewSet(schema.HashString, flattenStringListToInterface(config.SecurityGroupIds))

	if config.VpcId != "" {
		result["vpc_id"] = config.VpcId
	}

	return result
}
