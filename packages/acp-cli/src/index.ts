// Public re-exports for `@titular/acp-cli`. Programmatic consumers can use
// the same command runners the binary uses (e.g. `runConfigure`) without
// shelling out — useful in dashboards and integration tests.

export { runConfigure, ConfigureError } from "./commands/configure.js";
export type { ConfigureFlags, ConfigureDeps } from "./commands/configure.js";

export {
  runAgentCreate,
  runAgentList,
  runAgentStatus,
  AgentError,
  emitAgentsPage,
} from "./commands/agent.js";
export type {
  AgentCreateFlags,
  AgentListFlags,
  AgentStatusFlags,
  AgentDeps,
} from "./commands/agent.js";

export { runBrowse } from "./commands/browse.js";
export type { BrowseFlags, BrowseDeps } from "./commands/browse.js";

export { runJobCreate, runJobList, JobError } from "./commands/job.js";
export type { JobCreateFlags, JobListFlags, JobDeps } from "./commands/job.js";

export {
  runEventsListen,
  EventsError,
  parseEventBlock,
} from "./commands/events.js";
export type {
  EventsListenFlags,
  EventsDeps,
  SseEvent,
} from "./commands/events.js";

export {
  defaultConfigPath,
  isSupportedChain,
  loadConfig,
  mergeConfig,
  redactConfig,
  saveConfig,
  validateConfig,
  SUPPORTED_CHAINS,
} from "./config.js";

export {
  buildAcpAgent,
  buildGatewayClient,
  resolveAlchemyChain,
  resolveContracts,
} from "./sdk.js";
export type { ProviderFactory } from "./sdk.js";

export type { CliConfig, GlobalFlags, SupportedChain } from "./types.js";
