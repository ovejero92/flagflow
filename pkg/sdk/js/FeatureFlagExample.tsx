import { useEffect, useState } from "react";
import { FeatureFlagClient } from "./feature-flags-sdk";

import { DEFAULT_FLAG_SERVICE_URL } from "./feature-flags-sdk";

const client = new FeatureFlagClient(
  process.env.NEXT_PUBLIC_APP_ID ?? "your-app-uuid-here",
  process.env.NEXT_PUBLIC_FLAG_SERVICE_URL ?? DEFAULT_FLAG_SERVICE_URL
);

export function FeatureFlagExample() {
  const [darkMode, setDarkMode] = useState(false);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const userId = "user-123";

    client
      .isEnabled("dark-mode", userId)
      .then(setDarkMode)
      .catch(console.error)
      .finally(() => setLoading(false));
  }, []);

  if (loading) {
    return <p>Cargando feature flags...</p>;
  }

  return (
    <div style={{ padding: 16, background: darkMode ? "#111" : "#f5f5f5", color: darkMode ? "#fff" : "#111" }}>
      <h2>Feature Flag Demo</h2>
      <p>
        Dark mode está{" "}
        <strong>{darkMode ? "activado" : "desactivado"}</strong> para este usuario.
      </p>
    </div>
  );
}
