export type ProfileId =
  | "typescript-react"
  | "python-data"
  | "elixir-phoenix"
  | "dotnet-api"
  | "laravel"
  | "go-service";

export type Profile = {
  id: ProfileId;
  title: string;
  summary: string;
};

export const PROFILES: Profile[] = [
  {
    id: "typescript-react",
    title: "TypeScript + React Router + Tailwind",
    summary: "Strong typing, composition-first React, accessible UI, and motion with intent."
  },
  {
    id: "python-data",
    title: "Python Data Systems",
    summary: "Typed Python with Pydantic, explicit schemas, robust pipelines, and reliable data workflows."
  },
  {
    id: "elixir-phoenix",
    title: "Elixir + Phoenix",
    summary: "Functional boundaries, supervision, immutable state, and fault-tolerant services."
  },
  {
    id: "dotnet-api",
    title: "C# .NET API",
    summary: "Clean architecture in .NET with explicit contracts and testable application layers."
  },
  {
    id: "laravel",
    title: "PHP Laravel",
    summary: "Convention with discipline, clear domain services, and pragmatic, maintainable app structure."
  },
  {
    id: "go-service",
    title: "Go Service",
    summary: "Small composable packages, explicit interfaces, and clear operational behavior."
  }
];

export function findProfile(profileId: string): Profile | undefined {
  return PROFILES.find((profile) => profile.id === profileId);
}