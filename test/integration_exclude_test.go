package test_test

import (
	"github.com/gruntwork-io/terragrunt/test/helpers"
	"github.com/gruntwork-io/terragrunt/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"testing"
)

const (
	testExcludeByDefault = "fixtures/exclude/exclude-default"
	testExcludeDisabled  = "fixtures/exclude/exclude-disabled"
	testExcludeByAction  = "fixtures/exclude/exclude-by-action"
	testExcludeByFlags   = "fixtures/exclude/exclude-by-feature-flags"
)

func TestExcludeByDefault(t *testing.T) {
	t.Parallel()

	cleanupTerraformFolder(t, testExcludeByDefault)
	tmpEnvPath := helpers.CopyEnvironment(t, testExcludeByDefault)
	rootPath := util.JoinPath(tmpEnvPath, testExcludeByDefault)

	_, stderr, err := helpers.RunTerragruntCommandWithOutput(t, "terragrunt run-all apply -auto-approve --terragrunt-non-interactive --terragrunt-working-dir "+rootPath)

	require.NoError(t, err)

	assert.Contains(t, stderr, "app1")
	assert.NotContains(t, stderr, "app2")
}

func TestExcludeDisabled(t *testing.T) {
	t.Parallel()

	cleanupTerraformFolder(t, testExcludeDisabled)
	tmpEnvPath := helpers.CopyEnvironment(t, testExcludeDisabled)
	rootPath := util.JoinPath(tmpEnvPath, testExcludeDisabled)

	_, stderr, err := helpers.RunTerragruntCommandWithOutput(t, "terragrunt run-all apply -auto-approve --terragrunt-non-interactive --terragrunt-working-dir "+rootPath)

	require.NoError(t, err)

	assert.Contains(t, stderr, "app1")
	assert.Contains(t, stderr, "app2")
}

func TestExcludeApply(t *testing.T) {
	t.Parallel()

	cleanupTerraformFolder(t, testExcludeByAction)
	tmpEnvPath := helpers.CopyEnvironment(t, testExcludeByAction)
	rootPath := util.JoinPath(tmpEnvPath, testExcludeByAction)

	_, stderr, err := helpers.RunTerragruntCommandWithOutput(t, "terragrunt run-all plan --terragrunt-non-interactive --terragrunt-working-dir "+rootPath)

	require.NoError(t, err)

	assert.Contains(t, stderr, "exclude-apply")
	assert.NotContains(t, stderr, "exclude-plan")

	_, stderr, err = helpers.RunTerragruntCommandWithOutput(t, "terragrunt run-all apply -auto-approve --terragrunt-non-interactive --terragrunt-working-dir "+rootPath)

	require.NoError(t, err)

	// should be applied only exclude-plan
	assert.Contains(t, stderr, "exclude-plan")
	assert.NotContains(t, stderr, "exclude-apply")
}

func TestExcludeByFeatureFlagDefault(t *testing.T) {
	t.Parallel()

	cleanupTerraformFolder(t, testExcludeByFlags)
	tmpEnvPath := helpers.CopyEnvironment(t, testExcludeByFlags)
	rootPath := util.JoinPath(tmpEnvPath, testExcludeByFlags)

	_, stderr, err := helpers.RunTerragruntCommandWithOutput(t, "terragrunt run-all apply --terragrunt-non-interactive --terragrunt-working-dir "+rootPath)

	require.NoError(t, err)

	assert.Contains(t, stderr, "app2")
	assert.NotContains(t, stderr, "app1")
}

func TestExcludeByFeatureFlag(t *testing.T) {
	t.Parallel()

	cleanupTerraformFolder(t, testExcludeByFlags)
	tmpEnvPath := helpers.CopyEnvironment(t, testExcludeByFlags)
	rootPath := util.JoinPath(tmpEnvPath, testExcludeByFlags)

	_, stderr, err := helpers.RunTerragruntCommandWithOutput(t, "terragrunt run-all apply --feature exclude2=false --terragrunt-non-interactive --terragrunt-working-dir "+rootPath)

	require.NoError(t, err)

	assert.Contains(t, stderr, "app1")
	assert.Contains(t, stderr, "app2")
}

func TestExcludeAllByFeatureFlag(t *testing.T) {
	t.Parallel()

	cleanupTerraformFolder(t, testExcludeByFlags)
	tmpEnvPath := helpers.CopyEnvironment(t, testExcludeByFlags)
	rootPath := util.JoinPath(tmpEnvPath, testExcludeByFlags)

	_, stderr, err := helpers.RunTerragruntCommandWithOutput(t, "terragrunt run-all apply --feature exclude1=true --feature exclude2=true --terragrunt-non-interactive --terragrunt-working-dir "+rootPath)

	require.NoError(t, err)

	assert.NotContains(t, stderr, "app1")
	assert.NotContains(t, stderr, "app2")
}
