package firebase

import (
	"fmt"
)

func (fc *FirebaseClient) GetUserLiquidIntake(params QueryParams) (QueryBody, error) {
	liquidIntakePath := fmt.Sprintf("/liquidIntake/%s", fc.UserId())

	return fc.Query(liquidIntakePath, params)
}

func (fc *FirebaseClient) GetUserHydrationGoals(params QueryParams) (QueryBody, error) {
	goalsPath := fmt.Sprintf("/hydrationGoal/%s", fc.UserId())

	return fc.Query(goalsPath, params)
}

func (fc *FirebaseClient) GetDeviceInfo(params QueryParams, deviceId string) (QueryBody, error) {
	deviceInfoPath := fmt.Sprintf("/capUserSettings/%s/%s", fc.UserId(), deviceId)

	return fc.Query(deviceInfoPath, params)
}
