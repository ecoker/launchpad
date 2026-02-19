#!/usr/bin/env node

import { PROFILES } from "./profiles.js";
import { scaffoldProject } from "./scaffold.js";

type ParsedArgs = {
  command: "list" | "init" | "help";
  targetDir?: string;
  profile: string;
  force: boolean;
};

function printHelp() {
  console.log(`
agent-kit: opinionated AI instruction scaffolder

Usage:
  agent-kit list
  agent-kit init <targetDir> [--profile <id>] [--force]

Options:
  --profile <id>   Profile template to apply (default: typescript-react)
  --force          Overwrite files in a non-empty target directory

Examples:
  agent-kit list
  agent-kit init ./my-app --profile typescript-react
  agent-kit init ./analytics --profile python-data --force
`);
}

function parseArgs(argv: string[]): ParsedArgs {
  const [command = "help", ...rest] = argv;

  if (command === "list") {
    return { command: "list", profile: "typescript-react", force: false };
  }

  if (command !== "init") {
    return { command: "help", profile: "typescript-react", force: false };
  }

  let targetDir: string | undefined;
  let profile = "typescript-react";
  let force = false;

  for (let index = 0; index < rest.length; index += 1) {
    const token = rest[index];

    if (!token.startsWith("--") && !targetDir) {
      targetDir = token;
      continue;
    }

    if (token === "--force") {
      force = true;
      continue;
    }

    if (token === "--profile") {
      const next = rest[index + 1];
      if (!next) {
        throw new Error("Missing value for --profile");
      }
      profile = next;
      index += 1;
      continue;
    }

    throw new Error(`Unknown argument: ${token}`);
  }

  if (!targetDir) {
    throw new Error("Missing required <targetDir> for init command");
  }

  return { command: "init", targetDir, profile, force };
}

async function main() {
  const args = parseArgs(process.argv.slice(2));

  if (args.command === "help") {
    printHelp();
    return;
  }

  if (args.command === "list") {
    console.log("Available profiles:\n");
    for (const profile of PROFILES) {
      console.log(`- ${profile.id}`);
      console.log(`  ${profile.title}`);
      console.log(`  ${profile.summary}\n`);
    }
    return;
  }

  const result = await scaffoldProject({
    targetDir: args.targetDir!,
    profileId: args.profile,
    force: args.force
  });

  console.log(`Scaffolded '${result.profile.id}' into ${result.outputPath}`);
  console.log("Next steps:");
  console.log(`- cd ${args.targetDir}`);
  console.log("- Review .github/copilot-instructions.md and .github/instructions/*.instructions.md");
}

main().catch((error: unknown) => {
  const message = error instanceof Error ? error.message : "Unknown error";
  console.error(`Error: ${message}`);
  process.exitCode = 1;
});