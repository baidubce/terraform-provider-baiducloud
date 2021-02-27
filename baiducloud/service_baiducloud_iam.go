package baiducloud

import (
	"github.com/baidubce/bce-sdk-go/services/iam"
	"github.com/baidubce/bce-sdk-go/services/iam/api"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

type IamService struct {
	client *connectivity.BaiduClient
}

func (s IamService) ClearUserAttachedPolicy(userName string) error {
	raw, err := s.client.WithIamClient(func(iamClient *iam.Client) (i interface{}, e error) {
		return iamClient.ListUserAttachedPolicies(userName)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_iam", "list attached policy for user "+userName,
			BCESDKGoERROR)
	}
	listPolicyResult := raw.(*api.ListPolicyResult)
	for _, policy := range listPolicyResult.Policies {
		if err := s.DetachPolicyFromUser(userName, policy.Name, policy.Type); err != nil {
			return err
		}
	}
	return nil
}

func (s IamService) DetachPolicyFromUser(userName, policyName, policyType string) error {
	_, err := s.client.WithIamClient(func(iamClient *iam.Client) (i interface{}, e error) {
		return nil, iamClient.DetachPolicyFromUser(&api.DetachPolicyFromUserArgs{
			UserName:   userName,
			PolicyName: policyName,
			PolicyType: policyType,
		})
	})
	if err != nil && !NotFoundError(err) {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_iam", "detach attached policy for user "+userName,
			BCESDKGoERROR)
	}
	return nil
}

func (s IamService) ClearUserGroupMembership(userName string) error {
	raw, err := s.client.WithIamClient(func(iamClient *iam.Client) (i interface{}, e error) {
		return iamClient.ListGroupsForUser(userName)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_iam", "list groups for user "+userName,
			BCESDKGoERROR)
	}
	listGroupResult := raw.(*api.ListGroupsForUserResult)
	for _, group := range listGroupResult.Groups {
		if err := s.DeleteUserFromGroup(userName, group.Name); err != nil {
			return err
		}
	}
	return nil
}

func (s IamService) AddUserToGroup(userName, groupName string) error {
	_, err := s.client.WithIamClient(func(iamClient *iam.Client) (i interface{}, e error) {
		return nil, iamClient.AddUserToGroup(userName, groupName)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_iam",
			"add user "+userName+" to group "+groupName, BCESDKGoERROR)
	}
	return nil
}

func (s IamService) DeleteUserFromGroup(userName, groupName string) error {
	_, err := s.client.WithIamClient(func(iamClient *iam.Client) (i interface{}, e error) {
		return nil, iamClient.DeleteUserFromGroup(userName, groupName)
	})
	if err != nil && !NotFoundError(err) {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_iam",
			"delete user "+userName+" from group "+groupName, BCESDKGoERROR)
	}
	return nil
}

func (s IamService) ClearGroupAttachedPolicy(groupName string) error {
	raw, err := s.client.WithIamClient(func(iamClient *iam.Client) (i interface{}, e error) {
		return iamClient.ListGroupAttachedPolicies(groupName)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_iam", "list attached policy for group "+groupName,
			BCESDKGoERROR)
	}
	listPolicyResult := raw.(*api.ListPolicyResult)
	for _, policy := range listPolicyResult.Policies {
		if err := s.DetachPolicyFromGroup(groupName, policy.Name, policy.Type); err != nil {
			return err
		}
	}
	return nil
}

func (s IamService) DetachPolicyFromGroup(groupName, policyName, policyType string) error {
	_, err := s.client.WithIamClient(func(iamClient *iam.Client) (i interface{}, e error) {
		return nil, iamClient.DetachPolicyFromGroup(&api.DetachPolicyFromGroupArgs{
			GroupName:  groupName,
			PolicyName: policyName,
			PolicyType: policyType,
		})
	})
	if err != nil && !NotFoundError(err) {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_iam", "detach attached policy for group "+groupName,
			BCESDKGoERROR)
	}
	return nil
}

func (s IamService) ClearUserFromGroup(groupName string) error {
	listUsersResult, err := s.ListUserInGroup(groupName)
	if err != nil {
		return err
	}
	for _, user := range listUsersResult.Users {
		if err := s.DeleteUserFromGroup(user.Name, groupName); err != nil {
			return err
		}
	}
	return nil
}
func (s IamService) ListUserInGroup(groupName string) (*api.ListUsersInGroupResult, error) {
	raw, err := s.client.WithIamClient(func(iamClient *iam.Client) (i interface{}, e error) {
		return iamClient.ListUsersInGroup(groupName)
	})
	if err != nil {
		return nil, WrapErrorf(err, DefaultErrorMsg, "baiducloud_iam", "list users in group "+groupName,
			BCESDKGoERROR)
	}
	return raw.(*api.ListUsersInGroupResult), nil
}
