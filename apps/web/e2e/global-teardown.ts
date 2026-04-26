import { stopAnvil } from "./fixtures/anvil.js";

export default async function globalTeardown(): Promise<void> {
  stopAnvil();
}
