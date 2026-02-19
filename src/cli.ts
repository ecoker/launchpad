#!/usr/bin/env node

import path from "node:path";
import { PROFILES, type ProfileId } from "./profiles.js";
import { scaffoldProject, type InitOptions } from "./scaffold.js";
import { createPrompt, fmt, c, printBanner, printDone, printFileTree } from "./ui.js";

// ─── Arg parsing ─────────────────────────────────────────────────────────────

type ParsedArgs = {
  command: "list" | "init" | "help";
  targetDir?: string;
  profile?: string;
  force: boolean;
  interactive: boolean;
  addons: string[];
};

function parseArgs(argv: string[]): ParsedArgs {
  const [command = "help", ...rest] = argv;

  if (command === "list") {
    return { command: "list", profile: undefined, force: false, interactive: false, addons: [] };
  }

  if (command !== "init") {
    return { command: "help", profile: undefined, force: false, interactive: false, addons: [] };
  }

  let targetDir: string | undefined;
  let profile: string | undefined;
  let force = false;
  let interactive = true;
  const addons: string[] = [];

  for (let i = 0; i < rest.length; i += 1) {
    const token = rest[i];

    if (token === "--force" || token === "-f") { force = true; continue; }
    if (token === "--yes" || token === "-y") { interactive = false; continue; }

    if (token === "--profile" || token === "-p") {
      profile = rest[++i];
      if (!profile) throw new Error("Missing value for --profile");
      continue;
    }

    if (token === "--addon" || token === "-a") {
      const addon = rest[++i];
      if (!addon) throw new Error("Missing value for --addon");
      addons.push(addon);
      continue;
    }

    if (!token.startsWith("-") && !targetDir) {
      targetDir = token;
      continue;
    }

    throw new Error(`Unknown argument: ${token}`);
  }

  return { command: "init", targetDir, profile, force, interactive, addons };
}

// ─── Help ────────────────────────────────────────────────────────────────────

function printHelp() {
  printBanner();
  console.log(`${fmt.heading("Usage:")}`);
  console.log(`  agent-kit ${c.cyan}list${c.reset}                          List available profiles`);
  console.log(`  agent-kit ${c.cyan}init${c.reset} <dir> [options]           Scaffold a new project\n`);
  console.log(`${fmt.heading("Options:")}`);
  console.log(`  ${c.cyan}--profile, -p${c.reset} <id>    Profile template (default: interactive picker)`);
  console.log(`  ${c.cyan}--addon, -a${c.reset}   <id>    Include optional add-on (repeatable)`);
  console.log(`  ${c.cyan}--force, -f${c.reset}           Overwrite files in non-empty target`);
  console.log(`  ${c.cyan}--yes, -y${c.reset}             Skip interactive prompts (use defaults)\n`);
  console.log(`${fmt.heading("Examples:")}`);
  console.log(`  agent-kit init ./my-app`);
  console.log(`  agent-kit init ./api -p go-service -a data-intensive`);
  console.log(`  agent-kit init ./site -p typescript-react -y\n`);
}

// ─── List ────────────────────────────────────────────────────────────────────

function printList() {
  printBanner();
  console.log(`${fmt.heading("Available profiles:")}\n`);
  for (const p of PROFILES) {
    console.log(`  ${c.cyan}${c.bold}${p.id}${c.reset}`);
    console.log(`  ${p.title}`);
    console.log(`  ${c.dim}${p.summary}${c.reset}\n`);
  }

  console.log(`${fmt.heading("Optional add-ons:")}\n`);
  console.log(`  ${c.cyan}${c.bold}data-intensive${c.reset}    Postgres, NATS, Parquet, idempotent pipelines`);
  console.log(`  ${c.cyan}${c.bold}frontend-craft${c.reset}    Tailwind, shadcn/ui, Motion, visual polish\n`);
}

// ─── Interactive init ────────────────────────────────────────────────────────

async function interactiveInit(args: ParsedArgs): Promise<InitOptions> {
  printBanner();

  const prompt = await createPrompt();

  try {
    // 1. Target directory
    const targetDir = args.targetDir ?? await prompt.ask("Where should we create the project?", "./my-app");

    // 2. Profile
    const profileId = args.profile as ProfileId ?? await prompt.choose<ProfileId>(
      "What kind of project are you building?",
      PROFILES.map((p) => ({
        value: p.id,
        label: p.title,
        description: p.summary,
      }))
    );

    // 3. Add-ons
    const addons = args.addons.length > 0
      ? args.addons
      : await prompt.multiChoose("Which add-on instruction sets do you want?", [
          { value: "data-intensive", label: "Data-Intensive Systems — Postgres, NATS, Parquet, idempotent pipelines", defaultOn: true },
          { value: "frontend-craft", label: "Frontend Craft — Tailwind, shadcn/ui, Motion, visual polish", defaultOn: profileId === "typescript-react" },
        ]);

    // 4. Overwrite check
    const force = args.force || await prompt.confirm("Overwrite files if the directory isn't empty?", false);

    console.log("");

    return { targetDir, profileId, force, addons };
  } finally {
    prompt.close();
  }
}

// ─── Main ────────────────────────────────────────────────────────────────────

async function main() {
  const args = parseArgs(process.argv.slice(2));

  if (args.command === "help") {
    printHelp();
    return;
  }

  if (args.command === "list") {
    printList();
    return;
  }

  // Init command
  let options: InitOptions;

  if (args.interactive && !args.profile) {
    options = await interactiveInit(args);
  } else {
    // Non-interactive: require targetDir
    if (!args.targetDir) {
      throw new Error("Missing required <targetDir> for init command");
    }
    options = {
      targetDir: args.targetDir,
      profileId: args.profile ?? "typescript-react",
      force: args.force,
      addons: args.addons,
    };
  }

  const result = await scaffoldProject(options);

  const displayPath = (() => {
    const rel = path.relative(process.cwd(), result.outputPath);
    return rel.startsWith("..") ? result.outputPath : (rel || ".");
  })();

  printFileTree(result.createdFiles, result.outputPath);
  printDone(result.profile.title, displayPath);
}

main().catch((err: unknown) => {
  const message = err instanceof Error ? err.message : "Unknown error";
  console.error(`\n${fmt.error("Error:")} ${message}\n`);
  process.exitCode = 1;
});