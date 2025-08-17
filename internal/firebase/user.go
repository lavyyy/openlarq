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
