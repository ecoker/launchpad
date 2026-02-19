import { mkdir, readdir } from "node:fs/promises";
import path from "node:path";
import { fileURLToPath } from "node:url";
import { copyDirectory, exists, isDirectoryEmpty } from "./fs-utils.js";
import { findProfile, type Profile } from "./profiles.js";

const thisFilePath = fileURLToPath(import.meta.url);
const thisDir = path.dirname(thisFilePath);
const templateRoot = path.resolve(thisDir, "..", "templates");

export type InitOptions = {
  targetDir: string;
  profileId: string;
  force: boolean;
  addons: string[];
};

type ScaffoldResult = {
  profile: Profile;
  outputPath: string;
  createdFiles: string[];
};

async function collectFiles(dir: string, base = dir): Promise<string[]> {
  if (!(await exists(dir))) return [];
  const entries = await readdir(dir, { withFileTypes: true });
  const files: string[] = [];
  for (const entry of entries) {
    const full = path.join(dir, entry.name);
    if (entry.isDirectory()) {
      files.push(...(await collectFiles(full, base)));
    } else {
      files.push(full);
    }
  }
  return files;
}

const ADDON_TEMPLATES: Record<string, string> = {
  "data-intensive": "addons/data-intensive",
  "frontend-craft": "addons/frontend-craft",
};

export async function scaffoldProject(options: InitOptions): Promise<ScaffoldResult> {
  const profile = findProfile(options.profileId);
  if (!profile) {
    throw new Error(`Unknown profile '${options.profileId}'. Run 'agent-kit list' to see available profiles.`);
  }

  const outputPath = path.resolve(process.cwd(), options.targetDir);
  const outputExists = await exists(outputPath);
  if (outputExists && !(await isDirectoryEmpty(outputPath)) && !options.force) {
    throw new Error(`Target directory is not empty: ${outputPath}. Re-run with --force to overwrite files.`);
  }

  await mkdir(outputPath, { recursive: true });

  const replacements = {
    "{{PROJECT_NAME}}": path.basename(outputPath),
  };

  const copyOpts = { overwrite: options.force, replacements };

  // 1. Core templates (always)
  await copyDirectory(path.join(templateRoot, "core"), outputPath, copyOpts);

  // 2. Profile templates
  await copyDirectory(path.join(templateRoot, "profiles", profile.id), outputPath, copyOpts);

  // 3. Add-on templates
  for (const addon of options.addons) {
    const addonSubdir = ADDON_TEMPLATES[addon];
    if (!addonSubdir) continue;
    const addonPath = path.join(templateRoot, addonSubdir);
    if (await exists(addonPath)) {
      await copyDirectory(addonPath, outputPath, copyOpts);
    }
  }

  const createdFiles = await collectFiles(outputPath);

  return { profile, outputPath, createdFiles };
}