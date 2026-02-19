package ai

import "strings"

// ValidateSelectionCompatibility enforces hard selection constraints.
func ValidateSelectionCompatibility(selection Selection) []string {
	issues := make([]string, 0)

	if selection.ProfileID == "" {
		issues = append(issues, "profile_id is required")
	} else {
		validProfile := map[string]bool{
			"typescript-react": true,
			"python-data":      true,
			"elixir-phoenix":   true,
			"dotnet-api":       true,
			"laravel":          true,
			"go-service":       true,
		}
		if !validProfile[selection.ProfileID] {
			issues = append(issues, "profile_id is not supported by this Launchpad build")
		}
	}

	allowedAddonsByProfile := map[string]map[string]bool{
		"typescript-react": {
			"frontend-craft": true,
			"data-intensive": true,
		},
		"elixir-phoenix": {
			"frontend-craft": true,
			"data-intensive": true,
		},
		"laravel": {
			"frontend-craft": true,
			"data-intensive": true,
		},
		"python-data": {
			"data-intensive": true,
		},
		"dotnet-api": {
			"data-intensive": true,
		},
		"go-service": {
			"data-intensive": true,
		},
	}

	seenAddons := map[string]bool{}
	for _, addonID := range selection.AddonIDs {
		if addonID == "" {
			continue
		}
		if seenAddons[addonID] {
			issues = append(issues, "duplicate addon_id: "+addonID)
			continue
		}
		seenAddons[addonID] = true

		allowed, ok := allowedAddonsByProfile[selection.ProfileID]
		if !ok || !allowed[addonID] {
			issues = append(issues, "addon_id not compatible with selected profile: "+addonID)
		}
	}

	seenAssets := map[string]bool{}
	var paletteCount, fontCount, lintCount, testingCount, frameworkOpinionCount int
	for _, assetID := range selection.AssetIDs {
		if assetID == "" {
			continue
		}
		if seenAssets[assetID] {
			issues = append(issues, "duplicate asset_id: "+assetID)
			continue
		}
		seenAssets[assetID] = true

		switch {
		case strings.HasPrefix(assetID, "asset.palette."):
			paletteCount++
		case strings.HasPrefix(assetID, "asset.fonts."):
			fontCount++
		case strings.HasPrefix(assetID, "asset.lint"):
			lintCount++
		case strings.HasPrefix(assetID, "asset.testing."):
			testingCount++
		case strings.HasPrefix(assetID, "asset.framework."):
			frameworkOpinionCount++
		}
	}

	if paletteCount > 1 {
		issues = append(issues, "only one palette asset may be selected")
	}
	if fontCount > 1 {
		issues = append(issues, "only one font asset may be selected")
	}
	if lintCount > 1 {
		issues = append(issues, "only one linting asset may be selected")
	}
	if testingCount > 1 {
		issues = append(issues, "only one testing asset may be selected")
	}
	if frameworkOpinionCount > 1 {
		issues = append(issues, "only one framework-opinion asset may be selected")
	}

	return issues
}
