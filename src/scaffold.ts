import { mkdir } from "node:fs/promises";
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
};

export async function scaffoldProject(options: InitOptions): Promise<{ profile: Profile; outputPath: string }> {
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
    "{{PROJECT_NAME}}": path.basename(outputPath)
  };

  await copyDirectory(path.join(templateRoot, "core"), outputPath, {
    overwrite: options.force,
    replacements
  });

  await copyDirectory(path.join(templateRoot, "profiles", profile.id), outputPath, {
    overwrite: options.force,
    replacements
  });

  return { profile, outputPath };
}