// deploy.ts — Anchor migration script (runs on `anchor migrate`).
import * as anchor from "@coral-xyz/anchor";

module.exports = async (provider: anchor.AnchorProvider) => {
  anchor.setProvider(provider);
  // Deploy logic here once programs are production-ready.
  console.log("migration: no-op for stub workspace");
};
