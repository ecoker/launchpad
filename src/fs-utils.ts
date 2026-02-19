import { constants } from "node:fs";
import { access, copyFile, mkdir, readdir, readFile, stat, writeFile } from "node:fs/promises";
import path from "node:path";

export async function exists(targetPath: string): Promise<boolean> {
  try {
    await access(targetPath, constants.F_OK);
    return true;
  } catch {
    return false;
  }
}

export async function isDirectory(targetPath: string): Promise<boolean> {
  if (!(await exists(targetPath))) {
    return false;
  }

  const targetStat = await stat(targetPath);
  return targetStat.isDirectory();
}

export async function isDirectoryEmpty(targetPath: string): Promise<boolean> {
  if (!(await isDirectory(targetPath))) {
    return true;
  }

  const entries = await readdir(targetPath);
  return entries.length === 0;
}

async function copyTextWithReplacements(sourcePath: string, destinationPath: string, replacements: Record<string, string>) {
  const text = await readFile(sourcePath, "utf8");
  const replaced = Object.entries(replacements).reduce((accumulator, [token, value]) => {
    return accumulator.replaceAll(token, value);
  }, text);
  await writeFile(destinationPath, replaced, "utf8");
}

function shouldTreatAsText(filePath: string): boolean {
  const extension = path.extname(filePath).toLowerCase();
  return [".md", ".txt", ".json", ".yaml", ".yml", ".toml", ".env", ".ts", ".js"].includes(extension) || extension === "";
}

export async function copyDirectory(
  sourceDir: string,
  destinationDir: string,
  options: { overwrite: boolean; replacements: Record<string, string> }
): Promise<void> {
  await mkdir(destinationDir, { recursive: true });
  const entries = await readdir(sourceDir, { withFileTypes: true });

  for (const entry of entries) {
    const sourcePath = path.join(sourceDir, entry.name);
    const destinationPath = path.join(destinationDir, entry.name);

    if (entry.isDirectory()) {
      await copyDirectory(sourcePath, destinationPath, options);
      continue;
    }

    if (!options.overwrite && (await exists(destinationPath))) {
      throw new Error(`Refusing to overwrite existing file: ${destinationPath}`);
    }

    if (shouldTreatAsText(sourcePath)) {
      await copyTextWithReplacements(sourcePath, destinationPath, options.replacements);
      continue;
    }

    await copyFile(sourcePath, destinationPath);
  }
}