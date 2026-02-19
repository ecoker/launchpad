import { createInterface } from "node:readline/promises";
import { stdin, stdout } from "node:process";

// ─── ANSI helpers (zero deps) ────────────────────────────────────────────────

const esc = (code: string) => `\x1b[${code}m`;

export const c = {
  reset: esc("0"),
  bold: esc("1"),
  dim: esc("2"),
  italic: esc("3"),
  underline: esc("4"),

  // Colors
  cyan: esc("36"),
  green: esc("32"),
  yellow: esc("33"),
  magenta: esc("35"),
  red: esc("31"),
  blue: esc("34"),
  white: esc("37"),
  gray: esc("90"),
} as const;

export const fmt = {
  heading: (text: string) => `${c.bold}${c.cyan}${text}${c.reset}`,
  success: (text: string) => `${c.green}${text}${c.reset}`,
  warn: (text: string) => `${c.yellow}${text}${c.reset}`,
  error: (text: string) => `${c.red}${c.bold}${text}${c.reset}`,
  accent: (text: string) => `${c.magenta}${text}${c.reset}`,
  dim: (text: string) => `${c.dim}${text}${c.reset}`,
  label: (text: string) => `${c.bold}${text}${c.reset}`,
  file: (text: string) => `${c.blue}${c.underline}${text}${c.reset}`,
};

// ─── ASCII art banner ────────────────────────────────────────────────────────

export function printBanner() {
  console.log("");
  console.log(
    `${c.cyan}${c.bold}   ┌─────────────────────────────────────────┐${c.reset}`
  );
  console.log(
    `${c.cyan}${c.bold}   │         ${c.magenta}⚡ agent-kit ${c.cyan}                    │${c.reset}`
  );
  console.log(
    `${c.cyan}${c.bold}   │   ${c.reset}${c.dim}opinionated AI instruction scaffolder${c.cyan}${c.bold}  │${c.reset}`
  );
  console.log(
    `${c.cyan}${c.bold}   └─────────────────────────────────────────┘${c.reset}`
  );
  console.log("");
}

// ─── Interactive prompts ─────────────────────────────────────────────────────

export async function createPrompt() {
  const rl = createInterface({ input: stdin, output: stdout });

  const ask = async (question: string, fallback?: string): Promise<string> => {
    const suffix = fallback ? ` ${c.dim}(${fallback})${c.reset}` : "";
    const answer = (await rl.question(`${c.bold}${question}${suffix} ${c.cyan}› ${c.reset}`)).trim();
    return answer || fallback || "";
  };

  const confirm = async (question: string, defaultYes = true): Promise<boolean> => {
    const hint = defaultYes ? "Y/n" : "y/N";
    const answer = await ask(question, hint);
    if (answer === hint) return defaultYes;
    return ["y", "yes"].includes(answer.toLowerCase());
  };

  const choose = async <T extends string>(
    question: string,
    options: { value: T; label: string; description: string }[]
  ): Promise<T> => {
    console.log(`\n${c.bold}${question}${c.reset}\n`);
    options.forEach((opt, i) => {
      console.log(`  ${c.cyan}${c.bold}${i + 1}${c.reset}  ${fmt.label(opt.label)}`);
      console.log(`     ${c.dim}${opt.description}${c.reset}`);
    });
    console.log("");

    while (true) {
      const answer = await ask(`Pick a number (1-${options.length})`);
      const index = parseInt(answer, 10) - 1;
      if (index >= 0 && index < options.length) {
        console.log(`     ${fmt.success("✔")} ${options[index].label}\n`);
        return options[index].value;
      }
      console.log(fmt.warn("  Please enter a valid number."));
    }
  };

  const multiChoose = async (
    question: string,
    options: { value: string; label: string; defaultOn: boolean }[]
  ): Promise<string[]> => {
    console.log(`\n${c.bold}${question}${c.reset}\n`);
    options.forEach((opt, i) => {
      const tag = opt.defaultOn ? fmt.success(" (included by default)") : "";
      console.log(`  ${c.cyan}${c.bold}${i + 1}${c.reset}  ${opt.label}${tag}`);
    });
    console.log("");

    const answer = await ask("Enter numbers to toggle, comma-separated", "all");
    if (answer === "all") {
      return options.map((o) => o.value);
    }

    const selected = new Set(
      options.filter((o) => o.defaultOn).map((o) => o.value)
    );

    for (const token of answer.split(",")) {
      const index = parseInt(token.trim(), 10) - 1;
      if (index >= 0 && index < options.length) {
        const value = options[index].value;
        if (selected.has(value)) {
          selected.delete(value);
        } else {
          selected.add(value);
        }
      }
    }

    return [...selected];
  };

  const close = () => rl.close();

  return { ask, confirm, choose, multiChoose, close };
}

// ─── Summary printer ─────────────────────────────────────────────────────────

export function printFileTree(files: string[], rootDir: string) {
  console.log(`\n${fmt.heading("Created files:")}\n`);
  const sorted = [...files].sort();
  for (const file of sorted) {
    const relative = file.replace(rootDir + "/", "");
    console.log(`  ${c.dim}└─${c.reset} ${fmt.file(relative)}`);
  }
  console.log("");
}

export function printDone(profileLabel: string, targetDir: string) {
  console.log(
    `${fmt.success("✔")} Scaffolded ${fmt.accent(profileLabel)} into ${fmt.file(targetDir)}`
  );
  console.log("");
  console.log(`${fmt.heading("Next steps:")}`);
  console.log(`  ${c.dim}1.${c.reset} cd ${fmt.file(targetDir)}`);
  console.log(
    `  ${c.dim}2.${c.reset} Review ${fmt.file(".github/copilot-instructions.md")} — your always-on standards`
  );
  console.log(
    `  ${c.dim}3.${c.reset} Browse ${fmt.file(".github/instructions/")} — language & framework rules`
  );
  console.log(
    `  ${c.dim}4.${c.reset} Edit freely — these are ${c.italic}your${c.reset} opinions now`
  );
  console.log("");
  console.log(`${c.dim}Happy building. Write something beautiful. ✨${c.reset}`);
  console.log("");
}
