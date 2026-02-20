package cli

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadKeyFromDotEnv(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    string
	}{
		{
			name:    "simple value",
			content: "OPENAI_API_KEY=sk-test123\n",
			want:    "sk-test123",
		},
		{
			name:    "double-quoted value",
			content: "OPENAI_API_KEY=\"sk-quoted123\"\n",
			want:    "sk-quoted123",
		},
		{
			name:    "single-quoted value",
			content: "OPENAI_API_KEY='sk-single123'\n",
			want:    "sk-single123",
		},
		{
			name:    "export prefix",
			content: "export OPENAI_API_KEY=sk-exported123\n",
			want:    "sk-exported123",
		},
		{
			name:    "export with quotes",
			content: "export OPENAI_API_KEY=\"sk-both123\"\n",
			want:    "sk-both123",
		},
		{
			name:    "inline comment",
			content: "OPENAI_API_KEY=sk-commented123 # my key\n",
			want:    "sk-commented123",
		},
		{
			name:    "KEY alias",
			content: "KEY=sk-alias123\n",
			want:    "sk-alias123",
		},
		{
			name:    "skips comments and blanks",
			content: "# comment\n\nOPENAI_API_KEY=sk-afterblank\n",
			want:    "sk-afterblank",
		},
		{
			name:    "no matching key",
			content: "OTHER_KEY=value\n",
			want:    "",
		},
		{
			name:    "empty file",
			content: "",
			want:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			envPath := filepath.Join(dir, ".env")
			if err := os.WriteFile(envPath, []byte(tt.content), 0o644); err != nil {
				t.Fatalf("writing .env: %v", err)
			}

			orig, _ := os.Getwd()
			t.Cleanup(func() { os.Chdir(orig) })
			os.Chdir(dir)

			got := loadKeyFromDotEnv()
			if got != tt.want {
				t.Errorf("loadKeyFromDotEnv() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestLoadKeyFromDotEnv_NoFile(t *testing.T) {
	dir := t.TempDir()
	orig, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(orig) })
	os.Chdir(dir)

	got := loadKeyFromDotEnv()
	if got != "" {
		t.Errorf("expected empty string when no .env exists, got %q", got)
	}
}
