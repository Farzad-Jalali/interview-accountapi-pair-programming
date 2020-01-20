package security

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_should_allow_when_user_has_permission(t *testing.T) {
	organisationId := uuid.New()

	ctx := buildContextWithPermissions(organisationId, "Account", CREATE)

	checkFn := RestrictWithPermissions("Account", CREATE)

	err := checkFn(ctx, &organisationId)

	assert.Nil(t, err)
}

func Test_should_not_allow_when_user_has_permission_but_wrong_organisation(t *testing.T) {

	organisationId := uuid.New()

	ctx := buildContextWithPermissions(organisationId, "Account", CREATE)

	checkFn := RestrictWithPermissions("Account", CREATE)

	randomId := uuid.New()
	err := checkFn(ctx, &randomId)

	assert.NotNil(t, err)
}

func Test_should_not_allow_when_user_does_not_have_permission_for_record_type(t *testing.T) {
	organisationId := uuid.New()

	ctx := buildContextWithPermissions(organisationId, "Payment", CREATE)

	checkFn := RestrictWithPermissions("Account", CREATE)

	randomId := uuid.New()
	err := checkFn(ctx, &randomId)

	assert.NotNil(t, err)
}

func Test_should_allow_when_application_context(t *testing.T) {
	organisationId := uuid.New()

	ctx := ApplicationContext(context.Background())

	checkFn := RestrictWithPermissions("Account", CREATE)

	err := checkFn(ctx, &organisationId)

	assert.Nil(t, err)
}

func Test_check_permission_without_organisation_should_allow_when_user_has_permission(t *testing.T) {
	organisationId := uuid.New()

	ctx := buildContextWithPermissions(organisationId, "sort_codes", READ)

	err := CheckPermissionForResourceWithoutOrganisation(ctx, "sort_codes", READ)

	assert.Nil(t, err)
}

func Test_check_permission_without_organisation_should_not_allow_when_user_does_not_have_permission_for_record_type(t *testing.T) {
	organisationId := uuid.New()

	ctx := buildContextWithPermissions(organisationId, "Payment", CREATE)

	err := CheckPermissionForResourceWithoutOrganisation(ctx, "sort_codes", READ)

	assert.NotNil(t, err)
}

func Test_check_permission_without_organisation_should_allow_when_application_context(t *testing.T) {
	ctx := ApplicationContext(context.Background())

	err := CheckPermissionForResourceWithoutOrganisation(ctx, "sort_codes", READ)

	assert.Nil(t, err)
}

func Test_get_organisations_is_unlimited_when_application_context(t *testing.T) {
	ctx := ApplicationContext(context.Background())

	organisations, _ := GetOrganisationsWithPermission(ctx, "Account", CREATE)

	assert.True(t, organisations.IsUnlimited())
}

func Test_empty_filtered_organisations_is_unlimited_when_application_context(t *testing.T) {
	ctx := ApplicationContext(context.Background())

	organisations, _ := GetOrganisationsWithPermission(ctx, "Account", CREATE)

	assert.True(t, organisations.IntersectFilter(nil).IsUnlimited())
}

func Test_filtered_organisations_is_filtered_when_application_context(t *testing.T) {
	ctx := ApplicationContext(context.Background())

	organisations, _ := GetOrganisationsWithPermission(ctx, "Account", CREATE)

	assert.False(t, organisations.IntersectFilter([]uuid.UUID{uuid.New()}).IsUnlimited())
}

func buildContextWithPermissions(organisationId uuid.UUID, recordType string, allowedActions ...AuthoriseAction) *context.Context {
	acls := AccessControlList{}

	for _, allowedAction := range allowedActions {
		acls = append(acls, AccessControlListEntry{OrganisationId: organisationId, Action: allowedAction, RecordType: recordType})
	}

	ctx := context.WithValue(context.Background(), "acls", acls)
	return &ctx
}
