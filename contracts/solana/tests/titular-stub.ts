import * as anchor from "@coral-xyz/anchor";
import { Program } from "@coral-xyz/anchor";
import { TitularStub } from "../target/types/titular_stub";
import { assert } from "chai";

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
