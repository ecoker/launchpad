package scaffold

// Profile represents a language/framework profile that can be scaffolded.
type Profile struct {
	ID      string
	Title   string
	Summary string
	Dir     string // directory name inside templates/profiles/
}

// Addon represents an optional add-on instruction set.
type Addon struct {
	ID      string
	Title   string
	Summary string
	Dir     string // directory name inside templates/addons/
}

// Profiles lists every available profile.
var Profiles = []Profile{
	{
		ID:      "typescript-react",
		Title:   "TypeScript + React",
		Summary: "React Router v7, Tailwind, shadcn/ui, Motion – fullstack TypeScript",
		Dir:     "typescript-react",
	},
	{
		ID:      "python-data",
		Title:   "Python Data / AI",
		Summary: "Pydantic, functional style, data pipelines & AI agents",
		Dir:     "python-data",
	},
	{
		ID:      "elixir-phoenix",
		Title:   "Elixir + Phoenix",
		Summary: "Phoenix LiveView, Ecto, OTP – functional by nature",
		Dir:     "elixir-phoenix",
	},
	{
		ID:      "dotnet-api",
		Title:   ".NET API",
		Summary: "C# minimal APIs, Entity Framework, clean architecture",
		Dir:     "dotnet-api",
	},
	{
		ID:      "laravel",
		Title:   "Laravel",
		Summary: "Laravel + Inertia, Eloquent, queues – the PHP way",
		Dir:     "laravel",
	},
	{
		ID:      "go-service",
		Title:   "Go Service",
		Summary: "Idiomatic Go microservices, stdlib-first, minimal deps",
		Dir:     "go-service",
	},
}

// Addons lists every available add-on.
var Addons = []Addon{
	{
		ID:      "data-intensive",
		Title:   "Data-Intensive",
		Summary: "Postgres, NATS, Parquet, event-driven architecture",
		Dir:     "data-intensive",
	},
	{
		ID:      "frontend-craft",
		Title:   "Frontend Craft",
		Summary: "Advanced CSS, animation, accessibility, responsive design",
		Dir:     "frontend-craft",
	},
}

// FindProfile returns the profile with the given ID, or nil if not found.
func FindProfile(id string) *Profile {
	for i := range Profiles {
		if Profiles[i].ID == id {
			return &Profiles[i]
		}
	}
	return nil
}

// FindAddon returns the addon with the given ID, or nil if not found.
func FindAddon(id string) *Addon {
	for i := range Addons {
		if Addons[i].ID == id {
			return &Addons[i]
		}
	}
	return nil
}

// ProfileIDs returns a slice of all profile IDs.
func ProfileIDs() []string {
	ids := make([]string, len(Profiles))
	for i, p := range Profiles {
		ids[i] = p.ID
	}
	return ids
}

// AddonIDs returns a slice of all addon IDs.
func AddonIDs() []string {
	ids := make([]string, len(Addons))
	for i, a := range Addons {
		ids[i] = a.ID
	}
	return ids
}
