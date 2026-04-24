"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { useState } from "react";
import { useForm } from "react-hook-form";
import modulesFixture from "../../../fixtures/modules.json";
import { type LaunchResult, submitLaunch } from "./actions";
import {
  type IdentityValues,
  type ModuleMeta,
  type ModulesValues,
  type SoulValues,
  type WizardState,
  buildLaunchParams,
  computeBitmap,
  identitySchema,
  soulSchema,
} from "./types";

const STEPS = ["Agent Identity", "Soul", "Modules", "Review", "Submit"] as const;
type Step = 0 | 1 | 2 | 3 | 4;

const FIXTURE_IMAGE_URI = "https://fixtures.titular.xyz/images/agent-placeholder.png";
const FIXTURE_SOUL_URI = "https://fixtures.titular.xyz/soul/sample-persona.json";

function FieldError({ message }: { message?: string }) {
  if (!message) return null;
  return (
    <p role='alert' aria-live='polite' style={{ color: "#c0392b", fontSize: 13, marginTop: 4 }}>
      {message}
    </p>
  );
}

function StepIdentity({
  defaultValues,
  onNext,
}: {
  defaultValues: Partial<IdentityValues>;
  onNext: (values: IdentityValues) => void;
}) {
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<IdentityValues>({
    resolver: zodResolver(identitySchema),
    defaultValues: { imageURI: FIXTURE_IMAGE_URI, ...defaultValues },
  });

  return (
    <form onSubmit={handleSubmit(onNext)} noValidate aria-label='Agent identity form'>
      <div style={{ marginBottom: 16 }}>
        <label htmlFor='name' style={{ display: "block", marginBottom: 4, fontWeight: 600 }}>
          Agent name <span aria-hidden='true'>*</span>
        </label>
        <input
          id='name'
          type='text'
          aria-required='true'
          aria-describedby={errors.name ? "name-error" : undefined}
          style={{ width: "100%", padding: "8px 12px", boxSizing: "border-box" }}
          {...register("name")}
        />
        <span id='name-error'>
          <FieldError message={errors.name?.message} />
        </span>
      </div>

      <div style={{ marginBottom: 16 }}>
        <label htmlFor='symbol' style={{ display: "block", marginBottom: 4, fontWeight: 600 }}>
          Token symbol <span aria-hidden='true'>*</span>
        </label>
        <input
          id='symbol'
          type='text'
          aria-required='true'
          aria-describedby={errors.symbol ? "symbol-error" : undefined}
          placeholder='e.g. MYAGENT'
          style={{ width: "100%", padding: "8px 12px", boxSizing: "border-box" }}
          {...register("symbol")}
        />
        <span id='symbol-error'>
          <FieldError message={errors.symbol?.message} />
        </span>
      </div>

      <div style={{ marginBottom: 16 }}>
        <label htmlFor='description' style={{ display: "block", marginBottom: 4, fontWeight: 600 }}>
          Description <span aria-hidden='true'>*</span>
        </label>
        <textarea
          id='description'
          aria-required='true'
          aria-describedby={errors.description ? "description-error" : undefined}
          rows={4}
          style={{ width: "100%", padding: "8px 12px", boxSizing: "border-box" }}
          {...register("description")}
        />
        <span id='description-error'>
          <FieldError message={errors.description?.message} />
        </span>
      </div>

      <div style={{ marginBottom: 24 }}>
        <label htmlFor='imageURI' style={{ display: "block", marginBottom: 4, fontWeight: 600 }}>
          Image URL <span aria-hidden='true'>*</span>
        </label>
        <input
          id='imageURI'
          type='url'
          aria-required='true'
          aria-describedby={errors.imageURI ? "imageURI-error" : undefined}
          style={{ width: "100%", padding: "8px 12px", boxSizing: "border-box" }}
          {...register("imageURI")}
        />
        <p style={{ fontSize: 12, color: "#666", marginTop: 4 }}>
          Fixture URL pre-filled. Real upload wired in a future iteration.
        </p>
        <span id='imageURI-error'>
          <FieldError message={errors.imageURI?.message} />
        </span>
      </div>

      <button type='submit' style={{ padding: "10px 24px", cursor: "pointer" }}>
        Next
      </button>
    </form>
  );
}

function StepSoul({
  defaultValues,
  onBack,
  onNext,
}: {
  defaultValues: Partial<SoulValues>;
  onBack: () => void;
  onNext: (values: SoulValues) => void;
}) {
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<SoulValues>({
    resolver: zodResolver(soulSchema),
    defaultValues: { soulURI: FIXTURE_SOUL_URI, ...defaultValues },
  });

  return (
    <form onSubmit={handleSubmit(onNext)} noValidate aria-label='Soul configuration form'>
      <div style={{ marginBottom: 16 }}>
        <label htmlFor='soulURI' style={{ display: "block", marginBottom: 4, fontWeight: 600 }}>
          Soul persona URI <span aria-hidden='true'>*</span>
        </label>
        <input
          id='soulURI'
          type='url'
          aria-required='true'
          aria-describedby={errors.soulURI ? "soulURI-error" : undefined}
          style={{ width: "100%", padding: "8px 12px", boxSizing: "border-box" }}
          {...register("soulURI")}
        />
        <p style={{ fontSize: 12, color: "#666", marginTop: 4 }}>
          Sample fixture: name, role, traits, systemPrompt
        </p>
        <span id='soulURI-error'>
          <FieldError message={errors.soulURI?.message} />
        </span>
      </div>

      <div style={{ display: "flex", gap: 12 }}>
        <button type='button' onClick={onBack} style={{ padding: "10px 24px", cursor: "pointer" }}>
          Back
        </button>
        <button type='submit' style={{ padding: "10px 24px", cursor: "pointer" }}>
          Next
        </button>
      </div>
    </form>
  );
}

function StepModules({
  defaultValues,
  onBack,
  onNext,
}: {
  defaultValues: ModulesValues;
  onBack: () => void;
  onNext: (values: ModulesValues) => void;
}) {
  const modules = modulesFixture as ModuleMeta[];
  const [enabled, setEnabled] = useState<Set<string>>(new Set(defaultValues.enabled));
  const [configs, setConfigs] = useState<Record<string, Record<string, string | number>>>(
    Object.fromEntries(
      modules.map((m) => [
        m.id,
        Object.fromEntries(
          m.configFields.map((f) => [f.key, defaultValues.configs[m.id]?.[f.key] ?? f.default])
        ),
      ])
    )
  );

  function toggleModule(id: string) {
    setEnabled((prev) => {
      const next = new Set(prev);
      if (next.has(id)) next.delete(id);
      else next.add(id);
      return next;
    });
  }

  function setConfigField(moduleId: string, key: string, value: string | number) {
    setConfigs((prev) => ({
      ...prev,
      [moduleId]: { ...prev[moduleId], [key]: value },
    }));
  }

  const bitmap = computeBitmap(Array.from(enabled));

  function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    onNext({
      enabled: Array.from(enabled),
      configs: Object.fromEntries(Array.from(enabled).map((id) => [id, configs[id] ?? {}])),
    });
  }

  return (
    <form onSubmit={handleSubmit} aria-label='Modules selection form'>
      <p style={{ marginBottom: 16, color: "#555" }}>
        Toggle modules to enable. Configure options for each enabled module.
      </p>

      <fieldset style={{ border: "none", padding: 0, margin: 0 }} aria-label='Module toggles'>
        {modules.map((mod) => {
          const isEnabled = enabled.has(mod.id);
          return (
            <div
              key={mod.id}
              style={{
                border: "1px solid",
                borderColor: isEnabled ? "#2563eb" : "#d1d5db",
                borderRadius: 8,
                padding: 16,
                marginBottom: 12,
                background: isEnabled ? "#eff6ff" : "#fff",
              }}
            >
              <div style={{ display: "flex", alignItems: "flex-start", gap: 12 }}>
                <input
                  type='checkbox'
                  id={`module-${mod.id}`}
                  checked={isEnabled}
                  onChange={() => toggleModule(mod.id)}
                  aria-label={`Enable ${mod.label} module`}
                  style={{ marginTop: 3, flexShrink: 0 }}
                />
                <div style={{ flex: 1 }}>
                  <label
                    htmlFor={`module-${mod.id}`}
                    style={{ fontWeight: 600, cursor: "pointer", display: "block" }}
                  >
                    {mod.label}
                  </label>
                  <p style={{ fontSize: 13, color: "#555", marginTop: 4 }}>{mod.description}</p>

                  {isEnabled && mod.configFields.length > 0 && (
                    <div style={{ marginTop: 12, paddingTop: 12, borderTop: "1px solid #dbeafe" }}>
                      {mod.configFields.map((field) => (
                        <div key={field.key} style={{ marginBottom: 10 }}>
                          <label
                            htmlFor={`${mod.id}-${field.key}`}
                            style={{ display: "block", fontSize: 13, marginBottom: 4 }}
                          >
                            {field.label}
                          </label>
                          <input
                            id={`${mod.id}-${field.key}`}
                            type={field.type === "address" ? "text" : "number"}
                            value={String(configs[mod.id]?.[field.key] ?? field.default)}
                            onChange={(e) => {
                              const raw = e.target.value;
                              const val =
                                field.type === "number" ? (raw === "" ? 0 : Number(raw)) : raw;
                              setConfigField(mod.id, field.key, val);
                            }}
                            style={{
                              padding: "6px 10px",
                              width: "100%",
                              boxSizing: "border-box",
                            }}
                          />
                        </div>
                      ))}
                    </div>
                  )}
                </div>
              </div>
            </div>
          );
        })}
      </fieldset>

      <output
        aria-live='polite'
        style={{
          display: "block",
          marginTop: 16,
          marginBottom: 16,
          padding: "10px 14px",
          background: "#f1f5f9",
          borderRadius: 6,
          fontFamily: "monospace",
          fontSize: 13,
        }}
      >
        Module bitmap: <strong>0b{bitmap.toString(2).padStart(7, "0")}</strong> ({bitmap})
      </output>

      <div style={{ display: "flex", gap: 12 }}>
        <button type='button' onClick={onBack} style={{ padding: "10px 24px", cursor: "pointer" }}>
          Back
        </button>
        <button type='submit' style={{ padding: "10px 24px", cursor: "pointer" }}>
          Next
        </button>
      </div>
    </form>
  );
}

function StepReview({
  state,
  onBack,
  onSubmit,
}: {
  state: WizardState;
  onBack: () => void;
  onSubmit: () => void;
}) {
  const params = buildLaunchParams(state);
  return (
    <div aria-label='Launch review'>
      <p style={{ marginBottom: 16 }}>
        Review your launch parameters before submitting. Estimated gas: <strong>350,000</strong>{" "}
        (fixture).
      </p>

      <pre
        aria-label='LaunchParams JSON preview'
        style={{
          background: "#0f172a",
          color: "#e2e8f0",
          padding: 20,
          borderRadius: 8,
          overflowX: "auto",
          fontSize: 13,
          lineHeight: 1.6,
        }}
      >
        {JSON.stringify(params, null, 2)}
      </pre>

      <div style={{ display: "flex", gap: 12, marginTop: 24 }}>
        <button type='button' onClick={onBack} style={{ padding: "10px 24px", cursor: "pointer" }}>
          Back
        </button>
        <button
          type='button'
          onClick={onSubmit}
          style={{
            padding: "10px 24px",
            cursor: "pointer",
            background: "#2563eb",
            color: "#fff",
            border: "none",
            borderRadius: 6,
          }}
        >
          Launch agent
        </button>
      </div>
    </div>
  );
}

function StepSubmit({ state, onBack }: { state: WizardState; onBack: () => void }) {
  const [result, setResult] = useState<LaunchResult | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [pending, setPending] = useState(false);

  async function handleLaunch() {
    setPending(true);
    setError(null);
    try {
      const params = buildLaunchParams(state);
      const res = await submitLaunch(params);
      setResult(res);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Unknown error");
    } finally {
      setPending(false);
    }
  }

  if (result) {
    return (
      <output aria-live='polite' aria-label='Launch result' style={{ display: "block" }}>
        <h2 style={{ color: "#16a34a" }}>Agent launched!</h2>
        <dl style={{ lineHeight: 2 }}>
          <dt style={{ fontWeight: 600 }}>Transaction hash</dt>
          <dd>
            <code style={{ wordBreak: "break-all" }}>{result.txHash}</code>
          </dd>
          <dt style={{ fontWeight: 600 }}>Agent ID</dt>
          <dd>{result.agentId}</dd>
          <dt style={{ fontWeight: 600 }}>Token address</dt>
          <dd>
            <code>{result.tokenAddress}</code>
          </dd>
          <dt style={{ fontWeight: 600 }}>Curve address</dt>
          <dd>
            <code>{result.curveAddress}</code>
          </dd>
          <dt style={{ fontWeight: 600 }}>Gas used</dt>
          <dd>{result.gasUsed.toLocaleString()}</dd>
        </dl>
        <a
          href={`/agent/${result.agentSlug}`}
          style={{
            display: "inline-block",
            marginTop: 16,
            padding: "10px 24px",
            background: "#2563eb",
            color: "#fff",
            textDecoration: "none",
            borderRadius: 6,
          }}
        >
          View agent page
        </a>
      </output>
    );
  }

  return (
    <div aria-label='Submit launch'>
      <p style={{ marginBottom: 16 }}>
        Ready to launch. This calls the server action (fixture stub; real chain tx in later
        iterations).
      </p>
      {error && (
        <p role='alert' style={{ color: "#c0392b", marginBottom: 12 }}>
          {error}
        </p>
      )}
      <div style={{ display: "flex", gap: 12 }}>
        <button
          type='button'
          onClick={onBack}
          disabled={pending}
          style={{ padding: "10px 24px", cursor: pending ? "not-allowed" : "pointer" }}
        >
          Back
        </button>
        <button
          type='button'
          onClick={handleLaunch}
          disabled={pending}
          aria-busy={pending}
          style={{
            padding: "10px 24px",
            cursor: pending ? "not-allowed" : "pointer",
            background: "#2563eb",
            color: "#fff",
            border: "none",
            borderRadius: 6,
            opacity: pending ? 0.7 : 1,
          }}
        >
          {pending ? "Launching..." : "Confirm launch"}
        </button>
      </div>
    </div>
  );
}

export default function LaunchpadCreatePage() {
  const [step, setStep] = useState<Step>(0);
  const [wizardState, setWizardState] = useState<WizardState>({
    identity: {},
    soul: {},
    modules: { enabled: [], configs: {} },
  });

  function updateIdentity(values: IdentityValues) {
    setWizardState((prev) => ({ ...prev, identity: values }));
    setStep(1);
  }

  function updateSoul(values: SoulValues) {
    setWizardState((prev) => ({ ...prev, soul: values }));
    setStep(2);
  }

  function updateModules(values: ModulesValues) {
    setWizardState((prev) => ({ ...prev, modules: values }));
    setStep(3);
  }

  return (
    <main style={{ maxWidth: 680, margin: "0 auto", padding: "40px 24px" }}>
      <h1 style={{ marginBottom: 8 }}>Launch your agent</h1>

      <nav aria-label='Wizard progress' style={{ marginBottom: 32 }}>
        <ol
          style={{
            display: "flex",
            gap: 8,
            listStyle: "none",
            padding: 0,
            margin: 0,
            flexWrap: "wrap",
          }}
        >
          {STEPS.map((label, idx) => (
            <li
              key={label}
              aria-current={idx === step ? "step" : undefined}
              style={{
                padding: "4px 12px",
                borderRadius: 20,
                fontSize: 13,
                fontWeight: idx === step ? 700 : 400,
                background: idx === step ? "#2563eb" : idx < step ? "#dbeafe" : "#f1f5f9",
                color: idx === step ? "#fff" : idx < step ? "#1d4ed8" : "#6b7280",
              }}
            >
              {idx + 1}. {label}
            </li>
          ))}
        </ol>
      </nav>

      {step === 0 && <StepIdentity defaultValues={wizardState.identity} onNext={updateIdentity} />}
      {step === 1 && (
        <StepSoul defaultValues={wizardState.soul} onBack={() => setStep(0)} onNext={updateSoul} />
      )}
      {step === 2 && (
        <StepModules
          defaultValues={wizardState.modules}
          onBack={() => setStep(1)}
          onNext={updateModules}
        />
      )}
      {step === 3 && (
        <StepReview state={wizardState} onBack={() => setStep(2)} onSubmit={() => setStep(4)} />
      )}
      {step === 4 && <StepSubmit state={wizardState} onBack={() => setStep(3)} />}
    </main>
  );
}
