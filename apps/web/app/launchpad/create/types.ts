import { z } from "zod";
import modulesFixture from "../../../fixtures/modules.json";

// ─── ERC-20 symbol regex ────────────────────────────────────────────────────
const ERC20_SYMBOL_RE = /^[A-Z0-9]{1,11}$/;

// ─── Step 1 — Agent identity ─────────────────────────────────────────────────
export const identitySchema = z.object({
  name: z.string().min(1, "Name is required").max(64, "Max 64 characters"),
  symbol: z
    .string()
    .min(1, "Symbol is required")
    .max(11, "Max 11 characters")
    .regex(ERC20_SYMBOL_RE, "Symbol must be 1-11 uppercase letters/digits"),
  description: z.string().min(1, "Description is required").max(1000, "Max 1000 characters"),
  imageURI: z.string().url("Must be a valid URL"),
});
export type IdentityValues = z.infer<typeof identitySchema>;

// ─── Step 2 — Soul ───────────────────────────────────────────────────────────
export const soulSchema = z.object({
  soulURI: z.string().url("Must be a valid URL"),
});
export type SoulValues = z.infer<typeof soulSchema>;

// ─── Module metadata (from fixture) ──────────────────────────────────────────
export type ModuleMeta = (typeof modulesFixture)[number];

// ─── Step 3 — Modules ────────────────────────────────────────────────────────
export const moduleConfigSchema = z.record(z.string(), z.union([z.string(), z.number()]));

export const modulesSchema = z.object({
  enabled: z.array(z.string()),
  configs: z.record(z.string(), moduleConfigSchema),
});
export type ModulesValues = z.infer<typeof modulesSchema>;

// ─── Derived LaunchParams (mirrors LaunchpadFactory.LaunchParams) ─────────────
export interface LaunchParams {
  name: string;
  symbol: string;
  imageURI: string;
  soulURI: string;
  /** Bitmap: bit N set if module with bit=N is enabled */
  moduleFlags: number;
  /** ABI-encoded config bytes per enabled module (fixture: JSON strings) */
  moduleConfigs: string[];
}

// ─── Wizard state (all steps combined) ───────────────────────────────────────
export interface WizardState {
  identity: Partial<IdentityValues>;
  soul: Partial<SoulValues>;
  modules: ModulesValues;
}

export function computeBitmap(enabledIds: string[]): number {
  let flags = 0;
  for (const id of enabledIds) {
    const meta = (modulesFixture as ModuleMeta[]).find((m) => m.id === id);
    if (meta) flags |= 1 << meta.bit;
  }
  return flags;
}

export function buildLaunchParams(state: WizardState): LaunchParams {
  const { identity, soul, modules } = state;
  const moduleFlags = computeBitmap(modules.enabled);
  return {
    name: identity.name ?? "",
    symbol: identity.symbol ?? "",
    imageURI: identity.imageURI ?? "",
    soulURI: soul.soulURI ?? "",
    moduleFlags,
    moduleConfigs: modules.enabled.map((id) => {
      const cfg = modules.configs[id] ?? {};
      return JSON.stringify(cfg);
    }),
  };
}
