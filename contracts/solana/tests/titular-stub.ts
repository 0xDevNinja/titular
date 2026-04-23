import * as anchor from "@coral-xyz/anchor";
import type { Program } from "@coral-xyz/anchor";
import { assert } from "chai";
import type { TitularStub } from "../target/types/titular_stub";

describe("titular-stub", () => {
  const provider = anchor.AnchorProvider.env();
  anchor.setProvider(provider);

  const program = anchor.workspace.TitularStub as Program<TitularStub>;

  it("calls initialize without error", async () => {
    const tx = await program.methods.initialize().rpc();
    assert.isString(tx, "expected transaction signature");
    console.log("initialize tx:", tx);
  });
});
