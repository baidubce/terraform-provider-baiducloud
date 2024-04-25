package mongodb

import (
	"log"

	"github.com/baidubce/bce-sdk-go/services/mongodb"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func findInstance(conn *connectivity.BaiduClient, instanceID string) (*mongodb.InstanceDetail, error) {
	raw, err := conn.WithMongoDBClient(func(client *mongodb.Client) (interface{}, error) {
		return client.GetInstanceDetail(instanceID)
	})
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] Query MongoDB Instance (%s) detail: %+v", instanceID, raw)
	return raw.(*mongodb.InstanceDetail), nil
}

func findAllInstance(conn *connectivity.BaiduClient, args mongodb.ListMongodbArgs) ([]mongodb.InstanceModel, error) {
	result := make([]mongodb.InstanceModel, 0)
	nextMarker := ""
	for {
		realArgs := &mongodb.ListMongodbArgs{
			Marker:         nextMarker,
			MaxKeys:        1000,
			EngineVersion:  args.EngineVersion,
			StorageEngine:  args.StorageEngine,
			DbInstanceType: args.DbInstanceType,
		}
		raw, err := conn.WithMongoDBClient(func(client *mongodb.Client) (i interface{}, e error) {
			return client.ListMongodb(realArgs)
		})
		if err != nil {
			return nil, err
		}

		response := raw.(*mongodb.ListMongodbResult)
		result = append(result, response.DbInstances...)

		if response.IsTruncated {
			args.Marker = response.NextMarker
		} else {
			return result, nil
		}
	}
}
